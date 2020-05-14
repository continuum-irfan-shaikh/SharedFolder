package db

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

const zeroTTL = 0

var (
	defaultTime = time.Time{}

	execBatch = func(_ *gocql.Batch) error {
		return nil
	}
)

type (
	// AccessorMock contains data for mocks
	AccessorMock struct {
		tableName   string
		keys        []string
		columns     []string
		columnsName []string
		viewTables  map[string][]string

		CustomExecRelease func(q *gocqlx.Queryx) error

		rows       map[string]row
		viewRows   map[string]map[string]row
		SelectMock QueryMockFunc
		GetMock    QueryMockFunc
		observer   Observer
	}

	row struct {
		shouldExpire bool
		expire      time.Time
		value       interface{}
	}
)

// NewAccessorMock constructor
func NewAccessorMock(tableName string, tableKeys []string, item Model, viewTables map[string][]string) *AccessorMock {
	viewRows := make(map[string]map[string]row, len(viewTables))
	for table := range viewTables {
		viewRows[table] = make(map[string]row)
	}

	mock := &AccessorMock{
		tableName:  tableName,
		keys:       tableKeys,
		viewTables: viewTables,
		rows:       make(map[string]row),
		viewRows:   viewRows,
		SelectMock: DefaultSelectMock,
		GetMock:    DefaultGetMock,
		observer:   &DefaultObserver{},
	}

	mock.columns, mock.columnsName, _ = getColumnsAndColumnNames(item, mock)
	return mock
}

// MockExecuteBatch mock execute batch
func MockExecuteBatch(err error) func() {
	oldExecuteBatch := execBatch
	execBatch = func(_ *gocql.Batch) error {
		return err
	}
	return func() {
		execBatch = oldExecuteBatch
	}
}

// SetValueTo initialize variable with item value
func SetValueTo(variable, item interface{}) {
	val := reflect.ValueOf(variable)
	if val.Kind() != reflect.Ptr {
		panic("some: variable must be a pointer")
	}
	val.Elem().Set(reflect.ValueOf(item).Elem())
}

// SetSliceTo initialize variable (slice) with items value
func SetSliceTo(variable interface{}, items []interface{}) {
	val := reflect.ValueOf(variable)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		panic("some: variable must be a slice")
	}

	for _, item := range items {
		newItem := reflect.ValueOf(item)
		if !val.CanSet() {
			panic("cannot set value! some: variable must be a pointer")
		}
		val.Set(reflect.Append(val, newItem))
	}
}

// Keys mocks Keys
func (m *AccessorMock) Keys() []string {
	return m.keys
}

// Table mocks Table
func (m *AccessorMock) Table() string {
	return m.tableName
}

// ExecuteBatch execute batch
func (m *AccessorMock) ExecuteBatch(batch *gocql.Batch) error {
	return execBatch(batch)
}

// GetColumns mocks GetColumns
func (m *AccessorMock) GetColumns() []string {
	return m.columnsName
}

// Quote mocks Quote
func (m *AccessorMock) Quote(s string) string {
	return fmt.Sprintf("%q", s)
}

// All mocks All
func (m *AccessorMock) All(values interface{}, keyCols ...interface{}) error {
	filter := make(map[string]string)
	for i, v := range keyCols {
		filter[m.keys[i]] = fmt.Sprint(v)
	}
	rows := m.FindBy(filter)
	if len(rows) == 0 {
		return gocql.ErrNotFound
	}
	SetSliceTo(values, rows)
	return nil
}

// Get mocks Get
func (m *AccessorMock) Get(value interface{}, keyCols ...interface{}) error {
	hash := m.getKeyHash(keyCols)
	for k, r := range m.rows {
		if m.checkExpireTTL(m.tableName, k, r) {
			continue
		}
		if strings.HasPrefix(k, hash) {
			SetValueTo(value, r.value)
			return nil
		}
	}
	return gocql.ErrNotFound
}

