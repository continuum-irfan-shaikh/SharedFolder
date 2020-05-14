package db

import (
	"errors"
	"reflect"

	"github.com/gocql/gocql"

	ref "gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/db/reflect"
)

var (
	// ZeroUUID is default (empty) UUID
	ZeroUUID = gocql.UUID{}

	// ErrConvertStructToMap failed convert struct to map
	ErrConvertStructToMap = errors.New("could not convert entity struct to map")
)

// Model interface for acquire model's ID if it not present/unique
type Model interface {
	AcquireID() error
}

// GetQueryKeys returns keys
func GetQueryKeys(item interface{}, keyColumns []string) ([]interface{}, error) {
	m, ok := ref.StructToMap(item)
	if !ok {
		return nil, ErrConvertStructToMap
	}
	keys := make([]interface{}, 0, len(keyColumns))
	for _, k := range keyColumns {
		v, found := m[k]
		if !found {
			return nil, errors.New("model has no specified key " + k)
		}

		if reflect.Zero(reflect.TypeOf(v)).Interface() == v {
			break
		}
		keys = append(keys, v)
	}
	return keys, nil
}
