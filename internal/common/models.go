package common

import (
	"time"
)

type STDomain struct {
	Provider string
	AuthArgs map[string]string
	Prefixes []map[string]string
}

type STLookups struct {
	Timeout  time.Duration
	Interval time.Duration
	V4Enable bool
	V4Addr   string
	V4Path   string
	V6Enable bool
	V6Addr   string
	V6Path   string
}

type STConfig struct {
	Domains map[string]STDomain
	Lookups STLookups
}
