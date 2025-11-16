package interfaces

// OrderReadWriter combines read and write operations needed by services that
// mutate order aggregates.
type OrderReadWriter interface {
	OrderReader
	OrderWriter
}

// OrderRepository provides the full surface area of order persistence
// operations. Prefer injecting one of the smaller interfaces above whenever a
// service only needs a subset.
type OrderRepository interface {
	OrderReader
	OrderWriter
	OrderQuery
}
