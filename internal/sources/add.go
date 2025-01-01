package sources

import (
	"github.com/ducng99/gohole/internal/db"
	"github.com/ducng99/gohole/internal/logger"
)

func AddSource(url string) error {
	dbConn := db.New().Conn

	sourceID, err := db.GetSourceID(dbConn, url)
	if err != nil {
		logger.Printf(logger.LogError, "Failed when checking if source already added\n")
		return err
	}

	if sourceID > 0 {
		logger.Printf(logger.LogWarn, "Source already exists, update source instead\n")
		return UpdateSource(sourceID)
	}

	sourceID, err = db.AddSource(dbConn, url)
	if err != nil {
		logger.Printf(logger.LogError, "Failed when adding a source entry to database\n")
		return err
	}

	logger.Printf(logger.LogNormal, "Saved new source for %s\n", url)

	if err := ForceUpdateSource(sourceID); err != nil {
		logger.Printf(logger.LogError, "Could not finish setting up new source\n")
		return err
	}

	logger.Printf(logger.LogSuccess, "Added new source successfully\n")

	return nil
}
