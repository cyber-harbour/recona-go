package models

type AS struct {
	Number       int64       `json:"number,omitempty"`
	Organization string      `json:"organization,omitempty"`
	Ipv4Ranges   []*ASSubnet `json:"ipv4_ranges,omitempty"`
	Ipv6Ranges   []*ASSubnet `json:"ipv6_ranges,omitempty"`
	UpdatedAt    string      `json:"updated_at,omitempty"`
}

type ASSubnet struct {
	Cidr string `json:"cidr,omitempty"`
	Isp  string `json:"isp,omitempty"`
}

type ASResponse struct {
	PaginationResponse
	AutonomousSystems []*AS `json:"autonomous_systems"`
}
