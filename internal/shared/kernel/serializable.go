package kernel

type Serializer interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}
