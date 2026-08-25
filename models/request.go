package models

// HttpRequest represents a saved REST API request.
type HttpRequest struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	Method         string `json:"method"`
	Body           string `json:"body"`
	RequestHeaders string `json:"request_headers"`
	StatusCode     int    `json:"status_code"`
	ResponseID     string `json:"response_id"`
}
