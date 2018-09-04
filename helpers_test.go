package rest_test

import (
	"context"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	rest "github.com/ktnyt/go-rest"
)

func init() {
	gob.Register(&Todo{})
}

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

func (t *Todo) Merge(other interface{}) error {
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

func NewTodo() rest.Model {
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

func Filter(ctx context.Context) rest.Filter {
	return func(value interface{}) bool {
		switch done := ctx.Value("done").(type) {
		case string:
			if done == "true" {
				return value.(*Todo).Done
			}
			if done == "false" {
				return !value.(*Todo).Done
			}
			return true
		default:
			return true
		}
	}
}

func Convert(value interface{}) rest.Model {
	return value.(*Todo)
}

var emptyContext = context.Background()
var trueContext = context.WithValue(emptyContext, "done", "true")
var falseContext = context.WithValue(emptyContext, "done", "false")

func NewTodoDictService() rest.Service {
	return rest.NewDictService(NewTodo, Filter, Convert)
}