// AllFromTable mocks AllFromTable
func (m *AccessorMock) AllFromTable(table string, values interface{}, keyCols ...interface{}) error {
	err := m.checkTable(table)
	if err != nil {
		return err
	}
	keys := m.viewTables[table]
	filter := make(map[string]string)
	for i, v := range keyCols {
		filter[keys[i]] = fmt.Sprint(v)
	}
	rows := m.FindByTable(table, filter)
	if len(rows) == 0 {
		return gocql.ErrNotFound
	}
	SetSliceTo(values, rows)
	return nil
}

// GetFromTable mocks GetFromTable
func (m *AccessorMock) GetFromTable(table string, value interface{}, keyCols ...interface{}) error {
	err := m.checkTable(table)
	if err != nil {
		return err
	}
	tableData := m.viewRows[table]
	hash := m.getKeyHash(keyCols)
	for k, v := range tableData {
		if m.checkExpireTTL(table, k, v) {
			continue
		}
		if strings.HasPrefix(k, hash) {
			SetValueTo(value, v.value)
			return nil
		}
	}
	return gocql.ErrNotFound
}

// Add mocks Add
func (m *AccessorMock) Add(item Model) error {
	err := item.AcquireID()
	if err != nil {
		return err
	}
	return m.insert(item, zeroTTL, EventBeforeAdd, EventAfterAdd)
}

// AddWithTTL mocks AddWithTTL
func (m *AccessorMock) AddWithTTL(item Model, ttl time.Duration) error {
	err := item.AcquireID()
	if err != nil {
		return err
	}
	return m.insert(item, ttl, EventBeforeAdd, EventAfterAdd)
}

// Update mocks Update
func (m *AccessorMock) Update(item Model) error {
	return m.insert(item, zeroTTL, EventBeforeUpdate, EventAfterUpdate)
}

// UpdateWithTTL mocks UpdateWithTTL
func (m *AccessorMock) UpdateWithTTL(item Model, ttl time.Duration) error {
	return m.insert(item, ttl, EventBeforeUpdate, EventAfterUpdate)
}

// Delete mocks Delete
func (m *AccessorMock) Delete(item Model) error {
	keyCols, err := GetQueryKeys(item, m.Keys())
	if err != nil {
		return err
	}
	m.observer.OnNotify(EventBeforeDelete, item)
	id := m.getKeyHash(keyCols)
	delete(m.rows, id)
	for table, keys := range m.viewTables {
		tableID := m.getKeyHashWithKeysFrom(keys, item)
		delete(m.viewRows[table], tableID)
	}
	m.observer.OnNotify(EventAfterDelete, item)
	return nil
}

// AddWithBatch mocks AddWithBatch
func (m *AccessorMock) AddWithBatch(batch *gocql.Batch, item Model) error {
	return m.Add(item)
}

// AddWithBatchAndTTL mocks AddWithBatchAndTTL
func (m *AccessorMock) AddWithBatchAndTTL(batch *gocql.Batch, item Model, ttl time.Duration) error {
	return m.AddWithTTL(item, ttl)
}

// UpdateWithBatch mocks UpdateWithBatch
func (m *AccessorMock) UpdateWithBatch(batch *gocql.Batch, item Model) error {
	return m.Update(item)
}

// UpdateWithBatchAndTTL mocks UpdateWithBatchAndTTL
func (m *AccessorMock) UpdateWithBatchAndTTL(batch *gocql.Batch, item Model, ttl time.Duration) error {
	return m.UpdateWithTTL(item, ttl)
}

// DeleteWithBatch mocks DeleteWithBatch
func (m *AccessorMock) DeleteWithBatch(batch *gocql.Batch, item Model) error {
	return m.Delete(item)
}

// ExecRelease mocks ExecRelease
func (m *AccessorMock) ExecRelease(q *gocqlx.Queryx) error {
	if m.CustomExecRelease != nil {
		return m.CustomExecRelease(q)
	}
	return nil
}

// QuerySelect mocks QuerySelect
func (m *AccessorMock) QuerySelect(values interface{}, queryBuilder qb.Builder, params map[string]interface{}) error {
	return m.SelectMock(m, values, queryBuilder, params)
}

// QuerySelectPagination mocks QuerySelectPagination
func (m *AccessorMock) QuerySelectPagination(values interface{}, queryBuilder qb.Builder, params map[string]interface{}, pageSize, page int) error {
	return m.SelectMock(m, values, queryBuilder, params)
}

