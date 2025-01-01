package dns

import (
	"strings"
	"time"

	"github.com/ducng99/gohole/cmd/globalFlags"
	"github.com/ducng99/gohole/internal/db"
	"github.com/ducng99/gohole/internal/logger"
	"github.com/miekg/dns"
)

const (
	upstreamDnsServer = "1.1.1.1:53"
)

func resolver(domain string, qtype uint16) []dns.RR {
	dnsEntry := getDomain(domain)
	if dnsEntry.block {
		if globalFlags.Verbose {
			logger.Printf(logger.LogNormal, "[DNS]: Blocked %s\n", domain)
		}
		return nil
	}

	message := new(dns.Msg)
	message.SetQuestion(dns.Fqdn(domain), qtype)
	message.RecursionDesired = true

	client := &dns.Client{Timeout: 5 * time.Second}

	response, _, err := client.Exchange(message, upstreamDnsServer)
	if err != nil {
		logger.Printf(logger.LogError, "[DNS]: %v\n", err)
		return nil
	}

	if response == nil {
		logger.Printf(logger.LogError, "[DNS]: no response from %s\n", upstreamDnsServer)
		return nil
	}

	if response.Rcode != dns.RcodeSuccess {
		if globalFlags.Verbose {
			logger.Printf(logger.LogError, "[DNS]: no valid answer from %s for %s\n", upstreamDnsServer, domain)
		}
		return nil
	}

	if globalFlags.Verbose {
		for _, answer := range response.Answer {
			logger.Printf(logger.LogNormal, "[DNS]: %s\n", answer.String())
		}
	}

	return response.Answer
}

func getDomain(domain string) *cachedDnsEntry {
	domain = strings.TrimRight(domain, ".")

	cacheLock.RLock()
	entry, ok := cacheDNS[domain]
	cacheLock.RUnlock()

	if !ok {
		dbConn := db.New().Conn

		existsInDb, err := db.HasDomain(dbConn, domain)
		if err != nil {
			logger.Printf(logger.LogError, "%v\n", err)
			return nil
		}

		entry = &cachedDnsEntry{
			block: existsInDb,
			eat:   time.Now().Add(cacheExpireDuration),
		}

		cacheLock.Lock()
		cacheDNS[domain] = entry
		cacheLock.Unlock()
	}

	return entry
}
