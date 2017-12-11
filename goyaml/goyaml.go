package goyaml

import (
	"bytes"
	"go/build"
	"go/doc"
	"go/printer"
	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
	"gopkg.in/yaml.v2"
	"path/filepath"
	pathpkg "path"
	"os"
	"fmt"
)

//  GoYAMLGeneration Generate the YAML file for golang projects
func GoYAMLGeneration(packageSource string, packageName string, ymlOutput string) error {
	// initialization
	ns := vfs.NameSpace{}
	ns.Bind("/", vfs.OS("C:/Go"), "/", vfs.BindReplace)
	ns.Bind("/src", vfs.OS(packageSource), "/", vfs.BindAfter)
	c := godoc.NewCorpus(ns)
	p := godoc.NewPresentation(c)
	// Begin of get package info
	abspath, relpath := paths(ns, p, packageName)
	var mode godoc.PageInfoMode
	mode = godoc.NoHTML
	var info *godoc.PageInfo
	info = p.GetPkgPageInfo(abspath, relpath, mode)
	// end of get package info
	// to YAML struct
	docPackage := ToDocfx(info)
	// create YAML file
	yamlFile, err := os.Create(ymlOutput + "/" + packageName + ".yml")
	if err != nil {
		fmt.Errorf("Failed to create file: ", packageName)
		return nil
	}
	yamlBytes, err := yaml.Marshal(docPackage)
	if err != nil {
		fmt.Errorf("Failed to Marshal")
		return err
	}
	yamlFile.WriteString("#YamlMIME: GoLangPkg\n")
	yamlFile.Write(yamlBytes)
	yamlFile.Close()
	return nil
}


// paths determines the paths to use.
//
// If we are passed an operating system path like . or ./foo or /foo/bar or c:\mysrc,
// we need to map that path somewhere in the fs name space so that routines
// like getPageInfo will see it.  We use the arbitrarily-chosen virtual path "/target"
// for this.  That is, if we get passed a directory like the above, we map that
// directory so that getPageInfo sees it as /target.
// Returns the absolute and relative paths.
func paths(fs vfs.NameSpace, pres *godoc.Presentation, path string) (string, string) {
	if filepath.IsAbs(path) {
		fs.Bind(target, vfs.OS(path), "/", vfs.BindReplace)
		return target, target
	}
	if build.IsLocalImport(path) {
		cwd, _ := os.Getwd() // ignore errors
		path = filepath.Join(cwd, path)
		fs.Bind(target, vfs.OS(path), "/", vfs.BindReplace)
		return target, target
	}
	if bp, _ := build.Import(path, "", build.FindOnly); bp.Dir != "" && bp.ImportPath != "" {
		fs.Bind(target, vfs.OS(bp.Dir), "/", vfs.BindReplace)
		return target, bp.ImportPath
	}
	return pathpkg.Join(pres.PkgFSRoot(), path), path
}

const (
	target = "/target"
)

func nodeFunc(info *godoc.PageInfo, node interface{}) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, info.FSet, node)
	return buf.String()
}



type DocsPackage struct {
	IsMain      bool                  `json:"ismain"`
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	ImportPath  string                `json:"importPath"`
	Dir         string                `json:"dir"`
	Consts      []DocsValue           `json:"consts"`
	Types       []DocsType            `json:"types"`
	Vars        []DocsValue           `json:"vars"`
	Funcs       []DocsFunc            `json:"funcs"`
	Notes       map[string][]DocsNote `json:"notes"`
	Examples    []DocsExample         `json:"examples"`
	Dirs        []DocsDir             `json:"dirs"`
}

type DocsDir struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Summary string `json:"summary"`
	HasPkg  bool   `json:"haspkg"`
}

type DocsNote struct {
	UID         string `json:"uid"`
	Description string `json:"description"`
}

