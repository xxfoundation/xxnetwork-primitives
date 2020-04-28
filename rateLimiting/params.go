package rateLimiting

import "time"

// Params structure holds the default values for different Buckets.
type Params struct {
	// Capacity for newly created buckets
	Capacity uint

	// Leak rate for newly created buckets
	LeakRate float64

	// How often to look for and discard stale buckets
	CleanPeriod time.Duration

	// Age of stale buckets when discarded
	MaxDuration time.Duration

	// File path for whitelist file
	WhitelistFile string
}
