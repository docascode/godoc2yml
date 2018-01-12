package goyaml

import (
	"bytes"
	"fmt"
	"go/build"
	"go/doc"
	"go/printer"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
	"gopkg.in/yaml.v2"
)

var SourceRepo string
var SourceBranch string

//  GoYAMLGeneration Generate the YAML file for golang projects
func GoYAMLGeneration(packageSource string, ymlOutput string, packagePrefix string) error {
	// Split package name from packageSource
	packagePaths := strings.Split(packageSource, "/")
	packageName := packagePaths[len(packagePaths)-1]
	packagePaths = packagePaths[:len(packagePaths)-1]
	packageSource = strings.Join(packagePaths, "/")

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

	// Begin: get position
	// PrintPosition(info) - only for test purpose
	// End: get position

	// to YAML struct
	docPackage, types := ToDocfx(info, packageName, packagePrefix)

	// create package YAML file
	var packagePath string
	if packagePrefix == "" {
		packagePath = ymlOutput
	} else {
		packagePath = ymlOutput + "/" + packagePrefix
	}
	os.Mkdir(packagePath, os.ModePerm)
	yamlBytes, err := yaml.Marshal(docPackage)
	if err != nil {
		fmt.Errorf("Failed to Marshal")
		return err
	}
	PrintYaml(yamlBytes, packagePath, packageName)

	// create type Yaml files
	typePath := packagePath + "/" + packageName
	os.Mkdir(typePath, os.ModePerm)
	for _, t := range types {
		yamlBytes, err = yaml.Marshal(t)
		if err != nil {
			fmt.Errorf("Failed to Marshal")
			return err
		}
		PrintYaml(yamlBytes, typePath, t.Name)
	}
	return nil
}

func PrintYaml(yamlBytes []byte, outputPath string, fileName string) error {
	yamlFile, err := os.Create(outputPath + "/" + fileName + ".yml")
	if err != nil {
		fmt.Errorf("Failed to create file: ", fileName)
		return err
	}

	yamlFile.WriteString("#YamlMIME: GoLangPkg\n")
	yamlFile.Write(yamlBytes)
	yamlFile.Close()
	return nil
}

