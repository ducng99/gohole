package sources

import (
	"github.com/ducng99/gohole/internal/db"
	"github.com/ducng99/gohole/internal/logger"
)

func RemoveSource(sourceID int64) error {
	logger.Printf(logger.LogNormal, "Removing source %d\n", sourceID)

	dbConn := db.New().Conn

	tx, err := dbConn.Begin()
	if err != nil {
		logger.Printf(logger.LogError, "Failed when create db transaction\n")
		return err
	}

	if err := db.ClearSourceDomains(tx, sourceID); err != nil {
		logger.Printf(logger.LogError, "Failed when removing domains for source %d\n", sourceID)
		return err
	}

	if err := db.RemoveSource(tx, sourceID); err != nil {
		logger.Printf(logger.LogError, "Failed when removing source %d\n", sourceID)
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Printf(logger.LogError, "Failed when committing changes\n")
		return err
	}

	db.New().Vacuum()

	logger.Printf(logger.LogSuccess, "Removed source %d and its domains\n", sourceID)

	return nil
}
