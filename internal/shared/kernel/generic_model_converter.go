package kernel

type GenericModelConverter[model any, domain any] interface {
	ToDomain(model model) (domain, error)
	ToModel(domain domain) (model, error)
}
