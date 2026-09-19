package todo

import (
	"context"

	"gorm.io/gorm"
)

type ListQuery struct {
	Desc   bool
	Done   *bool
	Limit  int
	Offset int
}

type Repository interface {
	Create(ctx context.Context, todo Todo) (int64, error)
	Get(ctx context.Context, id int64) (Model, error)
	List(ctx context.Context, q ListQuery) ([]Model, error)
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

func (r *GormRepository) Get(ctx context.Context, id int64) (Model, error) {
	return gorm.G[Model](r.db).Where("id = ?", id).Take(ctx)
}

func (r *GormRepository) List(ctx context.Context, q ListQuery) ([]Model, error) {
	order := "id ASC"
	if q.Desc {
		order = "id DESC"
	}
	g := gorm.G[Model](r.db).Order(order)
	if q.Done != nil {
		g = g.Where("is_done = ?", *q.Done)
	}
	if q.Limit > 0 {
		g = g.Limit(q.Limit)
	}
	if q.Offset > 0 {
		g = g.Offset(q.Offset)
	}
	return g.Find(ctx)
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
