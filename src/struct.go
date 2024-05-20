package src

// Struct following CIM Web data model
type BrowserArtifact struct {
	ArtifactType    string `json:"artifact_type,omitempty"`
	Dest            string `json:"dest,omitempty"`
	DestPort        int    `json:"dest_port,omitempty"`
	Src             string `json:"src,omitempty"`
	User            string `json:"user,omitempty"`
	App             string `json:"app,omitempty"`
	Action          string `json:"action,omitempty"`
	Cached          bool   `json:"cached,omitempty"`
	Cookie          string `json:"cookie,omitempty"`
	Url             string `json:"url,omitempty"`
	UrlDomain       string `json:"url_domain,omitempty"`
	HttpMethod      string `json:"http_method,omitempty"`
	HttpReferrer    string `json:"http_referrer,omitempty"`
	HttpUserAgent   string `json:"http_user_agent,omitempty"`
	HttpContentType string `json:"http_content_type,omitempty"`
	Duration        int    `json:"duration,omitempty"`
	Status          string `json:"status,omitempty"`
	BytesIn         int    `json:"bytes_in,omitempty"`
	BytesOut        int    `json:"bytes_out,omitempty"`
	Timestamp       int    `json:"timestamp,omitempty"`
	TimestampType   string `json:"timestamp_type,omitempty"`

	// Additional fields from Firefox
	Typed         int    `json:"typed,omitempty"`
	VisitCount    int    `json:"visit_count,omitempty"`
	Title         string `json:"title,omitempty"`
	BookmarkTitle string `json:"bookmark_title,omitempty"`

	// Downloads additional fields
	Metadata string `json:"metadata,omitempty"`
	Filename string `json:"filename,omitempty"`
	MimeType string `json:"mime_type,omitempty"`

	// Formhistory additional fields
	Fieldname string `json:"fieldname,omitempty"`
	Value     string `json:"value,omitempty"`

	// Addons additional fields
	AddonName       string `json:"addon_name,omitempty"`
	AddonType       string `json:"addon_type,omitempty"`
	Active          bool   `json:"active,omitempty"`
	Visible         bool   `json:"visible,omitempty"`
	Description     string `json:"description,omitempty"`
	FullDescription string `json:"full_description,omitempty"`
	Version         string `json:"version,omitempty"`
	SourceURI       string `json:"source_uri,omitempty"`
	HomePageURL     string `json:"home_page_url,omitempty"`
	AboutURL        string `json:"about_url,omitempty"`
	ReviewURL       string `json:"review_url,omitempty"`
	CreatorName     string `json:"creator_name,omitempty"`
	CreatorURL      string `json:"creator_url,omitempty"`
	AverageRating   float32
	RatingCount     int
}
