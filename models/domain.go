package models

// Domain represents comprehensive domain information used for domain analysis,
// security scanning, and reconnaissance operations
type Domain struct {
	// Basic domain identification
	Name         string `json:"name,omitempty"`          // The primary domain name (e.g., "example.com")
	NameReversed string `json:"name_reversed,omitempty"` // Domain name in reverse order for indexing

	// WHOIS information and metadata
	WhoisParsed    *WhoisParsed `json:"whois_parsed,omitempty"`     // Structured WHOIS data
	WhoisError     string       `json:"whois_error,omitempty"`      // Error message if WHOIS lookup failed
	WhoisUpdatedAt string       `json:"whois_updated_at,omitempty"` // Timestamp of last WHOIS update
	UpdatedAt      string       `json:"updated_at,omitempty"`       // Last update timestamp for this record
	Whois          string       `json:"whois,omitempty"`            // Raw WHOIS response data

	// DNS and network configuration
	DNSRecords *DNSRecords `json:"dns_records,omitempty"` // Complete DNS record information

	// Web content and analysis
	Extract    *Extract    `json:"extract,omitempty"`    // Extracted content from domain's website
	Screenshot *Screenshot `json:"screenshot,omitempty"` // Screenshot of the domain's main page

	// SSL/TLS certificate information
	CertificateSummaries *CertSummary `json:"certificate_summaries,omitempty"` // SSL certificate details

	// DNS record type flags - indicate what types of DNS records exist
	IsNs        bool `json:"is_ns,omitempty"`        // Has Name Server records
	IsMx        bool `json:"is_mx,omitempty"`        // Has Mail Exchange records
	IsPtr       bool `json:"is_ptr,omitempty"`       // Has Pointer records (reverse DNS)
	IsCname     bool `json:"is_cname,omitempty"`     // Has Canonical Name records
	IsSubdomain bool `json:"is_subdomain,omitempty"` // This is a subdomain, not a root domain

	// Domain structure and parsing
	Suffix          string `json:"suffix,omitempty"`            // Top-level domain (TLD)
	NameFullReverse string `json:"name_full_reverse,omitempty"` // Complete reversed domain name
	NameWithoutTld  string `json:"name_without_tld,omitempty"`  // Domain name excluding the TLD
	SubdomainPart   string `json:"subdomain_part,omitempty"`    // The subdomain portion only

	// HTTP request/response data
	RequestAnswer *RequestAnswer `json:"request_answer,omitempty"` // HTTP response information

	// Technology detection and analysis
	Technologies []*Technology `json:"technologies,omitempty"` // Detected web technologies, frameworks, etc.

	// Geographic and ISP information
	Geo []*DomainGeoInfo `json:"geo,omitempty"` // Geographic location data for domain's IPs
	Isp []*DomainIspInfo `json:"isp,omitempty"` // Internet Service Provider information

	// Security and vulnerability data
	SeverityDetails *SeverityDetails `json:"severity_details,omitempty"` // Security severity assessment
	CveList         []*DomainCVE     `json:"cve_list,omitempty"`         // Common Vulnerabilities and Exposures

	// Processing and operational flags
	IsForceImport    bool   `json:"is_force_import,omitempty"`    // Force reimport of domain data
	IsDomainExtended bool   `json:"is_domain_extended,omitempty"` // Extended domain analysis performed
	UserScanAt       string `json:"user_scan_at,omitempty"`       // Timestamp of user-initiated scan
	OperationType    string `json:"operation_type,omitempty"`     // Type of operation performed on domain
}

