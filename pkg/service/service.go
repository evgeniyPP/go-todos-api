package service

import (
	"github.com/evgeniyPP/go-todos-api"
	"github.com/evgeniyPP/go-todos-api/pkg/repository"
)

type Authorization interface {
	CreateUser(user todos.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list todos.TodoList) (int, error)
	GetAll(userId int) ([]todos.TodoList, error)
	GetById(userId int, id int) (todos.TodoList, error)
	Update(userId int, id int, input todos.UpdateListInput) error
	Delete(userId int, id int) error
}

type TodoItem interface {
	Create(userId int, listId int, item todos.TodoItem) (int, error)
	GetAll(userId int, listId int) ([]todos.TodoItem, error)
	GetById(userId int, id int) (todos.TodoItem, error)
	Update(userId int, id int, input todos.UpdateItemInput) error
	Delete(userId int, id int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
