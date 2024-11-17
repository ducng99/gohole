package sources

import (
	"errors"
	"net/http"
	"time"

	"github.com/ducng99/gohole/internal/db"
	"github.com/ducng99/gohole/internal/hosts"
	"github.com/ducng99/gohole/internal/logger"
)

func UpdateAllSources() error {
	dbConn := db.New().Conn

	sources, err := db.GetSources(dbConn)
	if err != nil {
		logger.Printf(logger.LogError, "Failed when getting all sources from database\n")
		return err
	}

	var updateErr error = nil

	for _, source := range sources {
		if err := UpdateSource(source.ID); err != nil {
			updateErr = errors.Join(updateErr, err)
		}
	}

	return updateErr
}

func UpdateSource(sourceID int64) error {
	dbConn := db.New().Conn

	logger.Printf(logger.LogNormal, "Checking source %d if it needs to be updated\n", sourceID)

	lastFetch, _, err := db.GetSourceStats(dbConn, sourceID)
	if err != nil {
		if errors.Is(err, db.ErrSourceNotFound) {
			logger.Printf(logger.LogError, "Source not found in database\n")
			return nil
		}

		logger.Printf(logger.LogError, "Failed when getting source stats\n")
		return err
	}

	if time.Now().Unix()-lastFetch > maxStoreSourceSeconds {
		logger.Printf(logger.LogNormal, "Source %d is outdated\n", sourceID)
		return ForceUpdateSource(sourceID)
	}

	logger.Printf(logger.LogNormal, "Source %d is up-to-date\n", sourceID)

	return nil
}

func ForceUpdateSource(sourceID int64) error {
	logger.Printf(logger.LogNormal, "Begin updating source %d\n", sourceID)

	dbConn := db.New().Conn

	url, err := db.GetSourceUrl(dbConn, sourceID)
	if err != nil {
		if errors.Is(err, db.ErrSourceNotFound) {
			logger.Printf(logger.LogError, "Source not found in database\n")
			return nil
		}

		return err
	}

	tx, err := dbConn.Begin()
	if err != nil {
		logger.Printf(logger.LogError, "Failed when creating db transaction\n")
		return err
	}

	if err := db.ClearSourceDomains(tx, sourceID); err != nil {
		logger.Printf(logger.LogError, "Failed when clearing domains for source %d\n", sourceID)
		return err
	}

	// Get the file from the world
	logger.Printf(logger.LogNormal, "Fetching file from %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		logger.Printf(logger.LogError, "Failed when fetching %s\n", url)
		return err
	}

	// Parse the file to get the domains
	logger.Printf(logger.LogNormal, "Parsing file\n")

	domains, err := hosts.ParseFromReader(resp.Body)
	if err != nil {
		logger.Printf(logger.LogError, "Failed when parsing domains\n")
		return err
	}

	logger.Printf(logger.LogSuccess, "Parsed %d domains\n", len(domains))

	// Add all domains to database
	logger.Printf(logger.LogNormal, "Adding %d domains to database\n", len(domains))

	if err := db.AddDomains(tx, domains, sourceID); err != nil {
		logger.Printf(logger.LogError, "Failed when adding domains to database\n")
		return err
	}

	logger.Printf(logger.LogSuccess, "Added %d domains to database\n", len(domains))

	// Update source entry with stats
	if err := db.UpdateSource(tx, sourceID, len(domains)); err != nil {
		logger.Printf(logger.LogError, "Failed when updating source entry\n")
		return err
	}

	if err = tx.Commit(); err != nil {
		logger.Printf(logger.LogError, "Failed when committing changes to database\n")
		return err
	}

	db.New().Vacuum()

	logger.Printf(logger.LogSuccess, "Finished updating source %d\n", sourceID)

	return nil
}
