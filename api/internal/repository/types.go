package repository

import "time"

// NowFunc returns the current time. Can be overridden for testing.
var NowFunc = time.Now
