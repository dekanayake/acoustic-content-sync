package csv

import "github.com/wesovilabs/koazee"

type AcousticDataRecord struct {
	values []AcousticData
}

type AcousticData struct {
	Name string
	Type string
	Value string
}

func convert(columnHeader string, configTypeMapping *ContentTypeMapping, dataRow DataRow) (AcousticData, error) {
	data := AcousticData{}
	acousticFieldMapping, err := configTypeMapping.GetFieldMapping(columnHeader)
	if err != nil {
		return data, err
	}
	data.Name = acousticFieldMapping.AcousticProperty
	data.Type = acousticFieldMapping.PropertyType
	data.Value, err = dataRow.Get(columnHeader)
	if err != nil {
		return data, err
	}
	return data, nil
}

func  Transform(contentType string, dataFeedPath string, configPath string) (*[]AcousticDataRecord,error) {
	config, err := InitConfig(configPath)
	if err != nil {
		return nil,err
	}
	configTypeMapping, err := config.Get(contentType)
	if err != nil {
		return nil,err
	}
	dataFeed, err := LoadCSV(dataFeedPath)
	if err != nil {
		return nil,err
	}
	columnHeaders, err := dataFeed.Headers()
	if err != nil {
		return nil,err
	}
	acousticDataList := make([]AcousticDataRecord, 0, dataFeed.RecordCount())
	for ok := true; ok; ok = dataFeed.HasNext() {
		dataRow := dataFeed.Next()
		acousticDataOut := koazee.StreamOf(columnHeaders).
			Map(func (columnHeader string) (AcousticData,error){
				return convert(columnHeader, configTypeMapping, dataRow)
			}).Do().Out()

		err := acousticDataOut.Err()
		if err != nil {
			return nil,err
		}
		acousticData := acousticDataOut.Val().([]AcousticData)
		acousticDataList = append(acousticDataList, AcousticDataRecord{values: acousticData})
	}
	return &acousticDataList,nil
}
