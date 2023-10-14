package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/iruma-tea/go_todo_app/entity"
	"github.com/iruma-tea/go_todo_app/store"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		ResponseJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	err := at.Validator.Struct(b)
	if err != nil {
		ResponseJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title:   b.Title,
		Status:  entity.TaskStatusTodo,
		Created: time.Now(),
	}

	id, err := store.Tasks.Add(t)
	if err != nil {
		ResponseJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return

	}
	rsp := struct {
		ID int `json:"id"`
	}{ID: int(id)}

	ResponseJSON(ctx, w, rsp, http.StatusOK)
}
