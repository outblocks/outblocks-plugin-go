package communication

type Handshake struct {
	Protocol string `json:"protocol"`
	Addr     string `json:"addr"`
}
