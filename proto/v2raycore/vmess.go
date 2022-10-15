package v2raycore

type User struct {
	ID       string `json:"id"`
	Security string `json:"security"`
}

type VNext struct {
	Address string `json:"address"`
	Port    int64  `json:"port"`
	Users   []User `json:"users"`
}

type VmessSettings struct {
	VNext []VNext `json:"vnext"`
}
