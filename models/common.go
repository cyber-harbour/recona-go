package models

// PaginationResponse wraps search or query results with pagination metadata.
// This is typically returned from API endpoints that support paginated results,
// providing both the result count information and pagination details.
type PaginationResponse struct {
	// TotalItems provides information about the total number of items found
	// for the query, including relationship information for approximate counts.
	TotalItems TotalItems `json:"total_items"`

	// Embedded Pagination struct containing the current page parameters.
	Pagination
}

// TotalItems represents the total count of items found for a query along with
// metadata about the accuracy of that count. This is useful for large datasets
// where exact counts might be expensive to compute.
type TotalItems struct {
	// Value is the total number of items found. This could be exact or approximate
	// depending on the relation field.
	Value int64 `json:"value"`

	// Relation indicates the relationship between the actual total and the Value field.
	// Common values might be "eq" (exact), "gte" (greater than or equal), etc.
	// This helps clients understand if the count is precise or an estimate.
	Relation string `json:"relation"`
}

// Pagination contains the standard pagination parameters used across API endpoints.
// These parameters control which subset of results to return and are typically
// provided as query parameters in GET requests.
type Pagination struct {
	// Limit specifies the maximum number of items to return per page.
	// This helps control response size and API performance.
	// Value must be in the range between 1 and 500.
	// required: false
	// minimum: 1
	// exclusiveMinimum: false
	// maximum: 500
	// exclusiveMaximum: false
	Limit int `json:"limit" schema:"limit" validate:"limit"`

	// Offset specifies the number of items to skip before starting to return results.
	// Used in combination with Limit to implement pagination (page = offset/limit + 1).
	// Value must be in the range between 0 and 9999.
	// required: false
	// minimum: 0
	// exclusiveMinimum: false
	// maximum: 9999
	// exclusiveMaximum: false
	Offset int `json:"offset" schema:"offset" validate:"offset"`
}

// SearchRequest combines search parameters with pagination controls.
// This is typically used as the request body for search API endpoints,
// allowing clients to specify both what to search for and how to paginate results.
type SearchRequest struct {
	// Embedded Search struct containing the search criteria.
	Search

	// Embedded Pagination struct containing pagination parameters.
	Pagination
}

// Search defines the core search parameters for querying data.
// This struct contains the essential elements needed to perform searches
// across the system's data.
type Search struct {
	// Query is the main search term or query string. The exact format and syntax
	// depend on the search implementation (could be full-text, structured query, etc.).
	// This field may be required depending on the validation rules applied.
	Query string `json:"query,omitempty" validate:"query_required"`

	// Filters contains additional filtering criteria to narrow down search results.
	// This is typically a JSON string or structured query format that specifies
	// field-specific filters (e.g., date ranges, categories, status values).
	// This field is optional and can be omitted if no additional filtering is needed.
	Filters string `json:"filters,omitempty" validate:"omitempty"`
}
