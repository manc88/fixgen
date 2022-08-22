package main

import (
	"fmt"
	"strings"
)

const (
	_CONSTRUCT = `
func ({r} *{receiver}) {methodName}({p} {ptype}) *{receiver}{
	{r}.{fieldname} = {p}
	return {r}
}`

	_CONSTRUCTOR = `
func New{Ureceiver}() *{receiver}{
	return &{receiver}{}
}`

	_ADD_IN_SLICE = `
func ({r} *{receiver}) {methodName}({p} {ptype}) *{receiver}{
	{r}.{fieldname} = append({r}.{fieldname},{p})
	return {r}
}`

	_ADD_IN_MAP = `
func ({r} *{receiver}) {methodName}({k} {ktype},{v} {vtype}) *{receiver}{
	if {r}.{fieldname} == nil {
		{r}.{fieldname} = make(map[{ktype}]{vtype})
	}
	{r}.{fieldname}[{k}] = {v}
	return {r}
}`
)

type GeneratedData struct {
	fields map[string][]string
	pack   string
}

type UnitGenerator struct {
	data *ParsedData
}

func NewGenerator() *UnitGenerator {
	return &UnitGenerator{
		data: nil,
	}
}

func (u *UnitGenerator) Load(d *ParsedData) {
	u.data = d
}

func (u *UnitGenerator) Generate() (GeneratedData, error) {

	gd := GeneratedData{}
	if u.data == nil {
		return gd, fmt.Errorf("No data to parse")
	}

	gd.fields = make(map[string][]string)
	gd.pack = u.data.pack

	for structName, fields := range u.data.fields {
		params := make(map[string]string, len(fields)+2)
		params["receiver"] = structName
		params["Ureceiver"] = strings.ToUpper(structName[:1]) + structName[1:]
		params["r"] = strings.ToLower(structName[:1])
		for _, field := range fields {
			templ := ""
			switch field.t {
			case "ident":
				params["p"] = field.n[:1]
				params["ptype"] = field.v
				params["fieldname"] = field.n
				params["methodName"] = "With" + strings.ToUpper(field.n[:1]) + field.n[1:]
				templ = _CONSTRUCT
			case "slice":
				params["p"] = field.n[:1]
				params["ptype"] = field.v
				params["fieldname"] = field.n
				params["methodName"] = "AddIn" + strings.ToUpper(field.n[:1]) + field.n[1:]
				templ = _ADD_IN_SLICE
			case "map":
				params["k"] = "k"
				params["ktype"] = field.k
				params["v"] = "v"
				params["vtype"] = field.v
				params["fieldname"] = field.n
				params["methodName"] = "AddIn" + strings.ToUpper(field.n[:1]) + field.n[1:]
				templ = _ADD_IN_MAP
			default:
				fmt.Println("unimplemented type ", field.t)
				continue
			}

			gd.fields[structName] = append(gd.fields[structName], nPrint(templ, params))
		}

		gd.fields[structName] = append(gd.fields[structName], nPrint(_CONSTRUCTOR, params))

	}

	return gd, nil
}

func nPrint(format string, params map[string]string) string {
	for key, val := range params {
		format = strings.Replace(format, "{"+key+"}", val, -1)
	}
	return format
}
