package constants

import "time"

const (
	USER_AGENT          = "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/113.0.5672.127 Safari/537.36"
	DEFAULT_CRAWL_DELAY = 1 * time.Millisecond
	DEFAULT_MAX_DEPTH   = 10
	DEFAULT_MAX_PAGES   = 100
	ALPHA               = 0.85
	EPSILON             = 0.0001
	NUM_OF_RESULT       = 100
	DOMAIN              = "avash.in"
	MAX_PARALLELISM     = 5
)
