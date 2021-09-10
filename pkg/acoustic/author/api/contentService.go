package api

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/wesovilabs/koazee"
)

type ContentService interface {
	CreateContentWithRetry(record AcousticDataRecord, contentType string) (*ContentCreateResponse, error)
}

type contentService struct {
	acousticAuthApiUrl string
	acousticContentLib string
	contentClient      ContentClient
}

func NewContentService(acousticAuthApiUrl string, acousticContentLib string) ContentService {
	return &contentService{
		acousticAuthApiUrl: acousticAuthApiUrl,
		acousticContentLib: acousticContentLib,
		contentClient:      NewContentClient(acousticAuthApiUrl),
	}
}

func (service *contentService) CreateContentWithRetry(record AcousticDataRecord, contentType string) (*ContentCreateResponse, error) {
	response, err := service.create(record, contentType)
	if err != nil && errors.IsRetryableError(err) {
		ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
		times := 1
		for range ticker.C {
			if times == 3 {
				ticker.Stop()
				return response, err
			}
			response, err = service.create(record, contentType)
			if err != nil && errors.IsRetryableError(err) {
				times++
				continue
			} else {
				ticker.Stop()
				return response, err
			}
			ticker.Stop()
			break
		}
	}
	return response, err
}

func (service *contentService) create(record AcousticDataRecord, contentType string) (*ContentCreateResponse, error) {
	acousticContentDataOut := koazee.StreamOf(record.Values).
		Reduce(func(acc map[string]interface{}, columnData GenericData) (map[string]interface{}, error) {
			if columnData.Ignore {
				return acc, nil
			}
			if columnData.Value == nil {
				return acc, nil
			}
			if acc == nil {
				acc = make(map[string]interface{})
			}
			element, err := Build(columnData.Type)
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
	err := acousticContentDataOut.Err().UserError()
	if err != nil {
		return nil, err
	}
	acousticContentData := acousticContentDataOut.Val().(map[string]interface{})
	content := Content{
		Name:      record.Name(),
		TypeId:    contentType,
		Status:    env.ContentStatus(),
		LibraryID: service.acousticContentLib,
		Elements:  acousticContentData,
		Tags:      record.Tags,
	}
	response, createErr := service.contentClient.Create(content)
	if createErr != nil {
		return nil, createErr
	} else {
		return response, nil
	}
}
