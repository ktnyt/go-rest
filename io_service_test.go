package rest_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"testing"

	rest "github.com/ktnyt/go-rest"
	"github.com/stretchr/testify/require"
)

type BufferIOHandler struct {
	buffer []byte
	build  rest.ServiceBuilder
}

func NewBufferIOHandler(build rest.ServiceBuilder) rest.IOHandler {
	handler := &BufferIOHandler{build: build}
	service := handler.build()

	if err := handler.Save(service); err != nil {
		log.Fatal(err)
	}

	return handler
}

func (h *BufferIOHandler) Save(service rest.Service) error {
	writer := new(bytes.Buffer)
	if err := gob.NewEncoder(writer).Encode(service); err != nil {
		return err
	}
	h.buffer = writer.Bytes()
	return nil
}

func (h *BufferIOHandler) Load() (rest.Service, error) {
	service := h.build()
	reader := bytes.NewBuffer(h.buffer)
	if err := gob.NewDecoder(reader).Decode(service); err != nil {
		return nil, err
	}
	return service, nil
}

func TestIOService(t *testing.T) {
	handler := NewBufferIOHandler(NewTodoDictService)
	service := rest.NewIOService(handler)
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
