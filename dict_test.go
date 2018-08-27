package rest_test

import (
	"bytes"
	"encoding/json"
	"net/url"
	"testing"

	rest "github.com/ktnyt/go-rest"
	"github.com/stretchr/testify/require"
)

func TestDict(t *testing.T) {
	d := rest.NewDict()
	todo := RandomTodo()
	key := todo.MakeKey(d.Length())

	t.Run("is empty on creation", func(t *testing.T) {
		require.Equal(t, d.Length(), 0)
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

func filter(params url.Values) rest.Filter {
	return func(model rest.Model) bool {
		done := params.Get("done")
		if len(done) > 0 {
			if done == "true" {
				return model.(*Todo).Done
			}
			if done == "false" {
				return !model.(*Todo).Done
			}
		}
		return true
	}
}

var trueParams = url.Values{"done": []string{"true"}}
var falseParams = url.Values{"done": []string{"false"}}
var emptyParams = url.Values{}

func TestDictService(t *testing.T) {
	service := rest.NewDictService(Construct, filter)
	count := 10
	trueKeys := make([]string, 0)
	falseKeys := make([]string, 0)
	keys := make([]string, count)

	t.Run("can create data", func(t *testing.T) {
		for i := 0; i < count; i++ {
			todo := RandomTodo()

			data, err := json.Marshal(&todo)
			require.NoError(t, err)

			_, err = service.Create(bytes.NewReader(data))
			require.NoError(t, err)
		}
	})

	t.Run("can browse data", func(t *testing.T) {
		t.Run("with filters", func(t *testing.T) {
			trueModels, err := service.Browse(trueParams)
			require.NoError(t, err)

			for _, model := range trueModels {
				trueKeys = append(trueKeys, model.(*Todo).Key)
			}

			falseModels, err := service.Browse(falseParams)
			require.NoError(t, err)

			for _, model := range falseModels {
				falseKeys = append(falseKeys, model.(*Todo).Key)
			}

			models := append(trueModels, falseModels...)
			require.Len(t, models, count)
		})

		t.Run("without filters", func(t *testing.T) {
			models, err := service.Browse(emptyParams)
			require.NoError(t, err)
			require.Len(t, models, count)

			for i, model := range models {
				keys[i] = model.(*Todo).Key
			}
		})
	})

	t.Run("can select data", func(t *testing.T) {
		t.Run("for existing keys", func(t *testing.T) {
			for _, key := range keys {
				_, err := service.Select(key)
				require.NoError(t, err)
			}
		})

		t.Run("for existing keys only", func(t *testing.T) {
			_, err := service.Select("foo")
			require.Error(t, err)
		})
	})

	t.Run("can update data", func(t *testing.T) {
		t.Run("for existing keys", func(t *testing.T) {
			for _, key := range keys {
				todo := RandomTodo()

				data, err := json.Marshal(&todo)
				require.NoError(t, err)

				_, err = service.Update(key, bytes.NewReader(data))
				require.NoError(t, err)
			}
		})

		t.Run("for existing keys only", func(t *testing.T) {
			todo := RandomTodo()

			data, err := json.Marshal(&todo)
			require.NoError(t, err)

			_, err = service.Update("foo", bytes.NewReader(data))
			require.Error(t, err)
		})

		t.Run("for valid data only", func(t *testing.T) {
			for _, key := range keys {
				todo := InvalidTodo()

				data, err := json.Marshal(&todo)
				require.NoError(t, err)

				before, err := service.Select(key)
				require.NoError(t, err)

				_, err = service.Update(key, bytes.NewReader(data))
				require.Error(t, err)

				after, err := service.Select(key)
				require.NoError(t, err)
				require.Equal(t, before, after)
			}
		})
	})

	t.Run("can modify data", func(t *testing.T) {
		t.Run("for existing keys", func(t *testing.T) {
			for _, key := range keys {
				todo := RandomTodo()

				data, err := json.Marshal(&todo)
				require.NoError(t, err)

				_, err = service.Modify(key, bytes.NewReader(data))
				require.NoError(t, err)
			}
		})

		t.Run("for existing keys only", func(t *testing.T) {
			todo := RandomTodo()

			data, err := json.Marshal(&todo)
			require.NoError(t, err)

			_, err = service.Modify("foo", bytes.NewReader(data))
			require.Error(t, err)
		})

		t.Run("for valid data only", func(t *testing.T) {
			for _, key := range keys {
				todo := InvalidTodo()

				data, err := json.Marshal(&todo)
				require.NoError(t, err)

				before, err := service.Select(key)
				require.NoError(t, err)

				_, err = service.Modify(key, bytes.NewReader(data))
				require.Error(t, err)

				after, err := service.Select(key)
				require.NoError(t, err)
				require.Equal(t, before, after)
			}
		})
	})

	t.Run("can remove data", func(t *testing.T) {
		t.Run("for existing keys", func(t *testing.T) {
			for _, key := range trueKeys {
				_, err := service.Remove(key)
				require.NoError(t, err)
			}

			models, err := service.Browse(emptyParams)
			require.NoError(t, err)
			require.Len(t, models, len(falseKeys))
		})

		t.Run("for existing keys only", func(t *testing.T) {
			_, err := service.Remove("foo")
			require.Error(t, err)

			models, err := service.Browse(emptyParams)
			require.NoError(t, err)
			require.Len(t, models, len(falseKeys))
		})
	})

	t.Run("can delete data", func(t *testing.T) {
		t.Run("without filters", func(t *testing.T) {
			_, err := service.Delete(emptyParams)
			require.NoError(t, err)
			models, err := service.Browse(emptyParams)
			require.NoError(t, err)
			require.Len(t, models, 0)
		})

		t.Run("with filters", func(t *testing.T) {
			for i := 0; i < count; i++ {
				todo := RandomTodo()

				data, err := json.Marshal(&todo)
				require.NoError(t, err)

				_, err = service.Create(bytes.NewReader(data))
				require.NoError(t, err)
			}

			trueModels, err := service.Delete(trueParams)
			require.NoError(t, err)

			falseModels, err := service.Delete(falseParams)
			require.NoError(t, err)

			models := append(trueModels, falseModels...)
			require.Len(t, models, count)
		})
	})
}
