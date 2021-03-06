package orm

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micros/logger"
	"os"
	"reflect"
	"strings"
	"time"
)

// DB contains information for current db connection
type DB struct {
	Value        interface{}
	Error        error
	RowsAffected int64

	// single db
	db      SQLCommon
	logMode int
	logger  logger.Logger
	search  *search
	values  map[string]interface{}

	// global db
	parent    *DB
	callbacks *Callback
	dialect   Dialect
}

type closer interface {
	Close() error
}

func Open(dialect string, args ...interface{}) (db *DB, err error) {
	if len(args) == 0 {
		err = errors.New("invalid database source")
		return nil, err
	}
	var dbSQL SQLCommon

	switch value := args[0].(type) {
	case string:
		var driver = dialect
		var source = value
		dbSQL, err = sql.Open(driver, source)
	case SQLCommon:
		dbSQL = value
	}

	db = &DB{
		db:        dbSQL,
		logger:    logger.L,
		values:    map[string]interface{}{},
		callbacks: DefaultCallback,
		dialect:   newDialect(dialect, dbSQL),
	}
	db.parent = db
	if err != nil {
		return
	}
	// Send a ping to make sure the database connection is alive.
	if d, ok := dbSQL.(*sql.DB); ok {
		if err = d.Ping(); err != nil {
			d.Close()
		}
	}
	return
}

func (s *DB) New() *DB {
	clone := s.clone()
	clone.search = nil
	clone.Value = nil
	return clone
}

// Close close current db connection
func (s *DB) Close() error {
	//类型断言
	if db, ok := s.parent.db.(closer); ok {
		return db.Close()
	}
	return errors.New("can't close current db")
}

// DB get `*sql.DB` from current connection
func (s *DB) DB() *sql.DB {
	//类型断言
	db, _ := s.db.(*sql.DB)
	return db
}

// CommonDB return the underlying `*sql.DB` or `*sql.Tx` instance
func (s *DB) CommonDB() SQLCommon {
	return s.db
}

// SetLogger replace default logger
func (s *DB) SetLogger(log logger.Logger) {
	s.logger = log
}

// LogMode set log mode, `true` for detailed logs, `false` for no log, default, will only print error logs
func (s *DB) LogMode(enable bool) *DB {
	if enable {
		s.logMode = 2
	} else {
		s.logMode = 1
	}
	return s
}

// NewScope create a scope for current operation
func (s *DB) NewScope(value interface{}) *Scope {
	dbClone := s.clone()
	dbClone.Value = value
	return &Scope{db: dbClone, Search: dbClone.search.clone(), Value: value}
}

// QueryExpr returns the query as expr object
func (s *DB) QueryExpr() *expr {
	scope := s.NewScope(s.Value)
	scope.InstanceSet("skip_bindvar", true)
	scope.prepareQuerySQL()

	return Expr(scope.SQL, scope.SQLVars...)
}

// Where return a new relation, filter records with given conditions, accepts `map`, `struct` or `string` as conditions
func (s *DB) Where(query interface{}, args ...interface{}) *DB {
	return s.clone().search.Where(query, args...).db
}

// Or filter records that match before conditions or this one, similar to `Where`
func (s *DB) Or(query interface{}, args ...interface{}) *DB {
	return s.clone().search.Or(query, args...).db
}

// Not filter records that don't match current conditions, similar to `Where`
func (s *DB) Not(query interface{}, args ...interface{}) *DB {
	return s.clone().search.Not(query, args...).db
}

// Limit specify the number of records to be retrieved
func (s *DB) Limit(limit interface{}) *DB {
	return s.clone().search.Limit(limit).db
}

// Offset specify the number of records to skip before starting to return the records
func (s *DB) Offset(offset interface{}) *DB {
	return s.clone().search.Offset(offset).db
}

// Order specify order when retrieve records from database, set reorder to `true` to overwrite defined conditions
//     db.Order("name DESC")
//     db.Order("name DESC", true) // reorder
//     db.Order(micros.Expr("name = ? DESC", "first"))
func (s *DB) Order(value interface{}, reorder ...bool) *DB {
	return s.clone().search.Order(value, reorder...).db
}

// Select specify fields that you want to retrieve from database when querying, by default, will select all fields;
// When creating/updating, specify fields that you want to save to database
func (s *DB) Select(query interface{}, args ...interface{}) *DB {
	return s.clone().search.Select(query, args...).db
}

// Group specify the group method on the find
func (s *DB) Group(query string) *DB {
	return s.clone().search.Group(query).db
}

