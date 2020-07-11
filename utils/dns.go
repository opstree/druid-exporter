package utils

import (
	"net"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// ReverseDNSLookup returns hostname if no results are found
func ReverseDNSLookup(host string, dnsCache *cache.Cache) string {
	if cacheValue, found := dnsCache.Get(host); found { // Value found in cache
		return cacheValue.(string)
	}
	names, err := net.LookupAddr(host)
	if err != nil || len(names) == 0 { // DNSLookup failed, returning input value
		logrus.Tracef("DNS lookup failed: host=%s", host)
		dnsCache.Set(host, host, cache.DefaultExpiration)
		return host
	}
	dnsResult := strings.Split(names[0], ".")[0] // DNS Lookup succeeded, returning result
	logrus.Tracef("DNS lookup succeeded: host=(%s => %s)", host, dnsResult)
	dnsCache.Set(host, dnsResult, cache.DefaultExpiration)
	return dnsResult
}
