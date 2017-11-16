package orm

import "time"

type Model struct {
	ID        uint `gorm:"primary_key"`
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
	Birthday          *time.Time // Time
	CreatedAt         time.Time  // CreatedAt: Time of record is created, will be insert automatically
	UpdatedAt         time.Time  // UpdatedAt: Time of record is updated, will be updated automatically
	Latitude          float64
	PasswordHash      []byte
	IgnoreMe          int64    `sql:"-"`
	IgnoreStringSlice []string `sql:"-"`
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

func GetPreparedUser(name string) *User {

	return &User{
		Name: name,
		Age:  20,
	}
}

func runMigration() {
	if err := DB.DropTableIfExists(&User{}).Error; err != nil {
		fmt.Printf("Got error when try to delete table users, %+v\n", err)
	}

	for _, table := range []string{"animals", "user_languages"} {
		DB.Exec(fmt.Sprintf("drop table %v;", table))
	}

	values := []interface{}{&Short{}, &ReallyLongThingThatReferencesShort{}, &ReallyLongTableNameToTestMySQLNameLengthLimit{}, &NotSoLongTableName{}, &User{}}
	for _, value := range values {
		DB.DropTable(value)
	}
	if err := DB.AutoMigrate(values...).Error; err != nil {
		panic(fmt.Sprintf("No error should happen when create table, but got %+v", err))
	}
}

func TestIndexes(t *testing.T) {
	if err := DB.Model(&Email{}).AddIndex("idx_email_email", "email").Error; err != nil {
		t.Errorf("Got error when tried to create index: %+v", err)
	}

	scope := DB.NewScope(&Email{})
	if !scope.Dialect().HasIndex(scope.TableName(), "idx_email_email") {
		t.Errorf("Email should have index idx_email_email")
	}

	if err := DB.Model(&Email{}).RemoveIndex("idx_email_email").Error; err != nil {
		t.Errorf("Got error when tried to remove index: %+v", err)
	}

	if scope.Dialect().HasIndex(scope.TableName(), "idx_email_email") {
		t.Errorf("Email's index idx_email_email should be deleted")
	}

	if err := DB.Model(&Email{}).AddIndex("idx_email_email_and_user_id", "user_id", "email").Error; err != nil {
		t.Errorf("Got error when tried to create index: %+v", err)
	}

	if !scope.Dialect().HasIndex(scope.TableName(), "idx_email_email_and_user_id") {
		t.Errorf("Email should have index idx_email_email_and_user_id")
	}

	if err := DB.Model(&Email{}).RemoveIndex("idx_email_email_and_user_id").Error; err != nil {
		t.Errorf("Got error when tried to remove index: %+v", err)
	}

	if scope.Dialect().HasIndex(scope.TableName(), "idx_email_email_and_user_id") {
		t.Errorf("Email's index idx_email_email_and_user_id should be deleted")
	}

	if err := DB.Model(&Email{}).AddUniqueIndex("idx_email_email_and_user_id", "user_id", "email").Error; err != nil {
		t.Errorf("Got error when tried to create index: %+v", err)
	}

	if !scope.Dialect().HasIndex(scope.TableName(), "idx_email_email_and_user_id") {
		t.Errorf("Email should have index idx_email_email_and_user_id")
	}

	if DB.Save(&User{Name: "unique_indexes", Emails: []Email{{Email: "user1@example.comiii"}, {Email: "user1@example.com"}, {Email: "user1@example.com"}}}).Error == nil {
		t.Errorf("Should get to create duplicate record when having unique index")
	}

	var user = User{Name: "sample_user"}
	DB.Save(&user)
	if DB.Model(&user).Association("Emails").Append(Email{Email: "not-1duplicated@gmail.com"}, Email{Email: "not-duplicated2@gmail.com"}).Error != nil {
		t.Errorf("Should get no error when append two emails for user")
	}

	if DB.Model(&user).Association("Emails").Append(Email{Email: "duplicated@gmail.com"}, Email{Email: "duplicated@gmail.com"}).Error == nil {
		t.Errorf("Should get no duplicated email error when insert duplicated emails for a user")
	}

	if err := DB.Model(&Email{}).RemoveIndex("idx_email_email_and_user_id").Error; err != nil {
		t.Errorf("Got error when tried to remove index: %+v", err)
	}

	if scope.Dialect().HasIndex(scope.TableName(), "idx_email_email_and_user_id") {
		t.Errorf("Email's index idx_email_email_and_user_id should be deleted")
	}

	if DB.Save(&User{Name: "unique_indexes", Emails: []Email{{Email: "user1@example.com"}, {Email: "user1@example.com"}}}).Error != nil {
		t.Errorf("Should be able to create duplicated emails after remove unique index")
	}
}

func TestAutoMigration(t *testing.T) {
	if err := DB.AutoMigrate(&EmailWithIdx{}).Error; err != nil {
		t.Errorf("Auto Migrate should not raise any error")
	}

	now := time.Now()
	DB.Save(&EmailWithIdx{Email: "jinzhu@example.org", UserAgent: "pc", RegisteredAt: &now})

	scope := DB.NewScope(&EmailWithIdx{})
	if !scope.Dialect().HasIndex(scope.TableName(), "idx_email_agent") {
		t.Errorf("Failed to create index")
	}

	if !scope.Dialect().HasIndex(scope.TableName(), "uix_email_with_idxes_registered_at") {
		t.Errorf("Failed to create index")
	}

	var bigemail EmailWithIdx
	DB.First(&bigemail, "user_agent = ?", "pc")
	if bigemail.Email != "jinzhu@example.org" || bigemail.UserAgent != "pc" || bigemail.RegisteredAt.IsZero() {
		t.Error("Big Emails should be saved and fetched correctly")
	}
}

type MultipleIndexes struct {
	ID     int64
	UserID int64  `sql:"unique_index:uix_multipleindexes_user_name,uix_multipleindexes_user_email;index:idx_multipleindexes_user_other"`
	Name   string `sql:"unique_index:uix_multipleindexes_user_name"`
	Email  string `sql:"unique_index:,uix_multipleindexes_user_email"`
	Other  string `sql:"index:,idx_multipleindexes_user_other"`
}

func TestMultipleIndexes(t *testing.T) {
	if err := DB.DropTableIfExists(&MultipleIndexes{}).Error; err != nil {
		fmt.Printf("Got error when try to delete table multiple_indexes, %+v\n", err)
	}

	DB.AutoMigrate(&MultipleIndexes{})
	if err := DB.AutoMigrate(&EmailWithIdx{}).Error; err != nil {
		t.Errorf("Auto Migrate should not raise any error")
	}

	DB.Save(&MultipleIndexes{UserID: 1, Name: "jinzhu", Email: "jinzhu@example.org", Other: "foo"})

	scope := DB.NewScope(&MultipleIndexes{})
	if !scope.Dialect().HasIndex(scope.TableName(), "uix_multipleindexes_user_name") {
		t.Errorf("Failed to create index")
	}

	if !scope.Dialect().HasIndex(scope.TableName(), "uix_multipleindexes_user_email") {
		t.Errorf("Failed to create index")
	}

	if !scope.Dialect().HasIndex(scope.TableName(), "uix_multiple_indexes_email") {
		t.Errorf("Failed to create index")
	}

	if !scope.Dialect().HasIndex(scope.TableName(), "idx_multipleindexes_user_other") {
		t.Errorf("Failed to create index")
	}

	if !scope.Dialect().HasIndex(scope.TableName(), "idx_multiple_indexes_other") {
		t.Errorf("Failed to create index")
	}

	var mutipleIndexes MultipleIndexes
	DB.First(&mutipleIndexes, "name = ?", "jinzhu")
	if mutipleIndexes.Email != "jinzhu@example.org" || mutipleIndexes.Name != "jinzhu" {
		t.Error("MutipleIndexes should be saved and fetched correctly")
	}

	// Check unique constraints
	if err := DB.Save(&MultipleIndexes{UserID: 1, Name: "name1", Email: "jinzhu@example.org", Other: "foo"}).Error; err == nil {
		t.Error("MultipleIndexes unique index failed")
	}

	if err := DB.Save(&MultipleIndexes{UserID: 1, Name: "name1", Email: "foo@example.org", Other: "foo"}).Error; err != nil {
		t.Error("MultipleIndexes unique index failed")
	}

	if err := DB.Save(&MultipleIndexes{UserID: 2, Name: "name1", Email: "jinzhu@example.org", Other: "foo"}).Error; err == nil {
		t.Error("MultipleIndexes unique index failed")
	}

	if err := DB.Save(&MultipleIndexes{UserID: 2, Name: "name1", Email: "foo2@example.org", Other: "foo"}).Error; err != nil {
		t.Error("MultipleIndexes unique index failed")
	}
}
