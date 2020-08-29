package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/wesovilabs/koazee"
)

type ContentService interface {
	Create(contentType string, dataFeedPath string, configPath string) error
}

type contentService struct {
	acousticAuthApiUrl string
	acousticContentLib string
}

func NewContentService(acousticAuthApiUrl string, acousticContentLib string) ContentService {
	return &contentService{
		acousticAuthApiUrl: acousticAuthApiUrl,
		acousticContentLib: acousticContentLib,
	}
}

func (service *contentService) Create(contentType string, dataFeedPath string, configPath string) error {
	contentClient := api.NewContentClient(service.acousticAuthApiUrl)
	records, err := Transform(contentType, dataFeedPath, configPath)
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	err = koazee.StreamOf(records).
		ForEach(func(record api.AcousticDataRecord) error {
			acousticContentDataOut := koazee.StreamOf(record.Values).
				Reduce(func(acc map[string]interface{}, columnData api.GenericData) (map[string]interface{}, error) {
					if acc == nil {
						acc = make(map[string]interface{})
					}
					element, err := api.Build(columnData.Type)
					if err != nil {
						return nil, errors.ErrorWithStack(err)
					}
					element, err = element.Convert(columnData)
					if err != nil {
						return nil, errors.ErrorWithStack(err)
					}
					acc[columnData.Name] = element
					return acc, nil
				})
			err := acousticContentDataOut.Err()
			if err != nil {
				return errors.ErrorWithStack(err)
			}
			acousticContentData := acousticContentDataOut.Val().(map[string]interface{})
			content := api.Content{
				Name:      record.Name(),
				TypeId:    contentType,
				Status:    env.ContentStatus(),
				LibraryID: service.acousticContentLib,
				Elements:  acousticContentData,
			}
			_, createErr := contentClient.Create(content)
			if createErr != nil {
				return errors.ErrorWithStack(err)
			}
			return nil
		}).Do().Out().Err()
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	return nil
}