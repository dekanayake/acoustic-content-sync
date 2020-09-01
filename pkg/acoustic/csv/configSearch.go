package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/wesovilabs/koazee"
)

func (searchMapping SearchMapping) SearchRequest() api.SearchRequest {
	searchRequest := api.SearchRequest{}
	searchRequest.Term = searchMapping.SearchTerm
	searchRequest.ContentTypes = koazee.StreamOf([]string{searchMapping.ContentType}).
		Filter(func(contentType string) bool {
			return contentType != ""
		}).Out().Val().([]string)
	searchRequest.Classification = searchMapping.Classification
	return searchRequest
}
