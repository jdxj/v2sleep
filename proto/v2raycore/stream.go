package v2raycore

type StreamSettings struct {
	Network      string        `json:"network"`
	Security     string        `json:"security"`
	TlsSettings  *TlsSettings  `json:"tlsSettings,omitempty"`
	TcpSettings  *TcpSettings  `json:"tcpSettings,omitempty"`
	HttpSettings *HttpSettings `json:"httpSettings,omitempty"`
}

type TlsSettings struct {
	ServerName    string `json:"serverName"`
	AllowInsecure bool   `json:"allowInsecure"`
}

// <<< tcp settings begin

type Request struct {
	Version string              `json:"version"`
	Method  string              `json:"method"`
	Path    []string            `json:"path"`
	Headers map[string][]string `json:"headers"`
}

type Response struct {
	Version string              `json:"version"`
	Status  string              `json:"status"`
	Reason  string              `json:"reason"`
	Headers map[string][]string `json:"headers"`
}

type HttpHeader struct {
	Type     string    `json:"type"`
	Request  *Request  `json:"request,omitempty"`
	Response *Response `json:"response,omitempty"`
}

type TcpSettings struct {
	Header HttpHeader `json:"header"`
}

// tcp settings end >>>

type HttpSettings struct {
	Host    []string            `json:"host"`
	Path    string              `json:"path"`
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers,omitempty"`
}
