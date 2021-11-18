package api

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"github.com/wesovilabs/koazee"
)

type ContentService interface {
	CreateOrUpdateContentWithRetry(record AcousticDataRecord, contentType string) (*ContentAutheringResponse, error)
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

func (service *contentService) CreateOrUpdateContentWithRetry(record AcousticDataRecord, contentType string) (*ContentAutheringResponse, error) {
	response, err := service.createOrUpdate(record, contentType)
	if err != nil && errors.IsRetryableError(err) {
		ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
		times := 1
		for range ticker.C {
			if times == 3 {
				ticker.Stop()
				return response, err
			}
			response, err = service.createOrUpdate(record, contentType)
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

func handlePreContentCreateFunctions(content Content) (Content, error) {
	for fieldName, element := range content.Elements {
		preContentCreateFuncs := element.(Element).PreContentCreateFunctions()
		for _, preContentCreateFunc := range preContentCreateFuncs {
			element, err := preContentCreateFunc()
			if err != nil {
				return Content{}, err
			}
			content.Elements[fieldName] = element
		}
	}
	return content, nil
}

func handlePreContentUpdateFunctions(content Content) (Content, error) {
	for fieldName, element := range content.Elements {
		if elementInstance, ok := element.(Element); ok {
			preContentUpdateFuncs := elementInstance.PreContentUpdateFunctions()
			for _, preContentUpdateFunc := range preContentUpdateFuncs {
				element, err := preContentUpdateFunc(element.(Element))
				if err != nil {
					return Content{}, err
				}
				content.Elements[fieldName] = element
			}
		}
	}
	return content, nil
}

func handlePostContentUpdateFunctions(content Content) error {
	for _, element := range content.Elements {
		if elementInstance, ok := element.(Element); ok {
			postContentUpdateFuncs := elementInstance.PostContentUpdateFunctions()
			for _, postContentUpdateFunc := range postContentUpdateFuncs {
				err := postContentUpdateFunc(element.(Element))
				if err != nil {
					logrus.Warn("post update field with error ", err)
				}
			}
		}
	}
	return nil
}

func (service *contentService) createOrUpdate(record AcousticDataRecord, contentType string) (*ContentAutheringResponse, error) {
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
	if !record.Update && record.CreateNonExistingItems {
		query, err := record.searchQuerytoGetTheContentToUpdate()
		if err != nil {
			return nil, err
		}
		searchRequest := SearchRequest{
			Term:           query,
			ContentTypes:   []string{record.SearchType},
			Classification: "content",
		}
		searchResponse, err := NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), searchRequest, Pagination{Start: 0, Rows: 1})
		if err != nil {
			return nil, err
		}
		if searchResponse.Count == 0 {
			content, err := handlePreContentCreateFunctions(content)
			if err != nil {
				return nil, err
			}
			response, createErr := service.contentClient.Create(content)
			if createErr != nil {
				return nil, createErr
			} else {
				return response, nil
			}
		} else {
			return nil, nil
		}
	}
	if record.Update {
		query, err := record.searchQuerytoGetTheContentToUpdate()
		if err != nil {
			return nil, err
		}
		searchRequest := SearchRequest{
			Term:           query,
			ContentTypes:   []string{record.SearchType},
			Classification: "content",
		}
		searchResponse, err := NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), searchRequest, Pagination{Start: 0, Rows: 1})
		if err != nil {
			return nil, err
		}
		if searchResponse.Count > 0 {
			contentId := searchResponse.Documents[0].Document.ID
			existingContent, err := service.contentClient.Get(contentId)
			updatedContent := Content{}
			copier.Copy(&updatedContent, &existingContent)
			if err != nil {
				return nil, err
			}
			for newContentElementKey, newElement := range content.Elements {
				existingContentElement := updatedContent.Elements[newContentElementKey]
				if existingContentElement == nil {
					updatedContent.Elements[newContentElementKey] = newElement
				} else {
					existingElement, err := Convert(existingContentElement.(map[string]interface{}))
					existingContent.Elements[newContentElementKey] = existingElement
					if err != nil {
						return nil, err
					}
					updatedElement, err := existingElement.Update(newElement.(Element))
					if err != nil {
						return nil, err
					}
					if updatedElement != nil {
						updatedContent.Elements[newContentElementKey] = updatedElement
					}
				}
			}
			content, err := handlePreContentUpdateFunctions(updatedContent)
			defer handlePostContentUpdateFunctions(*existingContent)
			if err != nil {
				return nil, err
			}
			response, udpateError := service.contentClient.Update(content)
			if udpateError != nil {
				return nil, udpateError
			} else {
				return response, nil
			}

		} else {
			if !record.CreateNonExistingItems {
				return nil, errors.ErrorMessageWithStack("No existing items found for query :" + query + " search type :" + record.SearchType)
			}
		}
	}
	content, err = handlePreContentCreateFunctions(content)
	if err != nil {
		return nil, err
	}
	response, createErr := service.contentClient.Create(content)
	if createErr != nil {
		return nil, createErr
	} else {
		return response, nil
	}
}
