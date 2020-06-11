package create

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"github.com/henrylee2cn/erpc/v6"
	"github.com/henrylee2cn/goutil"
	"github.com/xiaoenai/xmodel/cmd/create/structtag"
)

const (
	// MYSQL_MODEL name of the struct used to create mysql model
	MYSQL_MODEL = "__MYSQL_MODEL__"
	// MONGO_MODEL name of the struct used to create mongo model
	MONGO_MODEL = "__MONGO_MODEL__"
)

type (
	tplInfo struct {
		src               []byte
		fileSet           *token.FileSet
		astFile           *ast.File
		doc               string
		models            *models
		realStructTypes   []*structType
		realStructTypeMap map[string]*structType
		aliasTypes        []*aliasType
		typeImports       []string
	}
	models struct {
		mysql []*structType
		mongo []*structType
	}
	structType struct {
		doc              string
		name             string
		fields           []*field
		primaryFields    []*field
		uniqueFields     []*field
		isDefaultPrimary bool
		modelStyle       string // mysql, mongo
		node             *ast.StructType
	}
	field struct {
		Name      string
		ModelName string
		Typ       string
		isQuery   bool
		queryName string
		anonymous bool
		tag       string
		doc       string
		comment   string
	}
	aliasType struct {
		doc         string
		name        string
		text        string
		rawTypeName string
		rawStruct   *structType
	}
)

func newTplInfo(tplBytes []byte) *tplInfo {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", tplBytes, parser.ParseComments)
	if err != nil {
		erpc.Fatalf("[XModel] %v", err)
	}
	return &tplInfo{
		src:     tplBytes,
		fileSet: fset,
		astFile: file,
		doc:     addSlash(file.Doc.Text()),
		models:            new(models),
		realStructTypeMap: make(map[string]*structType),
	}
}

func (t *tplInfo) Parse() *tplInfo {
	t.parseImports()
	ok := t.hasType(emptyStructType.name)
	if ok {
		erpc.Fatalf("[XModel] Keep structure name cannot be used: %s", emptyStructType.name)
	}
	t.collectStructs()
	t.initModelStructs()
	return t
}

func (t *tplInfo) TypeImportString() string {
	return strings.Join(t.typeImports, "\n")
}

func (t *tplInfo) TypesString() string {
	var a string
	for _, s := range t.aliasTypes {
		a += s.String()
	}
	for _, s := range t.realStructTypes {
		a += s.String()
	}
	return a
}

func (t *tplInfo) getCodeBlock(i interface{}) string {
	var dst bytes.Buffer
	err := format.Node(&dst, t.fileSet, i)
	if err != nil {
		erpc.Fatalf("[XModel] %v", err)
	}
	return dst.String()
}

func (t *tplInfo) parseImports() {
	for _, imp := range t.astFile.Imports {
		s := t.getCodeBlock(imp)
		t.typeImports = append(t.typeImports, s)
	}
}

// collectStructs collects and maps structType nodes to their positions
func (t *tplInfo) collectStructs() {
	collectStructs := func(n ast.Node) bool {
		decl, ok := n.(ast.Decl)
		if !ok {
			return true
		}
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			return true
		}
		var groupDoc string
		if len(genDecl.Specs) == 1 {
			groupDoc = genDecl.Doc.Text()
		}
		for _, spec := range genDecl.Specs {
			var e ast.Expr
			var structName string
			var doc = groupDoc

			switch x := spec.(type) {
			case *ast.TypeSpec:
				if x.Type == nil {
					continue
				}
				structName = x.Name.Name
				e = x.Type
				if s := x.Doc.Text(); s != "" {
					doc = x.Doc.Text()
				}
			}

			x, ok := e.(*ast.StructType)
			if !ok {
				continue
			}

			if len(x.Fields.List) == 0 {
				switch structName {
				case MYSQL_MODEL, MONGO_MODEL:
				default:
					if goutil.IsExportedName(structName) {
						a := &aliasType{
							doc:  addSlash(doc),
							name: structName,
							text: fmt.Sprintf("%s = struct{}", structName),
						}
						a.rawTypeName = a.text[strings.LastIndex(strings.TrimSpace(strings.Split(a.text, "//")[0]), " ")+1:]
						if a.doc == "" {
							a.doc = fmt.Sprintf("// %s alias of type %s\n", a.name, a.rawTypeName)
						}
						t.aliasTypes = append(t.aliasTypes, a)
					}
					continue
				}
			}

			t.realStructTypes = append(
				t.realStructTypes,
				structType{
					name: structName,
					doc:  addSlash(doc),
					node: x,
				}.init(t),
			)
		}
		return true
	}
	ast.Inspect(t.astFile, collectStructs)
	t.sortStructs()
}

