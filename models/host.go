package models

// Host represents comprehensive information about a network host/IP address,
// including network services, security vulnerabilities, and geographic data
type Host struct {
	// Network identification
	IP string `json:"ip,omitempty"` // IP address of the host (IPv4 or IPv6)

	// Geographic and network provider information
	Geo *Geo `json:"geo,omitempty"` // Geographic location data (country, city, coordinates)
	Isp *ISP `json:"isp,omitempty"` // Internet Service Provider and network details

	// Port and service information
	Ports []*Port `json:"ports,omitempty"` // List of open ports and running services

	// DNS reverse lookup
	PtrRecord *PTRRecord `json:"ptr_record,omitempty"` // PTR record for reverse DNS lookup (IP to hostname)

	// Security assessment and vulnerabilities
	SeverityDetails *SeverityDetails `json:"severity_details"` // Overall security risk assessment and severity scoring
	CVEList         []*CVE           `json:"cve_list"`         // List of Common Vulnerabilities and Exposures found

	// Technology detection and analysis
	Technologies []*Technology `json:"technologies"` // Detected technologies, software, and services running on the host

	// Security and abuse information
	Abuses *Abuse `json:"abuses,omitempty"` // Abuse reports and malicious activity associated with this IP

	// SSL/TLS certificate information
	// SSL certificates found on various ports
	CertificateSummaries []*CertificateSummary `json:"certificate_summaries,omitempty"`

	// Metadata
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of last update to this host record
}

type Technology struct {
	Name                  string `json:"name,omitempty"`
	Version               string `json:"version,omitempty"`
	VersionRepresentation int64  `json:"version_representation,omitempty"`
	Port                  int64  `json:"port,omitempty"`
	LogoBase64            string `json:"logo_base64,omitempty"`
}

type SeverityDetails struct {
	High   int32 `json:"high,omitempty"`
	Low    int32 `json:"low,omitempty"`
	Medium int32 `json:"medium,omitempty"`
}

type Abuse struct {
	Score             int32            `json:"score,omitempty"`
	ReportsNum        int64            `json:"reports_num,omitempty"`
	Reports           []*AbuseReport   `json:"reports,omitempty"`
	AllCategories     []*AbuseCategory `json:"all_categories,omitempty"`
	IsWhitelistWeak   bool             `json:"is_whitelist_weak,omitempty"`
	IsWhitelistStrong bool             `json:"is_whitelist_strong,omitempty"`
	UpdatedAt         string           `json:"updated_at,omitempty"`
}

type AbuseReport struct {
	ReportedAt string           `json:"reported_at,omitempty"`
	Comment    string           `json:"comment,omitempty"`
	Categories []*AbuseCategory `json:"categories,omitempty"`
}