// PrintPosition
// some Sample codes to print the source code position, for constants & functions & methods
func PrintPosition(info *godoc.PageInfo) {
	// print position info for constant
	fmt.Println("---Constant source info example--------------------")
	if len(info.PDoc.Consts) > 0 {
		constant := info.PDoc.Consts[0] // use the first one for example
		position := constant.Decl.Pos()
		fs := info.FSet.Position(position)
		fmt.Println(constant.Names)
		fmt.Println(fs.Filename, fs.Line)
	}
	fmt.Println("---Function source info example--------------------")
	// print position info for function
	if len(info.PDoc.Funcs) > 0 {
		function := info.PDoc.Funcs[0]
		position := function.Decl.Pos()
		fs := info.FSet.Position(position)
		fmt.Println(function.Name)
		fmt.Println(fs.Filename, fs.Line)
	}
	fmt.Println("---Type source info example--------------------")
	// print position info for type
	if len(info.PDoc.Types) > 0 {
		typea := info.PDoc.Types[0]
		position := typea.Decl.Pos()
		fs := info.FSet.Position(position)
		fmt.Println(typea.Name)
		fmt.Println(fs.Filename, fs.Line)
		if len(typea.Methods) > 0 {
			fmt.Println("---Method source info example--------------------")
			methoda := typea.Methods[0]
			position := methoda.Decl.Pos()
			fs := info.FSet.Position(position)
			fmt.Println(methoda.Name)
			fmt.Println(fs.Filename, fs.Line)
		}
	}
	fmt.Println("-----------------------")
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
	Uid         string                `json:"uid"`
	Name        string                `json:"name"`
	IsMain      bool                  `json:"ismain"`
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	ImportPath  string                `json:"importPath"`
	Dir         string                `json:"dir"`
	Consts      []DocsValue           `json:"consts"`
	Types       []string              `json:"types"`
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

type SourcePosition struct {
	Repo string `json:"repo"`
	Branch string `json:"branch"`
	File string `json:"file"`
	Line int    `json:"line"`
}

type DocsValue struct {
	Names       []string       `json:"names"`
	Summary     string         `json:"summary"`
	Description string         `json:"description"`
	Code        string         `json:"code"`
	Source      SourcePosition `json:"source"`
}

type DocsType struct {
	Uid         string `json:"uid"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Code        string `json:"code"`

	Consts  []DocsValue    `json:"consts"`
	Vars    []DocsValue    `json:"vars"`
	Funcs   []DocsFunc     `json:"funcs"`
	Methods []DocsFunc     `json:"methods"`
	Source  SourcePosition `json:"source"`
}

type DocsFunc struct {
	Uid         string         `json:"uid"`
	Name        string         `json:"name"`
	Summary     string         `json:"summary"`
	Description string         `json:"description"`
	Code        string         `json:"code"`
	Source      SourcePosition `json:"source"`
}

func ToDocfx(info *godoc.PageInfo, packageName string, packagePrefix string) (*DocsPackage, []DocsType) {
	pkg := &DocsPackage{
		IsMain: info.IsMain,
		Dir:    info.Dirname,
		Notes:  toDocsNotes(info.Notes),
		Dirs:   toDocsDirs(info.Dirs),
	}
	var types []DocsType
	if info.PDoc != nil {
		if packagePrefix == "" {
			pkg.Uid = packageName
		} else {
			pkg.Uid = packagePrefix + "." + packageName
		}
		pkg.Name = packageName
		pkg.ImportPath = info.PDoc.ImportPath
		pkg.Summary = summary(info.PDoc.Doc)
		pkg.Description = description(info.PDoc.Doc)
		pkg.Examples = toDocsExamples(info.Examples, info)
		pkg.Consts = toDocsValues(info.PDoc.Consts, info)
		pkg.Vars = toDocsValues(info.PDoc.Vars, info)
		pkg.Funcs = toDocsFuncs(info.PDoc.Funcs, info, pkg.Uid)
		types, pkg.Types = toDocsTypes(info.PDoc.Types, info, pkg.Uid)
	}
	return pkg, types
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

func toDocsTypes(types []*doc.Type, info *godoc.PageInfo, parentUid string) ([]DocsType, []string) {
	arr := make([]DocsType, len(types))
	uidArr := make([]string, len(types))
	for i, t := range types {
		position := t.Decl.Pos()
		fs := info.FSet.Position(position)
		uid := parentUid + "." + t.Name
		arr[i] = DocsType{
			Uid:         uid,
			Name:        t.Name,
			Summary:     summary(t.Doc),
			Description: description(t.Doc),
			Code:        nodeFunc(info, t.Decl),
			Consts:      toDocsValues(t.Consts, info),
			Vars:        toDocsValues(t.Vars, info),
			Funcs:       toDocsFuncs(t.Funcs, info, uid),
			Methods:     toDocsFuncs(t.Methods, info, uid),
			Source:      SourcePosition{Repo: SourceRepo, Branch: SourceBranch, File: fs.Filename, Line: fs.Line},
		}
		uidArr[i] = uid
	}
	return arr, uidArr
}

func toDocsFuncs(funcs []*doc.Func, info *godoc.PageInfo, parentUid string) []DocsFunc {
	arr := make([]DocsFunc, len(funcs))
	for i, f := range funcs {
		position := f.Decl.Pos()
		fs := info.FSet.Position(position)
		arr[i] = DocsFunc{
			Uid:         parentUid + "." + f.Name,
			Name:        f.Name,
			Summary:     summary(f.Doc),
			Description: description(f.Doc),
			Code:        nodeFunc(info, f.Decl),
			Source:      SourcePosition{Repo: SourceRepo, Branch: SourceBranch, File: fs.Filename, Line: fs.Line},
		}
	}
	return arr
}

func toDocsValues(values []*doc.Value, info *godoc.PageInfo) []DocsValue {
	arr := make([]DocsValue, len(values))
	for i, v := range values {
		position := v.Decl.Pos()
		fs := info.FSet.Position(position)
		arr[i] = DocsValue{
			Names:       v.Names,
			Summary:     summary(v.Doc),
			Description: description(v.Doc),
			Code:        nodeFunc(info, v.Decl),
			Source:      SourcePosition{Repo: SourceRepo, Branch: SourceBranch, File: fs.Filename, Line: fs.Line},
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
