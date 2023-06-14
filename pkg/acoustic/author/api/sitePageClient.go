package api

import (
	"encoding/json"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"gopkg.in/resty.v1"
)

type SitePageClient interface {
	Create(siteID string, sitePage SitePage) (*SitePageResponse, error)
	GetChildPages(siteID string, parentPageID string) ([]SitePageResponse, error)
}

type sitePageClient struct {
	c1             *resty.Client
	c2             *resty.Client
	acousticApiUrl string
}

func NewSitePageClient(acousticApiUrl string) SitePageClient {
	return &sitePageClient{
		c1:             Connect(),
		c2:             Connect(),
		acousticApiUrl: acousticApiUrl,
	}
}

func (sitePageClient sitePageClient) Create(siteID string, sitePage SitePage) (*SitePageResponse, error) {
	req := sitePageClient.c1.NewRequest().SetBody(sitePage).
		SetResult(&SitePageResponse{}).
		SetError(&ContentAuthoringErrorResponse{}).SetPathParams(map[string]string{"siteId": siteID})

	if resp, err := req.Post(sitePageClient.acousticApiUrl + "authoring/v1/sites/{siteId}/pages"); err != nil {
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

func (sitePageClient sitePageClient) GetChildPages(siteID string, parentPageID string) ([]SitePageResponse, error) {
	req := sitePageClient.c2.NewRequest().
		SetResult(&SitePageResponseList{}).
		SetError(&ContentAuthoringErrorResponse{}).SetPathParams(map[string]string{"siteId": siteID, "parentPageID": parentPageID})

	if resp, err := req.Get(sitePageClient.acousticApiUrl + "mydelivery/v1/sites/{siteId}/pages/by-parent/{parentPageID}"); err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*SitePageResponseList).Items, nil
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
