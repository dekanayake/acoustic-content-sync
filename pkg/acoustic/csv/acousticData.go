package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/wesovilabs/koazee"
)

func convert(columnHeader string, configTypeMapping *ContentTypeMapping, dataRow DataRow) (api.GenericData, error) {
	data := api.GenericData{}
	acousticFieldMapping, err := configTypeMapping.GetFieldMapping(columnHeader)
	if err != nil {
		return data, errors.ErrorWithStack(err)
	}
	data.Name = acousticFieldMapping.AcousticProperty
	data.Type = acousticFieldMapping.PropertyType
	value, err := dataRow.Get(columnHeader)
	if err != nil {
		return data, errors.ErrorWithStack(err)
	}
	data.Value = acousticFieldMapping.Value(value)
	context, err := acousticFieldMapping.Context(dataRow, configTypeMapping)
	if err != nil {
		return api.GenericData{}, errors.ErrorWithStack(err)
	}
	data.Context = context
	return data, nil
}

func Transform(contentType string, dataFeedPath string, configPath string) ([]api.AcousticDataRecord, error) {
	config, err := InitConfig(configPath)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	configTypeMapping, err := config.GetContentType(contentType)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	dataFeed, err := LoadCSV(dataFeedPath)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	columnHeaders, err := dataFeed.Headers()
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	acousticDataList := make([]api.AcousticDataRecord, 0, dataFeed.RecordCount())
	for ok := true; ok; ok = dataFeed.HasNext() {
		dataRow := dataFeed.Next()
		acousticDataOut := koazee.StreamOf(columnHeaders).
			Map(func(columnHeader string) (api.GenericData, error) {
				return convert(columnHeader, configTypeMapping, dataRow)
			}).Do().Out()

		err := acousticDataOut.Err()
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		acousticData := acousticDataOut.Val().([]api.GenericData)
		acousticDataList = append(acousticDataList, api.AcousticDataRecord{
			Values:       acousticData,
			NameFields:   configTypeMapping.Name,
			Tags:         configTypeMapping.Tags,
			CSVRecordKey: configTypeMapping.CsvRecordKey,
		})
	}
	return acousticDataList, nil
}
