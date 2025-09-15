package models

import "time"

// Profile represents a complete customer profile including their response data,
// permissions, and current request usage statistics. This is typically returned
// when fetching detailed customer information including their current quota status.
type Profile struct {
	CustomerResponse

	// Permissions defines the access control and limits for this customer.
	// These permissions control what the customer can do within the system.
	Permissions Permissions `json:"permissions"`

	// RequestCount tracks the number of requests made by the customer within
	// the current daily period (between StartAt and EndAt).
	RequestCount int `json:"request_count"`

	// RequestLimitPerDay specifies the maximum number of requests this customer
	// is allowed to make per day based on their subscription or custom limits.
	RequestLimitPerDay int `json:"request_limit_per_day"`

	// StartAt indicates when the current daily request counting period began.
	// This is used to track daily usage and typically resets at a fixed time each day.
	StartAt time.Time `json:"start_at"`

	// EndAt indicates when the current daily request counting period ends.
	// After this time, the request count should reset and a new period begins.
	// Note: Missing json tag - should probably be `json:"end_at"`
	EndAt time.Time
}

// Permissions defines the various limits and access controls applied to a customer.
// These settings determine what actions the customer can perform and at what scale.
type Permissions struct {
	// UIRowsLimit restricts the maximum number of rows that can be displayed
	// in the user interface for this customer's queries or data views.
	UIRowsLimit int `db:"ui_rows_limit" json:"ui_rows_limit"`

	// APIRowsLimit restricts the maximum number of rows that can be returned
	// in a single API response for this customer's requests.
	APIRowsLimit int `db:"api_rows_limit" json:"api_rows_limit"`

	// RequestLimitPerDay sets the maximum number of API requests this customer
	// can make in a 24-hour period. Used for quota management.
	RequestLimitPerDay int `db:"request_limit_per_day" json:"request_limit_per_day"`

	// FilterLimits controls how many filters the customer can apply simultaneously
	// in their queries or searches.
	FilterLimits int `db:"filter_limit" json:"filter_limit"`

	// RequestRateLimit defines the maximum number of requests per minute/second
	// this customer can make to prevent system overload.
	RequestRateLimit int `db:"request_rate_limit" json:"request_rate_limit"`
}

// CustomerResponse contains the core customer information and metadata.
// This struct represents the standard customer data returned in API responses
// and includes subscription, organization, and usage statistics.
type CustomerResponse struct {
	// ID is the unique identifier for the customer in the system.
	ID int64 `json:"id"`

	// Login is the customer's username used for authentication.
	Login string `json:"login"`

	// Status represents the customer's account status (active, suspended, etc.).
	// Consider documenting the possible status values.
	Status int `json:"status"`

	// Nickname is the customer's display name or preferred name.
	Nickname string `json:"nickname"`

	// SubscriptionID links the customer to their current subscription plan.
	SubscriptionID int `json:"subscription_id"`

	// SubscriptionName is the human-readable name of the customer's subscription plan.
	// This is optional and may be null for customers without active subscriptions.
	SubscriptionName *string `json:"subscription_name,omitempty"`

	// GroupID identifies which customer group this customer belongs to.
	// Groups are used for organizing customers and applying group-level permissions.
	GroupID int64 `json:"group_id"`

	// GroupTitle is the human-readable name of the customer's group.
	// This is optional and may be null.
	GroupTitle *string `json:"group_title,omitempty"`

	// RoleID defines the customer's permission level within their organization.
	// Possible roles: 1(super_admin), 2(admin), 3(user)
	RoleID int `json:"role_id"`

	// SubscriptionStartedAt indicates when the customer's current subscription began.
	SubscriptionStartedAt time.Time `json:"subscription_start_at"`

	// SubscriptionExpiresAt indicates when the customer's current subscription expires.
	// After this date, the customer may lose access to premium features.
	SubscriptionExpiresAt time.Time `json:"subscription_expires_at"`

	// OrganizationID identifies which organization this customer belongs to.
	// Multiple customers can belong to the same organization.
	OrganizationID int64 `json:"organization_id"`

	// OrganizationTitle is the human-readable name of the customer's organization.
	// This is optional and may be null.
	OrganizationTitle *string `json:"organization_title,omitempty"`

	// CreatedAt records when this customer account was first created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt records when this customer's information was last modified.
	UpdatedAt time.Time `json:"updated_at"`

	// LastSeen tracks when the customer was last active in the system.
	// This is optional and may be null for customers who have never logged in.
	LastSeen *time.Time `json:"last_seen"`

	// TotalRequestCount tracks the total number of API requests this customer
	// has made since their account was created.
	TotalRequestCount int64 `json:"total_request_count"`

	// DailyRequestCount tracks the number of requests made in the current day.
	// This counter typically resets daily.
	DailyRequestCount int64 `json:"daily_request_count"`

	// WeekRequestCount tracks the number of requests made in the current week.
	// This counter typically resets weekly.
	WeekRequestCount int64 `json:"week_request_count"`

	// RequestLimitPerDay defines the customer's daily request quota.
	// Note: This duplicates the field in Permissions - consider consolidating.
	RequestLimitPerDay int64 `json:"request_limit_per_day"`

	// EnabledTwoFA indicates whether the customer has enabled two-factor authentication
	// for enhanced account security.
	EnabledTwoFA bool `json:"enabled_two_fa"`

	// ProductsPermission defines which specific products or features this customer
	// has access to within the system.
	ProductsPermission *ProductsPermission `json:"products_permission"`
}

// ProductsPermission defines access control for specific products or features.
// This allows fine-grained control over which parts of the system a customer can use.
type ProductsPermission struct {
	// Recona indicates whether the customer has access to the Recona product/feature.
	// Consider expanding this struct as more products are added to the system.
	Recona bool `db:"recona" json:"recona"`
}
