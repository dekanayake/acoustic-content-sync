package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wesovilabs/koazee"
)

type DeleteService interface {
	Delete(libraryId string, deleteMappingName string, configPath string) error
}

type deleteService struct {
	acousticAuthApiUrl string
	assetClient        api.AssetClient
	contentClient      api.ContentClient
	searchClient       api.SearchClient
}

func NewDeleteService(acousticAuthApiUrl string) DeleteService {
	return &deleteService{
		acousticAuthApiUrl: acousticAuthApiUrl,
		assetClient:        api.NewAssetClient(acousticAuthApiUrl),
		contentClient:      api.NewContentClient(acousticAuthApiUrl),
		searchClient:       api.NewSearchClient(acousticAuthApiUrl),
	}
}

func delete(d deleteService, assetType api.AssetType, id string) error {
	if assetType == api.DOCUMENT {
		err := d.contentClient.Delete(id)
		log.WithField("type", api.DOCUMENT).WithField("id", id).Info("Deleted")
		if err != nil {
			log.WithField("type", api.DOCUMENT).WithField("id", id).Info("Delete Failed")
			return errors.ErrorWithStack(err)
		}
	} else {
		err := d.assetClient.Delete(id)
		log.WithField("type", api.DOCUMENT).WithField("id", id).Info("Deleted")
		if err != nil {
			log.WithField("type", api.DOCUMENT).WithField("id", id).Info("Delete Failed")
			return errors.ErrorWithStack(err)
		}
	}
	return nil
}

func (d deleteService) Delete(libraryId string, deleteMappingName string, configPath string) error {
	config, err := InitConfig(configPath)
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	deleteMapping, err := config.GetDeleteMapping(deleteMappingName)
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	searchRequest := deleteMapping.SearchMapping.SearchRequest()
	start := 0
	rows := 100
	for {
		searchResponse, err := d.searchClient.Search(libraryId, true, false, searchRequest, api.Pagination{Start: start, Rows: rows})
		if err != nil {
			return errors.ErrorWithStack(err)
		}
		if searchResponse.IsCountLessThanStart() {
			start, rows = searchResponse.NextPagination()
			searchResponse, err = d.searchClient.Search(libraryId, true, false, searchRequest, api.Pagination{Start: start, Rows: rows})
			if err != nil {
				return errors.ErrorWithStack(err)
			}
		}
		err = koazee.StreamOf(searchResponse.Documents).
			ForEach(func(documentItem api.DocumentItem) error {
				err := delete(d, deleteMapping.AssetType, documentItem.Document.ID)
				if err != nil {
					return errors.ErrorWithStack(err)
				}
				return nil
			}).Do().Out().Err().UserError()
		if err != nil {
			return errors.ErrorWithStack(err)
		}
		if !searchResponse.HasNext() {
			break
		} else {
			start, rows = searchResponse.NextPagination()
		}
	}
	return nil
}
