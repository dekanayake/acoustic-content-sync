package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/env"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	logruserror "github.com/dekanayake/acoustic-content-sync/pkg/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/wesovilabs/koazee"
	"os"
)

type ContentService interface {
	Create(contentType string, dataFeedPath string, configPath string) (ContentCreationStatus, error)
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

type contentService struct {
	acousticAuthApiUrl string
	acousticContentLib string
}

func NewContentService(acousticAuthApiUrl string, acousticContentLib string) ContentService {
	return &contentService{
		acousticAuthApiUrl: acousticAuthApiUrl,
		acousticContentLib: acousticContentLib,
	}
}

func (service *contentService) Create(contentType string, dataFeedPath string, configPath string) (ContentCreationStatus, error) {
	contentClient := api.NewContentClient(service.acousticAuthApiUrl)
	records, err := Transform(contentType, dataFeedPath, configPath)
	if err != nil {
		return ContentCreationStatus{}, errors.ErrorWithStack(err)
	}
	failed := make([]ContentCreationFailedStatus, 0)
	success := make([]ContentCreationSuccessStatus, 0)
	koazee.StreamOf(records).
		ForEach(func(record api.AcousticDataRecord) {
			acousticContentDataOut := koazee.StreamOf(record.Values).
				Reduce(func(acc map[string]interface{}, columnData api.GenericData) (map[string]interface{}, error) {
					if acc == nil {
						acc = make(map[string]interface{})
					}
					element, err := api.Build(columnData.Type)
					if err != nil {
						return nil, errors.ErrorWithStack(err)
					}
					element, err = element.Convert(columnData)
					if err != nil {
						return nil, errors.ErrorWithStack(err)
					}
					acc[columnData.Name] = element
					return acc, nil
				})
			err := acousticContentDataOut.Err().UserError()
			if err != nil {
				failed = append(failed, ContentCreationFailedStatus{
					CSVIDKey:   record.CSVRecordKey,
					CSVIDValue: record.CSVRecordKeyValue(),
					Error:      errors.ErrorWithStack(err),
				})
				return
			}
			acousticContentData := acousticContentDataOut.Val().(map[string]interface{})
			content := api.Content{
				Name:      record.Name(),
				TypeId:    contentType,
				Status:    env.ContentStatus(),
				LibraryID: service.acousticContentLib,
				Elements:  acousticContentData,
				Tags:      record.Tags,
			}
			response, createErr := contentClient.Create(content)
			if createErr != nil {
				failed = append(failed, ContentCreationFailedStatus{
					CSVIDKey:   record.CSVRecordKey,
					CSVIDValue: record.CSVRecordKeyValue(),
					Error:      errors.ErrorWithStack(createErr),
				})
			} else {
				success = append(success, ContentCreationSuccessStatus{
					CSVIDKey:   record.CSVRecordKey,
					CSVIDValue: record.CSVRecordKeyValue(),
					ContentID:  response.Id,
				})
			}
		}).Do()
	return ContentCreationStatus{Success: success, Failed: failed}, nil
}
