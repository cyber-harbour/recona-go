package models

import "time"

// NistCVEData represents comprehensive vulnerability information from NIST's
// Common Vulnerabilities and Exposures (CVE) database, including scoring,
// exploit availability, and mitigation data
type NistCVEData struct {
	// Basic CVE identification and status
	ID     string `json:"id"`     // CVE identifier (e.g., "CVE-2023-12345")
	Status string `json:"status"` // CVE status (e.g., "Published", "Modified", "Rejected")

	// Data availability flags - indicate what additional data is present
	HasPOC     bool `json:"has_poc"`     // Proof of Concept exploit code is available
	HasEPSS    bool `json:"has_epss"`    // Exploit Prediction Scoring System data is available
	HasCVSS    bool `json:"has_cvss"`    // Common Vulnerability Scoring System data is available
	HasTargets bool `json:"has_targets"` // Specific target configurations are defined

	// Critical vulnerability indicators
	IsKEVListed bool `json:"is_kev_listed"` // Listed in CISA's Known Exploited Vulnerabilities catalog

	// Classification and metadata
	Tags        []string `json:"tags"`        // Additional classification tags
	Description string   `json:"description"` // Detailed vulnerability description

	// External information and references
	References []*Reference `json:"references"` // External links and documentation

	// Threat intelligence and scoring data
	KEV  *KEV  `json:"kev"`  // Known Exploited Vulnerability catalog entry
	CVSS *CVSS `json:"cvss"` // Common Vulnerability Scoring System metrics and scores
	EPSS *EPSS `json:"epss"` // Exploit Prediction Scoring System probability data
	POC  *POC  `json:"poc"`  // Proof of Concept exploit information

	// Weakness classification
	CWES []string `json:"cwes"` // Common Weakness Enumeration identifiers

	// Affected systems and configurations
	Configurations []*Configuration `json:"configurations"` // Vulnerable software configurations and versions

	// Timeline information
	LastModifiedAt *time.Time `json:"last_modified_at"` // When the CVE was last updated in the database
	PublishedAt    *time.Time `json:"published_at"`     // When the CVE was first published
}

type Reference struct {
	Source string   `json:"source"`
	Tags   []string `json:"tags"`
	URL    string   `json:"url"`
}

type KEV struct {
	VulnerabilityName string     `json:"vulnerability_name"`
	ActionRequired    string     `json:"action_required"`
	ExploitAdded      *time.Time `json:"exploit_added"`
	ActionDue         *time.Time `json:"action_due"`
}

type Configuration struct {
	Nodes    []*Node `json:"nodes"`
	Operator string  `json:"operator"`
}

type Node struct {
	CPEMatch []*CPEMatch `json:"cpe_match"`
	Negate   bool        `json:"negate"`
	Operator string      `json:"operator"`
}

type CPEMatch struct {
	Criteria              string `json:"criteria"`
	MatchCriteriaID       string `json:"match_criteria_id"`
	VersionEndExcluding   string `json:"version_end_excluding"`
	VersionEndIncluding   string `json:"version_end_including"`
	VersionStartExcluding string `json:"version_start_excluding"`
	VersionStartIncluding string `json:"version_start_including"`
	Vulnerable            bool   `json:"vulnerable"`
}

type CVSS struct {
	Score    float64 `json:"score"`
	Severity string  `json:"severity"`
	Metrics  *Metric `json:"metrics"`
}

type Metric struct {
	V2  []*CVSSV2 `json:"v2,omitempty"`
	V3  []*CVSSV3 `json:"v3,omitempty"`
	V31 []*CVSSV3 `json:"v3_1,omitempty"`
	V4  []*CVSSV4 `json:"v4,omitempty"`
}

type CVSSV2 struct {
	ACInsufInfo             bool        `json:"ac_insuf_info"`
	BaseSeverity            string      `json:"base_severity"`
	CVSSData                *CVSSDataV2 `json:"cvss_data"`
	ExploitabilityScore     float64     `json:"exploitability_score"`
	ImpactScore             float64     `json:"impact_score"`
	ObtainAllPrivilege      bool        `json:"obtain_all_privilege"`
	ObtainOtherPrivilege    bool        `json:"obtain_other_privilege"`
	ObtainUserPrivilege     bool        `json:"obtain_user_privilege"`
	Source                  string      `json:"source"`
	Type                    string      `json:"type"`
	UserInteractionRequired bool        `json:"user_interaction_required"`
}

type CVSSDataV2 struct {
	AccessComplexity      string  `json:"access_complexity"`
	AccessVector          string  `json:"access_vector"`
	Authentication        string  `json:"authentication"`
	AvailabilityImpact    string  `json:"availability_impact"`
	BaseScore             float64 `json:"base_score"`
	ConfidentialityImpact string  `json:"confidentiality_impact"`
	IntegrityImpact       string  `json:"integrity_impact"`
	VectorString          string  `json:"vector_string"`
	Version               string  `json:"version"`
}

