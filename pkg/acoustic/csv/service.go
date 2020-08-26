package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/wesovilabs/koazee"
)

type Service interface {
	Create(contentType string,dataFeedPath string, configPath string) error
}

type service struct {
	acousticAuthApiUrl string
	acousticContentLib string
}




func NewService(acousticAuthApiUrl string, acousticContentLib string) Service {
	return &service{
		acousticAuthApiUrl :acousticAuthApiUrl,
		acousticContentLib:acousticContentLib,
	}
}


func (service *service) Create(contentType string, dataFeedPath string, configPath string) error {
	contentClient := api.NewContentClient(service.acousticAuthApiUrl)
	records,err := Transform(contentType,dataFeedPath,configPath)
	if err != nil {
		return err
	}
	err = koazee.StreamOf(records).
			ForEach(func(record api.AcousticDataRecord) error {
				acousticContentDataOut := koazee.StreamOf(record.Values).
					Reduce(func(acc map[string]interface{},columnData api.GenericData) (map[string]interface{},error){
						if acc == nil {
							acc = make(map[string]interface{})
						}
						element,err := api.Build(columnData.Type)
						if err != nil {
							return nil,err
						}
						element,err =  element.Convert(columnData)
						if err != nil {
							return nil,err
						}
						acc[columnData.Name] = element
						return acc,nil
				})
			err := acousticContentDataOut.Err()
			if err != nil {
				return err
			}
			acousticContentData := acousticContentDataOut.Val().(map[string]interface{})
			content := api.Content{
				Name: record.Name(),
				TypeId: contentType,
				Status: "draft",
				LibraryID: service.acousticContentLib,
				Elements: acousticContentData,
			}
			_ , createErr := contentClient.Create(content)
			if createErr != nil {
				return createErr
			}
			return nil
		}).Do().Out().Err()
	if err != nil {
		return err
	}
	return nil
}


