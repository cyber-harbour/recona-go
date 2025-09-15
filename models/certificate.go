package models

// Certificate represents SSL/TLS certificate information including
// parsed certificate data, raw content, and validation status
type Certificate struct {
	// Parsed certificate data
	Parsed *Parsed `json:"parsed,omitempty"` // Structured certificate information (subject, issuer, dates, etc.)

	// Raw certificate content
	Raw string `json:"raw,omitempty"` // Raw certificate data in PEM or DER format

	// Certificate identification
	FingerprintSha256 string `json:"fingerprint_sha256,omitempty"` // SHA-256 fingerprint for unique identification

	// Certificate validation and trust status
	Validation *Validation `json:"validation,omitempty"` // Certificate chain validation results and trust status

	// Metadata
	UpdatedAt string `json:"updated_at,omitempty"` // Timestamp of last update to this certificate record
}

type Validation struct {
	Valid  bool   `json:"valid,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type Parsed struct {
	Extensions             *Extensions         `json:"extensions,omitempty"`
	FingerprintMd5         string              `json:"fingerprint_md5,omitempty"`
	FingerprintSha1        string              `json:"fingerprint_sha1,omitempty"`
	FingerprintSha256      string              `json:"fingerprint_sha256,omitempty"`
	Issuer                 *Issuer             `json:"issuer,omitempty"`
	IssuerDn               string              `json:"issuer_dn,omitempty"`
	Names                  []string            `json:"names,omitempty"`
	Redacted               bool                `json:"redacted,omitempty"`
	SerialNumber           string              `json:"serial_number,omitempty"`
	Signature              *Signature          `json:"signature,omitempty"`
	SignatureAlgorithm     *SignatureAlgorithm `json:"signature_algorithm,omitempty"`
	SpkiSubjectFingerprint string              `json:"spki_subject_fingerprint,omitempty"`
	Subject                *Subject            `json:"subject,omitempty"`
	SubjectDn              string              `json:"subject_dn,omitempty"`
	SubjectKeyInfo         *SubjectKeyInfo     `json:"subject_key_info,omitempty"`
	TbsFingerprint         string              `json:"tbs_fingerprint,omitempty"`
	TbsNoctFingerprint     string              `json:"tbs_noct_fingerprint,omitempty"`
	ValidationLevel        string              `json:"validation_level,omitempty"`
	Validity               *Validity           `json:"validity,omitempty"`
	Version                int64               `json:"version,omitempty"`
}

type Extensions struct {
	AuthorityInfoAccess         *AuthorityInfoAccess           `json:"authority_info_access,omitempty"`
	AuthorityKeyID              string                         `json:"authority_key_id,omitempty"`
	BasicConstraints            *BasicConstraints              `json:"basic_constraints,omitempty"`
	CertificatePolicies         []*CertPolicies                `json:"certificate_policies,omitempty"`
	CrlDistributionPoints       []string                       `json:"crl_distribution_points,omitempty"`
	ExtendedKeyUsage            *ExtendedKeyUsage              `json:"extended_key_usage,omitempty"`
	KeyUsage                    *KeyUsage                      `json:"key_usage,omitempty"`
	SignedCertificateTimestamps []*SignedCertificateTimestamps `json:"signed_certificate_timestamps,omitempty"`
	SubjectAltName              *SubjectAltName                `json:"subject_alt_name,omitempty"`
	SubjectKeyID                string                         `json:"subject_key_id,omitempty"`
}

type AuthorityInfoAccess struct {
	IssuerUrls []string `json:"issuer_urls,omitempty"`
	Ocspurls   []string `json:"ocspurls,omitempty"`
}

type BasicConstraints struct {
	IsCa bool `json:"is_ca,omitempty"`
}

type CertPolicies struct {
	Cps        []string                  `json:"cps,omitempty"`
	ID         string                    `json:"id,omitempty"`
	UserNotice []*CertPoliciesUserNotice `json:"user_notice,omitempty"`
}

type ExtendedKeyUsage struct {
	ClientAuth bool `json:"client_auth,omitempty"`
	ServerAuth bool `json:"server_auth,omitempty"`
}

type KeyUsage struct {
	ContentCommitment bool  `json:"content_commitment,omitempty"`
	DigitalSignature  bool  `json:"digital_signature,omitempty"`
	KeyEncipherment   bool  `json:"key_encipherment,omitempty"`
	Value             int64 `json:"value,omitempty"`
}

type SubjectAltName struct {
	DNSNames    []string `json:"dns_names,omitempty"`
	DNSNamesV2  []string `json:"dns_names_v2,omitempty"`
	IPAddresses []string `json:"ip_addresses,omitempty"`
}

type Issuer struct {
	CommonName         []string `json:"common_name,omitempty"`
	Country            []string `json:"country,omitempty"`
	EmailAddress       []string `json:"email_address,omitempty"`
	Locality           []string `json:"locality,omitempty"`
	Organization       []string `json:"organization,omitempty"`
	OrganizationalUnit []string `json:"organizational_unit,omitempty"`
	Province           []string `json:"province,omitempty"`
}

type Signature struct {
	SelfSigned         bool                `json:"self_signed,omitempty"`
	SignatureAlgorithm *SignatureAlgorithm `json:"signature_algorithm,omitempty"`
	Valid              bool                `json:"valid,omitempty"`
	Value              string              `json:"value,omitempty"`
}

type SignatureAlgorithm struct {
	Name string `json:"name,omitempty"`
	Oid  string `json:"oid,omitempty"`
}

type Subject struct {
	CommonName           []string `json:"common_name,omitempty"`
	CommonNameLowercase  []string `json:"common_name_lowercase,omitempty"`
	Country              []string `json:"country,omitempty"`
	EmailAddress         []string `json:"email_address,omitempty"`
	JurisdictionCountry  []string `json:"jurisdiction_country,omitempty"`
	JurisdictionLocality []string `json:"jurisdiction_locality,omitempty"`
	JurisdictionProvince []string `json:"jurisdiction_province,omitempty"`
	Locality             []string `json:"locality,omitempty"`
	Organization         []string `json:"organization,omitempty"`
	OrganizationalUnit   []string `json:"organizational_unit,omitempty"`
	PostalCode           []string `json:"postal_code,omitempty"`
	Province             []string `json:"province,omitempty"`
	SerialNumber         []string `json:"serial_number,omitempty"`
	StreetAddress        []string `json:"street_address,omitempty"`
}

type SubjectKeyInfo struct {
	EcdsaPublicKey    *EcdsaPublicKey `json:"ecdsa_public_key,omitempty"`
	FingerprintSha256 string          `json:"fingerprint_sha256,omitempty"`
	KeyAlgorithm      *KeyAlgorithm   `json:"key_algorithm,omitempty"`
	RsapublicKey      *RSAPublicKey   `json:"rsapublic_key,omitempty"`
}

type Validity struct {
	End    string `json:"end,omitempty"`
	Length int64  `json:"length,omitempty"`
	Start  string `json:"start,omitempty"`
}

type DomainCertificateIssuerDN struct {
	Raw          string `json:"raw,omitempty"`
	C            string `json:"C,omitempty"`
	CN           string `json:"CN,omitempty"`
	L            string `json:"L,omitempty"`
	O            string `json:"O,omitempty"`
	OU           string `json:"OU,omitempty"`
	ST           string `json:"ST,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
}

