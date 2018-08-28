package rest_test

import (
	"testing"

	rest "github.com/ktnyt/go-rest"
	"github.com/stretchr/testify/require"
)

func TestDict(t *testing.T) {
	d := rest.NewDict()
	todo := RandomTodo()
	key := todo.MakeKey(d.Len())

	t.Run("is empty on creation", func(t *testing.T) {
		require.Equal(t, d.Len(), 0)
	})

	t.Run("can insert values", func(t *testing.T) {
		require.True(t, d.Insert(key, todo))
	})

	t.Run("can get a value", func(t *testing.T) {
		t.Run("for existing key", func(t *testing.T) {
			require.Equal(t, todo, d.Get(key))
		})
		t.Run("only for existing key", func(t *testing.T) {
			require.Nil(t, d.Get("foo"))
		})
	})

	t.Run("can set a value", func(t *testing.T) {
		t.Run("for existing key", func(t *testing.T) {
			require.True(t, d.Set(key, todo))
		})
		t.Run("only for existing key", func(t *testing.T) {
			require.False(t, d.Set("foo", todo))
		})
	})

	t.Run("can search values", func(t *testing.T) {
		require.NotZero(t, len(d.Search(func(model rest.Model) bool {
			return !model.(*Todo).Done
		})))
		require.Zero(t, len(d.Search(func(model rest.Model) bool {
			return model.(*Todo).Done
		})))
	})

	t.Run("can remove a value", func(t *testing.T) {
		t.Run("for existing key", func(t *testing.T) {
			require.Equal(t, todo, d.Remove(key))
		})
		t.Run("only for existing key", func(t *testing.T) {
			require.Nil(t, d.Remove(key))
		})
	})
}
