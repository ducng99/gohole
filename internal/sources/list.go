package sources

import (
	"fmt"
	"time"

	"github.com/ducng99/gohole/internal/db"
	"github.com/ducng99/gohole/internal/logger"
)

func ListSources() error {
	dbConn := db.New().Conn

	sources, err := db.GetSources(dbConn)
	if err != nil {
		logger.Printf(logger.LogError, "Failed when getting all sources details\n")
		return err
	}

	logger.Printf(logger.LogNormal, "Found %d sources:\n", len(sources))

	for _, source := range sources {
		lastFetch := time.Unix(source.LastFetch, 0)
		fmt.Printf("%d. %s - Last fetched: %s - No. entries: %d\n", source.ID, source.Url, lastFetch.Format(time.RFC822), source.NumEntries)
	}

	return nil
}
