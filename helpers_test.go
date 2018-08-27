package rest_test

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	rest "github.com/ktnyt/go-rest"
)

type Todo struct {
	Key       string
	Content   string
	CreatedAt time.Time
	Done      bool
}

func (t *Todo) Validate() error {
	if len(t.Content) == 0 {
		return fmt.Errorf("todo content is empty")
	}
	if t.CreatedAt.After(time.Now()) {
		return fmt.Errorf("todo is created in the future")
	}
	return nil
}

func (t *Todo) MakeKey(i int) string {
	t.Key = strconv.Itoa(i)
	t.Done = i&1 == 1
	return t.Key
}

func (t *Todo) Merge(other rest.Model) error {
	switch other := other.(type) {
	case *Todo:
		t.Content = other.Content
		t.CreatedAt = other.CreatedAt
		t.Done = other.Done
		return nil
	default:
		return fmt.Errorf("attempted to merge non-Todo object")
	}
}

func Construct() rest.Model {
	return &Todo{}
}

func RandomTodo() *Todo {
	return &Todo{
		Content:   uuid.Must(uuid.NewV4()).String(),
		CreatedAt: time.Now(),
		Done:      false,
	}
}

func InvalidTodo() *Todo {
	return &Todo{
		Content:   "",
		CreatedAt: time.Now(),
		Done:      false,
	}
}
