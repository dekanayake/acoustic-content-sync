package csv

import (
	"encoding/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	logruserror "github.com/dekanayake/acoustic-content-sync/pkg/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/wesovilabs/koazee"
	"os"
)

type DeleteService interface {
	DeleteByFeed(deleteMappingName string, contentType string, dataFeedPath string, configPath string) (ContentDeletionStatus, error)
	Delete(libraryId string, deleteMappingName string, configPath string) error
}

type ContentDeletionStatus struct {
	Success []ContentDeletionSuccessStatus
	Failed  []ContentDeletionFailedStatus
}

type ContentDeletionSuccessStatus struct {
	CSVIDKey   string
	CSVIDValue string
	ContentID  string
}

type ContentDeletionFailedStatus struct {
	CSVIDKey   string
	CSVIDValue string
	Error      error
}

func (ContentDeletionStatus ContentDeletionStatus) TotalCount() int {
	return len(ContentDeletionStatus.Failed) + len(ContentDeletionStatus.Success)
}

func (contentDeletionStatus ContentDeletionStatus) FailuresExist() bool {
	return len(contentDeletionStatus.Failed) > 0
}

func (contentDeletionStatus ContentDeletionStatus) PrintFailed() (error error) {
	if env.LogErrorsToFile() {
		f, err := os.OpenFile(env.ErrorLogFileLocation(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			error = errors.ErrorWithStack(err)
		}
		defer func() {
			cerr := f.Close()
			if error == nil {
				error = errors.ErrorWithStack(cerr)
			}
		}()

		if env.WriteFailedRecordIDToCSV() {
			failedIdsCSVFile, err := os.Create("failed_ids.csv")
			if err != nil {
				return errors.ErrorWithStack(err)
			}
			defer failedIdsCSVFile.Close()
			failedIdsCSVWriter := csv.NewWriter(failedIdsCSVFile)
			defer failedIdsCSVWriter.Flush()
			for _, failedRecord := range contentDeletionStatus.Failed {
				failedIdsCSVWriter.Write([]string{failedRecord.CSVIDValue})
			}
		}

		errorLog := log.New()
		Formatter := new(log.TextFormatter)
		errorLog.SetFormatter(Formatter)
		errorLog.SetOutput(f)
		errorLog.Error("------------Following content failed to delete in acoustic-----------------")
		koazee.StreamOf(contentDeletionStatus.Failed).
			ForEach(func(failed ContentDeletionFailedStatus) {
				errorHandling := logruserror.PkgErrorEntry{
					Entry: errorLog.WithField("CSV Key ", failed.CSVIDKey).
						WithField(" CSV Value ", failed.CSVIDValue)}
				errorHandling.WithError(failed.Error).Error("failed record")
			}).Do()
	} else {
		koazee.StreamOf(contentDeletionStatus.Failed).
			ForEach(func(failed ContentDeletionFailedStatus) {
				errorHandling := logruserror.PkgErrorEntry{
					Entry: log.WithField("CSV Key ", failed.CSVIDKey).
						WithField(" CSV Value ", failed.CSVIDValue)}
				errorHandling.WithError(failed.Error).Error("failed record")
			}).Do()
	}
	return error
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

func (d deleteService) DeleteByFeed(deleteMappingName string, contentType string, dataFeedPath string, configPath string) (ContentDeletionStatus, error) {
	records := []api.AcousticDataRecord{}
	if dataFeedPath != "" {
		var err error = nil
		records, err = TransformContent(contentType, dataFeedPath, configPath)
		if err != nil {
			return ContentDeletionStatus{}, errors.ErrorWithStack(err)
		}
	}
	config, err := InitContentTypeMappingConfig(configPath)
	if err != nil {
		return ContentDeletionStatus{}, errors.ErrorWithStack(err)
	}
	deleteMapping, err := config.GetDeleteMapping(deleteMappingName)
	if err != nil {
		return ContentDeletionStatus{}, errors.ErrorWithStack(err)
	}

	failed := make([]ContentDeletionFailedStatus, 0)
	success := make([]ContentDeletionSuccessStatus, 0)

	if len(records) > 0 {
		for _, record := range records {
			query, err := record.SearchQuerytoGetTheContent()
			if err != nil {
				return ContentDeletionStatus{}, err
			}
			searchRequest := api.SearchRequest{
				Terms:          query,
				ContentTypes:   []string{record.SearchType},
				Classification: "content",
			}
			searchResponse, err := api.NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), record.SearchOnLibrary, record.SearchOnDeliveryAPI, searchRequest, api.Pagination{Start: 0, Rows: 1})
			if err != nil {
				log.WithField(record.CSVRecordKey, record.CSVRecordKeyValue()).Error("Failed in deleting  the content ")
				failed = append(failed, ContentDeletionFailedStatus{
					CSVIDKey:   record.CSVRecordKey,
					CSVIDValue: record.CSVRecordKeyValue(),
					Error:      errors.ErrorWithStack(err),
				})
			}
			if searchResponse.Count > 0 {
				err := delete(d, deleteMapping.AssetType, searchResponse.Documents[0].Document.ID)
				if err != nil {
					log.WithField(record.CSVRecordKey, record.CSVRecordKeyValue()).Error("Failed in deleting  the content ")
					failed = append(failed, ContentDeletionFailedStatus{
						CSVIDKey:   record.CSVRecordKey,
						CSVIDValue: record.CSVRecordKeyValue(),
						Error:      errors.ErrorWithStack(err),
					})
				} else {
					log.WithField(record.CSVRecordKey, record.CSVRecordKeyValue()).Info("Successfully deleted the content ")
					success = append(success, ContentDeletionSuccessStatus{
						CSVIDKey:   record.CSVRecordKey,
						CSVIDValue: record.CSVRecordKeyValue(),
						ContentID:  searchResponse.Documents[0].Document.ID,
					})
				}
			} else {
				failed = append(failed, ContentDeletionFailedStatus{
					CSVIDKey:   record.CSVRecordKey,
					CSVIDValue: record.CSVRecordKeyValue(),
					Error:      errors.ErrorMessageWithStack("content is not available"),
				})
			}
		}
	}
	return ContentDeletionStatus{
		Success: success,
		Failed:  failed,
	}, nil
}

func (d deleteService) Delete(libraryId string, deleteMappingName string, configPath string) error {
	config, err := InitContentTypeMappingConfig(configPath)
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
