package gorm

import (
	"context"
	"errors"

	myErrors "gochat/internal/shared/errors"

	"gorm.io/gorm"
)

type Repository[
	GormModel any,
	DomainModel any,
] struct {
	db                        *gorm.DB
	gormGenericModelConverter GenericModelConverter[*GormModel, *DomainModel]
}

func NewRepository[
	GormModel any,
	DomainModel any,
](db *gorm.DB, converter GenericModelConverter[*GormModel, *DomainModel]) *Repository[GormModel, DomainModel] {
	return &Repository[GormModel, DomainModel]{
		db:                        db,
		gormGenericModelConverter: converter,
	}
}

func (repo *Repository[
	GormModel,
	DomainModel,
]) WithTx(tx *gorm.DB) *Repository[GormModel, DomainModel] {
	r := *repo
	r.db = tx
	return &r
}

func (repo *Repository[GormModel, DomainModel]) Save(
	ctx context.Context,
	domainModel *DomainModel,
) error {
	if domainModel == nil {
		return myErrors.ErrEmptyPointer
	}
	gormModelPointer := repo.gormGenericModelConverter.ToModel(domainModel)

	if err := repo.db.WithContext(ctx).Create(gormModelPointer).Error; err != nil {
		return TranslateError(err)
	}
	return nil
}

func (repo *Repository[GormModel, DomainModel]) Saves(
	ctx context.Context,
	domainModels []*DomainModel,
) error {
	if domainModels == nil {
		return myErrors.ErrEmptyPointer
	}

	if len(domainModels) == 0 {
		return nil
	}

	// 将领域模型切片逐个转换为 GORM 模型切片
	gormModels := make([]*GormModel, 0, len(domainModels))
	for i := range domainModels {
		if domainModels[i] == nil {
			return myErrors.ErrEmptyPointer
		}
		gormModels = append(gormModels, repo.gormGenericModelConverter.ToModel(domainModels[i]))
	}

	// 批量创建
	if err := repo.db.WithContext(ctx).Create(&gormModels).Error; err != nil {
		return TranslateError(err)
	}
	return nil
}

func (repo *Repository[
	GormModel,
	DomainModel,
]) Find(ctx context.Context, scopes ...Scope) (*DomainModel, error) {
	gormModelPointer := new(GormModel)
	if err := applyScopes(repo.db.WithContext(ctx), scopes...).First(gormModelPointer).Error; err != nil {
		return nil, TranslateError(err)
	}
	return repo.gormGenericModelConverter.ToDomain(gormModelPointer), nil
}

func (repo *Repository[GormModel, DomainModel]) Finds(
	ctx context.Context,
	scopes ...Scope,
) ([]*DomainModel, error) {
	var gormModels []GormModel
	if err := applyScopes(repo.db.WithContext(ctx), scopes...).Find(&gormModels).Error; err != nil {
		return nil, TranslateError(err)
	}
	domainModels := make([]*DomainModel, 0, len(gormModels))
	for i := range gormModels {
		gormModel := &gormModels[i]
		domainModels = append(domainModels, repo.gormGenericModelConverter.ToDomain(gormModel))
	}
	return domainModels, nil
}

func (repo *Repository[GormModel, DomainModel]) Count(
	ctx context.Context,
	scopes ...Scope,
) (int64, error) {
	var count int64
	if err := applyScopes(repo.db.WithContext(ctx).Model(new(GormModel)), scopes...).Count(&count).Error; err != nil {
		return 0, TranslateError(err)
	}
	return count, nil
}

func (repo *Repository[GormModel, DomainModel]) Exists(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (bool, error) {
	gm := new(GormModel)
	err := repo.db.WithContext(ctx).Select("1").Where(query, args...).Take(gm).Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func (repo *Repository[
	GormModel,
	DomainModel,
]) Delete(ctx context.Context, query interface{}, args ...interface{}) error {
	return TranslateError(repo.db.WithContext(ctx).Where(query, args...).Delete(new(GormModel)).Error)
}

func (repo *Repository[
	GormModel,
	DomainModel,
]) HardDelete(ctx context.Context, query interface{}, args ...interface{}) error {
	return TranslateError(repo.db.WithContext(ctx).Unscoped().Where(query, args...).Delete(new(GormModel)).Error)
}

func (repo *Repository[GormModel, DomainModel]) Update(
	ctx context.Context,
	query interface{},
	column string,
	value interface{},
	args ...interface{},
) error {
	exists, err := repo.Exists(ctx, query, args...)
	if err != nil {
		return err
	}
	if !exists {
		return myErrors.ErrNotFound
	}
	result := repo.db.WithContext(ctx).Model(new(GormModel)).Where(query, args...).Update(column, value)
	return TranslateError(result.Error)
}

func (repo *Repository[GormModel, DomainModel]) Updates(
	ctx context.Context,
	query interface{},
	values map[string]any,
	args ...interface{},
) error {
	exists, err := repo.Exists(ctx, query, args...)
	if err != nil {
		return err
	}
	if !exists {
		return myErrors.ErrNotFound
	}
	result := repo.db.WithContext(ctx).Model(new(GormModel)).Where(query, args...).Updates(values)
	return TranslateError(result.Error)
}
