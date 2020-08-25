package csv

import (
	"encoding/csv"
	"errors"
	"github.com/wesovilabs/koazee"
	"github.com/wesovilabs/koazee/stream"
	"io"
	"os"
)

type contentData struct {
	rows []dataRow
}

type dataRow struct {
	columns map[string]string
}

func  load(csvFile *os.File) (*contentData,error){
	records := csv.NewReader(csvFile)
	headerRecord, err := records.Read()
	if err == io.EOF {
		return nil,errors.New("CSV file is empty. ")
	}
	if err != nil {
		return nil,err
	}
	if headerRecord != nil {
		headersMap := make(map[int]string)
		for index, columnHeader := range headerRecord {
			headersMap[index] = columnHeader
		}
		dataRows := make([]dataRow, 0,len(headersMap))
		for {
			contentRecord, err := records.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil,err
			}
			row := make(map[string]string)
			for index, columnValue := range contentRecord {
				row[headersMap[index]] = columnValue
			}
			dataRow := dataRow{columns:row}
			dataRows = append(dataRows, dataRow)
		}
		contentData := &contentData{rows : dataRows}
		return contentData,nil
	} else {
		return nil, errors.New("No error nor record found!")
	}
}

type DataFeed interface {
	HasNext() bool
	Next() DataRow
	Headers() ([]string,error)
	RecordCount() int
}

type DataRow interface {
	Get(columnName string) (string,error)
}

type dataFeed struct {
	rowIndex   int
	rowSize    int
	rows []dataRow
	rowsStream stream.Stream
}

func  LoadCSV(csvFilePath string) (DataFeed,error) {
	csvFile, err := os.Open(csvFilePath)
	defer csvFile.Close()
	if err != nil {
		return nil,err
	} else {
		if contentData,err := load(csvFile) ; err != nil {
			return nil,err
		} else {
			return &dataFeed{
						rowIndex:   0,
						rowSize:    len(contentData.rows),
						rowsStream: koazee.StreamOf(contentData.rows),
						rows : contentData.rows,
				},nil
		}
	}
}

func (dataFeed *dataFeed) RecordCount() int{
	return dataFeed.rowSize
}

func (dataFeed *dataFeed) HasNext() bool{
	return dataFeed.rowIndex < dataFeed.rowSize
}

func (dataFeed *dataFeed) Next() DataRow{
	val,remainingRows :=  dataFeed.rowsStream.Pop()
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
		return keys,nil
	} else {
		return nil, errors.New("No contents in datafeed")
	}
}

func (dataRow *dataRow) Get(columnName string) (string,error) {
		if _, ok := dataRow.columns[columnName]; !ok {
			return "",errors.New("No value found for column name :" + columnName)
		} else {
			return dataRow.columns[columnName],nil
		}
}