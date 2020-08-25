package api

import (
	"encoding/json"
	"errors"
	"github.com/dekanayake/acoustic-content-sync/pkg/context"
	"gopkg.in/resty.v1"
	"net/http"
)


type ContentClient interface {
	Create(content Content) (*ContentCreateResponse,error)
}

type contentClient struct {
	c resty.Client
	acousticApiUrl string
}

func NewContentClient(acousticApiUrl string, authToken *http.Cookie) ContentClient {
	return &contentClient{
		c: *resty.New().SetCookie(authToken).SetDebug(context.IsDebugEnabled()),
		acousticApiUrl: acousticApiUrl,
	}
}

func (contentClient *contentClient) Create (content Content) (*ContentCreateResponse,error) {
	req := contentClient.c.NewRequest().SetBody(content).
		SetResult(&ContentCreateResponse{}).
		SetError(&ContentCreateErrorResponse{})

	if resp, err := req.Post(contentClient.acousticApiUrl + "/authoring/v1/content") ; err != nil {
		return nil,err
	}   else if resp.IsSuccess() {
		return resp.Result().(*ContentCreateResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentCreateErrorResponse)
		errorString,_ := json.MarshalIndent(error, "", "\t")
		return nil, errors.New("error in creating content : " + resp.Status() + "  " + string(errorString))
	} else {
		return nil, errors.New("error in creating content : " + resp.Status())
	}
}
