package orm

import "time"

//Model base model definition
//type User struct {
//      orm.Model
//}
type Model struct {
	ID        uint `gorm:"primary_key"`
	Ver       int  `sql:"version"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
