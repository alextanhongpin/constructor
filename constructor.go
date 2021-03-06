package constructor

import (
	"flag"
	"fmt"
	"go/types"
	"os"
	"strings"

	"github.com/alextanhongpin/constructor/internal/loader"
	"github.com/alextanhongpin/constructor/internal/set"
	"github.com/alextanhongpin/pkg/stringcase"
	jen "github.com/dave/jennifer/jen"
)

const GeneratorName = "constructor"

func Run() error {
	includep := flag.String("include", "*", "the fields to include. Use comma to specify multiple fields. * means all fields")
	excludep := flag.String("exclude", "", "the fields to exclude. Use comma to specify multiple fields.")
	inp := flag.String("in", os.Getenv("GOFILE"), "the input file, defaults to the file with the go:generate comment")
	outp := flag.String("out", "", "the output file name, defaults to the same location as the input file, and the name generated from the struct type if none is provided")
	typep := flag.String("type", "", "the type of the struct to generate.")
	flag.Parse()

	in := loader.RelativePath(*inp)
	pkg := loader.LoadPackage(in)
	obj := pkg.Types.Scope().Lookup(*typep)
	if obj == nil {
		return fmt.Errorf("type %s not found", *typep)
	}
	str, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("type %s is not a struct", *typep)
	}

	fields := set.NewString()
	for i := 0; i < str.NumFields(); i++ {
		fields.Add(str.Field(i).Name())
	}

	var includes *set.String
	switch *includep {
	case "*":
		includes = set.NewString(fields.List()...)
	default:
		includes = set.NewString(strings.Split(*includep, ",")...)
		difference := includes.Difference(fields)
		if difference.Len() > 0 {
			return fmt.Errorf("includes contains unknown fields: %s", strings.Join(difference.List(), ", "))
		}
		includes = includes.Intersect(fields)
	}
	if *excludep != "" {
		excludes := set.NewString(strings.Split(*excludep, ",")...)
		difference := excludes.Difference(fields)
		if difference.Len() > 0 {
			return fmt.Errorf("excludes contains unknown fields: %s", strings.Join(difference.List(), ", "))
		}
		includes = includes.Difference(excludes)
	}
	included := includes.Map()

	var params []jen.Code
	values := make(jen.Dict)
	for i := 0; i < str.NumFields(); i++ {
		f := str.Field(i)
		name := f.Name()
		if !included[name] {
			continue
		}
		params = append(params, newVar(f))
		values[jen.Id(name)] = jen.Id(stringcase.CamelCase(name))
	}

	entity := *typep
	newEntity := "New" + entity
	newFile := *outp
	if newFile == "" {
		newFile = stringcase.KebabCase(entity) + "_gen.go"
	}

	f := jen.NewFilePathName(pkg.PkgPath, pkg.Name)
	f.HeaderComment(fmt.Sprintf("Code generated by %s, DO NOT EDIT.", GeneratorName))
	f.Func().Id(newEntity).Params(params...).Op("*").Id(entity).Block(
		jen.Return(jen.Op("&").Id(entity).Values(values)),
	)

	if err := f.Save(newFile); err != nil {
		return err
	}

	return nil
}

func newVar(t *types.Var) *jen.Statement {
	v := newCodeVisitor(jen.Id(stringcase.CamelCase(t.Name())))
	_ = Walk(v, t.Type())
	return v.stmt
}

type codeVisitor struct {
	stmt *jen.Statement
}

func newCodeVisitor(stmt *jen.Statement) *codeVisitor {
	return &codeVisitor{
		stmt: stmt,
	}
}

func (v *codeVisitor) Visit(t types.Type) bool {
	switch u := t.(type) {
	case *types.Pointer:
		v.stmt = v.stmt.Op("*")
	case *types.Slice:
		v.stmt = v.stmt.Index()
	case *types.Array:
		v.stmt = v.stmt.Index(jen.Lit(u.Len()))
	case *types.Map:
		k := newCodeVisitor(jen.Null())
		_ = Walk(k, u.Key())
		v.stmt = v.stmt.Map(k.stmt)
	case *types.Named:
		o := u.Obj()
		p := o.Pkg()
		v.stmt = v.stmt.Qual(p.Path(), o.Name())
		return false
	default:
		v.stmt = v.stmt.Id(u.String())
	}
	return true
}

type Visitor interface {
	Visit(t types.Type) bool
}

func Walk(visitor Visitor, t types.Type) bool {
	if !visitor.Visit(t) {
		return false
	}

	switch u := t.(type) {
	case *types.Pointer:
		return Walk(visitor, u.Elem())
	case *types.Named:
		return Walk(visitor, u.Underlying())
	case *types.Slice:
		return Walk(visitor, u.Elem())
	case *types.Array:
		return Walk(visitor, u.Elem())
	case *types.Map:
		return Walk(visitor, u.Elem())
	default:
		return types.IdenticalIgnoreTags(t, u)
	}
}
