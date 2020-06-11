package create

import (
	"testing"
)

var tInfo = newTplInfo([]byte(`
package create
import (
	"testing"
)
import f "fmt"

type __MYSQL_MODEL__ struct {
	A
}

type __MONGO_MODEL__ struct {
	B
}

// A comment ...
type A struct{
	// X doc ...
	X string // X comment ...
	// Y doc ...
	Y int // Y comment ...
}

// B comment ...
type B struct{
	// X doc ...
	X string // X comment ...
	// Y doc ...
	Y int // Y comment ...
}
`))

func TestParse(t *testing.T) {
	tInfo.Parse()
	t.Logf("TypeImportString: %s", tInfo.TypeImportString())
	t.Logf("TypesString:\n%v", tInfo.TypesString())
	for _, m := range tInfo.models.mysql {
		t.Logf("mysql:\n%v", m)
	}
	for _, m := range tInfo.models.mongo {
		t.Logf("mongo:\n%v", m)
	}
}
