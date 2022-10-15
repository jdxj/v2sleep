package v2raycore

type Outbound struct {
	Tag            string          `json:"tag"`
	Protocol       string          `json:"protocol"`
	Settings       any             `json:"settings"`
	StreamSettings *StreamSettings `json:"streamSettings,omitempty"`
}

type OutboundConfig struct {
	Outbounds []*Outbound `json:"outbounds"`
}

type Outbounder interface {
	Outbound() (*Outbound, error)
}
