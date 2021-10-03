package repository

import (
	"fmt"
	"strings"

	"github.com/evgeniyPP/go-todos-api"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todos.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	createUsersListQuery := fmt.Sprintf(
		"INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	if _, err := tx.Exec(createUsersListQuery, userId, id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todos.TodoList, error) {
	var lists []todos.TodoList
	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description FROM %s tl 
		 INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1`,
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

func (r *TodoListPostgres) GetById(userId int, id int) (todos.TodoList, error) {
	var list todos.TodoList
	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description FROM %s tl 
		 INNER JOIN %s ul on tl.id = ul.list_id 
		 WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, id)
	return list, err
}

func (r *TodoListPostgres) Update(userId int, id int, input todos.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s tl SET %s FROM %s ul 
		 WHERE tl.id = ul.list_id AND ul.list_id = $%d AND ul.user_id = $%d`,
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, id, userId)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoListPostgres) Delete(userId int, id int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s tl USING %s ul 
		 WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, id)
	return err
}
