package api

import (
	"encoding/json"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/thoas/go-funk"
	"gopkg.in/resty.v1"
	"strconv"
)

type SitePageClient interface {
	Create(siteID string, sitePage SitePage) (*SitePageResponse, error)
	GetChildPages(siteID string, parentPageID string) ([]string, error)
	Update(siteID string, pageId string, sitePage SitePage) (*SitePageResponse, error)
	Move(siteID string, sourcePageID string, sourcePageRev string, targetPageID string, targetPageRev string, targetPagePosition int) (*SitePageResponse, error)
	Delete(siteID string, pageId string, deleteContent bool) (*SitePageResponse, error)
	Get(siteID string, pageId string) (*SitePageResponse, error)
}

type sitePageClient struct {
	c1             *resty.Client
	c2             *resty.Client
	acousticApiUrl string
}

func NewSitePageClient(acousticApiUrl string) SitePageClient {
	apiKey := env.AcousticAPIKey()
	return &sitePageClient{
		c1:             Connect(apiKey),
		c2:             Connect(apiKey),
		acousticApiUrl: acousticApiUrl,
	}
}

type childSitePageResponseList struct {
	Items []childSitePageResponse `json:"items"`
}

type childSitePageResponse struct {
	ID string `json:"id"`
}

func (sitePageClient sitePageClient) Create(siteID string, sitePage SitePage) (*SitePageResponse, error) {
	req := sitePageClient.c1.NewRequest().SetBody(sitePage).
		SetResult(&SitePageResponse{}).
		SetError(&ContentAuthoringErrorResponse{}).SetPathParams(map[string]string{"siteId": siteID})

	if resp, err := req.Post(sitePageClient.acousticApiUrl + "/authoring/v1/sites/{siteId}/pages"); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*SitePageResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in creating content : " + resp.Status() + "  " + string(errorString))
	} else {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in creating content : " + resp.Status() + "  " + string(errorString))
	}
}

func (sitePageClient sitePageClient) Update(siteID string, pageId string, sitePage SitePage) (*SitePageResponse, error) {
	req := sitePageClient.c1.NewRequest().SetBody(sitePage).
		SetResult(&SitePageResponse{}).
		SetError(&ContentAuthoringErrorResponse{}).SetPathParams(map[string]string{"siteId": siteID, "pageID": pageId}).SetQueryParam("forceOverride", "true")

	if resp, err := req.Put(sitePageClient.acousticApiUrl + "/authoring/v1/sites/{siteId}/pages/{pageID}"); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*SitePageResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in creating content : " + resp.Status() + "  " + string(errorString))
	} else {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in creating content : " + resp.Status() + "  " + string(errorString))
	}
}

func (sitePageClient sitePageClient) Move(siteID string, sourcePageID string, sourcePageRev string, targetPageID string, targetPageRev string, targetPagePosition int) (*SitePageResponse, error) {
	req := sitePageClient.c1.NewRequest().
		SetResult(&SitePageResponse{}).
		SetError(&ContentAuthoringErrorResponse{}).SetPathParams(map[string]string{"siteId": siteID}).
		SetQueryParams(map[string]string{
			"sourceId":       sourcePageID,
			"sourceRev":      sourcePageRev,
			"targetId":       targetPageID,
			"targetRev":      targetPageRev,
			"targetPosition": strconv.Itoa(targetPagePosition),
		})

	if resp, err := req.Put(sitePageClient.acousticApiUrl + "/authoring/v1/sites/{siteId}/pages/move"); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*SitePageResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in moving content : " + resp.Status() + "  " + string(errorString))
	} else {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in moving content : " + resp.Status() + "  " + string(errorString))
	}
}

func (sitePageClient sitePageClient) Get(siteID string, pageId string) (*SitePageResponse, error) {
	req := sitePageClient.c1.NewRequest().
		SetResult(&SitePageResponse{}).
		SetError(&ContentAuthoringErrorResponse{}).SetPathParams(map[string]string{"siteId": siteID, "pageId": pageId})

	if resp, err := req.Get(sitePageClient.acousticApiUrl + "/authoring/v1/sites/{siteId}/pages/{pageId}"); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*SitePageResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in moving content : " + resp.Status() + "  " + string(errorString))
	} else {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in moving content : " + resp.Status() + "  " + string(errorString))
	}
}

func (sitePageClient sitePageClient) Delete(siteID string, pageId string, deleteContent bool) (*SitePageResponse, error) {
	req := sitePageClient.c1.NewRequest().
		SetResult(&SitePageResponse{}).
		SetError(&ContentAuthoringErrorResponse{}).SetPathParams(map[string]string{"siteId": siteID, "pageId": pageId}).
		SetQueryParams(map[string]string{
			"delete-content": strconv.FormatBool(deleteContent),
		})

	if resp, err := req.Delete(sitePageClient.acousticApiUrl + "/authoring/v1/sites/{siteId}/pages/{pageId}"); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*SitePageResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in moving content : " + resp.Status() + "  " + string(errorString))
	} else {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in moving content : " + resp.Status() + "  " + string(errorString))
	}
}

func (sitePageClient sitePageClient) GetChildPages(siteID string, parentPageID string) ([]string, error) {
	req := sitePageClient.c2.NewRequest().
		SetResult(&childSitePageResponseList{}).
		SetError(&ContentAuthoringErrorResponse{}).SetPathParams(map[string]string{"siteId": siteID, "parentPageID": parentPageID})

	if resp, err := req.Get(sitePageClient.acousticApiUrl + "/mydelivery/v1/sites/{siteId}/pages/by-parent/{parentPageID}"); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return funk.Map(resp.Result().(*childSitePageResponseList).Items, func(x childSitePageResponse) string {
			return x.ID
		}).([]string), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in creating content : " + resp.Status() + "  " + string(errorString))
	} else {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in creating content : " + resp.Status() + "  " + string(errorString))
	}
}