type DomainCertificateSubjectDN struct {
	Raw                         string `json:"raw,omitempty"`
	C                           string `json:"C,omitempty"`
	CN                          string `json:"CN,omitempty"`
	L                           string `json:"L,omitempty"`
	O                           string `json:"O,omitempty"`
	OU                          string `json:"OU,omitempty"`
	ST                          string `json:"ST,omitempty"`
	BusinessCategory            string `json:"businessCategory,omitempty"`
	JurisdictionCountry         string `json:"jurisdictionCountry,omitempty"`
	JurisdictionStateOrProvince string `json:"jurisdictionStateOrProvince,omitempty"`
	PostalCode                  string `json:"postalCode,omitempty"`
	SerialNumber                string `json:"serialNumber,omitempty"`
	Street                      string `json:"street,omitempty"`
	EmailAddress                string `json:"emailAddress,omitempty"`
}

type EcdsaPublicKey struct {
	B      string `json:"b,omitempty"`
	Curve  string `json:"curve,omitempty"`
	Gx     string `json:"gx,omitempty"`
	Gy     string `json:"gy,omitempty"`
	Length int64  `json:"length,omitempty"`
	N      string `json:"n,omitempty"`
	P      string `json:"p,omitempty"`
	Pub    string `json:"pub,omitempty"`
	X      string `json:"x,omitempty"`
	Y      string `json:"y,omitempty"`
}

type CertPoliciesUserNotice struct {
	ExplicitText string `json:"explicit_text,omitempty"`
}

type KeyAlgorithm struct {
	Name string `json:"name,omitempty"`
}

type RSAPublicKey struct {
	Exponent int64  `json:"exponent,omitempty"`
	Length   int64  `json:"length,omitempty"`
	Modulus  string `json:"modulus,omitempty"`
}

type SignedCertificateTimestamps struct {
	LogID     string `json:"log_id,omitempty"`
	Signature string `json:"signature,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Version   int64  `json:"version,omitempty"`
}

type CertificatesResponse struct {
	PaginationResponse
	Certificates []*Certificate `json:"certificates"`
}
