package proto

type ConfType int

const (
	ShareLink    ConfType = 1
	ClashSubAddr ConfType = 2
	V2raySubAddr ConfType = 3
)

type V2rayNG interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}