// QueryGet mocks QueryGet
func (m *AccessorMock) QueryGet(value interface{}, queryBuilder qb.Builder, params map[string]interface{}) error {
	return m.GetMock(m, value, queryBuilder, params)
}

// Register set observer for repo
func (m *AccessorMock) Register(observer Observer) {
	m.observer = observer
}

// Deregister revert observer to default (empty)
func (m *AccessorMock) Deregister() {
	m.observer = &DefaultObserver{}
}

// FindBy returns items with filters
func (m *AccessorMock) FindBy(condition map[string]string) []interface{} {
	var results []interface{}
	for k, v := range m.rows {
		if m.checkExpireTTL(m.tableName, k, v) {
			continue
		}
		if m.accept(v.value, condition) {
			results = append(results, v.value)
		}
	}
	return results
}

// FindByTable returns items with filters
func (m *AccessorMock) FindByTable(table string, condition map[string]string) []interface{} {
	var results []interface{}
	tableData, ok := m.viewRows[table]
	if !ok {
		return results
	}

	for k, v := range tableData {
		if m.checkExpireTTL(table, k, v) {
			continue
		}
		if m.accept(v.value, condition) {
			results = append(results, v.value)
		}
	}
	return results
}

// Dump prints dump
func (m *AccessorMock) Dump() {
	fmt.Println("Dump [" + m.tableName + "]")
	for key, value := range m.rows {
		fmt.Println(key, " --> ", value)
	}
	fmt.Println()
}

func (m *AccessorMock) insert(item Model, ttl time.Duration, beforeEvent, afterEvent EventType) error {
	m.observer.OnNotify(beforeEvent, item)
	id := m.getKeyHashFrom(item)

	rowItem := row{
		shouldExpire: ttl != zeroTTL,
		expire:      time.Now().Add(ttl),
		value:       item,
	}

	m.rows[id] = rowItem
	for table, keys := range m.viewTables {
		tableID := m.getKeyHashWithKeysFrom(keys, item)
		m.viewRows[table][tableID] = rowItem
	}
	m.observer.OnNotify(afterEvent, item)
	return nil
}

func (m *AccessorMock) getKeyHash(keyCols ...interface{}) string {
	s := ""
	for _, v := range keyCols {
		s += fmt.Sprintf("%v", v)
	}
	return strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
}

func (m *AccessorMock) getKeyHashFrom(item interface{}) string {
	return m.getKeyHashWithKeysFrom(m.keys, item)
}

func (m *AccessorMock) getKeyHashWithKeysFrom(tableKeys []string, item interface{}) string {
	s := reflect.ValueOf(item).Elem()
	keys := make([]interface{}, 0, len(tableKeys))
	for _, key := range tableKeys {
		val := s.FieldByName(key).Interface()
		if isDefaultValue(val) {
			break
		}
		keys = append(keys, val)
	}

	return m.getKeyHash(keys)
}

func (m *AccessorMock) checkTable(table string) error {
	_, ok := m.viewTables[table]
	if !ok {
		return errors.Errorf("not configured table %s", table)
	}
	return nil
}

func (m *AccessorMock) accept(item interface{}, condition map[string]string) bool {
	s := reflect.ValueOf(item).Elem()
	for key, value := range condition {
		val := s.FieldByName(key).Interface()
		realVal := fmt.Sprintf("%v", val)
		if realVal != value {
			return false
		}
	}
	return true
}

func isDefaultValue(value interface{}) bool {
	switch val := value.(type) {
	case gocql.UUID:
		return val == ZeroUUID
	case string:
		return val == ""
	case int, int8, int16, int32, int64:
		return val == 0
	case bool:
		return false
	case time.Time:
		return val == defaultTime
	default:
		panic(fmt.Sprintf("unknown value: %#v type: %T", val, val))
	}
}

func (m *AccessorMock) checkExpireTTL(table, item string, row row) bool {
	if !row.shouldExpire || !time.Now().After(row.expire) {
		return false
	}

	if table == m.tableName {
		delete(m.rows, item)
	} else {
		delete(m.viewRows[table], item)
	}
	return true
}
