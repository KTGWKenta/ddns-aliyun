package cueAST

import (
	"reflect"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/token"
	"go.uber.org/zap"

	"github.com/kentalee/errors"

	"github.com/kentalee/ddns/common/defines"
)

func FillNestedField(r *cue.Runtime, source []ast.Decl, value interface{}, path ...string) (newSource []ast.Decl, err error) {
	var valueExpr ast.Expr
	switch _value := value.(type) {
	case string:
		valueExpr = ast.NewString(_value)
	case bool:
		valueExpr = ast.NewBool(_value)
	case nil:
		valueExpr = ast.NewNull()
	case ast.Expr:
		valueExpr = _value
	default:
		var instance *cue.Instance
		if instance, err = r.CompileFile(&ast.File{}); err != nil {
			return nil, errors.Note(err)
		}
		if instance, err = instance.Fill(value, "_temp_val"); err != nil {
			return nil, errors.Note(err)
		}
		var ok bool
		if valueExpr, ok = instance.Lookup("_temp_val").Syntax(cue.Final()).(ast.Expr); !ok {
			return source, errors.Note(defines.ErrInvalidDataType,
				zap.String(defines.ErrKeyField, strings.Join(path, defines.FieldPathSep)),
				zap.Stringer(defines.ErrKeyType, reflect.TypeOf(value)))
		}
	}
	var target *ast.Field
	if newSource, target, err = MakeNestedField(source, path...); err == nil {
		target.Value = valueExpr
	}
	return newSource, err
}

func MakeNestedField(source []ast.Decl, path ...string) (newSource []ast.Decl, underlying *ast.Field, err error) {
	var last, rest = FindNestedField(source, nil, path...)
	if len(rest) == 0 {
		return source, last, nil
	}
	var root *ast.Field
	root, underlying = NestedField(rest...)
	if last == nil {
		newSource = append(source, root)
	} else {
		newSource = source
		if structVal, ok := last.Value.(*ast.StructLit); ok {
			structVal.Elts = append(structVal.Elts, root)
		} else {
			return newSource, nil, errors.Note(defines.ErrInvalidDataType)
		}
	}
	return newSource, underlying, nil
}

func FindNestedField(source []ast.Decl, parent []string, path ...string) (last *ast.Field, rest []string) {
	rest = path
	for _, decl := range source {
		var _last *ast.Field
		var _rest []string
		switch _decl := decl.(type) {
		case *ast.Field:
			switch _label := _decl.Label.(type) {
			case *ast.Ident:
				if _label.Name != path[0] {
					continue
				}
			case *ast.BasicLit:
				if _label.Kind != token.STRING || _label.Value != ast.NewString(path[0]).Value {
					continue
				}
			default:
				continue
			}
			if len(path) == 1 { // last part, matched!
				_last, _rest = _decl, nil // current
			} else if _decl.Value == nil {
				_last, _rest, _decl.Value = _decl, path[1:], ast.NewStruct() // current
			} else if subStruct, ok := _decl.Value.(*ast.StructLit); ok {
				_last, _rest = FindNestedField(subStruct.Elts, append(parent, path[0]), path[1:]...) // sub
				if _last == nil {
					_last, _rest = _decl, path[1:] // current
				}
			} else {
				continue
			}
		case *ast.StructLit:
			_last, _rest = FindNestedField(_decl.Elts, parent, path...)
		}
		if _last != nil && len(_rest) < len(rest) {
			last, rest = _last, _rest
		}
	}
	return last, rest
}

func NestedField(path ...string) (root, underlying *ast.Field) {
	root = &ast.Field{Label: ast.NewString(path[0])}
	underlying = root
	for _, pathPart := range path[1:] {
		subField := &ast.Field{Label: ast.NewString(pathPart)}
		underlying.Value = ast.NewStruct(subField)
		underlying = subField
	}
	return root, underlying
}

var ErrCantParseCUEInstance = errors.New("c0f236d900020001", "can't parse cue instance as ast file node")

func MergeRuntimeInstance(decl []ast.Decl, instance *cue.Instance) (file *ast.File, err error) {
	file = &ast.File{Filename: "", Decls: decl}
	if instance != nil {
		if eNodes, ok := instance.Lookup().Syntax(cue.Concrete(true), cue.Final()).(ast.Expr); ok {
			var instanceFile *ast.File
			if instanceFile, err = astutil.ToFile(eNodes); err != nil {
				return nil, errors.Note(err)
			} else {
				file.Decls = append(file.Decls, instanceFile.Decls...)
			}
		} else {
			return nil, errors.Note(ErrCantParseCUEInstance)
		}
	}
	return file, nil
}
