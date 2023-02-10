package fs

type Client interface {
	Write(id string, data []byte) error
	Read(id string) ([]byte, error)
	Remove(id string) error
}