// Having specify HAVING conditions for GROUP BY
func (s *DB) Having(query interface{}, values ...interface{}) *DB {
	return s.clone().search.Having(query, values...).db
}

// Scopes pass current database connection to arguments `func(*DB) *DB`, which could be used to add conditions dynamically
//     func AmountGreaterThan1000(db *micros.DB) *micros.DB {
//         return db.Where("amount > ?", 1000)
//     }
//
//     func OrderStatus(status []string) func (db *micros.DB) *micros.DB {
//         return func (db *micros.DB) *micros.DB {
//             return db.Scopes(AmountGreaterThan1000).Where("status in (?)", status)
//         }
//     }
//
//     db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
func (s *DB) Scopes(funcs ...func(*DB) *DB) *DB {
	for _, f := range funcs {
		s = f(s)
	}
	return s
}

// Unscoped return all record including deleted record
func (s *DB) Unscoped() *DB {
	return s.clone().search.unscoped().db
}

// Attrs initialize struct with argument if record not found with `FirstOrInit`
func (s *DB) Attrs(attrs ...interface{}) *DB {
	return s.clone().search.Attrs(attrs...).db
}

// Assign assign result with argument regardless it is found or not with `FirstOrInit`
func (s *DB) Assign(attrs ...interface{}) *DB {
	return s.clone().search.Assign(attrs...).db
}

// First find first record that match given conditions, order by primary key
func (s *DB) First(out interface{}, where ...interface{}) *DB {
	newScope := s.clone().NewScope(out)
	newScope.Search.Limit(1)
	return newScope.Set("micros:order_by_primary_key", "ASC").
		inlineCondition(where...).callCallbacks(s.parent.callbacks.queries).db
}

// Last find last record that match given conditions, order by primary key
func (s *DB) Last(out interface{}, where ...interface{}) *DB {
	newScope := s.clone().NewScope(out)
	newScope.Search.Limit(1)
	return newScope.Set("micros:order_by_primary_key", "DESC").
		inlineCondition(where...).callCallbacks(s.parent.callbacks.queries).db
}

// Find find records that match given conditions
func (s *DB) Find(out interface{}, where ...interface{}) *DB {
	return s.clone().NewScope(out).inlineCondition(where...).callCallbacks(s.parent.callbacks.queries).db
}

// Scan scan value to a struct
func (s *DB) Scan(dest interface{}) *DB {
	return s.clone().NewScope(s.Value).Set("micros:query_destination", dest).callCallbacks(s.parent.callbacks.queries).db
}

// Row return `*sql.Row` with given conditions
func (s *DB) Row() *sql.Row {
	return s.NewScope(s.Value).row()
}

// Rows return `*sql.Rows` with given conditions
func (s *DB) Rows() (*sql.Rows, error) {
	return s.NewScope(s.Value).rows()
}

// ScanRows scan `*sql.Rows` to give struct
func (s *DB) ScanRows(rows *sql.Rows, result interface{}) error {
	var (
		clone        = s.clone()
		scope        = clone.NewScope(result)
		columns, err = rows.Columns()
	)

	if clone.AddError(err) == nil {
		scope.scan(rows, columns, scope.Fields())
	}

	return clone.Error
}

// Pluck used to query single column from a model as a map
//     var ages []int64
//     db.Find(&users).Pluck("age", &ages)
func (s *DB) Pluck(column string, value interface{}) *DB {
	return s.NewScope(s.Value).pluck(column, value).db
}

// Count get how many records for a model
func (s *DB) Count(value interface{}) *DB {
	return s.NewScope(s.Value).count(value).db
}

// FirstOrInit find first matched record or initialize a new one with given conditions (only works with struct, map conditions)
func (s *DB) FirstOrInit(out interface{}, where ...interface{}) *DB {
	c := s.clone()
	if result := c.First(out, where...); result.Error != nil {
		if !result.RecordNotFound() {
			return result
		}
		c.NewScope(out).inlineCondition(where...).initialize()
	} else {
		c.NewScope(out).updatedAttrsWithValues(c.search.assignAttrs)
	}
	return c
}

// FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
func (s *DB) FirstOrCreate(out interface{}, where ...interface{}) *DB {
	c := s.clone()
	if result := s.First(out, where...); result.Error != nil {
		if !result.RecordNotFound() {
			return result
		}
		return c.NewScope(out).inlineCondition(where...).initialize().callCallbacks(c.parent.callbacks.creates).db
	} else if len(c.search.assignAttrs) > 0 {
		return c.NewScope(out).InstanceSet("micros:update_interface", c.search.assignAttrs).callCallbacks(c.parent.callbacks.updates).db
	}
	return c
}

