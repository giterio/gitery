package prototypes

// ContextKey is a type used for key of context
type ContextKey int

const (
	// RouteKey ...
	RouteKey ContextKey = iota
	// UserKey ...
	UserKey
)
