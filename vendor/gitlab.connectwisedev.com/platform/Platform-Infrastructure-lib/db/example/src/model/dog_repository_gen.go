// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package model

import (
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"
	"gitlab.connectwisedev.com/platform/Platform-Infrastructure-lib/db"
)

const dogBaseTableName = "dogs"

var (
	dogKeyColumns = []string{NameColumn}
	dogViewTables = map[string][]string{}

	dogsRepository *DogRepository
	initDogs       sync.Once
)

// DogRepository interface
type DogRepository struct {
	base db.Base
}

// Dogs singleton, thread-safe, returns pointer to Dogs repository
func Dogs() *DogRepository {
	initDogs.Do(func() {
		dogsRepository = &DogRepository{base: db.NewBase(
			&Dog{},
			dogBaseTableName,
			dogKeyColumns,
			dogViewTables,
		)}
	})
	return dogsRepository
}

// Get return item
func (r *DogRepository) Get(keyCols ...interface{}) (*Dog, error) {
	dog := new(Dog)
	if err := r.base.Get(dog, keyCols...); err != nil {
		return nil, err
	}
	return dog, nil
}

// All returns slice of Dogs
func (r *DogRepository) All() ([]*Dog, error) {
	var dogs []*Dog
	if err := r.base.All(&dogs); err != nil {
		return nil, err
	}
	return dogs, nil
}

// Add check/generates ID and inserts item
// nolint:interfacer
func (r *DogRepository) Add(item *Dog) error {
	return r.base.Add(item)
}

// AddWithTTL check/generates ID and inserts item with ttl
// nolint:interfacer
func (r *DogRepository) AddWithTTL(item *Dog, ttl time.Duration) error {
	return r.base.AddWithTTL(item, ttl)
}

// Update item in repository
// nolint:interfacer
func (r *DogRepository) Update(item *Dog) error {
	return r.base.Update(item)
}

// UpdateWithTTL item in repository with ttl
// nolint:interfacer
func (r *DogRepository) UpdateWithTTL(item *Dog, ttl time.Duration) error {
	return r.base.UpdateWithTTL(item, ttl)
}

// Delete item from repository
// nolint:interfacer
func (r *DogRepository) Delete(item *Dog) error {
	return r.base.Delete(item)
}

// AddWithBatch adds all queries to batch
// nolint:interfacer
func (r *DogRepository) AddWithBatch(batch *gocql.Batch, item *Dog) error {
	return r.base.AddWithBatch(batch, item)
}

// AddWithBatchAndTTL adds all queries with ttl to batch
// nolint:interfacer
func (r *DogRepository) AddWithBatchAndTTL(batch *gocql.Batch, item *Dog, ttl time.Duration) error {
	return r.base.AddWithBatchAndTTL(batch, item, ttl)
}

// UpdateWithBatch adds all queries to batch
// nolint:interfacer
func (r *DogRepository) UpdateWithBatch(batch *gocql.Batch, item *Dog) error {
	return r.base.UpdateWithBatch(batch, item)
}

// UpdateWithBatchAndTTL adds all queries with ttl to batch
// nolint:interfacer
func (r *DogRepository) UpdateWithBatchAndTTL(batch *gocql.Batch, item *Dog, ttl time.Duration) error {
	return r.base.UpdateWithBatchAndTTL(batch, item, ttl)
}

// DeleteWithBatch adds all queries to batch
// nolint:interfacer
func (r *DogRepository) DeleteWithBatch(batch *gocql.Batch, item *Dog) error {
	return r.base.DeleteWithBatch(batch, item)
}

// GetByID returns item by @id or error if item doesn't exists in repository
func (r *DogRepository) GetByID(id string) (*Dog, error) {
	return r.Get(id)
}

// GetByIDs returns slice of dogs by slice of ids
func (r *DogRepository) GetByIDs(ids ...string) ([]*Dog, error) {
	var dogs []*Dog
	queryBuilder := qb.Select(r.base.Table()).Where(
		qb.In(r.base.Quote(r.base.Keys()[0])),
	)

	err := r.base.QuerySelect(&dogs, queryBuilder, map[string]interface{}{
		r.base.Quote(r.base.Keys()[0]): ids,
	})
	return dogs, err
}

// GetRows returns slice of Dogs filtered by received keys
func (r *DogRepository) GetRows(keyCols ...interface{}) ([]*Dog, error) {
	var dogs []*Dog
	if err := r.base.All(&dogs, keyCols...); err != nil {
		return nil, err
	}
	return dogs, nil
}
