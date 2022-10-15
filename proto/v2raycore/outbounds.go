package v2raycore

type Outbound struct {
	Tag      string `json:"tag"`
	Protocol string `json:"protocol"`
	// 看看能不能改成接口
	Settings       any            `json:"settings"`
	StreamSettings StreamSettings `json:"streamSettings"`
}

type OutboundConfig struct {
	Outbounds []Outbound `json:"outbounds"`
}
