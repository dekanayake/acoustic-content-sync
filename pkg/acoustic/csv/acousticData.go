package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/wesovilabs/koazee"
)

func convert(acousticField string, configTypeMapping *ContentTypeMapping, dataRow DataRow) (api.GenericData, error) {
	acousticFieldMapping, err := configTypeMapping.GetFieldMappingByAcousticField(acousticField)
	if err != nil {
		return api.GenericData{}, errors.ErrorWithStack(err)
	}
	data, err := acousticFieldMapping.ConvertToGenericData(dataRow, configTypeMapping)
	if err != nil {
		return data, errors.ErrorWithStack(err)
	}
	return data, nil
}

func TransformContent(contentType string, dataFeedPath string, configPath string) ([]api.AcousticDataRecord, error) {
	config, err := InitContentTypeMappingConfig(configPath)
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
	acousticFields := configTypeMapping.GetAcousticFields()
	acousticDataList := make([]api.AcousticDataRecord, 0, dataFeed.RecordCount())
	for ok := true; ok; ok = dataFeed.HasNext() {
		dataRow := dataFeed.Next()
		acousticDataOut := koazee.StreamOf(acousticFields).
			Map(func(acousticField string) (api.GenericData, error) {
				return convert(acousticField, configTypeMapping, dataRow)
			}).Do().Out()

		err := acousticDataOut.Err()
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		acousticData := acousticDataOut.Val().([]api.GenericData)

		searchValues := make(map[string]string)
		for _, searchKey := range configTypeMapping.SearchKeys {
			for _, acousticDataItem := range acousticData {
				if acousticDataItem.Name == searchKey {
					searchValues[searchKey] = acousticDataItem.Value.(api.AcousticValue).Value
				}
			}
		}

		acousticDataList = append(acousticDataList, api.AcousticDataRecord{
			Values:                 acousticData,
			NameFields:             configTypeMapping.Name,
			Tags:                   configTypeMapping.Tags,
			Update:                 configTypeMapping.Update,
			CreateNonExistingItems: configTypeMapping.CreateNonExistingItems,
			SearchTerm:             configTypeMapping.SearchTerm,
			SearchTerms:            configTypeMapping.SearchTerms,
			SearchOnLibrary:        configTypeMapping.SearchOnLibrary,
			SearchOnDeliveryAPI:    configTypeMapping.SearchOnDeliveryAPI,
			SearchValues:           searchValues,
			SearchKeys:             configTypeMapping.SearchKeys,
			SearchType:             configTypeMapping.SearchType,
			CSVRecordKey:           configTypeMapping.CsvRecordKey,
			FilterRecords:          configTypeMapping.FilterRecords,
			FilterFileLocation:     configTypeMapping.FilterFileLocation,
			FilterType:             configTypeMapping.FilterType,
			FilterColumns:          configTypeMapping.FilterColumns,
		})
	}
	return acousticDataList, nil
}

func TransformSite(contentType string, dataFeedPath string, configPath string) ([]api.AcousticDataRecord, error) {
	config, err := InitContentTypeMappingConfig(configPath)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	siteMapping, err := config.GetSiteMapping(contentType)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	dataFeed, err := LoadCSV(dataFeedPath)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	acousticFields := siteMapping.GetAcousticFields()
	acousticDataList := make([]api.AcousticDataRecord, 0, dataFeed.RecordCount())
	for ok := true; ok; ok = dataFeed.HasNext() {
		dataRow := dataFeed.Next()
		acousticDataOut := koazee.StreamOf(acousticFields).
			Map(func(acousticField string) (api.GenericData, error) {
				return convert(acousticField, &siteMapping.ContentTypeMapping, dataRow)
			}).Do().Out()

		err := acousticDataOut.Err()
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		acousticData := acousticDataOut.Val().([]api.GenericData)

		searchValues := make(map[string]string)
		for _, searchKey := range siteMapping.SearchKeys {
			for _, acousticDataItem := range acousticData {
				if acousticDataItem.Name == searchKey {
					searchValues[searchKey] = acousticDataItem.Value.(api.AcousticValue).Value
				}
			}
		}

		acousticDataList = append(acousticDataList, api.AcousticDataRecord{
			Values:                 acousticData,
			NameFields:             siteMapping.Name,
			Tags:                   siteMapping.Tags,
			Update:                 siteMapping.Update,
			CreateNonExistingItems: siteMapping.CreateNonExistingItems,
			SearchTerm:             siteMapping.SearchTerm,
			SearchTerms:            siteMapping.SearchTerms,
			SearchOnLibrary:        siteMapping.SearchOnLibrary,
			SearchOnDeliveryAPI:    siteMapping.SearchOnDeliveryAPI,
			SearchValues:           searchValues,
			SearchKeys:             siteMapping.SearchKeys,
			SearchType:             siteMapping.SearchType,
			CSVRecordKey:           siteMapping.CsvRecordKey,
			FilterRecords:          siteMapping.FilterRecords,
			FilterFileLocation:     siteMapping.FilterFileLocation,
			FilterType:             siteMapping.FilterType,
			FilterColumns:          siteMapping.FilterColumns,
			SiteConfig: api.SiteConfig{
				DontCreatePageIfExist: siteMapping.DontCreatePageIfExist,
				UpdatePageIfExists:    siteMapping.UpdatePageIfExist,
			},
		})
	}
	return acousticDataList, nil
}