type DocsExample struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type DocsValue struct {
	Names       []string `json:"names"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Code        string   `json:"code"`
}

type DocsType struct {
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Code        string `json:"code"`

	Consts  []DocsValue `json:"consts"`
	Vars    []DocsValue `json:"vars"`
	Funcs   []DocsFunc  `json:"funcs"`
	Methods []DocsFunc  `json:"methods"`
}

type DocsFunc struct {
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

func ToDocfx(info *godoc.PageInfo) *DocsPackage {
	pkg := &DocsPackage{
		IsMain: info.IsMain,
		Dir:    info.Dirname,
		Notes:  toDocsNotes(info.Notes),
		Dirs:   toDocsDirs(info.Dirs),
	}
	if info.PDoc != nil {
		pkg.ImportPath = info.PDoc.ImportPath
		pkg.Summary = summary(info.PDoc.Doc)
		pkg.Description = description(info.PDoc.Doc)
		pkg.Examples = toDocsExamples(info.Examples, info)
		pkg.Consts = toDocsValues(info.PDoc.Consts, info)
		pkg.Vars = toDocsValues(info.PDoc.Vars, info)
		pkg.Funcs = toDocsFuncs(info.PDoc.Funcs, info)
		pkg.Types = toDocsTypes(info.PDoc.Types, info)
	}
	return pkg
}

func toDocsDirs(dirs *godoc.DirList) []DocsDir {
	if dirs == nil {
		return []DocsDir{}
	}

	arr := make([]DocsDir, len(dirs.List))
	for i, d := range dirs.List {
		arr[i] = DocsDir{
			Name:    d.Name,
			Path:    d.Path,
			Summary: d.Synopsis,
			HasPkg:  d.HasPkg,
		}
	}
	return arr
}

func toDocsTypes(types []*doc.Type, info *godoc.PageInfo) []DocsType {
	arr := make([]DocsType, len(types))
	for i, t := range types {
		arr[i] = DocsType{
			Name:        t.Name,
			Summary:     summary(t.Doc),
			Description: description(t.Doc),
			Code:        nodeFunc(info, t.Decl),
			Consts:      toDocsValues(t.Consts, info),
			Vars:        toDocsValues(t.Vars, info),
			Funcs:       toDocsFuncs(t.Funcs, info),
			Methods:     toDocsFuncs(t.Methods, info),
		}
	}
	return arr
}

func toDocsFuncs(funcs []*doc.Func, info *godoc.PageInfo) []DocsFunc {
	arr := make([]DocsFunc, len(funcs))
	for i, f := range funcs {
		arr[i] = DocsFunc{
			Name:        f.Name,
			Summary:     summary(f.Doc),
			Description: description(f.Doc),
			Code:        nodeFunc(info, f.Decl),
		}
	}
	return arr
}

func toDocsValues(values []*doc.Value, info *godoc.PageInfo) []DocsValue {
	arr := make([]DocsValue, len(values))
	for i, v := range values {
		arr[i] = DocsValue{
			Names:       v.Names,
			Summary:     summary(v.Doc),
			Description: description(v.Doc),
			Code:        nodeFunc(info, v.Decl),
		}
	}
	return arr
}

func toDocsExamples(examples []*doc.Example, info *godoc.PageInfo) []DocsExample {
	arr := make([]DocsExample, len(examples))
	for i, eg := range examples {
		cnode := &printer.CommentedNode{Node: eg.Code, Comments: eg.Comments}
		arr[i] = DocsExample{
			Name: eg.Name,
			Code: nodeFunc(info, cnode),
		}
	}
	return arr
}

func toDocsNotes(notes map[string][]*doc.Note) map[string][]DocsNote {
	m := map[string][]DocsNote{}
	for k, v := range notes {
		arr := make([]DocsNote, len(v))
		for i, n := range v {
			arr[i] = DocsNote{
				UID:         n.UID,
				Description: n.Body,
			}
		}
		m[k] = arr
	}
	return m
}

func summary(d string) string {
	return doc.Synopsis(d)
}

func description(d string) string {
	var buf bytes.Buffer
	doc.ToText(&buf, d, "", "    ", 999999)
	return buf.String()
}