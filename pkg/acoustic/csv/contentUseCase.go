package csv

import (
	"encoding/csv"
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	logruserror "github.com/dekanayake/acoustic-content-sync/pkg/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"github.com/wesovilabs/koazee"
	"os"
)

type ContentUseCase interface {
	CreateBatch(contentType string, dataFeedPath string, configPath string) (ContentCreationStatus, error)
	ReadBatch(contentType string, dataFeedPath string, configPath string) error
}

type ContentCreationStatus struct {
	Success []ContentCreationSuccessStatus
	Failed  []ContentCreationFailedStatus
}

type ContentCreationFailedStatus struct {
	CSVIDKey   string
	CSVIDValue string
	Error      error
}

type ContentCreationSuccessStatus struct {
	CSVIDKey   string
	CSVIDValue string
	ContentID  string
}

type contentUseCase struct {
	acousticAuthApiUrl string
	acousticContentLib string
	contentService     api.ContentService
}

func (contentCreationStatus ContentCreationStatus) TotalCount() int {
	return len(contentCreationStatus.Failed) + len(contentCreationStatus.Success)
}

func (contentCreationStatus ContentCreationStatus) FailuresExist() bool {
	return len(contentCreationStatus.Failed) > 0
}

func (contentCreationStatus ContentCreationStatus) PrintFailed() (error error) {
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

		errorLog := log.New()
		Formatter := new(log.TextFormatter)
		errorLog.SetFormatter(Formatter)
		errorLog.SetOutput(f)
		errorLog.Error("------------Following content failed to crate in acoustic-----------------")
		koazee.StreamOf(contentCreationStatus.Failed).
			ForEach(func(failed ContentCreationFailedStatus) {
				errorHandling := logruserror.PkgErrorEntry{
					Entry: errorLog.WithField("CSV Key ", failed.CSVIDKey).
						WithField(" CSV Value ", failed.CSVIDValue)}
				errorHandling.WithError(failed.Error).Error("failed record")
			}).Do()
	} else {
		koazee.StreamOf(contentCreationStatus.Failed).
			ForEach(func(failed ContentCreationFailedStatus) {
				errorHandling := logruserror.PkgErrorEntry{
					Entry: log.WithField("CSV Key ", failed.CSVIDKey).
						WithField(" CSV Value ", failed.CSVIDValue)}
				errorHandling.WithError(failed.Error).Error("failed record")
			}).Do()
	}
	return error
}

func NewContentUseCase(acousticAuthApiUrl string, acousticContentLib string) ContentUseCase {
	return &contentUseCase{
		acousticAuthApiUrl: acousticAuthApiUrl,
		acousticContentLib: acousticContentLib,
		contentService:     api.NewContentService(acousticAuthApiUrl, acousticContentLib),
	}
}

func getFilterValues(fileLocationPath string, columns []string) ([]map[string]string, error) {
	filterValuesFeed, err := LoadCSV(fileLocationPath)
	if err != nil {
		return nil, errors.ErrorWithStack(err)
	}
	filterValues := make([]map[string]string, 0)
	for ok := true; ok; ok = filterValuesFeed.HasNext() {
		filterValueMap := make(map[string]string)
		filterValueRecord := filterValuesFeed.Next()
		for _, column := range columns {
			filterValue, err := filterValueRecord.Get(column)
			if err != nil {
				return nil, errors.ErrorWithStack(err)
			}
			filterValueMap[column] = filterValue
			filterValues = append(filterValues, filterValueMap)
		}
	}
	return filterValues, nil
}

