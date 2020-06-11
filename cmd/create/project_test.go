package create

import (
	"strings"
	"testing"

	"github.com/xiaoenai/xmodel/cmd/info"
)

func TestGenerator(t *testing.T) {
	info.Init("test")
	proj := NewProject([]byte(__tpl__))
	proj.gen()
	t.Logf("args/const.gen.go:\n%s", proj.codeFiles["args/const.gen.go"])
	t.Logf("args/type.gen.go:\n%s", proj.codeFiles["args/type.gen.go"])
	t.Logf("model/init.go:\n%s", proj.codeFiles["model/init.go"])
	for k, v := range proj.codeFiles {
		if strings.HasPrefix(k, "model") {
			t.Logf("%s:\n%s", k, v)
		}
	}
}
