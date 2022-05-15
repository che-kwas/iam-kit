// Package codegen is a tool to automate the creation of code init function
// and generate document.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"go/types"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

var errCodeDocPrefix = `## 错误码列表

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
`

var (
	doc    = flag.Bool("doc", false, "if true only generate error code documentation in markdown format")
	output = flag.String("output", "", "output file name; default srcdir/src_filename_gen.go")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of codegen:\n")
	fmt.Fprintf(os.Stderr, "\tcodegen                            # generate code\n")
	fmt.Fprintf(os.Stderr, "\tcodegen -doc -output <output-path> # generate document\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("codegen: ")
	flag.Usage = Usage
	flag.Parse()

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	g := Generator{}
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		dir = filepath.Dir(args[0])
	}

	g.parsePackage(args)

	if !*doc {
		g.Printf("// Code generated by \"codegen %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
		g.Printf("\n")
		g.Printf("package %s", g.pkg.name)
		g.Printf("\n")
		g.Printf("import \"github.com/che-kwas/iam-kit/errcoder\"")
		g.Printf("\n")
	}

	// Run generate for each type.
	var src []byte
	if *doc {
		g.generateDocs()
		src = g.buf.Bytes()
	} else {
		g.generate()
		// Format the output.
		src = g.format()
	}

	// Write to file.
	outputName := *output
	if outputName == "" {
		absDir, _ := filepath.Abs(dir)
		baseName := fmt.Sprintf("%s_gen.go", strings.ReplaceAll(filepath.Base(absDir), "-", "_"))
		if len(flag.Args()) == 1 {
			baseName = fmt.Sprintf(
				"%s_gen.go",
				strings.ReplaceAll(filepath.Base(strings.TrimSuffix(flag.Args()[0], ".go")), "-", "_"),
			)
		}

		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, src, 0o600)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.
}

// Printf like fmt.Printf, but add the string to g.buf.
func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg    *Package  // Package to which this file belongs.
	file   *ast.File // Parsed AST.
	values []Value   // Accumulator for constant values of that type.
}

// Package defines options for package.
type Package struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string) {
	// nolint: staticcheck
	cfg := &packages.Config{Mode: packages.LoadSyntax}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file: file,
			pkg:  g.pkg,
		}
	}
}

// generate produces the register calls for the named type.
func (g *Generator) generate() {
	values := make([]Value, 0, 100)
	// Set the state for this run of the walker.
	for _, file := range g.pkg.files {
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		return
	}

	// Generate code that will fail if the constants change value.
	g.Printf("\t// init register error codes defines in `github.com/marmotedu/errors`\n")
	g.Printf("func init() {\n")
	for _, v := range values {
		code, description := v.ParseComment()
		g.Printf("\terrcoder.Register(%s, %s, \"%s\")\n", v.name, code, description)
	}
	g.Printf("}\n")
}

// generateDocs produces error code markdown document for the named type.
func (g *Generator) generateDocs() {
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		return
	}

	tmpl, _ := template.New("doc").Parse(errCodeDocPrefix)
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, "`")

	// Generate code that will fail if the constants change value.
	g.Printf(buf.String())
	for _, v := range values {
		code, description := v.ParseComment()
		// g.Printf("\tregister(%s, %s, \"%s\")\n", v.name, code, description)
		g.Printf("| %s | %d | %s | %s |\n", v.name, v.value, code, description)
	}
	g.Printf("\n")
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")

		return g.buf.Bytes()
	}

	return src
}

// Value represents a declared constant.
type Value struct {
	comment string
	name    string // The name of the constant.
	// The value is stored as a bit pattern alone. The boolean tells us
	// whether to interpret it as an int64 or a uint64; the only place
	// this matters is when sorting.
	// Much of the time the str field is all we need; it is printed
	// by Value.String.
	value  uint64 // Will be converted to int64 when needed.
	signed bool   // Whether the constant is a signed type.
	str    string // The string representation given by the "go/constant" package.
}

func (v *Value) String() string {
	return v.str
}

// ParseComment parse comment to http code and error code description.
func (v *Value) ParseComment() (string, string) {
	reg := regexp.MustCompile(`\w\s*-\s*(\d{3})\s*:\s*([A-Z].*)\s*\.\n*`)
	if !reg.MatchString(v.comment) {
		log.Printf("constant '%s' have wrong comment format, register with 500 as default", v.name)

		return "500", "Internal server error"
	}

	groups := reg.FindStringSubmatch(v.comment)
	if len(groups) != 3 {
		return "500", "Internal server error"
	}

	return groups[1], groups[2]
}

// nolint: gocognit
// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// We only care about const declarations.
		return true
	}
	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		vspec, _ := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		if vspec.Type == nil && len(vspec.Values) > 0 {
			// "X = 1". With no type but a value. If the constant is untyped,
			// skip this vspec and reset the remembered type.
			typ = ""

			// If this is a simple type conversion, remember the type.
			// We don't mind if this is actually a call; a qualified call won't
			// be matched (that will be SelectorExpr, not Ident), and only unusual
			// situations will result in a function call that appears to be
			// a type conversion.
			ce, ok := vspec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			id, ok := ce.Fun.(*ast.Ident)
			if !ok {
				continue
			}
			typ = id.Name
		}
		if vspec.Type != nil {
			// "X T". We have a type. Remember it.
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != "int" {
			// This is not the type we're looking for.
			continue
		}
		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}
			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the name, find its
			// types.Const, and extract its value.
			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				log.Fatalf("can't handle non-integer constant type %s", typ)
			}
			value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
			if value.Kind() != constant.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
			}
			i64, isInt := constant.Int64Val(value)
			u64, isUint := constant.Uint64Val(value)
			if !isInt && !isUint {
				log.Fatalf("internal error: value of %s is not an integer: %s", name, value.String())
			}
			if !isInt {
				u64 = uint64(i64)
			}
			v := Value{
				name:   name.Name,
				value:  u64,
				signed: info&types.IsUnsigned == 0,
				str:    value.String(),
			}
			if vspec.Doc != nil && vspec.Doc.Text() != "" {
				v.comment = vspec.Doc.Text()
			} else if c := vspec.Comment; c != nil && len(c.List) == 1 {
				v.comment = c.Text()
			}

			f.values = append(f.values, v)
		}
	}

	return false
}
