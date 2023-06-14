package api

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/thoas/go-funk"
	"github.com/wesovilabs/koazee"
	"strings"
)

type SiteService interface {
	CreatePageWithRetry(siteId string, parentPageId string, record AcousticDataRecord) (*SitePageResponse, error)
}

type siteService struct {
	acousticAuthApiUrl string
	sitePageClient     SitePageClient
}

func NewSiteService(acousticAuthApiUrl string) SiteService {
	return &siteService{
		acousticAuthApiUrl: acousticAuthApiUrl,
		sitePageClient:     NewSitePageClient(acousticAuthApiUrl),
	}
}

func (service *siteService) CreatePageWithRetry(siteId string, parentPageId string, record AcousticDataRecord) (*SitePageResponse, error) {
	response, err := service.createPage(siteId, parentPageId, record)
	if err != nil && errors.IsRetryableError(err) {
		ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
		times := 1
		for range ticker.C {
			if times == 3 {
				ticker.Stop()
				return response, err
			}
			response, err = service.createPage(siteId, parentPageId, record)
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

func (service *siteService) createParentPages(siteId string, parentPageID string, url string, contentID string) (string, error) {
	segments := strings.Split(url, "/")
	var currentParentPageId = parentPageID
	for _, segment := range segments[:len(segments)-1] {
		childPages, err := service.sitePageClient.GetChildPages(siteId, currentParentPageId)
		if err != nil {
			return "", errors.ErrorWithStack(err)
		}
		matchedChildPage := funk.Find(childPages, func(childPage SitePageResponse) bool {
			return childPage.Segment == segment
		})
		if matchedChildPage != nil {
			currentParentPageId = matchedChildPage.(SitePageResponse).ID
		} else {
			sitePage := SitePage{
				Name:      segment,
				Segment:   segment,
				ParentId:  currentParentPageId,
				ContentId: contentID,
			}
			createdSite, err := service.sitePageClient.Create(siteId, sitePage)
			if err != nil {
				return "", errors.ErrorWithStack(err)
			}
			currentParentPageId = createdSite.ID
		}
	}
	return currentParentPageId, nil
}

func (service *siteService) createPage(siteId string, parentPageId string, record AcousticDataRecord) (*SitePageResponse, error) {
	acousticContentDataOut := koazee.StreamOf(record.Values).
		Reduce(func(acc map[string]string, columnData GenericData) (map[string]string, error) {
			if acc == nil {
				acc = make(map[string]string)
			}
			var value string
			switch columnData.Value.(type) {
			default:
				value = columnData.Value.(string)
			case string:
				value = columnData.Value.(string)
			case AcousticValue:
				value = columnData.Value.(AcousticValue).Value
			}
			acc[columnData.Name] = value
			return acc, nil
		})
	err := acousticContentDataOut.Err().UserError()
	if err != nil {
		return nil, err
	}
	acousticContentData := acousticContentDataOut.Val().(map[string]string)
	if acousticContentData["name"] == "" {
		return nil, errors.ErrorMessageWithStack("No value for the name")
	}
	if acousticContentData["url"] == "" {
		return nil, errors.ErrorMessageWithStack("No value for the url")
	}
	query, err := record.SearchQuerytoGetTheContent()
	if err != nil {
		return nil, err
	}
	searchRequest := SearchRequest{
		Terms:          query,
		ContentTypes:   []string{record.SearchType},
		Classification: "content",
	}
	searchResponse, err := NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), record.SearchOnLibrary, record.SearchOnDeliveryAPI, searchRequest, Pagination{Start: 0, Rows: 1})
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	if searchResponse.Count == 0 {
		return nil, errors.ErrorMessageWithStack("The content provided is not available")
	} else {
		currentParentPageId, err := service.createParentPages(siteId, parentPageId, acousticContentData["url"], searchResponse.Documents[0].Document.ID)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		segments := strings.Split(acousticContentData["url"], "/")
		lastPageSegment := segments[len(segments)-1]
		pageToCreate := SitePage{
			Name:      acousticContentData["name"],
			ContentId: searchResponse.Documents[0].Document.ID,
			ParentId:  currentParentPageId,
			Segment:   lastPageSegment,
		}
		createdPage, err := service.sitePageClient.Create(siteId, pageToCreate)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		return createdPage, nil
	}

}
