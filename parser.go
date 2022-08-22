package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type field struct {
	n string
	t string
	k string
	v string
}

type ParsedData struct {
	fields map[string]map[string]field
	pack   string
}

func NewParsedData() *ParsedData {
	return &ParsedData{
		fields: make(map[string]map[string]field),
		pack:   "",
	}
}

func (p *ParsedData) AddStructFiled(structName, fieldName, fieldType, k, v string) {
	if fileds, ok := p.fields[structName]; ok {
		fileds[fieldName] = field{n: fieldName, t: fieldType, k: k, v: v}
		return
	}

	p.fields[structName] = make(map[string]field)
	p.fields[structName][fieldName] = field{n: fieldName, t: fieldType, k: k, v: v}

}

func (p *ParsedData) SetPack(pack string) {
	p.pack = pack
}

type UnitParser struct {
	data *ParsedData
}

func NewParser() *UnitParser {
	return &UnitParser{
		data: nil,
	}
}

func (p *UnitParser) Pop() *ParsedData {
	return p.data
}

func (p *UnitParser) Parse(fileName string) error {

	p.data = NewParsedData()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}

	p.data.SetPack(f.Name.Name)

	ast.Inspect(f, func(n ast.Node) bool {
		t, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if t.Type == nil {
			return true
		}

		x, ok := t.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range x.Fields.List {

			mp, ok := field.Type.(*ast.MapType)
			if ok {
				p.data.AddStructFiled(t.Name.Name, field.Names[0].Name, "map", fmt.Sprintf("%s", mp.Key), fmt.Sprintf("%s", mp.Value))
			}

			s, ok := field.Type.(*ast.ArrayType)
			if ok {
				p.data.AddStructFiled(t.Name.Name, field.Names[0].Name, "slice", "", fmt.Sprintf("%s", s.Elt))
			}

			_, ok = field.Type.(*ast.Ident)
			if ok {
				p.data.AddStructFiled(t.Name.Name, field.Names[0].Name, "ident", "", fmt.Sprintf("%s", field.Type))
			}
		}
		return true
	})
	return nil
}