type AbuseCategory struct {
	ID          int32  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type IPProxyModel struct {
	IP        string       `json:"ip,omitempty"`
	ProxyData []*ProxyData `json:"proxy_data,omitempty"`
}

type ProxyData struct {
	IsProxy   bool   `json:"is_proxy,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Port      int64  `json:"port,omitempty"`
	Type      string `json:"type,omitempty"`
}

type IPGeoISP struct {
	IP  string `json:"ip,omitempty"`
	Geo *Geo   `json:"geo,omitempty"`
	ISP *ISP   `json:"isp,omitempty"`
}

type IPAbuse struct {
	IP    string `json:"ip,omitempty"`
	Abuse *Abuse `json:"abuse,omitempty"`
}

type IPPTR struct {
	IP        string     `json:"ip,omitempty"`
	PtrRecord *PTRRecord `json:"ptr_record,omitempty"`
}

type Extract struct {
	Links                  []*Link              `json:"links,omitempty"`
	Emails                 []string             `json:"emails,omitempty"`
	Errors                 []string             `json:"errors,omitempty"`
	FaviconURI             *URI                 `json:"favicon_uri,omitempty"`
	FaviconSha256          string               `json:"favicon_sha256,omitempty"`
	MetaTags               []*MetaTag           `json:"meta_tags,omitempty"`
	Description            string               `json:"description,omitempty"`
	ResponseChain          []*ResponseChainLink `json:"response_chain,omitempty"`
	StatusCode             int64                `json:"status_code,omitempty"`
	Headers                []*HTTPHeader        `json:"headers,omitempty"`
	RobotsTxt              string               `json:"robots_txt,omitempty"`
	Scripts                []string             `json:"scripts,omitempty"`
	Styles                 []string             `json:"styles,omitempty"`
	Title                  string               `json:"title,omitempty"`
	RawResponse            string               `json:"raw_response,omitempty"`
	ExternalRedirectURI    *URI                 `json:"external_redirect_uri,omitempty"`
	ExtractedAt            string               `json:"extracted_at,omitempty"`
	Cookies                []*Cookies           `json:"cookies,omitempty"`
	AdsenseID              string               `json:"adsense_id,omitempty"`
	RobotsDisallow         []string             `json:"robots_disallow,omitempty"`
	GoogleAnalyticsKey     string               `json:"google_analytics_key,omitempty"`
	GoogleSiteVerification string               `json:"google_site_verification,omitempty"`
	GooglePlayApp          string               `json:"google_play_app,omitempty"`
	AppleItunesApp         string               `json:"apple_itunes_app,omitempty"`
}

type Cookies struct {
	Key      string `json:"key,omitempty"`
	Value    string `son:"value,omitempty"`
	Expire   string `json:"expire,omitempty"`
	MaxAge   int64  `json:"max_age,omitempty"`
	Path     string `json:"path,omitempty"`
	HTTPOnly bool   `json:"http_only,omitempty"`
	Security bool   `json:"security,omitempty"`
}

type MetaTag struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type ResponseChainLink struct {
	StatusCode int64         `json:"status_code,omitempty"`
	Headers    []*HTTPHeader `json:"headers,omitempty"`
}

type HTTPHeader struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type URI struct {
	FullURI string `json:"full_uri,omitempty"`
	Host    string `json:"host,omitempty"`
	Path    string `json:"path,omitempty"`
}

type Link struct {
	Anchor     string          `json:"anchor,omitempty"`
	Attributes *LinkAttributes `json:"attributes,omitempty"`
}

type LinkAttributes struct {
	NoFollow bool `json:"no_follow,omitempty"`
	URI      *URI `json:"uri,omitempty"`
}

type Port struct {
	Banner             string   `json:"banner,omitempty"`
	CpeApplication     string   `json:"cpe_application,omitempty"`
	CpeHardware        string   `json:"cpe_hardware,omitempty"`
	CpeOs              string   `json:"cpe_os,omitempty"`
	DeviceType         string   `json:"device_type,omitempty"`
	Extract            *Extract `json:"extract,omitempty"`
	Hostname           string   `json:"hostname,omitempty"`
	Info               string   `json:"info,omitempty"`
	MasscanServiceName string   `json:"masscan_service_name,omitempty"`
	OperationSystem    string   `json:"operation_system,omitempty"`
	Port               int64    `json:"port,omitempty"`
	Product            string   `json:"product,omitempty"`
	Service            string   `json:"service,omitempty"`
	Version            string   `json:"version,omitempty"`
	UpdatedAt          string   `json:"updated_at,omitempty"`
	IsSsl              bool     `json:"is_ssl,omitempty"`
}

type Geo struct {
	CityName       string    `json:"city_name,omitempty"`
	Country        string    `json:"country,omitempty"`
	CountryIsoCode string    `json:"country_iso_code,omitempty"`
	Location       *Location `json:"location,omitempty"`
}

type Location struct {
	Lon float64 `json:"lon,omitempty"`
	Lat float64 `json:"lat,omitempty"`
}

type ISP struct {
	AsNum   uint32 `json:"as_num,omitempty"`
	AsOrg   string `json:"as_org,omitempty"`
	Isp     string `json:"isp,omitempty"`
	Network string `json:"network,omitempty"`
}

type PTRRecord struct {
	Value     string `json:"value,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CVE struct {
	BaseScore    float32  `json:"base_score,omitempty"`
	ID           string   `json:"id,omitempty"`
	Ports        []int64  `json:"ports,omitempty"`
	Severity     string   `json:"severity,omitempty"`
	Vector       string   `json:"vector,omitempty"`
	Description  string   `json:"description,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
	EPSS         *EPSS    `json:"epss,omitempty"`
	HasPOC       bool     `json:"has_poc,omitempty"`
}

type CertificateSummary struct {
	FingerprintSha256 string                      `json:"fingerprint_sha256,omitempty"`
	IssuerDn          *DomainCertificateIssuerDN  `json:"issuer_dn,omitempty"`
	SubjectDn         *DomainCertificateSubjectDN `json:"subject_dn,omitempty"`
	TLSVersion        string                      `json:"tls_version,omitempty"`
	ValidityEnd       string                      `json:"validity_end,omitempty"`
	DNSNames          []string                    `json:"dns_names,omitempty"`
	Port              int64                       `json:"port,omitempty"`
	UpdatedAt         string                      `json:"updated_at,omitempty"`
}

type HostsResponse struct {
	PaginationResponse
	Hosts []*Host `json:"hosts"`
}