type CVSSV3 struct {
	CVSSData            *CVSSDataV3 `json:"cvss_data"`
	ExploitabilityScore float64     `json:"exploitability_score"`
	ImpactScore         float64     `json:"impact_score"`
	Source              string      `json:"source"`
	Type                string      `json:"type"`
}

type CVSSDataV3 struct {
	AttackComplexity      string  `json:"attack_complexity"`
	AttackVector          string  `json:"attack_vector"`
	AvailabilityImpact    string  `json:"availability_impact"`
	BaseScore             float64 `json:"base_score"`
	BaseSeverity          string  `json:"base_severity"`
	ConfidentialityImpact string  `json:"confidentiality_impact"`
	IntegrityImpact       string  `json:"integrity_impact"`
	PrivilegesRequired    string  `json:"privileges_required"`
	Scope                 string  `json:"scope"`
	UserInteraction       string  `json:"user_interaction"`
	VectorString          string  `json:"vector_string"`
	Version               string  `json:"version"`
}

type CVSSV4 struct {
	CVSSData *CVSSDataV4 `json:"cvss_data"`
	Source   string      `json:"source"`
	Type     string      `json:"type"`
}

type CVSSDataV4 struct {
	AttackComplexity                        string  `json:"attack_complexity"`
	AttackRequirements                      string  `json:"attack_requirements"`
	AttackVector                            string  `json:"attack_vector"`
	Automatable                             string  `json:"automatable"`
	AvailabilityRequirements                string  `json:"availability_requirements"`
	BaseScore                               float64 `json:"base_score"`
	BaseSeverity                            string  `json:"base_severity"`
	ConfidentialityRequirements             string  `json:"confidentiality_requirements"`
	ExploitMaturity                         string  `json:"exploit_maturity"`
	IntegrityRequirements                   string  `json:"integrity_requirements"`
	ModifiedAttackComplexity                string  `json:"modified_attack_complexity"`
	ModifiedAttackRequirements              string  `json:"modified_attack_requirements"`
	ModifiedAttackVector                    string  `json:"modified_attack_vector"`
	ModifiedPrivilegesRequired              string  `json:"modified_privileges_required"`
	ModifiedSubsequentSystemAvailability    string  `json:"modified_subsequent_system_availability"`
	ModifiedSubsequentSystemConfidentiality string  `json:"modified_subsequent_system_confidentiality"`
	ModifiedSubsequentSystemIntegrity       string  `json:"modified_subsequent_system_integrity"`
	ModifiedUserInteraction                 string  `json:"modified_user_interaction"`
	ModifiedVulnerableSystemAvailability    string  `json:"modified_vulnerable_system_availability"`
	ModifiedVulnerableSystemConfidentiality string  `json:"modified_vulnerable_system_confidentiality"`
	ModifiedVulnerableSystemIntegrity       string  `json:"modified_vulnerable_system_integrity"`
	PrivilegesRequired                      string  `json:"privileges_required"`
	ProviderUrgency                         string  `json:"provider_urgency"`
	Recovery                                string  `json:"recovery"`
	Safety                                  string  `json:"safety"`
	SubsequentSystemAvailability            string  `json:"subsequent_system_availability"`
	SubsequentSystemConfidentiality         string  `json:"subsequent_system_confidentiality"`
	SubsequentSystemIntegrity               string  `json:"subsequent_system_integrity"`
	UserInteraction                         string  `json:"user_interaction"`
	ValueDensity                            string  `json:"value_density"`
	VectorString                            string  `json:"vector_string"`
	Version                                 string  `json:"version"`
	VulnerabilityResponseEffort             string  `json:"vulnerability_response_effort"`
	VulnerableSystemAvailability            string  `json:"vulnerable_system_availability"`
	VulnerableSystemConfidentiality         string  `json:"vulnerable_system_confidentiality"`
	VulnerableSystemIntegrity               string  `json:"vulnerable_system_integrity"`
}

type EPSS struct {
	Score      float64   `json:"score"`
	Percentile float64   `json:"percentile"`
	Date       time.Time `json:"date"`
}

type POC struct {
	References []string `json:"references"`
}

type CWE struct {
	Code                string `json:"code"`
	Name                string `json:"name"`
	Abstraction         string `json:"abstraction"`
	Structure           string `json:"structure"`
	Status              string `json:"status"`
	Description         string `json:"description"`
	ExtendedDescription string `json:"extended_description"`
}

type CVEResponse struct {
	PaginationResponse
	CVEList []*NistCVEData `json:"cve_list"`
}

type CWEParams struct {
	IDs []string `json:"ids"` // List of CWE IDs to filter by (e.g., ["CWE-79", "CWE-89"])
}

type CWEResponse struct {
	Items []*CWE `json:"items"`
}
