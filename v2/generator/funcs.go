// Copyright (c) 2020 Mercari, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package generator

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/kenshaw/snaker"

	"go.mercari.io/yo/v2/internal"
	"go.mercari.io/yo/v2/models"
)

// newTemplateFuncs returns a set of template funcs bound to the supplied args.
func (g *Generator) newTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"filterFields": g.filterFields,
		"shortName":    g.shortName,
		"nullcheck":    g.nullcheck,

		"hasColumn":         g.hasColumn,
		"columnNames":       g.columnNames,
		"columnNamesQuery":  g.columnNamesQuery,
		"columnPrefixNames": g.columnPrefixNames,

		"hasField":   g.hasField,
		"fieldNames": g.fieldNames,

		"goParam":         g.goParam,
		"goEncodedParam":  g.goEncodedParam,
		"goParams":        g.goParams,
		"goParamDefs":     g.goParamDefs,
		"goEncodedParams": g.goEncodedParams,

		"escape":    g.escape,
		"toLower":   g.toLower,
		"pluralize": g.pluralize,

		"newPackage":     newPackage,
		"presetPackages": presetPackages,
	}
}

func ignoreFromMultiTypes(ignoreNames []interface{}) map[string]bool {
	ignore := map[string]bool{}
	for _, f := range ignoreNames {
		switch k := f.(type) {
		case string:
			ignore[k] = true

		case []*models.Field:
			for _, f := range k {
				ignore[f.Name] = true
			}
		}
	}

	return ignore
}

func (g *Generator) filterFields(fields []*models.Field, ignoreNames ...interface{}) []*models.Field {
	ignore := ignoreFromMultiTypes(ignoreNames)

	filtered := make([]*models.Field, 0, len(fields))
	for _, f := range fields {
		if ignore[f.Name] {
			continue
		}
		filtered = append(filtered, f)
	}

	return filtered
}

// columnNames creates a list of the column names found in fields.
//
// When escaped is true, if the column name is a reserved word, it's escaped with backquotes.
//
// Used to present a comma separated list of column names, that can be used in
// a SELECT, or UPDATE, or other SQL clause requiring an list of identifiers
// (ie, "field_1, field_2, field_3, ...").
func (g *Generator) columnNames(fields []*models.Field) string {
	str := ""
	i := 0
	for _, f := range fields {
		if i != 0 {
			str = str + ", "
		}
		str = str + internal.EscapeColumnName(f.ColumnName)
		i++
	}
	return str
}

// columnNamesQuery creates a list of the column names in fields as a query and
// joined by sep, excluding any models.Field with Name contained in ignoreNames.
//
// Used to create a list of column names in a WHERE clause (ie, "field_1 = $1
// AND field_2 = $2 AND ...") or in an UPDATE clause (ie, "field = $1, field =
// $2, ...").
func (g *Generator) columnNamesQuery(fields []*models.Field, sep string) string {
	str := ""
	i := 0
	for _, f := range fields {
		if i != 0 {
			str = str + sep
		}
		str = str + internal.EscapeColumnName(f.ColumnName) + " = " + g.loader.NthParam(i)
		i++
	}

	return str
}

// shortName generates a safe Go identifier for typ. typ is first checked
// against ShortNameTypeMap, and if not found, then the value is
// calculated and stored in the ShortNameTypeMap for future use.
//
// A shortname is the concatenation of the lowercase of the first character in
// the words comprising the name. For example, "MyCustomName" will have
// the shortname of "mcn".
//
// If a generated shortname conflicts with a Go reserved name, then the
// corresponding value in goReservedNames map will be used.
//
// Generated shortnames that have conflicts with any scopeConflicts member will
// have nameConflictSuffix appended.
//
// Note: recognized types for scopeConflicts are string, []*models.Field.
func (g *Generator) shortName(typ string, scopeConflicts ...interface{}) string {
	var v string
	var ok bool

	// check short name map
	if v, ok = ShortNameTypeMap[typ]; !ok {
		// calc the short name
		u := []string{}
		for _, s := range strings.Split(strings.ToLower(snaker.CamelToSnake(typ)), "_") {
			if len(s) > 0 && s != "id" {
				u = append(u, s[:1])
			}
		}
		v = strings.Join(u, "")

		// check go reserved names
		if n, ok := goReservedNames[v]; ok {
			v = n
		}

		// store back to short name map
		ShortNameTypeMap[typ] = v
	}

	// add scopeConflicts to conflicts
	for _, c := range scopeConflicts {
		switch k := c.(type) {
		case string:
			if k == v {
				v = v + g.nameConflictSuffix
			}

		case []*models.Field:
			for _, f := range k {
				if f.Name == v {
					v = v + g.nameConflictSuffix
				}
			}

		default:
			panic("shortName: supported type")
		}
	}

	// append suffix if conflict exists
	if _, ok := ConflictedShortNames[v]; ok {
		v = v + g.nameConflictSuffix
	}

	return v
}