func (t *tplInfo) lookupAliasType(name string) (*aliasType, bool) {
	for _, a := range t.aliasTypes {
		if a.name == name {
			return a, true
		}
	}
	return nil, false
}

func (t *tplInfo) lookupStructType(name string) (*structType, bool) {
	s, ok := t.realStructTypeMap[name]
	return s, ok
}

func (t *tplInfo) hasType(name string) bool {
	_, ok := t.lookupStructType(name)
	if !ok {
		_, ok = t.lookupAliasType(name)
	}
	return ok
}

func (t *tplInfo) lookupTypeFields(name string) ([]*field, bool) {
	s, ok := t.lookupStructType(name)
	if ok {
		return s.fields, true
	}
	a, ok := t.lookupAliasType(name)
	if ok {
		if a.rawStruct != nil {
			return a.rawStruct.fields, true
		}
	}
	return nil, false
}


func (t *tplInfo) initModelStructs() {
	for _, s := range t.models.mysql {
		s.initModel()
	}
	if len(t.models.mongo) > 0 {
		for _, s := range t.models.mongo {
			s.initModel()
		}
		var hasMongo bool
		const mongoImp = `"github.com/xiaoenai/xmodel/mongo"`
		for _, imp := range t.typeImports {
			if imp == mongoImp {
				hasMongo = true
				break
			}
		}
		if !hasMongo {
			t.typeImports = append(t.typeImports, mongoImp)
		}
	}
}

func (s *structType) initModel() {
	s.setFields(false, &field{
		Name: "UpdatedAt",
		Typ:  "int64",
		tag:  "`" + `json:"updated_at"` + "`",
	}, &field{
		Name: "CreatedAt",
		Typ:  "int64",
		tag:  "`" + `json:"created_at"` + "`",
	}, &field{
		Name: "DeletedTs",
		Typ:  "int64",
		tag:  "`" + `json:"deleted_ts"` + "`",
	})

	switch s.modelStyle {
	case "mysql":
		var hasPrimary bool
		s.rangeTags(func(tags *structtag.Tags, f *field, anonymous bool) bool {
			tag, _ := tags.Get("json")
			f.ModelName = tag.Name
			tag, err := tags.Get("key")
			if err == nil {
				if tag.Name == "pri" {
					s.primaryFields = append(s.primaryFields, f)
					hasPrimary = true
				} else {
					tags.Set(&structtag.Tag{
						Key:  "key",
						Name: "uni",
					})
					s.uniqueFields = append(s.uniqueFields, f)
				}
			}
			return true
		})
		if !hasPrimary {
			s.setFields(true, &field{
				Name:      "Id",
				ModelName: "id",
				Typ:       "int64",
				tag:       "`" + `json:"id" key:"pri"` + "`",
			})
			s.primaryFields = append(s.primaryFields, s.fields[0])
		}
		if len(s.primaryFields) == 1 && s.primaryFields[0].Typ == "int64" {
			s.isDefaultPrimary = true
		}

	case "mongo":
		var hasObjectId bool
		s.rangeTags(func(tags *structtag.Tags, f *field, anonymous bool) bool {
			if f.Typ != "mongo.ObjectId" {
				f.ModelName = goutil.SnakeString(f.Name)
				tags.Set(&structtag.Tag{
					Key:  "bson",
					Name: f.ModelName,
				})
				_, err := tags.Get("key")
				if err == nil {
					tags.Set(&structtag.Tag{
						Key:  "key",
						Name: "uni",
					})
					s.uniqueFields = append(s.uniqueFields, f)
				}
			} else if !hasObjectId {
				hasObjectId = true
				s.primaryFields = append(s.primaryFields, f)
				f.ModelName = "_id"
				tags.Set(&structtag.Tag{
					Key:  "bson",
					Name: "_id",
				})
				tags.Set(&structtag.Tag{
					Key:  "json",
					Name: "_id",
				})
				tags.Set(&structtag.Tag{
					Key:  "key",
					Name: "pri",
				})
			}
			return true
		})
		if !hasObjectId {
			s.setFields(true, &field{
				Name:      "Id",
				ModelName: "_id",
				Typ:       "mongo.ObjectId",
				tag:       "`" + `json:"_id" bson:"_id" key:"pri"` + "`",
			})
			s.primaryFields = append(s.primaryFields, s.fields[0])
		}
		if len(s.primaryFields) == 1 && s.primaryFields[0].Typ == "mongo.ObjectId" {
			s.isDefaultPrimary = true
		}
	}
}

