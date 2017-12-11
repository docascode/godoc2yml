# godoc2yml
This subrepository holds a POC for generating YAML format of GoLang API DOC, leveraging the godoc functions.


### Current Functions:
Accepts a package name and the directory containing its source codes
Generates a YAML file containing the entity(Variables, Constants, Functions, Types) and its corresponding document


### Current Gaps:
1. All documents are generated in a single YAML file <br/>
	further requirement of seperating documents by types should be supported <br/>

2. Current ast.Node for golang entity(Variables, Constants, Functions, Types) is directly printed via go/printer <br/>
	the format is not friendly(contains \n,\t .etc) and this contains source code details <br/>
	Format refinement, source code hiding, document seperation should be supported <br/>


### Other:
If finer manipulation of golang entity needed, kind of AST walker may need to be implemented.
In that case, can refer to https://github.com/nirasan/ast-walker/blob/master/lib/walk.go 