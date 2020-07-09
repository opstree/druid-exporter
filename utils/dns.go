package utils

import (
	"net"
	"strings"

	"github.com/patrickmn/go-cache"
)

// ReverseDNSLookup returns hostname if no results are found
func ReverseDNSLookup(host string, dnsCache *cache.Cache) string {
	if cacheValue, found := dnsCache.Get(host); found {
		return cacheValue.(string)
	}
	names, err := net.LookupAddr(host)
	if err != nil || len(names) == 0 {
		dnsCache.Set(host, host, cache.DefaultExpiration)
		return host
	}
	dnsResult := strings.Split(names[0], ".")[0]
	dnsCache.Set(host, dnsResult, cache.DefaultExpiration)
	return dnsResult
}