// Update update attributes
func (s *DB) Update(attrs ...interface{}) *DB {
	return s.Updates(toSearchableMap(attrs...), true)
}

// Updates update attributes
func (s *DB) Updates(values interface{}, ignoreProtectedAttrs ...bool) *DB {
	return s.clone().NewScope(s.Value).
		Set("micros:ignore_protected_attrs", len(ignoreProtectedAttrs) > 0).
		InstanceSet("micros:update_interface", values).
		callCallbacks(s.parent.callbacks.updates).db
}

// UpdateColumn update attributes
func (s *DB) UpdateColumn(attrs ...interface{}) *DB {
	return s.UpdateColumns(toSearchableMap(attrs...))
}

// UpdateColumns update attributes
func (s *DB) UpdateColumns(values interface{}) *DB {
	return s.clone().NewScope(s.Value).
		Set("micros:update_column", true).
		Set("micros:save_associations", false).
		InstanceSet("micros:update_interface", values).db
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (s *DB) Save(value interface{}) *DB {
	scope := s.clone().NewScope(value)
	if !scope.PrimaryKeyZero() {
		newDB := scope.callCallbacks(s.parent.callbacks.updates).db
		if newDB.Error == nil && newDB.RowsAffected == 0 {
			return s.New().FirstOrCreate(value)
		}
		return newDB
	}
	return scope.callCallbacks(s.parent.callbacks.creates).db
}

// Create insert the value into database
func (s *DB) Create(value interface{}) *DB {
	scope := s.clone().NewScope(value)
	return scope.callCallbacks(s.parent.callbacks.creates).db
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (s *DB) Delete(value interface{}, where ...interface{}) *DB {
	return s.clone().NewScope(value).inlineCondition(where...).callCallbacks(s.parent.callbacks.deletes).db
}

// Raw use raw sql as conditions, won't run it unless invoked by other methods
//    db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
func (s *DB) Raw(sql string, values ...interface{}) *DB {
	return s.clone().search.Raw(true).Where(sql, values...).db
}

// Exec execute raw sql
func (s *DB) Exec(sql string, values ...interface{}) *DB {
	scope := s.clone().NewScope(nil)
	generatedSQL := scope.buildWhereCondition(map[string]interface{}{"query": sql, "args": values})
	generatedSQL = strings.TrimSuffix(strings.TrimPrefix(generatedSQL, "("), ")")
	scope.Raw(generatedSQL)
	return scope.Exec().db
}

// Model specify the model you would like to run db operations
//    // update all users's name to `hello`
//    db.Model(&User{}).Update("name", "hello")
//    // if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
//    db.Model(&user).Update("name", "hello")
func (s *DB) Model(value interface{}) *DB {
	c := s.clone()
	c.Value = value
	return c
}

// Table specify the table you would like to run db operations
func (s *DB) Table(name string) *DB {
	clone := s.clone()
	clone.search.Table(name)
	clone.Value = nil
	return clone
}

// Debug start debug mode
func (s *DB) Debug() *DB {
	return s.clone().LogMode(true)
}

// Begin begin a transaction
func (s *DB) Begin() *DB {
	c := s.clone()
	if db, ok := c.db.(sqlDb); ok && db != nil {
		tx, err := db.Begin()
		c.db = interface{}(tx).(SQLCommon)
		c.AddError(err)
	} else {
		c.AddError(ErrCantStartTransaction)
	}
	return c
}

// Commit commit a transaction
func (s *DB) Commit() *DB {
	if db, ok := s.db.(sqlTx); ok && db != nil {
		s.AddError(db.Commit())
	} else {
		s.AddError(ErrInvalidTransaction)
	}
	return s
}

// Rollback rollback a transaction
func (s *DB) Rollback() *DB {
	if db, ok := s.db.(sqlTx); ok && db != nil {
		s.AddError(db.Rollback())
	} else {
		s.AddError(ErrInvalidTransaction)
	}
	return s
}

// NewRecord check if value's primary key is blank
func (s *DB) NewRecord(value interface{}) bool {
	return s.clone().NewScope(value).PrimaryKeyZero()
}

// RecordNotFound check if returning ErrRecordNotFound error
func (s *DB) RecordNotFound() bool {
	for _, err := range s.GetErrors() {
		if err == ErrRecordNotFound {
			return true
		}
	}
	return false
}

// CreateTable create table for models
func (s *DB) CreateTable(models ...interface{}) *DB {
	db := s.Unscoped()
	for _, model := range models {
		db = db.NewScope(model).createTable().db
	}
	return db
}

// DropTable drop table for models
func (s *DB) DropTable(values ...interface{}) *DB {
	db := s.clone()
	for _, value := range values {
		if tableName, ok := value.(string); ok {
			db = db.Table(tableName)
		}

		db = db.NewScope(value).dropTable().db
	}
	return db
}

// HasTable check has table or not
func (s *DB) HasTable(value interface{}) bool {
	var (
		scope     = s.clone().NewScope(value)
		tableName string
	)

	if name, ok := value.(string); ok {
		tableName = name
	} else {
		tableName = scope.TableName()
	}

	has := scope.Dialect().HasTable(tableName)
	s.AddError(scope.db.Error)
	return has
}

// AddIndex add index for columns with given name
func (s *DB) AddIndex(indexName string, columns ...string) *DB {
	scope := s.Unscoped().NewScope(s.Value)
	scope.addIndex(false, indexName, columns...)
	return scope.db
}

// AddUniqueIndex add unique index for columns with given name
func (s *DB) AddUniqueIndex(indexName string, columns ...string) *DB {
	scope := s.Unscoped().NewScope(s.Value)
	scope.addIndex(true, indexName, columns...)
	return scope.db
}

// Set set setting by name, will clone a new db, and update its setting
func (s *DB) Set(name string, value interface{}) *DB {
	return s.clone().InstantSet(name, value)
}

// InstantSet instant set setting, will affect current db
func (s *DB) InstantSet(name string, value interface{}) *DB {
	s.values[name] = value
	return s
}

// Get get setting by name
func (s *DB) Get(name string) (value interface{}, ok bool) {
	value, ok = s.values[name]
	return
}

// AddError add error to the db
func (s *DB) AddError(err error) error {
	if err != nil {
		if err != ErrRecordNotFound {
			if s.logMode == 0 {
				go s.print(fileWithLineNum(), err)
			} else {
				s.log(err)
			}

			errors := Errors(s.GetErrors())
			errors = errors.Add(err)
			if len(errors) > 1 {
				err = errors
			}
		}

		s.Error = err
	}
	return err
}

// GetErrors get happened errors from the db
func (s *DB) GetErrors() []error {
	if errs, ok := s.Error.(Errors); ok {
		return errs
	} else if s.Error != nil {
		return []error{s.Error}
	}
	return []error{}
}

////////////////////////////////////////////////////////////////////////////////
// Private Methods For DB
////////////////////////////////////////////////////////////////////////////////

func (s *DB) clone() *DB {
	db := &DB{
		db:      s.db,
		parent:  s.parent,
		logger:  s.logger,
		logMode: s.logMode,
		values:  map[string]interface{}{},
		Value:   s.Value,
		Error:   s.Error,
	}

	for key, value := range s.values {
		db.values[key] = value
	}

	if s.search == nil {
		db.search = &search{limit: -1, offset: -1}
	} else {
		db.search = s.search.clone()
	}

	db.search.db = db
	return db
}

func (s *DB) print(v ...interface{}) {
	s.logger.Info(v...)
}

func (s *DB) log(v ...interface{}) {
	if s != nil && s.logMode == 2 {
		s.print(append([]interface{}{"log", fileWithLineNum()}, v...)...)
	}
}

func (s *DB) slog(sql string, t time.Time, vars ...interface{}) {
	for _, seatVar := range vars {
		v := reflect.ValueOf(seatVar)
		str := fmt.Sprintf("'%v'", v)
		sql = strings.Replace(sql, "?", str, 1)
	}
	s.logger.Sql(fileWithLineNum()+" ", sql+" [", NowFunc().Sub(t), "] [", s.RowsAffected, "]")
}

//todo add config
func OpenConnection() (db *DB) {

	//	dbhost := os.Getenv("micros_DBADDRESS")
	//todo  修改读配置
	dbhost := "127.0.0.1:3306"
	dbhost = fmt.Sprintf("tcp(%v)", dbhost)
	db, err := Open("mysql", fmt.Sprintf("root:@%v/micros?charset=utf8&parseTime=True", dbhost))
	if err != nil {
		panic("can't not open connection," + err.Error())
	}

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}
	db.DB().SetMaxIdleConns(10)
	return
}
