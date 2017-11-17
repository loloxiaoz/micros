package orm

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/now"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var (
	db                 *DB
	t1, t2, t3, t4, t5 time.Time
)

func init() {
	var err error

	if db, err = OpenTestConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to test database, but got err=%+v", err))
	}
}

func OpenTestConnection() (db *DB, err error) {

	//	dbhost := os.Getenv("micros_DBADDRESS")
	//todo  修改读配置
	dbhost := "127.0.0.1:3306"
	dbhost = fmt.Sprintf("tcp(%v)", dbhost)
	db, err = Open("mysql", fmt.Sprintf("root:@%v/micros?charset=utf8&parseTime=True", dbhost))
	if err != nil {
		panic("can't not open connection," + err.Error())
	}

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(10)
	return
}

func TestStringPrimaryKey(t *testing.T) {
	type UUIDStruct struct {
		ID   string `micros:"primary_key"`
		Name string
	}

	db.DropTable(&UUIDStruct{})
	err := db.CreateTable(&UUIDStruct{}).Error
	if err != nil {
		t.Errorf("craete table UUIDStruct error")
	}

	data := UUIDStruct{ID: "uuid", Name: "hello"}
	if err := db.Save(&data).Error; err != nil || data.ID != "uuid" || data.Name != "hello" {
		t.Errorf("string primary key should not be populated")
	}

	data = UUIDStruct{ID: "uuid", Name: "hello world"}
	if err := db.Save(&data).Error; err != nil || data.ID != "uuid" || data.Name != "hello world" {
		t.Errorf("string primary key should not be populated")
	}
}

func TestCreateUser(t *testing.T) {
	db.DropTable(&User{})
	db.CreateTable(&User{})
	user := User{Name: "user"}
	db.Save(&user)

	var count1, count2 int64
	db.Model(&User{}).Count(&count1)
	if count1 <= 0 {
		t.Errorf("Should find some users")
	}

	if db.Where("name = ?", "jinzhu; delete * from user").First(&User{}).Error == nil {
		t.Errorf("Should got error with invalid SQL")
	}

	db.Model(&User{}).Count(&count2)
	if count1 != count2 {
		t.Errorf("No user should not be deleted by invalid SQL")
	}
}

func TestExceptionsWithInvalidSql(t *testing.T) {
	var columns []string
	if db.Model(&User{}).Where("sdsd.zaaa = ?", "sd;;;aa").Pluck("aaa", &columns).Error == nil {
		t.Errorf("Should got error with invalid SQL")
	}

	if db.Where("sdsd.zaaa = ?", "sd;;;aa").Pluck("aaa", &columns).Error == nil {
		t.Errorf("Should got error with invalid SQL")
	}

	if db.Where("sdsd.zaaa = ?", "sd;;;aa").Find(&User{}).Error == nil {
		t.Errorf("Should got error with invalid SQL")
	}

}

func TestSetTable(t *testing.T) {
	db.Create(getPreparedUser("pluck_user1", "pluck_user"))
	db.Create(getPreparedUser("pluck_user2", "pluck_user"))
	db.Create(getPreparedUser("pluck_user3", "pluck_user"))

	if err := db.Table("user").Where("role = ?", "pluck_user").Pluck("age", &[]int{}).Error; err != nil {
		t.Error("No errors should happen if set table for pluck", err)
	}

	var users []User
	if db.Table("user").Find(&[]User{}).Error != nil {
		t.Errorf("No errors should happen if set table for find")
	}

	if db.Table("invalid_table").Find(&users).Error == nil {
		t.Errorf("Should got error when table is set to an invalid table")
	}

	db.Exec("drop table deleted_users;")
	if db.Table("deleted_users").CreateTable(&User{}).Error != nil {
		t.Errorf("Create table with specified table")
	}

	db.Table("deleted_users").Save(&User{Name: "DeletedUser"})

	var deletedUsers []User
	db.Table("deleted_users").Find(&deletedUsers)
	if len(deletedUsers) != 1 {
		t.Errorf("Query from specified table")
	}

	db.Save(getPreparedUser("normal_user", "reset_table"))
	db.Table("deleted_users").Save(getPreparedUser("deleted_user", "reset_table"))
	var user1, user2, user3 User
	db.Where("role = ?", "reset_table").First(&user1).Table("deleted_users").First(&user2).Table("").First(&user3)
	if (user1.Name != "normal_user") || (user2.Name != "deleted_user") || (user3.Name != "normal_user") {
		t.Errorf("unset specified table with blank string")
	}
}

