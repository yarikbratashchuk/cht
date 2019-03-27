package fname

import "testing"

func TestCurrent(t *testing.T) {
	name := Current()
	if name != "fname.TestCurrent" {
		t.Error("did not work")
	}

	tt := test{}
	mname := tt.testMethod()
	if mname != "fname.test.testMethod" {
		t.Error("did not work")
	}
}

type test struct{}

func (t test) testMethod() string {
	return Current()
}
