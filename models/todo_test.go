package models

import (
	"testing"
)

func TestTodoCreation(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		item      string
		completed int
		want      Todo
	}{
		{
			name:      "Create basic todo",
			id:        1,
			item:      "Buy groceries",
			completed: 0,
			want: Todo{
				Id:        1,
				Item:      "Buy groceries",
				Completed: 0,
			},
		},
		{
			name:      "Create completed todo",
			id:        2,
			item:      "Walk the dog",
			completed: 1,
			want: Todo{
				Id:        2,
				Item:      "Walk the dog",
				Completed: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Todo{
				Id:        tt.id,
				Item:      tt.item,
				Completed: tt.completed,
			}

			if got.Id != tt.want.Id {
				t.Errorf("Todo.Id = %v, want %v", got.Id, tt.want.Id)
			}
			if got.Item != tt.want.Item {
				t.Errorf("Todo.Item = %v, want %v", got.Item, tt.want.Item)
			}
			if got.Completed != tt.want.Completed {
				t.Errorf("Todo.Completed = %v, want %v", got.Completed, tt.want.Completed)
			}
		})
	}
}

func TestTodoValidFields(t *testing.T) {
	todo := Todo{
		Id:        1,
		Item:      "Test todo",
		Completed: 0,
	}

	if todo.Id <= 0 {
		t.Error("Todo ID should be positive")
	}

	if len(todo.Item) == 0 {
		t.Error("Todo Item should not be empty")
	}

	if todo.Completed != 0 && todo.Completed != 1 {
		t.Error("Todo Completed should be either 0 or 1")
	}
}