type Order struct {
}

type Cart struct {
}

func (c Cart) TableName() string {
	return "shopping_cart"
}

func TestHasTable(t *testing.T) {
	type Foo struct {
		Id    int
		Stuff string
	}
	db.DropTable(&Foo{})

	// Table should not exist at this point, HasTable should return false
	if ok := db.HasTable("foo"); ok {
		t.Errorf("Table should not exist, but does")
	}
	if ok := db.HasTable(&Foo{}); ok {
		t.Errorf("Table should not exist, but does")
	}

	// We create the table
	if err := db.CreateTable(&Foo{}).Error; err != nil {
		t.Errorf("Table should be created")
	}

	// And now it should exits, and HasTable should return true
	if ok := db.HasTable("foo"); !ok {
		t.Errorf("Table should exist, but HasTable informs it does not")
	}
	if ok := db.HasTable(&Foo{}); !ok {
		t.Errorf("Table should exist, but HasTable informs it does not")
	}
}

func TestTableName(t *testing.T) {
	db := db.Model("")
	if db.NewScope(Order{}).TableName() != "order" {
		t.Errorf("Order's table name should be orders")
	}

	if db.NewScope(&Order{}).TableName() != "order" {
		t.Errorf("&Order's table name should be orders")
	}

	if db.NewScope([]Order{}).TableName() != "order" {
		t.Errorf("[]Order's table name should be orders")
	}

	if db.NewScope(&[]Order{}).TableName() != "order" {
		t.Errorf("&[]Order's table name should be orders")
	}

	if db.NewScope(Order{}).TableName() != "order" {
		t.Errorf("Order's singular table name should be order")
	}

	if db.NewScope(&Order{}).TableName() != "order" {
		t.Errorf("&Order's singular table name should be order")
	}

	if db.NewScope([]Order{}).TableName() != "order" {
		t.Errorf("[]Order's singular table name should be order")
	}

	if db.NewScope(&[]Order{}).TableName() != "order" {
		t.Errorf("&[]Order's singular table name should be order")
	}

	if db.NewScope(&Cart{}).TableName() != "shopping_cart" {
		t.Errorf("&Cart's singular table name should be shopping_cart")
	}

	if db.NewScope(Cart{}).TableName() != "shopping_cart" {
		t.Errorf("Cart's singular table name should be shopping_cart")
	}

	if db.NewScope(&[]Cart{}).TableName() != "shopping_cart" {
		t.Errorf("&[]Cart's singular table name should be shopping_cart")
	}

	if db.NewScope([]Cart{}).TableName() != "shopping_cart" {
		t.Errorf("[]Cart's singular table name should be shopping_cart")
	}
}

func TestNullValues(t *testing.T) {
	db.DropTable(&NullValue{})
	db.CreateTable(&NullValue{})

	if err := db.Save(&NullValue{
		Name:    sql.NullString{String: "hello", Valid: true},
		Gender:  &sql.NullString{String: "M", Valid: true},
		Age:     sql.NullInt64{Int64: 18, Valid: true},
		Male:    sql.NullBool{Bool: true, Valid: true},
		Height:  sql.NullFloat64{Float64: 100.11, Valid: true},
		AddedAt: NullTime{Time: time.Now(), Valid: true},
	}).Error; err != nil {
		t.Errorf("Not error should raise when test null value")
	}

	var nv NullValue
	db.First(&nv, "name = ?", "hello")

	if nv.Name.String != "hello" || nv.Gender.String != "M" || nv.Age.Int64 != 18 || nv.Male.Bool != true || nv.Height.Float64 != 100.11 || nv.AddedAt.Valid != true {
		t.Errorf("Should be able to fetch null value")
	}

	if err := db.Save(&NullValue{
		Name:    sql.NullString{String: "hello-2", Valid: true},
		Gender:  &sql.NullString{String: "F", Valid: true},
		Age:     sql.NullInt64{Int64: 18, Valid: false},
		Male:    sql.NullBool{Bool: true, Valid: true},
		Height:  sql.NullFloat64{Float64: 100.11, Valid: true},
		AddedAt: NullTime{Time: time.Now(), Valid: false},
	}).Error; err != nil {
		t.Errorf("Not error should raise when test null value")
	}

	var nv2 NullValue
	db.First(&nv2, "name = ?", "hello-2")
	if nv2.Name.String != "hello-2" || nv2.Gender.String != "F" || nv2.Age.Int64 != 0 || nv2.Male.Bool != true || nv2.Height.Float64 != 100.11 || nv2.AddedAt.Valid != false {
		t.Errorf("Should be able to fetch null value")
	}

	if err := db.Save(&NullValue{
		Name:    sql.NullString{String: "hello-3", Valid: false},
		Gender:  &sql.NullString{String: "M", Valid: true},
		Age:     sql.NullInt64{Int64: 18, Valid: false},
		Male:    sql.NullBool{Bool: true, Valid: true},
		Height:  sql.NullFloat64{Float64: 100.11, Valid: true},
		AddedAt: NullTime{Time: time.Now(), Valid: false},
	}).Error; err == nil {
		t.Errorf("Can't save because of name can't be null")
	}
}

