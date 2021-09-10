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

type ContentUseCase interface {
	CreateBatch(contentType string, dataFeedPath string, configPath string) (ContentCreationStatus, error)
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

func (contentUseCase *contentUseCase) CreateBatch(contentType string, dataFeedPath string, configPath string) (ContentCreationStatus, error) {
	records, err := Transform(contentType, dataFeedPath, configPath)
	if err != nil {
		return ContentCreationStatus{}, errors.ErrorWithStack(err)
	}
	failed := make([]ContentCreationFailedStatus, 0)
	success := make([]ContentCreationSuccessStatus, 0)
	koazee.StreamOf(records).
		ForEach(func(record api.AcousticDataRecord) {
			response, err := contentUseCase.contentService.CreateContentWithRetry(record, contentType)
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
