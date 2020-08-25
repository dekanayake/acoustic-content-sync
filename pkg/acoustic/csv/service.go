package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/wesovilabs/koazee"
	"net/http"
)

type Service interface {
	Create(contentType string,dataFeedPath string, configPath string) error
}

type service struct {
	authToken *http.Cookie
	acousticAuthApiUrl string
	acousticContentLib string
}




func NewService(authToken *http.Cookie, acousticAuthApiUrl string, acousticContentLib string) Service {
	return &service{
		authToken : authToken,
		acousticAuthApiUrl :acousticAuthApiUrl,
		acousticContentLib:acousticContentLib,
	}
}


func (service *service) Create(contentType string, dataFeedPath string, configPath string) error {
	contentClient := api.NewContentClient(service.acousticAuthApiUrl,service.authToken)
	records,err := Transform(contentType,dataFeedPath,configPath)
	if err != nil {
		return err
	}
	err = koazee.StreamOf(records).
			ForEach(func(record AcousticDataRecord) error {
				acousticContentDataOut := koazee.StreamOf(record.values).
					Reduce(func(acc map[string]interface{},columnData api.GenericData) (map[string]interface{},error){
						if acc == nil {
							acc = make(map[string]interface{})
						}
						element,err := api.Build(columnData.Type)
						if err != nil {
							return nil,err
						}
						element,err =  element.Enrich(columnData, api.TextConverter)
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
				Name: "test",
				TypeId: contentType,
				Status: "ready",
				LibraryID: service.acousticContentLib,
				Elements: acousticContentData,
			}
			contentClient.Create(content)
			return nil
		}).Out().Err()
	if err != nil {
		return err
	}
	return nil
}


