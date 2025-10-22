package kernel

type HashEncryptor interface {
	Encrypt(origin string) (hash string, err error)
}

type HashComparator interface {
	Compare(hash string, origin string) error
}
