package components

type WebhookView struct {
	ID                  string
	Name                string
	Description         string
	URL                 string
	RequestCount        int64
	CreatedAt           string
	ResponseStatus      int
	ResponseContentType string
	ResponseBody        string
}

type RequestView struct {
	ID            string
	WebhookID     string
	Method        string
	Path          string
	SourceIP      string
	ContentType   string
	ContentLength int64
	CreatedAt     string
}

type DetailRequestView struct {
	ID            string
	WebhookID     string
	Method        string
	Path          string
	Headers       map[string]string
	QueryParams   map[string]string
	Body          string
	ContentType   string
	SourceIP      string
	ContentLength int64
	CreatedAt     string
	HeadersJSON   string
	QueryJSON     string
	BodyFormatted string
	BodyLanguage  string

	ResponseStatus        int
	ResponseHeaders       map[string]string
	ResponseHeadersJSON   string
	ResponseBody          string
	ResponseBodyFormatted string
	ResponseBodyLanguage  string
}
