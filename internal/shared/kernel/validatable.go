package kernel

type Validatable interface {
	Validate() error
}
