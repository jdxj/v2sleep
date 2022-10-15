package v2raycore

type SS struct {
	Email    string `json:"email,omitempty"`
	Address  string `json:"address"`
	Port     int64  `json:"port"`
	Method   string `json:"method"`
	Password string `json:"password"`
	Level    int    `json:"level,omitempty"`
	IvCheck  bool   `json:"ivCheck,omitempty"`
}

type SSSettings struct {
	Servers []SS `json:"servers"`
}