func (s *structType) getField(fieldName string) *field {
	for _, f := range s.fields {
		if f.Name == fieldName {
			return f
		}
	}
	return nil
}

func (s *structType) setFields(toLeader bool, fields ...*field) {
	for _, f := range fields {
		for i, ff := range s.fields {
			if ff.Name == f.Name {
				s.fields = append(s.fields[:i], s.fields[i+1:]...)
				break
			}
		}
		if toLeader {
			s.fields = append([]*field{f}, s.fields...)
		} else {
			s.fields = append(s.fields, f)
		}
	}
}

func (s *structType) isInvildName() bool {
	switch s.name {
	case MYSQL_MODEL, MONGO_MODEL:
		return true
	default:
		return goutil.IsExportedName(s.name)
	}
}

func (s structType) init(t *tplInfo) *structType {
	if !s.isInvildName() {
		erpc.Fatalf("[XModel] Unexported struct name: %s", s.name)
	}
	for _, v := range s.node.Fields.List {
		f := new(field)
		if len(v.Names) > 0 {
			f.Name = v.Names[0].Name
			if !goutil.IsExportedName(f.Name) {
				erpc.Fatalf("[XModel] Unexported field name: %s.%s", s.name, f.Name)
			}
		}
		f.Typ = t.getCodeBlock(v.Type)
		if len(f.Name) == 0 {
			f.anonymous = true
			f.Name = strings.TrimPrefix(f.Typ, "*")
			if !goutil.IsExportedName(f.Name) {
				erpc.Fatalf("[XModel] Unexported anonymous field: %s.%s", s.name, f.Typ)
			}
		}
		f.doc = addSlash(v.Doc.Text())
		f.comment = addSlash(v.Comment.Text())
		if v.Tag != nil {
			f.tag = v.Tag.Value
			f.queryName, f.isQuery = getQueryField(f.tag)
			if len(f.queryName) == 0 {
				f.queryName = goutil.SnakeString(f.Name)
			}
		}
		s.fields = append(s.fields, f)
	}
	s.rangeTags(
		addJsonTag,
	)
	if s.doc == "" {
		s.doc = fmt.Sprintf("// %s comment...\n", s.name)
	}
	return &s
}

var queryRegexp = regexp.MustCompile("<\\s*query\\s*(:[^:>]*)?>")

func getQueryField(tag string) (queryName string, isQuery bool) {
	a := queryRegexp.FindStringSubmatch(tag)
	if len(a) != 2 {
		return
	}
	isQuery = true
	queryName = strings.TrimLeft(a[1], ":")
	queryName = strings.TrimSpace(queryName)
	return
}

