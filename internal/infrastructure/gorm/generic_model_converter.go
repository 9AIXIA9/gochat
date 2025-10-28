package gorm

type GenericModelConverter[model any, domain any] interface {
	ToDomain(model model) domain
	ToModel(domain domain) model
}
