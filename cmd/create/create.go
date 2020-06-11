package create

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/henrylee2cn/erpc/v6"
	"github.com/henrylee2cn/goutil"
	"github.com/xiaoenai/xmodel/cmd/create/tpl"
	"github.com/xiaoenai/xmodel/cmd/info"
)

// ModelTpl template file name
const ModelTpl = "__model__tpl__.go"

// ModelGenLock the file is used to markup generated project
const ModelGenLock = "__model__gen__.lock"

// CreateProject creates a project.
func CreateProject(force bool) {
	erpc.Infof("Generating project: %s", info.ProjPath())

	os.MkdirAll(info.AbsPath(), os.FileMode(0755))
	err := os.Chdir(info.AbsPath())
	if err != nil {
		erpc.Fatalf("[XModel] Jump working directory failed: %v", err)
	}

	force = force || !goutil.FileExists(ModelGenLock)

	// creates base files
	if force {
		tpl.Create()
	}

	// read temptale file
	b, err := ioutil.ReadFile(ModelTpl)
	if err != nil {
		b = []byte(strings.Replace(__tpl__, "__PROJ_NAME__", info.ProjName(), -1))
	}

	// new project code
	proj := NewProject(b)
	proj.Generator(force)

	// write template file
	f, err := os.OpenFile(ModelTpl, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		erpc.Fatalf("[XModel] Create files error: %v", err)
	}
	defer f.Close()
	f.Write(formatSource(b))

	tpl.RestoreAsset("./", ModelGenLock)

	erpc.Infof("Completed code generation!")
}
