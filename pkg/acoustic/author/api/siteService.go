package api

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"github.com/wesovilabs/koazee"
	"strings"
)

type PageCreationStatus string

const (
	PAGE_CREATED PageCreationStatus = "PAGE_CREATED"
	PAGE_UPDATED                    = "PAGE_UPDATED"
	PAGE_EXIST                      = "PAGE_EXIST"
)

type SiteService interface {
	CreatePageWithRetry(siteId string, parentPageId string, record AcousticDataRecord) (PageCreationStatus, *SitePageResponse, error)
	CreatePageForContent(siteId string, parentPageId string, contentID string, relativePath string) (string, error)
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

func (service *siteService) CreatePageForContent(siteId string, parentPageId string, contentID string, relativePath string) (string, error) {
	currentParentPageId, err := service.createParentPages(siteId, parentPageId, relativePath)
	if err != nil {
		return "", err
	}
	segments := strings.Split(relativePath, "/")
	lastPageSegment := segments[len(segments)-1]
	pageToCreate := SitePage{
		Name:      &lastPageSegment,
		ContentId: &contentID,
		ParentId:  &currentParentPageId,
		Segment:   &lastPageSegment,
	}
	createdPage, err := service.sitePageClient.Create(siteId, pageToCreate)
	if err != nil {
		return "", errors.ErrorWithStack(err)
	}
	return createdPage.ID, nil
}

func (service *siteService) CreatePageWithRetry(siteId string, parentPageId string, record AcousticDataRecord) (PageCreationStatus, *SitePageResponse, error) {
	status, response, err := service.createPage(siteId, parentPageId, record)
	if err != nil && errors.IsRetryableError(err) {
		ticker := backoff.NewTicker(backoff.NewExponentialBackOff())
		times := 1
		for range ticker.C {
			if times == 3 {
				ticker.Stop()
				return status, response, err
			}
			status, response, err = service.createPage(siteId, parentPageId, record)
			if err != nil && errors.IsRetryableError(err) {
				times++
				continue
			} else {
				ticker.Stop()
				return status, response, err
			}
			ticker.Stop()
			break
		}
	}
	return status, response, err
}

func (service *siteService) createParentPages(siteId string, parentPageID string, url string) (string, error) {
	segments := strings.Split(url, "/")
	var currentParentPageId = parentPageID
	for _, segment := range segments[:len(segments)-1] {
		childPages, err := service.getChildPages(siteId, currentParentPageId)
		if err != nil {
			return "", errors.ErrorWithStack(err)
		}
		matchedChildPage := funk.Find(childPages, func(childPage *SitePageResponse) bool {
			return childPage.Segment == segment
		})
		if matchedChildPage != nil {
			currentParentPageId = matchedChildPage.(*SitePageResponse).ID
		} else {
			contentTypeID := env.GetOrPanic("ParentPageContentTypeID")
			sitePage := SitePage{
				Name:          &segment,
				Segment:       &segment,
				ParentId:      &currentParentPageId,
				ContentTypeId: &contentTypeID,
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

func (service *siteService) getPage(siteId string, parentPageId string, segment string) (*SitePageResponse, bool, error) {
	childPages, err := service.getChildPages(siteId, parentPageId)
	if err != nil {
		return nil, false, errors.ErrorWithStack(err)
	}
	matchedChildPage := funk.Find(childPages, func(childPage *SitePageResponse) bool {
		return childPage.Segment == segment
	})
	if matchedChildPage != nil {
		value := matchedChildPage.(*SitePageResponse)
		return value, true, nil
	}
	return nil, false, nil
}

func (service *siteService) getChildPages(siteID string, parentPageID string) ([]*SitePageResponse, error) {
	childIds, err := service.sitePageClient.GetChildPages(siteID, parentPageID)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	childPageList := make([]*SitePageResponse, 0)
	for _, childId := range childIds {
		childPage, err := service.sitePageClient.Get(siteID, childId)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
		childPageList = append(childPageList, childPage)
	}
	return childPageList, nil
}

func (service *siteService) updatePage(contentID string, siteID string, parentPageID string, page *SitePageResponse) (*SitePageResponse, error) {
	var err error
	defer func() {
		if err != nil {
			pageToUpdate := SitePage{
				Name:          &page.Name,
				ContentId:     &page.ContentId,
				ParentId:      &page.ParentID,
				Segment:       &page.Segment,
				ContentTypeId: &page.ContentTypeId,
				Rev:           &page.Rev,
				LayoutId:      &page.LayoutId,
				Position:      &page.Position,
				Title:         &page.Title,
				Description:   &page.Description,
			}
			_, err := service.sitePageClient.Update(siteID, page.ID, pageToUpdate)
			if err != nil {
				log.Error("Error occured while reverting the current site page", err)
			}
		}
	}()
	movedName := page.Name + "_MOVED"
	movedSegment := page.Segment + "_MOVED"
	pageToUpdate := SitePage{
		Name:          &movedName,
		ContentId:     &page.ContentId,
		ParentId:      &page.ParentID,
		Segment:       &movedSegment,
		ContentTypeId: &page.ContentTypeId,
		Rev:           &page.Rev,
		LayoutId:      &page.LayoutId,
		Position:      &page.Position,
		Title:         &page.Title,
		Description:   &page.Description,
	}
	updatePage, err := service.sitePageClient.Update(siteID, page.ID, pageToUpdate)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}

	pageToCreate := SitePage{
		Name:      &page.Name,
		ContentId: &contentID,
		ParentId:  &parentPageID,
		Segment:   &page.Segment,
	}
	createdPage, err := service.sitePageClient.Create(siteID, pageToCreate)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}

	childPages, err := service.getChildPages(siteID, updatePage.ID)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}

	for _, childPage := range childPages {
		childPageToUpdate := SitePage{
			Name:          &childPage.Name,
			ContentId:     &childPage.ContentId,
			ParentId:      &createdPage.ID,
			Segment:       &childPage.Segment,
			ContentTypeId: &childPage.ContentTypeId,
			Rev:           &childPage.Rev,
			LayoutId:      &childPage.LayoutId,
			Position:      &childPage.Position,
			Title:         &childPage.Title,
			Description:   &childPage.Description,
		}
		_, err := service.sitePageClient.Update(siteID, childPage.ID, childPageToUpdate)
		if err != nil {
			return nil, errors.ErrorWithStack(err)
		}
	}

	_, err = service.sitePageClient.Delete(siteID, updatePage.ID, true)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	return createdPage, nil
}

func (service *siteService) createPage(siteId string, parentPageId string, record AcousticDataRecord) (PageCreationStatus, *SitePageResponse, error) {

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
		return "", nil, err
	}
	acousticContentData := acousticContentDataOut.Val().(map[string]string)
	if acousticContentData["name"] == "" {
		return "", nil, errors.ErrorMessageWithStack("No value for the name")
	}
	if acousticContentData["url"] == "" {
		return "", nil, errors.ErrorMessageWithStack("No value for the url")
	}
	query, err := record.SearchQuerytoGetTheContent()
	if err != nil {
		return "", nil, err
	}
	searchRequest := SearchRequest{
		Terms:          query,
		ContentTypes:   []string{record.SearchType},
		Classification: "content",
	}
	searchResponse, err := NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), record.SearchOnLibrary, record.SearchOnDeliveryAPI, searchRequest, Pagination{Start: 0, Rows: 1})
	if err != nil {
		return "", nil, errors.ErrorWithStack(err)
	}
	if searchResponse.Count == 0 {
		return "", nil, errors.ErrorMessageWithStack("The content provided is not available")
	} else {
		currentParentPageId, err := service.createParentPages(siteId, parentPageId, acousticContentData["url"])
		if err != nil {
			return "", nil, errors.ErrorWithStack(err)
		}
		segments := strings.Split(acousticContentData["url"], "/")
		lastPageSegment := segments[len(segments)-1]
		if record.SiteConfig.DontCreatePageIfExist {
			page, pageExist, err := service.getPage(siteId, currentParentPageId, lastPageSegment)
			if err != nil {
				return "", nil, errors.ErrorWithStack(err)
			}
			if pageExist {
				return PAGE_EXIST, page, nil
			}
		}
		if record.SiteConfig.UpdatePageIfExists {
			page, pageExist, err := service.getPage(siteId, currentParentPageId, lastPageSegment)
			if err != nil {
				return "", nil, errors.ErrorWithStack(err)
			}
			if pageExist {
				if page.ContentId != searchResponse.Documents[0].Document.ID {
					updatedPage, err := service.updatePage(searchResponse.Documents[0].Document.ID, siteId, currentParentPageId, page)
					if err != nil {
						return "", nil, errors.ErrorWithStack(err)
					}
					return PAGE_UPDATED, updatedPage, nil
				} else {
					return PAGE_EXIST, page, nil
				}
			}
		}
		pageToCreate := SitePage{
			Name:      &lastPageSegment,
			ContentId: &searchResponse.Documents[0].Document.ID,
			ParentId:  &currentParentPageId,
			Segment:   &lastPageSegment,
		}
		createdPage, err := service.sitePageClient.Create(siteId, pageToCreate)
		if err != nil {
			return "", nil, errors.ErrorWithStack(err)
		}
		return PAGE_CREATED, createdPage, nil

	}

}