func TestNullValuesWithFirstOrCreate(t *testing.T) {
	var nv1 = NullValue{
		Name:   sql.NullString{String: "first_or_create", Valid: true},
		Gender: &sql.NullString{String: "M", Valid: true},
	}

	var nv2 NullValue
	result := db.Where(nv1).FirstOrCreate(&nv2)

	if result.RowsAffected != 1 {
		t.Errorf("RowsAffected should be 1 after create some record")
	}

	if result.Error != nil {
		t.Errorf("Should not raise any error, but got %v", result.Error)
	}

	if nv2.Name.String != "first_or_create" || nv2.Gender.String != "M" {
		t.Errorf("first or create with nullvalues")
	}

	if err := db.Where(nv1).Assign(NullValue{Age: sql.NullInt64{Int64: 18, Valid: true}}).FirstOrCreate(&nv2).Error; err != nil {
		t.Errorf("Should not raise any error, but got %v", err)
	}

	//	if nv2.Age.Int64 != 18 {
	//		t.Errorf("should update age to 18")
	//	}
}

func TestTransaction(t *testing.T) {
	tx := db.Begin()
	u := User{Name: "transcation"}
	if err := tx.Save(&u).Error; err != nil {
		t.Errorf("No error should raise")
	}

	if err := tx.First(&User{}, "name = ?", "transcation").Error; err != nil {
		t.Errorf("Should find saved record")
	}

	if sqlTx, ok := tx.CommonDB().(*sql.Tx); !ok || sqlTx == nil {
		t.Errorf("Should return the underlying sql.Tx")
	}

	tx.Rollback()

	if err := tx.First(&User{}, "name = ?", "transcation").Error; err == nil {
		t.Errorf("Should not find record after rollback")
	}

	tx2 := db.Begin()
	u2 := User{Name: "transcation-2"}
	if err := tx2.Save(&u2).Error; err != nil {
		t.Errorf("No error should raise")
	}

	if err := tx2.First(&User{}, "name = ?", "transcation-2").Error; err != nil {
		t.Errorf("Should find saved record")
	}

	tx2.Commit()

	if err := db.First(&User{}, "name = ?", "transcation-2").Error; err != nil {
		t.Errorf("Should be able to find committed record")
	}
}

func TestRow(t *testing.T) {
	user1 := User{Name: "RowUser1", Age: 1, Birthday: parseTime("2000-1-1")}
	user2 := User{Name: "RowUser2", Age: 10, Birthday: parseTime("2010-1-1")}
	user3 := User{Name: "RowUser3", Age: 20, Birthday: parseTime("2020-1-1")}
	db.Save(&user1).Save(&user2).Save(&user3)

	row := db.Table("user").Where("name = ?", user2.Name).Select("age").Row()
	var age int64
	row.Scan(&age)
	if age != 10 {
		t.Errorf("Scan with Row")
	}
}

