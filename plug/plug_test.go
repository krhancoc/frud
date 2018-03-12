package plug_test

import (
	"testing"

	"github.com/krhancoc/frud/plug"
)

type Check interface {
	Get()
	Hello()
	There()
}

type ObjOne struct{}

func (*ObjOne) Get()   {}
func (*ObjOne) There() {}
func (*ObjOne) Hello() {}

type ObjTwo struct{}

func (*ObjTwo) Get()   {}
func (*ObjTwo) hello() {}

type ObjThree struct{}

var unimpimentedTests = []struct {
	value    interface{}
	expected int
}{
	{&ObjOne{}, 0},
	{&ObjTwo{}, 2},
	{&ObjThree{}, 3},
}

func TestCheckUnimplimented(t *testing.T) {

	for _, test := range unimpimentedTests {
		out := plug.CheckUnimplimented(test.value, (*Check)(nil))
		if len(out) != test.expected {
			t.Error("Expected", test.expected, "got", len(out))
		}
	}
}
