package todo

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, todo Todo) (int64, error)
	First(ctx context.Context, desc bool) (Model, error)
	List(ctx context.Context, desc bool) ([]Model, error)
	Update(ctx context.Context, id int, todo Todo) (int, error)
	Delete(ctx context.Context, id int) (int, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Migrate() error {
	return r.db.AutoMigrate(&Model{})
}

func (r *GormRepository) Create(ctx context.Context, todo Todo) (int64, error) {
	m := Model{
		Title:       todo.Title,
		Description: todo.Description,
		IsDone:      todo.IsDone,
	}
	if err := gorm.G[Model](r.db).Create(ctx, &m); err != nil {
		return 0, err
	}
	return int64(m.ID), nil
}

func (r *GormRepository) First(ctx context.Context, desc bool) (Model, error) {
	order := "id ASC"
	if desc {
		order = "id DESC"
	}
	return gorm.G[Model](r.db).Order(order).Take(ctx)
}

func (r *GormRepository) List(ctx context.Context, desc bool) ([]Model, error) {
	order := "id ASC"
	if desc {
		order = "id DESC"
	}
	return gorm.G[Model](r.db).Order(order).Find(ctx)
}

func (r *GormRepository) Update(ctx context.Context, id int, todo Todo) (int, error) {
	return gorm.G[Model](r.db).
		Where("id = ?", id).
		Updates(ctx, Model{
			Title:       todo.Title,
			Description: todo.Description,
			IsDone:      todo.IsDone,
		})
}

func (r *GormRepository) Delete(ctx context.Context, id int) (int, error) {
	return gorm.G[Model](r.db).
		Where("id = ?", id).
		Delete(ctx)
}
