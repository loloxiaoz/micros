package orm

import (
	"testing"
)

type CalculateField struct {
	Model
	Name     string
	Children []CalculateFieldChild
	Category CalculateFieldCategory
	TestTag  string `sql:"NOT NULL;DEFAULT:'hello'"`
}

type CalculateFieldChild struct {
	Model
	CalculateFieldID uint
	Name             string
}

type CalculateFieldCategory struct {
	Model
	CalculateFieldID uint
	Name             string
}

func TestCalculateField(t *testing.T) {
	var field CalculateField
	var scope = db.NewScope(&field)
	if _, ok := scope.FieldByName("Children"); !ok {
		t.Errorf("Should calculate fields correctly for the first time")
	}

	if _, ok := scope.FieldByName("Category"); !ok {
		t.Errorf("Should calculate fields correctly for the first time")
	}

	if field, ok := scope.FieldByName("TestTag"); !ok {
		t.Errorf("should find TestTag ")
	} else if _, ok := field.TagSettings["NOT NULL"]; !ok {
		t.Errorf("should find TestTag's tag settings")
	}
}
