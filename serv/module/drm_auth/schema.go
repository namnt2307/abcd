package drm_auth

const (
	COLLECTION_DRM_RESP_LOG = "drm_resp_log"
)

type CastlabLicenseRequest struct {
	Asset      string
	Variant    string
	User       string
	Session    string
	Client     string
	DrmScheme  string
	ClientInfo struct {
		Manufacturer string
		Model        string
		Version      string
		CertType     string
		DrmVersion   string
		SecLevel     interface{}
	}
	RequestMetadata struct {
		RemoteAddr string
		UserAgent  string
	}
}
type CastlabConfig struct {
	Vod_output_protect             bool
	Vod_store_licence              bool
	Vod_output_protect_analogue    bool
	Vod_output_protect_digital     bool
	Vod_output_protect_enforce     bool
	Livetv_output_protect          bool
	Livetv_store_licence           bool
	Livetv_output_protect_analogue bool
	Livetv_output_protect_digital  bool
	Livetv_output_protect_enforce  bool
	Save_request                   bool
}
type CastlabLicenseResponseSuccess struct {
	// Ref          []string `json:"ref"`
	AccountingID     string         `json:"accountingId"`
	AssetID          string         `json:"assetId"`
	VariantID        string         `json:"variantId"`
	StoreLicense     bool           `json:"storeLicense"`
	Profile          CastlabProfile `json:"profile"`
	OutputProtection interface{}    `json:"outputProtection,omitempty"`
	// Op struct {
	// 	Config interface{} `json:"config"`
	// } `json:"op"`
	// Csl struct {
	// 	CslTrackingID   string        `json:"cslTrackingId"`
	// 	SessionDuration int64         `json:"sessionDuration"`
	// 	Rules           []interface{} `json:"rules"`
	// } `json:"csl"`
}
type OutputProtection struct {
	Digital  bool `json:"digital"`
	Analogue bool `json:"analogue"`
	Enforce  bool `json:"enforce"`
}
type CastlabLicenseResponseError struct {
	Message     string `json:"message"`
	RedirectUrl string `json:"redirectUrl"`
}

type CastlabProfile struct {
	Rental struct {
		AbsoluteExpiration string `json:"absoluteExpiration,omitempty"`
		RelativeExpiration string `json:"relativeExpiration,omitempty"`
		PlayDuration       int64  `json:"playDuration"`
	} `json:"rental"`
}

type DrmLogInfo struct {
	Log_data struct {
		Request  interface{} `json:"request"`
		Response interface{} `json:"response"`
	} `json:"log_data"`
	Drm_service_name string `json:"drm_service_name"`
	Created_at       int64  `json:"created_at"`
}

type VieONLicenseRequest struct {
	Id      string
	Token   string
	Type    string
	User_id string
}