func (contentUseCase *contentUseCase) CreateBatch(contentType string, dataFeedPath string, configPath string) (ContentCreationStatus, error) {
	records, err := Transform(contentType, dataFeedPath, configPath)
	if err != nil {
		return ContentCreationStatus{}, errors.ErrorWithStack(err)
	}
	filterRecords := records[0].FilterRecords
	if filterRecords {
		filterValues, err := getFilterValues(records[0].FilterFileLocation, records[0].FilterColumns)
		if err != nil {
			return ContentCreationStatus{}, errors.ErrorWithStack(err)
		}
		records = funk.Filter(records, func(record api.AcousticDataRecord) bool {
			return funk.Contains(filterValues, func(filterValueMap map[string]string) bool {
				contains := true
				for filterKey, filterValue := range filterValueMap {
					value := record.GetValue(filterKey)
					contains = contains && value != nil && value == filterValue
				}
				return contains
			})
		}).([]api.AcousticDataRecord)
	}
	failed := make([]ContentCreationFailedStatus, 0)
	success := make([]ContentCreationSuccessStatus, 0)
	koazee.StreamOf(records).
		ForEach(func(record api.AcousticDataRecord) {
			response, err := contentUseCase.contentService.CreateOrUpdateContentWithRetry(record, contentType)
			if err != nil {
				log.WithField(record.CSVRecordKey, record.CSVRecordKeyValue()).Error("Failed in creating  the content ")
				failed = append(failed, ContentCreationFailedStatus{
					CSVIDKey:   record.CSVRecordKey,
					CSVIDValue: record.CSVRecordKeyValue(),
					Error:      errors.ErrorWithStack(err),
				})
			} else if response != nil {
				log.WithField(record.CSVRecordKey, record.CSVRecordKeyValue()).Info("Successfully created the content ")
				success = append(success, ContentCreationSuccessStatus{
					CSVIDKey:   record.CSVRecordKey,
					CSVIDValue: record.CSVRecordKeyValue(),
					ContentID:  response.Id,
				})
			}
		}).Do()
	return ContentCreationStatus{Success: success, Failed: failed}, nil
}

func (contentUseCase contentUseCase) ReadBatch(contentType string, dataFeedPath string, configPath string) error {
	csvFile, err := os.Create(dataFeedPath)
	defer csvFile.Close()
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	csvFileWriter := csv.NewWriter(csvFile)
	defer csvFileWriter.Flush()

	config, err := InitConfig(configPath)
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	configTypeMapping, err := config.GetContentType(contentType)
	if err != nil {
		return errors.ErrorWithStack(err)
	}
	acousticFields := configTypeMapping.GetAcousticFields()
	rowHeaders := make([]string, 0)
	for _, acousticField := range acousticFields {
		rowHeaders = append(rowHeaders, acousticField)
	}
	if err := csvFileWriter.Write(rowHeaders); err != nil {
		return errors.ErrorWithStack(err)
	}

	searchRequest := api.NewSearchRequest(configTypeMapping.SearchTerm, configTypeMapping.SearchTerms)
	searchRequest.ContentTypes = []string{configTypeMapping.SearchType}
	searchRequest.Classification = "content"

	start := 0
	rows := 100
	if configTypeMapping.PaginationRows > 0 {
		rows = configTypeMapping.PaginationRows
	}
	documents := make([]api.DocumentItem, 0)

	for {
		searchResponse, err := api.NewSearchClient(env.AcousticAPIUrl()).Search(env.LibraryID(), configTypeMapping.SearchOnLibrary, configTypeMapping.SearchOnDeliveryAPI, searchRequest, api.Pagination{Start: start, Rows: rows})
		if err != nil {
			return errors.ErrorWithStack(err)
		}
		documents = append(documents, searchResponse.Documents...)
		if !searchResponse.HasNext() {
			break
		} else {
			start, rows = searchResponse.NextPagination()
		}
	}
	contentClient := api.NewContentClient(env.AcousticAPIUrl())
	if len(documents) > 0 {
		for _, document := range documents {
			contentId := document.Document.ID
			existingContent, err := contentClient.Get(contentId)
			if err != nil {
				return errors.ErrorWithStack(err)
			}
			row := make([]string, 0)
			for _, acousticField := range acousticFields {
				if element, ok := existingContent.Elements[acousticField]; ok {
					existingElement, err := api.Convert(element.(map[string]interface{}))
					value, err := existingElement.ToCSV()
					if err != nil {
						return errors.ErrorWithStack(err)
					}
					row = append(row, value)
				}
			}
			if err := csvFileWriter.Write(row); err != nil {
				return errors.ErrorWithStack(err)
			}
		}
	} else {
		return errors.ErrorMessageWithStack("No records for the match with the search term")
	}
	return nil
}
