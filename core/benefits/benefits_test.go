package benefits

import (
	"testing"

	"github.com/malkhamis/quantax/core/human"
)

func Test_getChildCount(t *testing.T) {

	children := []*human.Person{
		&human.Person{},
		nil,
	}

	actual := getChildCount(children)
	if actual != 1 {
		t.Fatalf("expected nil child to not be counted, got: %d", actual)
	}
}
