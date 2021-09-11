package api

import (
	"encoding/json"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"gopkg.in/resty.v1"
	"io"
)

type AssetCreateRequest struct {
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Tags        Tags     `json:"tags"`
	Status      string   `json:"status"`
	Profiles    []string `json:"profiles,omitempty"`
	LibraryID   string   `json:"libraryId"`
}

type AssetCreateResponse struct {
	Id        string `json:"id"`
	AssetType string `json:"assetType"`
	MediaType string `json:"mediaType"`
	IsManaged bool   `json:"isManaged"`
}

type AssetClient interface {
	Create(
		reader io.Reader,
		resourceFileName string,
		tags []string,
		path string, status string, profiles []string, libraryID string) (*AssetCreateResponse, error)
	Delete(id string) error
}

type assetClient struct {
	c              *resty.Client
	acousticApiUrl string
}

func NewAssetClient(acousticApiUrl string) AssetClient {
	return &assetClient{
		c:              Connect(),
		acousticApiUrl: acousticApiUrl,
	}
}

func (assetClient *assetClient) Create(
	reader io.Reader,
	resourceFileName string,
	tags []string,
	path string, status string, profiles []string, libraryID string) (*AssetCreateResponse, error) {
	resourceCreateReq := AssetCreateRequest{
		Name:        resourceFileName,
		Description: resourceFileName,
		Path:        path,
		Status:      status,
		Tags:        Tags{Values: tags},
		Profiles:    profiles,
		LibraryID:   libraryID,
	}
	resourceCreateReqJson, err := json.Marshal(resourceCreateReq)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}

	resp, err := assetClient.c.NewRequest().
		SetHeader("Content-Type", "multipart/form-data").
		SetFileReader("resource", resourceFileName, reader).
		SetFormData(map[string]string{
			"data": string(resourceCreateReqJson),
		}).SetResult(&AssetCreateResponse{}).
		SetError(&ContentAuthoringErrorResponse{}).
		Post(assetClient.acousticApiUrl + "/authoring/v1/assets")

	if err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*AssetCreateResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in creating asset : " + "  " + string(errorString))
	} else {
		return nil, errors.ErrorMessageWithStack("error in creating asset : " + resp.Status())
	}
}

func (assetClient assetClient) Delete(id string) error {
	req := assetClient.c.NewRequest().SetPathParams(map[string]string{"id": id})

	if resp, err := req.Delete(assetClient.acousticApiUrl + "/authoring/v1/assets/{id}"); err != nil {
		return errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return errors.ErrorMessageWithStack("error in deleting asset : " + resp.Status() + "  " + string(errorString))
	} else {
		return errors.ErrorMessageWithStack("error in deleting asset : " + resp.Status())
	}
}
