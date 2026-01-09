package gin

type RequestPointers[Request any] interface {
	*Request
	Bindable
}
