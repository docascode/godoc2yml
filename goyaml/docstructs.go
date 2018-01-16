package goyaml

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
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	File   string `json:"file"`
	Line   int    `json:"line"`
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
