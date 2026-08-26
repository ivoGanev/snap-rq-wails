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
	if req.CollectionID == 0 {
		return models.HttpRequest{}, fmt.Errorf("collection id is required")
	}

	result, err := s.db.Exec(
		`INSERT INTO http_requests (collection_id, name, url, method, body, request_headers, status_code, response_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.CollectionID, req.Name, req.URL, req.Method, req.Body, req.RequestHeaders, req.StatusCode, req.ResponseID,
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
		`SELECT id, collection_id, name, url, method, body, request_headers, status_code, response_id
		 FROM http_requests WHERE id = ?`,
		id,
	)
	err := row.Scan(&req.ID, &req.CollectionID, &req.Name, &req.URL, &req.Method, &req.Body, &req.RequestHeaders, &req.StatusCode, &req.ResponseID)
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
		`SELECT id, collection_id, name, url, method, body, request_headers, status_code, response_id
		 FROM http_requests ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing requests: %w", err)
	}
	defer rows.Close()

	var requests []models.HttpRequest
	for rows.Next() {
		var req models.HttpRequest
		if err := rows.Scan(&req.ID, &req.CollectionID, &req.Name, &req.URL, &req.Method, &req.Body, &req.RequestHeaders, &req.StatusCode, &req.ResponseID); err != nil {
			return nil, fmt.Errorf("scanning request: %w", err)
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating requests: %w", err)
	}

	return requests, nil
}

// GetRequestsForCollection returns all requests belonging to a collection.
func (s *RequestService) GetRequestsForCollection(collectionID int64) ([]models.HttpRequest, error) {
	rows, err := s.db.Query(
		`SELECT id, collection_id, name, url, method, body, request_headers, status_code, response_id
		 FROM http_requests
		 WHERE collection_id = ?
		 ORDER BY name`,
		collectionID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing requests: %w", err)
	}
	defer rows.Close()

	var requests []models.HttpRequest
	for rows.Next() {
		var req models.HttpRequest
		if err := rows.Scan(&req.ID, &req.CollectionID, &req.Name, &req.URL, &req.Method, &req.Body, &req.RequestHeaders, &req.StatusCode, &req.ResponseID); err != nil {
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
	if req.CollectionID == 0 {
		return models.HttpRequest{}, fmt.Errorf("collection id is required")
	}

	_, err := s.db.Exec(
		`UPDATE http_requests
		 SET collection_id = ?, name = ?, url = ?, method = ?, body = ?, request_headers = ?, status_code = ?, response_id = ?
		 WHERE id = ?`,
		req.CollectionID, req.Name, req.URL, req.Method, req.Body, req.RequestHeaders, req.StatusCode, req.ResponseID, req.ID,
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

// CreateResponse saves a new response for a request and returns it with its generated ID.
func (s *RequestService) CreateResponse(resp models.HttpResponse) (models.HttpResponse, error) {
	result, err := s.db.Exec(
		`INSERT INTO responses (request_id, headers, status_code, body)
		 VALUES (?, ?, ?, ?)`,
		resp.RequestID, resp.Headers, resp.StatusCode, resp.Body,
	)
	if err != nil {
		return models.HttpResponse{}, fmt.Errorf("creating response: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.HttpResponse{}, fmt.Errorf("getting last insert id: %w", err)
	}

	resp.ID = id
	return resp, nil
}

// GetResponse retrieves a single response by ID.
func (s *RequestService) GetResponse(id int64) (models.HttpResponse, error) {
	var resp models.HttpResponse
	row := s.db.QueryRow(
		`SELECT id, request_id, headers, status_code, body FROM responses WHERE id = ?`,
		id,
	)
	err := row.Scan(&resp.ID, &resp.RequestID, &resp.Headers, &resp.StatusCode, &resp.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.HttpResponse{}, fmt.Errorf("response not found")
		}
		return models.HttpResponse{}, fmt.Errorf("getting response: %w", err)
	}
	return resp, nil
}

// GetResponsesForRequest returns all responses for a given request ID, newest first.
func (s *RequestService) GetResponsesForRequest(requestID int64) ([]models.HttpResponse, error) {
	rows, err := s.db.Query(
		`SELECT id, request_id, headers, status_code, body
		 FROM responses
		 WHERE request_id = ?
		 ORDER BY id DESC`,
		requestID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing responses: %w", err)
	}
	defer rows.Close()

	var responses []models.HttpResponse
	for rows.Next() {
		var resp models.HttpResponse
		if err := rows.Scan(&resp.ID, &resp.RequestID, &resp.Headers, &resp.StatusCode, &resp.Body); err != nil {
			return nil, fmt.Errorf("scanning response: %w", err)
		}
		responses = append(responses, resp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating responses: %w", err)
	}

	return responses, nil
}

// DeleteResponse removes a response by ID.
func (s *RequestService) DeleteResponse(id int64) error {
	_, err := s.db.Exec(`DELETE FROM responses WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting response: %w", err)
	}
	return nil
}
