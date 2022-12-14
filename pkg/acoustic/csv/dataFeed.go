package csv

import (
	"encoding/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/dimchansky/utfbom"
	"github.com/wesovilabs/koazee"
	"github.com/wesovilabs/koazee/stream"
	"io"
	"os"
	"strings"
)

type contentData struct {
	rows []dataRow
}

type dataRow struct {
	columns map[string]string
}

func load(csvFile *os.File) (*contentData, error) {
	reader, _ := utfbom.Skip(csvFile)
	records := csv.NewReader(reader)

	headerRecord, err := records.Read()
	var unparsedRecordsWriter *csv.Writer = nil
	if env.WriteUnParsedRecordsToCSV() {
		unparsedRecordsFile, err := os.Create("unparsed_records.csv")
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		defer unparsedRecordsFile.Close()
		unparsedRecordsWriter = csv.NewWriter(unparsedRecordsFile)
		defer unparsedRecordsWriter.Flush()
		if err := unparsedRecordsWriter.Write(headerRecord); err != nil {
			return nil, errors.ErrorWithStack(err)
		}
	}
	if err == io.EOF {
		return nil, errors.ErrorMessageWithStack("CSV file is empty. ")
	}
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	if headerRecord != nil {
		headersMap := make(map[int]string)
		for index, columnHeader := range headerRecord {
			headersMap[index] = columnHeader
		}
		dataRows := make([]dataRow, 0, len(headersMap))
		for {
			contentRecord, err := records.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				if env.WriteUnParsedRecordsToCSV() {
					if err := unparsedRecordsWriter.Write(contentRecord); err != nil {
						return nil, errors.ErrorWithStack(err)
					}
					continue
				} else {
					return nil, errors.ErrorWithStack(err)
				}
			}
			row := make(map[string]string)
			for index, columnValue := range contentRecord {
				row[headersMap[index]] = columnValue
			}
			dataRow := dataRow{columns: row}
			dataRows = append(dataRows, dataRow)
		}
		contentData := &contentData{rows: dataRows}
		return contentData, nil
	} else {
		return nil, errors.ErrorMessageWithStack("No error nor record found!")
	}
}

type DataFeed interface {
	HasNext() bool
	Next() DataRow
	Headers() ([]string, error)
	RecordCount() int
}

type DataRow interface {
	Get(columnName string) (string, error)
}

type dataFeed struct {
	rowIndex   int
	rowSize    int
	rows       []dataRow
	rowsStream stream.Stream
}

func LoadCSV(csvFilePath string) (DataFeed, error) {
	csvFile, err := os.Open(csvFilePath)
	defer csvFile.Close()
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	} else {
		if contentData, err := load(csvFile); err != nil {
			return nil, err
		} else {
			return &dataFeed{
				rowIndex:   0,
				rowSize:    len(contentData.rows),
				rowsStream: koazee.StreamOf(contentData.rows),
				rows:       contentData.rows,
			}, nil
		}
	}
}

func (dataFeed *dataFeed) RecordCount() int {
	return dataFeed.rowSize
}

func (dataFeed *dataFeed) HasNext() bool {
	return dataFeed.rowIndex < dataFeed.rowSize
}

func (dataFeed *dataFeed) Next() DataRow {
	val, remainingRows := dataFeed.rowsStream.Pop()
	dataFeed.rowsStream = remainingRows
	dataFeed.rowIndex += 1
	dataRow := val.Val().(dataRow)
	return &dataRow
}

func (dataFeed *dataFeed) Headers() ([]string, error) {
	dataRows := dataFeed.rows
	if len(dataRows) > 0 {
		columns := dataRows[0].columns
		keys := make([]string, 0, len(columns))
		for k := range columns {
			keys = append(keys, k)
		}
		return keys, nil
	} else {
		return nil, errors.ErrorMessageWithStack("No contents in datafeed")
	}
}

func (dataRow *dataRow) Get(columnName string) (string, error) {
	if _, ok := dataRow.columns[columnName]; !ok {
		return "", errors.ErrorMessageWithStack("No value found for column name :" + columnName)
	} else {
		return strings.TrimSpace(dataRow.columns[columnName]), nil
	}
}