type DomainCVE struct {
	BaseScore    float32  `json:"base_score,omitempty"`
	ID           string   `json:"id,omitempty"`
	Severity     string   `json:"severity,omitempty"`
	Vector       string   `json:"vector,omitempty"`
	Description  string   `json:"description,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
	EPSS         *EPSS    `json:"epss,omitempty"`
	HasPOC       bool     `json:"has_poc,omitempty"`
}

type SData struct {
	Env      []*ExposedEnv `json:"env,omitempty"`
	Git      []*ExposedGit `json:"git,omitempty"`
	PhpFiles []*ExposedPhp `json:"php_files,omitempty"`
}

type ExposedPhp struct {
	Path string `json:"path,omitempty"`
}

type DomainGeoInfo struct {
	CityName       string             `json:"city_name,omitempty"`
	Country        string             `json:"country,omitempty"`
	CountryIsoCode string             `json:"country_iso_code,omitempty"`
	Location       *DomainGeoLocation `json:"location,omitempty"`
	IP             string             `json:"ip,omitempty"`
}

type DomainGeoLocation struct {
	Lon float64 `json:"lon,omitempty"`
	Lat float64 `json:"lat,omitempty"`
}

type DomainIspInfo struct {
	AsNum   uint32 `json:"as_num,omitempty"`
	AsOrg   string `json:"as_org,omitempty"`
	AsName  string `json:"as_name,omitempty"`
	IP      string `json:"ip,omitempty"`
	Network string `json:"network,omitempty"`
}

type ExposedEnv struct {
	Path string             `json:"path,omitempty"`
	Data []*StringsKeyValue `json:"data,omitempty"`
	Raw  string             `json:"raw,omitempty"`
}

type ExposedGit struct {
	Path string `json:"path,omitempty"`
	Raw  string `json:"raw,omitempty"`
}

type StringsKeyValue struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type BugBounty struct {
	Name        string `json:"name,omitempty"`
	ProgramURL  string `json:"program_url,omitempty"`
	Count       int32  `json:"count,omitempty"`
	Change      int32  `json:"change,omitempty"`
	IsNew       bool   `json:"is_new,omitempty"`
	Platform    string `json:"platform,omitempty"`
	Bounty      bool   `json:"bounty,omitempty"`
	LastUpdated string `json:"last_updated,omitempty"`
}

type Screenshot struct {
	IsScreenshotted bool   `json:"IsScreenshotted,omitempty"`
	ScreenshotTime  string `json:"ScreenshotTime,omitempty"`
	ScreenshotError string `json:"ScreenshotError,omitempty"`
}

type CveLists struct {
	HTTPCVEList []*DomainCveList `json:"http_cve_list,omitempty"`
}

type DomainCveList struct {
	BaseScore  float32 `json:"base_score,omitempty"`
	ID         string  `json:"id,omitempty"`
	Severity   string  `json:"severity,omitempty"`
	Vector     string  `json:"vector,omitempty"`
	Technology string  `json:"technology,omitempty"`
}

type Files struct {
	Env *Env `json:"env,omitempty"`
}

type Env struct {
	Text string `json:"text,omitempty"`
	Path string `json:"path,omitempty"`
}

type WhoisParsed struct {
	ErrorCode  int32       `json:"error_code,omitempty"`
	Registrar  *Registrar  `json:"registrar,omitempty"`
	Registrant *Registrant `json:"registrant,omitempty"`
	Admin      *Registrant `json:"admin,omitempty"`
	Tech       *Registrant `json:"tech,omitempty"`
	Bill       *Registrant `json:"bill,omitempty"`
	UpdatedAt  string      `json:"updated_at,omitempty"`
}

type Registrar struct {
	CreatedDate    string `json:"created_date,omitempty"`
	DomainDnssec   string `json:"domain_dnssec,omitempty"`
	DomainID       string `json:"domain_id,omitempty"`
	DomainName     string `json:"domain_name,omitempty"`
	DomainStatus   string `json:"domain_status,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty"`
	NameServers    string `json:"name_servers,omitempty"`
	ReferralURL    string `json:"referral_url,omitempty"`
	RegistrarID    string `json:"registrar_id,omitempty"`
	RegistrarName  string `json:"registrar_name,omitempty"`
	UpdatedDate    string `json:"updated_date,omitempty"`
	WhoisServer    string `json:"whois_server,omitempty"`
	Emails         string `json:"emails,omitempty"`
}

type Registrant struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Organization string `json:"organization,omitempty"`
	Street       string `json:"street,omitempty"`
	StreetExt    string `json:"street_ext,omitempty"`
	City         string `json:"city,omitempty"`
	Province     string `json:"province,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
	Phone        string `json:"phone,omitempty"`
	PhoneExt     string `json:"phone_ext,omitempty"`
	Fax          string `json:"fax,omitempty"`
	FaxExt       string `json:"fax_ext,omitempty"`
	Email        string `json:"email,omitempty"`
}

type DNSRecords struct {
	A         []string   `json:"A,omitempty"`
	AAAA      []string   `json:"AAAA,omitempty"`
	CNAME     []string   `json:"CNAME,omitempty"`
	TXT       []string   `json:"TXT,omitempty"`
	NS        []string   `json:"NS,omitempty"`
	MX        []string   `json:"MX,omitempty"`
	SPF       []*SPF     `json:"SPF,omitempty"`
	SOA       *SOARecord `json:"SOA,omitempty"`
	CAA       []string   `json:"CAA,omitempty"`
	UpdatedAt string     `json:"updated_at,omitempty"`
}

type SPF struct {
	Version          string                `json:"version,omitempty"`
	ValidationErrors []*SpfValidationError `json:"validation_errors,omitempty"`
	Mechanisms       []*SpfMechanism       `json:"mechanisms,omitempty"`
	Modifiers        []*SpfModifier        `json:"modifiers,omitempty"`
	Raw              string                `json:"raw,omitempty"`
}

type SpfModifier struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type SpfMechanism struct {
	Name      string `json:"name,omitempty"`
	Qualifier string `json:"qualifier,omitempty"`
	Value     string `json:"value,omitempty"`
}

type SpfValidationError struct {
	Description string `json:"description,omitempty"`
	Target      string `json:"target,omitempty"`
}

type DNSA struct {
	A string `json:"A,omitempty"`
}

type SOARecord struct {
	Ns      string `json:"ns,omitempty"`
	Email   string `json:"email,omitempty"`
	Serial  int64  `json:"serial,omitempty"`
	Refresh int64  `json:"refresh,omitempty"`
	Retry   int64  `json:"retry,omitempty"`
	Expire  int64  `json:"expire,omitempty"`
	MinTTL  int64  `json:"min_ttl,omitempty"`
}

type CAARecord struct {
	Tag   string `json:"tag,omitempty"`
	Value string `json:"value,omitempty"`
	Flag  int64  `json:"flag,omitempty"`
}

type CertSummary struct {
	FingerprintSha256 string                      `json:"fingerprint_sha256,omitempty"`
	IssuerDn          *DomainCertificateIssuerDN  `json:"issuer_dn,omitempty"`
	SubjectDn         *DomainCertificateSubjectDN `json:"subject_dn,omitempty"`
	TLSVersion        string                      `json:"tls_version,omitempty"`
	ValidityEnd       string                      `json:"validity_end,omitempty"`
	DNSNames          []string                    `json:"dns_names,omitempty"`
	UpdatedAt         string                      `json:"updated_at,omitempty"`
}

type DomainsResponse struct {
	PaginationResponse
	Domains []*Domain `json:"domains"`
}
