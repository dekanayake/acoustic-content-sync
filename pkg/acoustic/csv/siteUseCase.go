package csv

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/author/api"
	"github.com/dekanayake/acoustic-content-sync/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wesovilabs/koazee"
)

type SiteUseCase interface {
	CreatePages(siteId string, parentPageId string, contentType string, dataFeedPath string, configPath string) (ContentCreationStatus, error)
}

type siteUseCase struct {
	acousticAuthApiUrl string
	contentService     api.ContentService
	siteService        api.SiteService
}

func NewSiteUseCase(acousticAuthApiUrl string) SiteUseCase {
	return &siteUseCase{
		acousticAuthApiUrl: acousticAuthApiUrl,
		siteService:        api.NewSiteService(acousticAuthApiUrl),
	}
}

func (s siteUseCase) CreatePages(siteId string, parentPageId string, contentType string, dataFeedPath string, configPath string) (ContentCreationStatus, error) {
	pageRecords, err := TransformSite(contentType, dataFeedPath, configPath)
	if err != nil {
		return ContentCreationStatus{}, errors.ErrorWithStack(err)
	}
	failed := make([]ContentCreationFailedStatus, 0)
	success := make([]ContentCreationSuccessStatus, 0)
	koazee.StreamOf(pageRecords).
		ForEach(func(record api.AcousticDataRecord) {
			response, err := s.siteService.CreatePageWithRetry(siteId, parentPageId, record)
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
					ContentID:  response.ID,
				})
			}
		}).Do()
	return ContentCreationStatus{Success: success, Failed: failed}, nil
}
