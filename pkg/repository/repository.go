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
	GetAll(userId int) ([]todos.TodoList, error)
	GetById(userId int, id int) (todos.TodoList, error)
	Update(userId int, id int, input todos.UpdateListInput) error
	Delete(userId int, id int) error
}

type TodoItem interface {
	Create(listId int, item todos.TodoItem) (int, error)
	GetAll(listId int) ([]todos.TodoItem, error)
	GetById(listId int, id int) (todos.TodoItem, error)
	Update(listId int, id int, input todos.UpdateItemInput) error
	Delete(listId int, id int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
