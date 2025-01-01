package dns

import (
	"sync"
	"time"
)

const cacheExpireDuration = 1 * time.Hour

type cachedDnsEntry struct {
	block bool
	eat   time.Time
}

var cacheDNS = map[string]*cachedDnsEntry{}
var cacheLock sync.RWMutex

func StartCacheCleaner() {
	go func() {
		for {
			clearCache()
			time.Sleep(cacheExpireDuration)
		}
	}()
}

func clearCache() {
	cacheLock.Lock()

	for domain, entry := range cacheDNS {
		if entry.eat.Before(time.Now()) {
			delete(cacheDNS, domain)
		}
	}

	cacheLock.Unlock()
}
