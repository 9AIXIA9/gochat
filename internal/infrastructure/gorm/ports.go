package gorm

type Model interface {
	TableName() string
}

type GenericModelConverter[model any, domain any] interface {
	ToDomain(model model) domain
	ToModel(domain domain) model
}
