package schema

type APIs struct {
	APIs []API `json:"apis"`
}

type API struct {
	Name        string   `json:"name"`
	URIs        string   `json:"uris"`
	UpstreamUrl string   `json:"upstream_url"`
	Checks      []Checks `json:"checks"`
}

type Checks struct {
	CheckType string `json:"checkType,omitempty"`
	Key       string `json:"key,omitempty"`
	Token     string `json:"token,omitempty"`
	Body      string `json:"body,omitempty"`
	Header    string `json:"header,omitempty"`
}