// columnPrefixNames creates a list of the column names found in fields with the
// supplied prefix.
//
// Used to present a comma separated list of column names with a prefix. Used in
// a SELECT, or UPDATE (ie, "t.field_1, t.field_2, t.field_3, ...").
func (g *Generator) columnPrefixNames(fields []*models.Field, prefix string) string {
	str := ""
	i := 0
	for _, f := range fields {
		if i != 0 {
			str = str + ", "
		}
		str = str + prefix + "." + internal.EscapeColumnName(f.ColumnName)
		i++
	}

	return str
}

// fieldNames creates a list of field names from fields of the adding the
// provided prefix, and excluding any models.Field with Name contained in ignoreNames.
//
// Used to present a comma separated list of field names, ie in a Go statement
// (ie, "t.Field1, t.Field2, t.Field3 ...")
func (g *Generator) fieldNames(fields []*models.Field, prefix string) string {
	str := ""
	i := 0
	for _, f := range fields {
		if i != 0 {
			str = str + ", "
		}

		str = str + prefix + "." + f.Name
		i++
	}

	return str
}

// goReservedNames is a map of go reserved names to "safe" names.
var goReservedNames = map[string]string{
	"break":       "brk",
	"case":        "cs",
	"chan":        "chn",
	"const":       "cnst",
	"continue":    "cnt",
	"default":     "def",
	"defer":       "dfr",
	"else":        "els",
	"fallthrough": "flthrough",
	"for":         "fr",
	"func":        "fn",
	"go":          "goVal",
	"goto":        "gt",
	"if":          "ifVal",
	"import":      "imp",
	"interface":   "iface",
	"map":         "mp",
	"package":     "pkg",
	"range":       "rnge",
	"return":      "ret",
	"select":      "slct",
	"struct":      "strct",
	"switch":      "swtch",
	"type":        "typ",
	"var":         "vr",

	// go types
	"error":      "e",
	"bool":       "b",
	"string":     "str",
	"byte":       "byt",
	"rune":       "r",
	"uintptr":    "uptr",
	"int":        "i",
	"int8":       "i8",
	"int16":      "i16",
	"int32":      "i32",
	"int64":      "i64",
	"uint":       "u",
	"uint8":      "u8",
	"uint16":     "u16",
	"uint32":     "u32",
	"uint64":     "u64",
	"float32":    "z",
	"float64":    "f",
	"complex64":  "c",
	"complex128": "c128",
}

// goParam make the first word of name to lowercase
func (g *Generator) goParam(name string) string {
	ns := strings.Split(snaker.CamelToSnake(name), "_")
	name = strings.ToLower(ns[0]) + name[len(ns[0]):]

	// check go reserved names
	if r, ok := goReservedNames[strings.ToLower(name)]; ok {
		name = r
	}

	return name
}

// goEncodedParam make the first word of name to lowercase
func (g *Generator) goEncodedParam(name string) string {
	return fmt.Sprintf("yoEncode(%s)", g.goParam(name))
}

