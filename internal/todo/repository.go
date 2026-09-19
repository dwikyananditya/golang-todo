package todo

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, todo Todo) error
	Get(ctx context.Context, id int) (Model, error)
	List(ctx context.Context) ([]Model, error)
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

func (r *GormRepository) Create(ctx context.Context, todo Todo) error {
	return gorm.G[Model](r.db).Create(ctx, &Model{
		Title:       todo.Title,
		Description: todo.Description,
		IsDone:      todo.IsDone,
	})
}

func (r *GormRepository) Get(ctx context.Context, id int) (Model, error) {
	return gorm.G[Model](r.db).Where("id = ?", id).Take(ctx)
}

func (r *GormRepository) List(ctx context.Context) ([]Model, error) {
	return gorm.G[Model](r.db).Order("id ASC").Find(ctx)
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
