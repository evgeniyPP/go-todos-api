package repository

import (
	"fmt"
	"strings"

	"github.com/evgeniyPP/go-todos-api"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todos.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}
	createListItemsQuery := fmt.Sprintf(
		"INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	if _, err := tx.Exec(createListItemsQuery, listId, itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(listId int) ([]todos.TodoItem, error) {
	var items []todos.TodoItem
	query := fmt.Sprintf(
		`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
		 INNER JOIN %s li on ti.id = li.item_id 
		 WHERE li.list_id = $1`,
		todoItemsTable, listsItemsTable)
	err := r.db.Select(&items, query, listId)
	return items, err
}

func (r *TodoItemPostgres) GetById(listId int, id int) (todos.TodoItem, error) {
	var item todos.TodoItem
	query := fmt.Sprintf(
		`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
		 INNER JOIN %s li on ti.id = li.item_id 
		 WHERE li.list_id = $1 AND li.item_id = $2`,
		todoItemsTable, listsItemsTable)
	err := r.db.Get(&item, query, listId, id)
	return item, err
}

func (r *TodoItemPostgres) Update(listId int, id int, input todos.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s ti SET %s FROM %s li
		 WHERE ti.id = li.item_id AND li.list_id=$%d AND li.item_id=$%d`,
		todoItemsTable, setQuery, listsItemsTable, argId, argId+1)
	args = append(args, listId, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoItemPostgres) Delete(listId int, id int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s ti USING %s li 
		 WHERE ti.id = li.item_id AND li.list_id = $1 AND li.item_id = $2`,
		todoItemsTable, listsItemsTable)
	_, err := r.db.Exec(query, listId, id)
	return err
}
