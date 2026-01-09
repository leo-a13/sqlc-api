package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	db "todo-app/internal/db/sqlc"
)

func RegisterTodoRoutes(r chi.Router, q *db.Queries) {
	r.Post("/todos", createTodo(q))
	r.Get("/todos", listTodos(q))
	r.Get("/todos/{id}", getTodo(q))
	r.Patch("/todos/{id}", updateTodo(q))
	r.Delete("/todos/{id}", deleteTodo(q))
}

func createTodo(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		todo, err := q.CreateTodo(r.Context(), db.CreateTodoParams{
			Title:       body.Title,
			Description: sql.NullString{},
			Priority:    3,
		})
		if err != nil {
			http.Error(w, "failed to create todo", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func listTodos(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := q.ListTodos(r.Context())
		if err != nil {
			http.Error(w, "failed to list todos", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(todos); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func getTodo(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		todoID, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "invalid todo id", http.StatusBadRequest)
			return
		}

		todo, err := q.GetTodo(r.Context(), todoID)
		if err != nil {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func updateTodo(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var body struct {
			Completed bool `json:"completed"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		todoID, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "invalid todo id", http.StatusBadRequest)
			return
		}

		existing, err := q.GetTodo(r.Context(), todoID)
		if err != nil {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}

		updated, err := q.UpdateTodo(r.Context(), db.UpdateTodoParams{
			ID:          todoID,
			Title:       existing.Title,
			Description: existing.Description,
			Completed:   body.Completed,
			Priority:    existing.Priority,
		})
		if err != nil {
			http.Error(w, "failed to update todo", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(updated); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func deleteTodo(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		todoID, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "invalid todo id", http.StatusBadRequest)
			return
		}

		if err := q.DeleteTodo(r.Context(), todoID); err != nil {
			http.Error(w, "failed to delete todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