func TestRows(t *testing.T) {
	user1 := User{Name: "RowsUser1", Age: 1, Birthday: parseTime("2000-1-1")}
	user2 := User{Name: "RowsUser2", Age: 10, Birthday: parseTime("2010-1-1")}
	user3 := User{Name: "RowsUser3", Age: 20, Birthday: parseTime("2020-1-1")}
	db.Save(&user1).Save(&user2).Save(&user3)

	rows, err := db.Table("user").Where("name = ? or name = ?", user2.Name, user3.Name).Select("name, age").Rows()
	if err != nil {
		t.Errorf("Not error should happen, got %v", err)
	}

	count := 0
	for rows.Next() {
		var name string
		var age int64
		rows.Scan(&name, &age)
		count++
	}

	if count != 2 {
		t.Errorf("Should found two records")
	}
}

func TestScanRows(t *testing.T) {
	user1 := User{Name: "ScanRowsUser1", Age: 1, Birthday: parseTime("2000-1-1")}
	user2 := User{Name: "ScanRowsUser2", Age: 10, Birthday: parseTime("2010-1-1")}
	user3 := User{Name: "ScanRowsUser3", Age: 20, Birthday: parseTime("2020-1-1")}
	db.Save(&user1).Save(&user2).Save(&user3)

	rows, err := db.Table("user").Where("name = ? or name = ?", user2.Name, user3.Name).Select("name, age").Rows()
	if err != nil {
		t.Errorf("Not error should happen, got %v", err)
	}

	type Result struct {
		Name string
		Age  int
	}

	var results []Result
	for rows.Next() {
		var result Result
		if err := db.ScanRows(rows, &result); err != nil {
			t.Errorf("should get no error, but got %v", err)
		}
		results = append(results, result)
	}

	if !reflect.DeepEqual(results, []Result{{Name: "ScanRowsUser2", Age: 10}, {Name: "ScanRowsUser3", Age: 20}}) {
		t.Errorf("Should find expected results")
	}
}

func TestScan(t *testing.T) {
	user1 := User{Name: "ScanUser1", Age: 1, Birthday: parseTime("2000-1-1")}
	user2 := User{Name: "ScanUser2", Age: 10, Birthday: parseTime("2010-1-1")}
	user3 := User{Name: "ScanUser3", Age: 20, Birthday: parseTime("2020-1-1")}
	db.Save(&user1).Save(&user2).Save(&user3)

	type result struct {
		Name string
		Age  int
	}

	var res result
	db.Table("user").Select("name, age").Where("name = ?", user3.Name).Scan(&res)
	if res.Name != user3.Name {
		t.Errorf("Scan into struct should work")
	}

	var doubleAgeRes = &result{}
	if err := db.Table("user").Select("age + age as age").Where("name = ?", user3.Name).Scan(&doubleAgeRes).Error; err != nil {
		t.Errorf("Scan to pointer of pointer")
	}
	if doubleAgeRes.Age != res.Age*2 {
		t.Errorf("Scan double age as age")
	}

	var ress []result
	db.Table("user").Select("name, age").Where("name in (?)", []string{user2.Name, user3.Name}).Scan(&ress)
	if len(ress) != 2 || ress[0].Name != user2.Name || ress[1].Name != user3.Name {
		t.Errorf("Scan into struct map")
	}
}

func TestRaw(t *testing.T) {
	user1 := User{Name: "ExecRawSqlUser1", Age: 1, Birthday: parseTime("2000-1-1")}
	user2 := User{Name: "ExecRawSqlUser2", Age: 10, Birthday: parseTime("2010-1-1")}
	user3 := User{Name: "ExecRawSqlUser3", Age: 20, Birthday: parseTime("2020-1-1")}
	db.Save(&user1).Save(&user2).Save(&user3)

	type result struct {
		Name  string
		Email string
	}

	var ress []result
	db.Raw("SELECT name, age FROM user WHERE name = ? or name = ?", user2.Name, user3.Name).Scan(&ress)
	if len(ress) != 2 || ress[0].Name != user2.Name || ress[1].Name != user3.Name {
		t.Errorf("Raw with scan")
	}

	rows, _ := db.Raw("select name, age from user where name = ?", user3.Name).Rows()
	count := 0
	for rows.Next() {
		count++
	}
	if count != 1 {
		t.Errorf("Raw with Rows should find one record with name 3")
	}

	db.Exec("update user set name=? where name in (?)", "jinzhu", []string{user1.Name, user2.Name, user3.Name})
	if db.Where("name in (?)", []string{user1.Name, user2.Name, user3.Name}).First(&User{}).Error != ErrRecordNotFound {
		t.Error("Raw sql to update records")
	}
}

