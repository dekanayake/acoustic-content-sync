package api

import (
	"encoding/json"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	"github.com/wesovilabs/koazee"
	"gopkg.in/resty.v1"
	"net/url"
	"strconv"
)

type SearchResponse struct {
	Start     int
	Rows      int
	Count     int            `json:"numFound"`
	Documents []DocumentItem `json:"documents"`
}

type DocumentItem struct {
	Document Document `json:"document"`
}

type Document struct {
	ID string `json:"id"`
}

type SearchRequest struct {
	Term           string
	ContentTypes   []string
	CriteriaList   []FilterCriteria
	Classification string
	AssetType      AssetType
}

type FilterCriteria interface {
	Query() string
}

type GenericFilterCriteria struct {
	Field string
	Value string
}

type Pagination struct {
	Start int
	Rows  int
}

func (searchResponse SearchResponse) HasNext() bool {
	return (searchResponse.Start + searchResponse.Rows) < searchResponse.Count
}

func (filterCriteria GenericFilterCriteria) Query() string {
	return filterCriteria.Field + ":" + filterCriteria.Value
}

func (searchRequest SearchRequest) ContentTypesQuery() string {
	contentTypes := koazee.StreamOf(searchRequest.ContentTypes).
		Reduce(func(acc string, contentType string) string {
			if acc != "" {
				acc += "OR"
			}
			acc += "\"" + contentType + "\""
			return acc
		}).String()
	if contentTypes == "" {
		return ""
	} else {
		return "type:" + contentTypes
	}
}

func (searchRequest SearchRequest) ClassificationQuery() string {
	return "classification:" + searchRequest.Classification
}

func (searchRequest SearchRequest) AssetTypeQuery() string {
	if searchRequest.AssetType != "" {
		return "assetType:" + string(searchRequest.AssetType)
	} else {
		return ""
	}

}

func (searchRequest SearchRequest) SearchTerm() string {
	if searchRequest.Term == "" {
		return "*"
	} else {
		return searchRequest.Term
	}
}

type SearchClient interface {
	Search(libraryId string, searchRequest SearchRequest, pagination Pagination) (SearchResponse, error)
}

type searchClient struct {
	c              *resty.Client
	acousticApiUrl string
}

func NewSearchClient(acousticApiUrl string) SearchClient {
	return &searchClient{
		c:              Connect(),
		acousticApiUrl: acousticApiUrl,
	}
}

func (searchClient searchClient) Search(libraryId string, searchRequest SearchRequest, pagination Pagination) (SearchResponse, error) {
	req := searchClient.c.NewRequest().SetResult(&SearchResponse{}).SetError(&ContentAuthoringErrorResponse{})
	req.SetQueryParam("q", searchRequest.SearchTerm())
	req.SetQueryParam("fl", "document:[json]")
	queryParams := make([]string, 0)
	queryParams = append(queryParams, "libraryId:(\""+libraryId+"\")")
	queryParams = append(queryParams, searchRequest.ContentTypesQuery())
	queryParams = append(queryParams, searchRequest.AssetTypeQuery())
	queryParams = append(queryParams, searchRequest.ClassificationQuery())
	koazee.StreamOf(searchRequest.CriteriaList).
		ForEach(func(criteria FilterCriteria) {
			queryParams = append(queryParams, criteria.Query())
		})
	queryParams = koazee.StreamOf(queryParams).
		Filter(func(queryParam string) bool {
			return queryParam != ""
		}).Out().Val().([]string)
	req.SetMultiValueQueryParams(url.Values{
		"fq": queryParams,
	})
	req.SetQueryParam("rows", strconv.Itoa(pagination.Rows))
	req.SetQueryParam("start", strconv.Itoa(pagination.Start))

	if resp, err := req.Get(searchClient.acousticApiUrl + "/authoring/v1/search"); err != nil {
		return SearchResponse{}, errors.ErrorWithStack(err)
	} else if resp.IsSuccess() {
		searchResponse := *resp.Result().(*SearchResponse)
		searchResponse.Start = pagination.Start
		searchResponse.Rows = pagination.Rows
		return searchResponse, nil
	} else if resp.IsError() && resp.StatusCode() == 400 {
		error := resp.Error().(*ContentAuthoringErrorResponse)
		errorString, _ := json.MarshalIndent(error, "", "\t")
		return SearchResponse{}, errors.ErrorMessageWithStack("error in searching  : " + resp.Status() + "  " + string(errorString))
	} else {
		return SearchResponse{}, errors.ErrorMessageWithStack("error in searching  : " + resp.Status())
	}
}
