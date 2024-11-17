package sources

import (
	"github.com/ducng99/gohole/internal/db"
	"github.com/ducng99/gohole/internal/hosts"
	"github.com/ducng99/gohole/internal/logger"
)

func WriteDomainsToHosts() error {
	logger.Printf(logger.LogNormal, "Writing domains to hosts file\n")

	dbConn := db.New().Conn

	domainEntries, err := db.GetDomains(dbConn)
	if err != nil {
		logger.Printf(logger.LogError, "Failed when getting all domains from database\n")
		return err
	}

	domains := make([]string, 0, len(domainEntries))

	for _, entry := range domainEntries {
		domains = append(domains, entry.Domain)
	}

	if err := hosts.AddDomainsToHosts(domains); err != nil {
		return err
	}

	return nil
}
