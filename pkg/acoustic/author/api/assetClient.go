package api

import (
	"encoding/json"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/thoas/go-funk"
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
	Path      string `json:"path"`
}

type AssetResponse struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type AssetClient interface {
	Create(
		reader io.Reader,
		resourceFileName string,
		tags []string,
		path string, status string, profiles []string, libraryID string) (*AssetCreateResponse, error)
	Delete(id string) error
	Get(id string) (*AssetResponse, error)
	GetByPath(path string) (bool, *AssetResponse, error)
}

type assetClient struct {
	c              *resty.Client
	acousticApiUrl string
}

func NewAssetClient(acousticApiUrl string) AssetClient {
	apiKey := env.AcousticAPIKey()
	return &assetClient{
		c:              Connect(apiKey),
		acousticApiUrl: acousticApiUrl,
	}
}

func (assetClient assetClient) Get(id string) (*AssetResponse, error) {
	resp, err := assetClient.c.NewRequest().SetResult(&AssetResponse{}).
		SetError(&ContentAuthoringErrorResponse{}).
		Get(assetClient.acousticApiUrl + "/authoring/v1/assets/" + id)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*AssetResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return nil, errors.ErrorMessageWithStack("error in retrieving asset : " + "  " + string(errorString))
	} else {
		return nil, errors.ErrorMessageWithStack("error in retrieving asset : " + resp.Status())
	}
}

func (assetClient assetClient) GetByPath(path string) (bool, *AssetResponse, error) {
	req := assetClient.c.NewRequest().SetResult(&AssetResponse{}).
		SetError(&ContentAuthoringErrorResponse{})
	req.SetQueryParam("path", path)
	resp, err := req.Get(assetClient.acousticApiUrl + "/authoring/v1/assets/record")
	if err != nil {
		return false, nil, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		return resp.Result().(*AssetResponse).ID != "", resp.Result().(*AssetResponse), nil
	} else if resp.IsError() && resp.StatusCode() == 404 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		matchedError := funk.Find(error.Errors, func(error ContentAuthoringError) bool {
			return error.Key == "error.asset.not.found.at.path.3011"
		})
		if matchedError != nil {
			return false, nil, nil
		} else {
			errorString, _ := json.MarshalIndent(error, "", "\t")
			return false, nil, errors.ErrorMessageWithStack("error in retrieving asset : " + "  " + string(errorString))
		}

	} else {
		return false, nil, errors.ErrorMessageWithStack("error in creating asset : " + resp.Status())
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
	req := assetClient.c.NewRequest().SetPathParams(map[string]string{"id": id}).SetError(&ContentAuthoringErrorResponse{})

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