// goParams converts a list of fields into their named Go parameters,
// skipping any models.Field with Name contained in ignoreNames. addType will cause
// the go Type to be added after each variable name. addPrefix will cause the
// returned string to be prefixed with ", " if the generated string is not
// empty.
//
// Any field name encountered will be checked against goReservedNames, and will
// have its name substituted by its corresponding looked up value.
//
// Used to present a comma separated list of Go variable names for use with as
// either a Go func parameter list, or in a call to another Go func.
// (ie, ", a, b, c, ..." or ", a T1, b T2, c T3, ...").
func (g *Generator) goParams(fields []*models.Field, addPrefix bool) string {
	return g.handleGoParams(fields, addPrefix, func(_ *models.Field, param string) string {
		return param
	})
}

func (g *Generator) goEncodedParams(fields []*models.Field, addPrefix bool) string {
	return g.handleGoParams(fields, addPrefix, func(_ *models.Field, param string) string {
		return fmt.Sprintf("yoEncode(%s)", param)
	})
}

func (g *Generator) goParamDefs(registry *PackageRegistry, fields []*models.Field, addPrefix bool) string {
	return g.handleGoParams(fields, addPrefix, func(field *models.Field, param string) string {
		return fmt.Sprintf("%s %s", param, field.Type.GetType(registry))
	})
}

func (g *Generator) handleGoParams(fields []*models.Field, addPrefix bool, handle func(field *models.Field, param string) string) string {
	i := 0
	vals := []string{}
	for _, f := range fields {
		s := "v" + strconv.Itoa(i)
		if len(f.Name) > 0 {
			s = g.goParam(f.Name)
		}

		// add to vals
		vals = append(vals, handle(f, s))

		i++
	}

	// concat generated values
	str := strings.Join(vals, ", ")
	if addPrefix && str != "" {
		return ", " + str
	}

	return str
}

// hascolumn takes a list of fields and determines if field with the specified
// column name is in the list.
func (g *Generator) hasColumn(fields []*models.Field, name string) bool {
	for _, f := range fields {
		if f.ColumnName == name {
			return true
		}
	}

	return false
}

// hasfield takes a list of fields and determines if field with the specified
// field name is in the list.
func (g *Generator) hasField(fields []*models.Field, name string) bool {
	for _, f := range fields {
		if f.Name == name {
			return true
		}
	}

	return false
}

// nullcheck generates a code to check the field value is null.
func (g *Generator) nullcheck(field *models.Field) string {
	paramName := g.goParam(field.Name)

	if t, ok := field.Type.(models.PlainFieldType); ok {
		if t.Pkg == models.GoSpannerPackage {
			switch t.Type {
			case "NullInt64",
				"NullString",
				"NullFloat64",
				"NullBool",
				"NullTime",
				"NullDate":
				return fmt.Sprintf("%s.IsNull()", paramName)
			}
		}
	}

	return fmt.Sprintf("yo, ok := %s.(yoIsNull); ok && yo.IsNull()", paramName)
}

// escaped returns the ColumnName of col. It is escaped for query.
func (g *Generator) escape(col string) string {
	return internal.EscapeColumnName(col)
}

// toLower converts s to lower case.
func (g *Generator) toLower(s string) string {
	return strings.ToLower(s)
}

// pluralize converts s to plural.
func (g *Generator) pluralize(s string) string {
	return g.inflector.Pluralize(s)
}

// newPackage returns a new models.Package object
func newPackage(path, name string) models.Package {
	if name == "" {
		return models.Package{
			Path: path,
			Name: filepath.Base(path),
		}
	}

	return models.Package{
		Path: path,
		Name: name,
	}
}

// presetPackages returns preset Go packages
func presetPackages() map[string]models.Package {
	return map[string]models.Package{
		"big":     models.MathBigPackage,
		"context": models.ContextPackage,
		"errors":  models.ErrorsPackage,
		"fmt":     models.FmtPackage,
		"strconv": models.StrconvPackage,
		"strings": models.StringsPackage,
		"time":    models.TimePackage,

		"apiIterator":               models.APIIteratorPackage,
		"googleApisGaxGoV2ApiError": models.GoogleApisGaxGoV2ApiErrorPackage,
		"goSpanner":                 models.GoSpannerPackage,
		"goCivil":                   models.GoCivilPackage,
		"gRPCCodes":                 models.GRPCCodesPackage,
		"gRPCStatus":                models.GRPCStatusPackage,
	}
}
