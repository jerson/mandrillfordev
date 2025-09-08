package types

import "time"

type MandrillRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"` // to|cc|bcc
}

type MandrillAttachment struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"` // base64-encoded
}

// Merge variables
type MandrillMergeVar struct {
	Name    string `json:"name"`
	Content any    `json:"content"`
}

type MandrillRcptMergeVars struct {
	Rcpt string             `json:"rcpt"`
	Vars []MandrillMergeVar `json:"vars"`
}

type MandrillMessage struct {
	HTML                    string                  `json:"html,omitempty"`
	Text                    string                  `json:"text,omitempty"`
	Subject                 string                  `json:"subject,omitempty"`
	FromEmail               string                  `json:"from_email"`
	FromName                string                  `json:"from_name,omitempty"`
	To                      []MandrillRecipient     `json:"to"`
	Headers                 map[string]string       `json:"headers,omitempty"`
	Important               bool                    `json:"important,omitempty"`
	TrackOpens              bool                    `json:"track_opens,omitempty"`
	TrackClicks             bool                    `json:"track_clicks,omitempty"`
	AutoText                bool                    `json:"auto_text,omitempty"`
	AutoHTML                bool                    `json:"auto_html,omitempty"`
	InlineCSS               bool                    `json:"inline_css,omitempty"`
	URLStripQS              bool                    `json:"url_strip_qs,omitempty"`
	PreserveRecipients      bool                    `json:"preserve_recipients,omitempty"`
	ViewContentLink         bool                    `json:"view_content_link,omitempty"`
	BccAddress              string                  `json:"bcc_address,omitempty"`
	TrackingDomain          string                  `json:"tracking_domain,omitempty"`
	SigningDomain           string                  `json:"signing_domain,omitempty"`
	ReturnPathDomain        string                  `json:"return_path_domain,omitempty"`
	Merge                   bool                    `json:"merge,omitempty"`
	MergeLanguage           string                  `json:"merge_language,omitempty"`
	GlobalMergeVars         []MandrillMergeVar      `json:"global_merge_vars,omitempty"`
	MergeVars               []MandrillRcptMergeVars `json:"merge_vars,omitempty"`
	Tags                    []string                `json:"tags,omitempty"`
	Subaccount              string                  `json:"subaccount,omitempty"`
	GoogleAnalyticsDomains  []string                `json:"google_analytics_domains,omitempty"`
	GoogleAnalyticsCampaign string                  `json:"google_analytics_campaign,omitempty"`
	Metadata                map[string]string       `json:"metadata,omitempty"`
	RecipientMetadata       []struct {
		Rcpt   string            `json:"rcpt"`
		Values map[string]string `json:"values"`
	} `json:"recipient_metadata,omitempty"`
	Attachments []MandrillAttachment `json:"attachments,omitempty"`
	Images      []MandrillAttachment `json:"images,omitempty"`
}

// Requests
type SendRequest struct {
	Key     string          `json:"key"`
	Message MandrillMessage `json:"message"`
	Async   bool            `json:"async,omitempty"`
	IPPool  string          `json:"ip_pool,omitempty"`
	SendAt  string          `json:"send_at,omitempty"`
}

type TemplateContent struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type SendTemplateRequest struct {
	Key             string            `json:"key"`
	TemplateName    string            `json:"template_name"`
	TemplateContent []TemplateContent `json:"template_content"`
	Message         MandrillMessage   `json:"message"`
	Async           bool              `json:"async,omitempty"`
	IPPool          string            `json:"ip_pool,omitempty"`
	SendAt          string            `json:"send_at,omitempty"`
}

type SendRawRequest struct {
	Key              string   `json:"key"`
	RawMessage       string   `json:"raw_message"`
	FromEmail        string   `json:"from_email,omitempty"`
	FromName         string   `json:"from_name,omitempty"`
	To               []string `json:"to,omitempty"`
	ReturnPathDomain string   `json:"return_path_domain,omitempty"`
	Async            bool     `json:"async,omitempty"`
	IPPool           string   `json:"ip_pool,omitempty"`
	SendAt           string   `json:"send_at,omitempty"`
}

type ParseRequest struct {
	Key        string `json:"key"`
	RawMessage string `json:"raw_message"`
}

type InfoRequest struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type ContentRequest struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type SearchRequest struct {
	Key      string   `json:"key"`
	Query    string   `json:"query,omitempty"`
	DateFrom string   `json:"date_from,omitempty"`
	DateTo   string   `json:"date_to,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Senders  []string `json:"senders,omitempty"`
	Limit    int      `json:"limit,omitempty"`
}

type ListScheduledRequest struct {
	Key string `json:"key"`
	To  string `json:"to,omitempty"`
}

type CancelScheduledRequest struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type RescheduleRequest struct {
	Key    string `json:"key"`
	Id     string `json:"id"`
	SendAt string `json:"send_at"`
}

// Responses & record types

type SendResult struct {
	Email        string `json:"email"`
	Status       string `json:"status"` // sent|rejected|queued|scheduled|invalid
	ID           string `json:"_id"`
	RejectReason string `json:"reject_reason,omitempty"`
}

type MessageRecord struct {
	ID           string
	CreatedAt    time.Time
	ScheduledAt  *time.Time
	SentAt       *time.Time
	Status       string
	RejectReason string
	Message      MandrillMessage
	From         string
	To           []string
	Subject      string
	Tags         []string
	Raw          []byte
	TemplateName string
}

// Template management
type Template struct {
	Name          string     `json:"name"`
	FromEmail     string     `json:"from_email,omitempty"`
	FromName      string     `json:"from_name,omitempty"`
	Subject       string     `json:"subject,omitempty"`
	Code          string     `json:"code,omitempty"`         // draft HTML
	Text          string     `json:"text,omitempty"`         // draft text
	PublishedCode string     `json:"publish_code,omitempty"` // published HTML
	PublishedText string     `json:"publish_text,omitempty"` // published text
	Labels        []string   `json:"labels,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
}

type TemplateAddRequest struct {
	Key       string   `json:"key"`
	Name      string   `json:"name"`
	FromEmail string   `json:"from_email,omitempty"`
	FromName  string   `json:"from_name,omitempty"`
	Subject   string   `json:"subject,omitempty"`
	Code      string   `json:"code,omitempty"`
	Text      string   `json:"text,omitempty"`
	Publish   bool     `json:"publish,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}

type TemplateInfoRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Update supports partial updates via pointer fields (nil means unchanged)
type TemplateUpdateRequest struct {
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	FromEmail *string   `json:"from_email"`
	FromName  *string   `json:"from_name"`
	Subject   *string   `json:"subject"`
	Code      *string   `json:"code"`
	Text      *string   `json:"text"`
	Publish   *bool     `json:"publish"`
	Labels    *[]string `json:"labels"`
}

type TemplatePublishRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type TemplateDeleteRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type TemplateListRequest struct {
	Key   string `json:"key"`
	Label string `json:"label,omitempty"`
}

type TemplateTimeSeriesRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type TemplateRenderRequest struct {
	Key             string            `json:"key"`
	TemplateName    string            `json:"template_name"`
	TemplateContent []TemplateContent `json:"template_content"`
}
