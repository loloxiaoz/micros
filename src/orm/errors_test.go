package orm

import (
	"errors"
	"testing"
)

func TestErrorsCanBeUsedOutside(t *testing.T) {
	errs := []error{errors.New("First"), errors.New("Second")}

	gErrs := Errors(errs)
	gErrs = gErrs.Add(errors.New("Third"))
	gErrs = gErrs.Add(gErrs)

	if gErrs.Error() != "First; Second; Third" {
		t.Fatalf("Gave wrong error, got %s", gErrs.Error())
	}
}
