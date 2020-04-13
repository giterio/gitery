package views

// Error is a custom error
type Error struct {
	StatusCode  int    `json:"status_code"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
	trace       string
}
