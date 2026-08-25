package services

import (
	"database/sql"
	"fmt"

	"snap-rq/models"
)

// RequestService provides CRUD operations for saved HTTP requests.
type RequestService struct {
	db *sql.DB
}

// NewRequestService returns a RequestService backed by the given database.
func NewRequestService(db *sql.DB) *RequestService {
	return &RequestService{db: db}
}

// CreateRequest saves a new HTTP request and returns it with its generated ID.
func (s *RequestService) CreateRequest(req models.HttpRequest) (models.HttpRequest, error) {
	result, err := s.db.Exec(
		`INSERT INTO http_requests (name, url, method, body, request_headers, status_code, response_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.Name, req.URL, req.Method, req.Body, req.RequestHeaders, req.StatusCode, req.ResponseID,
	)
	if err != nil {
		return models.HttpRequest{}, fmt.Errorf("creating request: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.HttpRequest{}, fmt.Errorf("getting last insert id: %w", err)
	}

	req.ID = id
	return req, nil
}

// GetRequest retrieves a single HTTP request by ID.
func (s *RequestService) GetRequest(id int64) (models.HttpRequest, error) {
	var req models.HttpRequest
	row := s.db.QueryRow(
		`SELECT id, name, url, method, body, request_headers, status_code, response_id
		 FROM http_requests WHERE id = ?`,
		id,
	)
	err := row.Scan(&req.ID, &req.Name, &req.URL, &req.Method, &req.Body, &req.RequestHeaders, &req.StatusCode, &req.ResponseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.HttpRequest{}, fmt.Errorf("request not found")
		}
		return models.HttpRequest{}, fmt.Errorf("getting request: %w", err)
	}
	return req, nil
}

// GetAllRequests returns all saved HTTP requests ordered by name.
func (s *RequestService) GetAllRequests() ([]models.HttpRequest, error) {
	rows, err := s.db.Query(
		`SELECT id, name, url, method, body, request_headers, status_code, response_id
		 FROM http_requests ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing requests: %w", err)
	}
	defer rows.Close()

	var requests []models.HttpRequest
	for rows.Next() {
		var req models.HttpRequest
		if err := rows.Scan(&req.ID, &req.Name, &req.URL, &req.Method, &req.Body, &req.RequestHeaders, &req.StatusCode, &req.ResponseID); err != nil {
			return nil, fmt.Errorf("scanning request: %w", err)
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating requests: %w", err)
	}

	return requests, nil
}

// UpdateRequest updates an existing HTTP request.
func (s *RequestService) UpdateRequest(req models.HttpRequest) (models.HttpRequest, error) {
	if req.ID == 0 {
		return models.HttpRequest{}, fmt.Errorf("request id is required")
	}

	_, err := s.db.Exec(
		`UPDATE http_requests
		 SET name = ?, url = ?, method = ?, body = ?, request_headers = ?, status_code = ?, response_id = ?
		 WHERE id = ?`,
		req.Name, req.URL, req.Method, req.Body, req.RequestHeaders, req.StatusCode, req.ResponseID, req.ID,
	)
	if err != nil {
		return models.HttpRequest{}, fmt.Errorf("updating request: %w", err)
	}

	return req, nil
}

// DeleteRequest removes an HTTP request by ID.
func (s *RequestService) DeleteRequest(id int64) error {
	_, err := s.db.Exec(`DELETE FROM http_requests WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting request: %w", err)
	}
	return nil
}