func TestGroup(t *testing.T) {
	rows, err := db.Select("name").Table("user").Group("name").Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			rows.Scan(&name)
		}
	} else {
		t.Errorf("Should not raise any error")
	}
}

func TestHaving(t *testing.T) {
	rows, err := db.Select("name, count(*) as total").Table("user").Group("name").Having("name IN (?)", []string{"2", "3"}).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var total int64
			rows.Scan(&name, &total)

			if name == "2" && total != 1 {
				t.Errorf("Should have one user having name 2")
			}
			if name == "3" && total != 2 {
				t.Errorf("Should have two user having name 3")
			}
		}
	} else {
		t.Errorf("Should not raise any error")
	}
}

func TestQueryBuilderSubselectInWhere(t *testing.T) {
	user := User{Name: "query_expr_select_ruser1", Email: "root@user1.com", Age: 32}
	db.Save(&user)
	user = User{Name: "query_expr_select_ruser2", Email: "nobody@user2.com", Age: 16}
	db.Save(&user)
	user = User{Name: "query_expr_select_ruser3", Email: "root@user3.com", Age: 64}
	db.Save(&user)
	user = User{Name: "query_expr_select_ruser4", Email: "somebody@user3.com", Age: 128}
	db.Save(&user)

	var users []User
	db.Select("*").Where("name IN (?)", db.
		Select("name").Table("user").Where("name LIKE ?", "query_expr_select%").QueryExpr()).Find(&users)

	if len(users) != 4 {
		t.Errorf("Four users should be found, instead found %d", len(users))
	}

	db.Select("*").Where("name LIKE ?", "query_expr_select%").Where("age >= (?)", db.
		Select("AVG(age)").Table("user").Where("name LIKE ?", "query_expr_select%").QueryExpr()).Find(&users)

	if len(users) != 2 {
		t.Errorf("Two users should be found, instead found %d", len(users))
	}
}

func TestQueryBuilderSubselectInHaving(t *testing.T) {
	user := User{Name: "query_expr_having_ruser1", Email: "root@user1.com", Age: 64}
	db.Save(&user)
	user = User{Name: "query_expr_having_ruser2", Email: "root@user2.com", Age: 128}
	db.Save(&user)
	user = User{Name: "query_expr_having_ruser3", Email: "root@user1.com", Age: 64}
	db.Save(&user)
	user = User{Name: "query_expr_having_ruser4", Email: "root@user2.com", Age: 128}
	db.Save(&user)

	var users []User
	db.Select("AVG(age) as avgage").Where("name LIKE ?", "query_expr_having_%").Group("email").Having("AVG(age) > (?)", db.
		Select("AVG(age)").Where("name LIKE ?", "query_expr_having_%").Table("user").QueryExpr()).Find(&users)

	if len(users) != 1 {
		t.Errorf("Two user group should be found, instead found %d", len(users))
	}
}

func DialectHasTzSupport() bool {
	// NB: mssql and FoundationDB do not support time zones.
	if dialect := os.Getenv("micros_DIALECT"); dialect == "foundation" {
		return false
	}
	return true
}

