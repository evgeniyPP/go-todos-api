package repository

import (
	"github.com/evgeniyPP/go-todos-api"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todos.User) (int, error)
	GetUser(username, password string) (todos.User, error)
}

type TodoList interface {
	Create(userId int, list todos.TodoList) (int, error)
}

type TodoItem interface{}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
	}
}
