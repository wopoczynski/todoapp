package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *TodoModel) (*uuid.UUID, error)
	GetAll(ctx context.Context) *[]TodoModel
	Get(ctx context.Context, id uuid.UUID) (*TodoModel, error)
	Update(ctx context.Context, model *TodoModel) error
	DeleteAll(ctx context.Context) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MysqlTodoRepository struct {
	db *gorm.DB
}

func NewMysqlTodoRepository(db *gorm.DB) TodoRepository {
	return &MysqlTodoRepository{
		db: db,
	}
}

func (r *MysqlTodoRepository) Create(ctx context.Context, todo *TodoModel) (*uuid.UUID, error) {
	err := r.db.WithContext(ctx).Create(todo).Error
	if err != nil {
		return nil, err
	}

	return &todo.ID, nil
}

func (r *MysqlTodoRepository) GetAll(ctx context.Context) *[]TodoModel {
	var todos []TodoModel
	result := r.db.WithContext(ctx).Find(&todos)
	if err := result.Error; err != nil {
		log.Error().Err(err).Msg("unable to fetch todos")
		return nil
	}

	return &todos
}

func (r *MysqlTodoRepository) Get(ctx context.Context, id uuid.UUID) (*TodoModel, error) {
	var todo TodoModel
	err := r.db.WithContext(ctx).First(&todo, id).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to fetch todos")
		return nil, err
	}

	return &todo, nil
}

func (r *MysqlTodoRepository) Update(ctx context.Context, todo *TodoModel) error {
	err := r.db.WithContext(ctx).Save(todo).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *MysqlTodoRepository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("1 = 1").Delete(&TodoModel{}).Error
}

func (r *MysqlTodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&TodoModel{}, id).Error
}
