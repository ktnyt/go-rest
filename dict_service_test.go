package rest_test

import (
	"bytes"
	"encoding/json"
	"testing"

	rest "github.com/ktnyt/go-rest"
	"github.com/stretchr/testify/require"
)

func TestDictService(t *testing.T) {
	service := rest.NewDictService(NewTodo, Filter, Convert)
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
			trueValues, err := service.Browse(trueContext)
			require.NoError(t, err)

			for _, value := range trueValues {
				trueKeys = append(trueKeys, value.(*Todo).Key)
			}

			falseValues, err := service.Browse(falseContext)
			require.NoError(t, err)

			for _, value := range falseValues {
				falseKeys = append(falseKeys, value.(*Todo).Key)
			}

			values := append(trueValues, falseValues...)
			require.Len(t, values, count)
		})

		t.Run("without filters", func(t *testing.T) {
			values, err := service.Browse(emptyContext)
			require.NoError(t, err)
			require.Len(t, values, count)

			for i, value := range values {
				keys[i] = value.(*Todo).Key
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

			values, err := service.Browse(emptyContext)
			require.NoError(t, err)
			require.Len(t, values, len(falseKeys))
		})

		t.Run("for existing keys only", func(t *testing.T) {
			_, err := service.Remove("foo")
			require.Error(t, err)

			values, err := service.Browse(emptyContext)
			require.NoError(t, err)
			require.Len(t, values, len(falseKeys))
		})
	})

	t.Run("can delete data", func(t *testing.T) {
		t.Run("without filters", func(t *testing.T) {
			_, err := service.Delete(emptyContext)
			require.NoError(t, err)
			values, err := service.Browse(emptyContext)
			require.NoError(t, err)
			require.Len(t, values, 0)
		})

		t.Run("with filters", func(t *testing.T) {
			for i := 0; i < count; i++ {
				todo := RandomTodo()

				data, err := json.Marshal(&todo)
				require.NoError(t, err)

				_, err = service.Create(bytes.NewReader(data))
				require.NoError(t, err)
			}

			trueValues, err := service.Delete(trueContext)
			require.NoError(t, err)

			falseValues, err := service.Delete(falseContext)
			require.NoError(t, err)

			values := append(trueValues, falseValues...)
			require.Len(t, values, count)
		})
	})
}
