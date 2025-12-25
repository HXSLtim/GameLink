package db

// ConfigProvider defines the configuration interface needed by db package
type ConfigProvider interface {
	GetDatabaseDSN() string
	IsSeedEnabled() bool
	GetMaxOpenConns() int
}

// MetricsProvider defines the metrics interface needed by db package
type MetricsProvider interface {
	InstrumentGorm(db interface{}) error
}