func TestTimeWithZone(t *testing.T) {
	var format = "2006-01-02 15:04:05 -0700"
	var times []time.Time
	GMT8, _ := time.LoadLocation("Asia/Shanghai")
	times = append(times, time.Date(2013, 02, 19, 1, 51, 49, 123456789, GMT8))
	times = append(times, time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.UTC))

	for index, vtime := range times {
		name := "time_with_zone_" + strconv.Itoa(index)
		user := User{Name: name, Birthday: &vtime}

		if !DialectHasTzSupport() {
			// If our driver dialect doesn't support TZ's, just use UTC for everything here.
			utcBirthday := user.Birthday.UTC()
			user.Birthday = &utcBirthday
		}

		db.Save(&user)
		expectedBirthday := "2013-02-18 17:51:49 +0000"
		foundBirthday := user.Birthday.UTC().Format(format)
		if foundBirthday != expectedBirthday {
			t.Errorf("User's birthday should not be changed after save for name=%s, expected bday=%+v but actual value=%+v", name, expectedBirthday, foundBirthday)
		}

		var findUser, findUser2 User
		db.First(&findUser, "name = ?", name)
		foundBirthday = findUser.Birthday.UTC().Format(format)
		if foundBirthday != expectedBirthday {
			t.Errorf("User's birthday should not be changed after find for name=%s, expected bday=%+v but actual value=%+v", name, expectedBirthday, foundBirthday)
		}

		if db.Where("id = ? AND birthday >= ?", findUser.Id, user.Birthday.Add(-time.Minute)).First(&findUser2).RecordNotFound() {
			t.Errorf("User should be found")
		}
	}
}

func TestSetAndGet(t *testing.T) {
	if value, ok := db.Set("hello", "world").Get("hello"); !ok {
		t.Errorf("Should be able to get setting after set")
	} else {
		if value.(string) != "world" {
			t.Errorf("Setted value should not be changed")
		}
	}

	if _, ok := db.Get("non_existing"); ok {
		t.Errorf("Get non existing key should return error")
	}
}

func TestOpenExistingDB(t *testing.T) {
	db.Save(&User{Name: "jnfeinstein"})
	dialect := os.Getenv("micros_DIALECT")

	ndb, err := Open(dialect, db.DB())
	if err != nil {
		t.Errorf("Should have wrapped the existing DB connection")
	}

	var user User
	if ndb.Where("name = ?", "jnfeinstein").First(&user).Error == ErrRecordNotFound {
		t.Errorf("Should have found existing record")
	}
}

func TestDdlErrors(t *testing.T) {
	var err error

	if err = db.Close(); err != nil {
		t.Errorf("Closing DDL test db connection err=%s", err)
	}
	defer func() {
		// Reopen DB connection.
		if db, err = OpenTestConnection(); err != nil {
			t.Fatalf("Failed re-opening db connection: %s", err)
		}
	}()

	if err := db.Find(&User{}).Error; err == nil {
		t.Errorf("Expected operation on closed db to produce an error, but err was nil")
	}
}

func TestOpenWithOneParameter(t *testing.T) {
	ndb, err := Open("dialect")
	if ndb != nil {
		t.Error("Open with one parameter returned non nil for db")
	}
	if err == nil {
		t.Error("Open with one parameter returned err as nil")
	}
}

func TestBlockGlobalUpdate(t *testing.T) {
	db.DropTable(&Toy{})
	db.CreateTable(&Toy{})
	ndb := db.New()
	ndb.Create(&Toy{Name: "Stuffed Animal", OwnerType: "Nobody"})

	err := ndb.Model(&Toy{}).Update("OwnerType", "Human").Error
	if err != nil {
		t.Error("Unexpected error on global update")
	}

	err = ndb.Delete(&Toy{}).Error
	if err != nil {
		t.Error("Unexpected error on global delete")
	}

	ndb.Create(&Toy{Name: "Stuffed Animal", OwnerType: "Nobody"})

	err = ndb.Model(&Toy{}).Where(&Toy{OwnerType: "Martian"}).Update("OwnerType", "Astronaut").Error
	if err != nil {
		t.Error("Unxpected error on conditional update")
	}

	err = ndb.Where(&Toy{OwnerType: "Martian"}).Delete(&Toy{}).Error
	if err != nil {
		t.Error("Unexpected error on conditional delete")
	}
}

func BenchmarkORM(b *testing.B) {
	b.N = 2000
	for x := 0; x < b.N; x++ {
		e := strconv.Itoa(x) + "benchmark@example.org"
		now := time.Now()
		email := EmailWithIdx{Email: e, UserAgent: "pc", RegisteredAt: &now}
		// Insert
		db.Save(&email)
		// Query
		db.First(&EmailWithIdx{}, "email = ?", e)
		// Update
		db.Model(&email).UpdateColumn("email", "new-"+e)
		// Delete
		db.Delete(&email)
	}
}

func parseTime(str string) *time.Time {
	t := now.New(time.Now().UTC()).MustParse(str)
	return &t
}