func (s *structType) rangeTags(fns ...func(tags *structtag.Tags, f *field, anonymous bool) bool) {
	for _, fn := range fns {
		for _, v := range s.fields {
			logName := v.Name
			if len(logName) == 0 {
				logName = v.Typ
			}
			tags, err := structtag.Parse(strings.TrimSpace(strings.Trim(v.tag, "`")))
			if err != nil {
				erpc.Fatalf("[XModel] %s.%s: %s", s.name, logName, err.Error())
			}
			if !fn(tags, v, len(v.Name) == 0) {
				break
			}
			v.tag = "`" + tags.String() + "`"
		}
	}
}

var addJsonTag = func(tags *structtag.Tags, f *field, anonymous bool) bool {
	tag, _ := tags.Get("json")
	if tag != nil {
		return true
	}
	tags.Set(&structtag.Tag{
		Key:  "json",
		Name: goutil.SnakeString(f.Name),
	})
	return true
}

func (s *structType) String() string {
	r := fmt.Sprintf("%stype %s struct {\n", s.doc, s.name)
	for _, f := range s.fields {
		if f.anonymous {
			r += fmt.Sprintf("%s  %s  %s  %s", f.doc, f.Typ, f.tag, f.comment)
		} else {
			r += fmt.Sprintf("%s  %s  %s  %s", f.doc+f.Name, f.Typ, f.tag, f.comment)
		}
		if r[len(r)-1] != '\n' {
			r += "\n"
		}
	}
	r += "}\n\n"
	return r
}

func (t *tplInfo) sortStructs() {
	var lastList []*structType
	var mysqlList []*structType
	var mongoList []*structType
	for _, v := range t.realStructTypes {
		switch v.name {
		case MYSQL_MODEL:
			mysqlList = append(mysqlList, v)
		case MONGO_MODEL:
			mongoList = append(mongoList, v)
		default:
			lastList = append(lastList, v)
			t.realStructTypeMap[v.name] = v
		}
	}
	t.realStructTypes = lastList
	for _, name := range getStructFieldNames(mysqlList) {
		for i := 0; i < len(t.realStructTypes); i++ {
			v := t.realStructTypes[i]
			if v.name == name {
				if len(v.modelStyle) > 0 {
					erpc.Fatalf("[XModel] %s: multiple specified model style", v.name)
				}
				t.models.mysql = append(t.models.mysql, v)
				v.modelStyle = "mysql"
				break
			}
		}
	}
	for _, name := range getStructFieldNames(mongoList) {
		for i := 0; i < len(t.realStructTypes); i++ {
			v := t.realStructTypes[i]
			if v.name == name {
				if len(v.modelStyle) > 0 {
					erpc.Fatalf("[XModel] %s: multiple specified model style", v.name)
				}
				t.models.mongo = append(t.models.mongo, v)
				v.modelStyle = "mongo"
				break
			}
		}
	}
}

func getStructFieldNames(v []*structType) (a []string) {
	for _, s := range v {
		for _, f := range s.fields {
			a = append(a, f.Name)
		}
	}
	return
}

var emptyStructType = &aliasType{
	doc:  "// EmptyStruct alias of type struct {}\n",
	name: "EmptyStruct",
	text: "EmptyStruct = struct{}",
}

func (a *aliasType) String() string {
	return fmt.Sprintf("%stype %s\n", a.doc, a.text)
}

func joinName(a, b string) string {
	a = strings.Trim(a, "_")
	b = strings.Trim(b, "_")
	if a == "" {
		return b
	}
	a = strings.ToUpper(a[:1]) + a[1:]
	if b == "" {
		return a
	}
	b = strings.ToUpper(b[:1]) + b[1:]
	return a + "_" + b
}

func firstLowerLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(string(s[0]))
}

func addSlash(txt string) (comment string) {
	r := bufio.NewReader(strings.NewReader(txt))
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		comment += "// " + string(line) + "\n"
	}
	return
}
