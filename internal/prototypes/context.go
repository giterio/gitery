package prototypes

// ContextKey is a type used for key of context
type ContextKey int

const (
	// RouteKey is the key of Route stored in request context
	RouteKey ContextKey = iota
	// UserKey is the key of JWT token payload stored in request context
	UserKey
)
