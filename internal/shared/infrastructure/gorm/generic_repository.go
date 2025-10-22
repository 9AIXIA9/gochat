package gorm

import (
	"context"
	"errors"
	"gochat/internal/shared/kernel"

	myErrors "gochat/internal/shared/errors"
	"gorm.io/gorm"
)

type Repository[
	GormModel any,
	DomainModel any,
] struct {
	db                        *gorm.DB
	gormGenericModelConverter kernel.GenericModelConverter[*GormModel, *DomainModel]
}

func NewRepository[
	GormModel any,
	DomainModel any,
](db *gorm.DB, converter kernel.GenericModelConverter[*GormModel, *DomainModel]) *Repository[GormModel, DomainModel] {
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
	gormModelPointer, err := repo.gormGenericModelConverter.ToModel(domainModel)
	if err != nil {
		return err
	}

	if err := repo.db.WithContext(ctx).Create(gormModelPointer).Error; err != nil {
		return TranslateError(err)
	}

	updated, err := repo.gormGenericModelConverter.ToDomain(gormModelPointer)
	if err != nil {
		return err
	}

	if updated != nil {
		*domainModel = *updated // Write back to domain model (回写领域模型)}
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
	return repo.gormGenericModelConverter.ToDomain(gormModelPointer)
}

func (repo *Repository[GormModel, DomainModel]) FindMany(
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
		domainModel, err := repo.gormGenericModelConverter.ToDomain(gormModel)
		if err != nil {
			return nil, err
		}
		domainModels = append(domainModels, domainModel)
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
	result := repo.db.WithContext(ctx).Where(query, args...).Delete(new(GormModel))
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return myErrors.ErrNotFound
	}
	return nil
}

func (repo *Repository[
	GormModel,
	DomainModel,
]) HardDelete(ctx context.Context, query interface{}, args ...interface{}) error {
	result := repo.db.WithContext(ctx).Unscoped().Where(query, args...).Delete(new(GormModel))
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return myErrors.ErrNotFound
	}
	return nil
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
