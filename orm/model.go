package orm

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"time"
)

type Model struct {
	ID        uint `micros:"primary_key"`
	Ver       int  `sql:"version"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type User struct {
	Id                int64
	Age               int64
	UserNum           Num
	Name              string `sql:"size:255"`
	Email             string
	Role              string
	Birthday          *time.Time // Time
	CreatedAt         time.Time  // CreatedAt: Time of record is created, will be insert automatically
	UpdatedAt         time.Time  // UpdatedAt: Time of record is updated, will be updated automatically
	Latitude          float64
	PasswordHash      []byte
	IgnoreMe          int64    `sql:"-"`
	IgnoreStringSlice []string `sql:"-"`
}

type Toy struct {
	Id        int
	Name      string
	OwnerId   int
	OwnerType string
}

type EmailWithIdx struct {
	Id           int64
	UserId       int64
	Email        string     `sql:"index:idx_email_agent"`
	UserAgent    string     `sql:"index:idx_email_agent"`
	RegisteredAt *time.Time `sql:"unique_index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type NotSoLongTableName struct {
	Id                int64
	ReallyLongThingID int64
	ReallyLongThing   ReallyLongTableNameToTestMySQLNameLengthLimit
}

type ReallyLongTableNameToTestMySQLNameLengthLimit struct {
	Id int64
}

type ReallyLongThingThatReferencesShort struct {
	Id      int64
	ShortID int64
	Short   Short
}

type Short struct {
	Id int64
}

type Num int64

func (i *Num) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
	case int64:
		*i = Num(s)
	default:
		return errors.New("Cannot scan NamedInt from " + reflect.ValueOf(src).String())
	}
	return nil
}

// Scanner
type NullValue struct {
	Id      int64
	Name    sql.NullString  `sql:"not null"`
	Gender  *sql.NullString `sql:"not null"`
	Age     sql.NullInt64
	Male    sql.NullBool
	Height  sql.NullFloat64
	AddedAt NullTime
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	nt.Time, nt.Valid = value.(time.Time), true
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func getPreparedUser(name, role string) *User {

	return &User{
		Name: name,
		Role: role,
		Age:  20,
	}
}
