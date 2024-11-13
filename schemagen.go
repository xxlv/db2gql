package main

import (
	"fmt"
)

type SchemaGenerator struct {
	Name       string
	RawColumns []Column
}

func (sg *SchemaGenerator) genObject() string {
	tf := TypeFormatter{
		Kind: "type",
		Name: sg.Name,
		NameTypes: asNameTypeFormatter(sg.RawColumns, func(c Column) bool {
			return true
		}),
	}
	return tf.Format()
}

var needNotRepeatedField []string = []string{
	"updated_at",
	"created_at",
}

func asNameTypeFormatter(cols []Column, skipFilter func(Column) bool) []*NameTypeFormatter {
	result := []*NameTypeFormatter{}
	if len(cols) <= 0 {
		return result
	}
	addTypes := map[string]any{}
	for _, v := range cols {
		if !skipFilter(v) {
			continue
		}
		name := v.Name
		// if this col's name already loaded , use $table.name instead
		if _, ok := addTypes[name]; ok {
			needSkip := false
			// if in this map ,just skip
			for _, v := range needNotRepeatedField {
				if v == name {
					needSkip = true
					break
				}
			}
			if needSkip {
				continue
			}
			name = asCamStyle(v.Table) + asCamStyle(v.Name)
		}
		typ := mapMySQLTypeToGraphQL(v.Type, v.Null)
		addTypes[name] = nil

		// hack
		if name == "id" {
			typ = "ID"
			if v.Null == "YES" {
				typ = typ + "!"
			}
		}
		result = append(result, &NameTypeFormatter{
			Name:    name,
			Type:    typ,
			Comment: v.Comment,
			Args:    []*ArgsFormatter{},
		})
	}

	return result
}

// buildin payload obejct
func (sg *SchemaGenerator) genPayload() string {
	tf := &TypeFormatter{
		Kind: "type",
		Name: getTypePayloadObject(sg.Name),
		NameTypes: []*NameTypeFormatter{
			{
				Name:    asLowCaseCamStyle(sg.Name),
				Type:    asCamStyle(sg.Name),
				Comment: "FIXME: please add comment.",
			},
			{
				Name:    asCamStyleWithoutUnderline(sg.Name) + "_UserErrors",
				Type:    fmt.Sprintf("[%s!]!", asCamStyle(sg.Name)+"UserErrors"),
				Comment: "The list of errors that occurred from executing the mutation.",
			},
		},
	}
	return tf.Format()
}

func (sg *SchemaGenerator) genUserErrors() string {
	tf := TypeFormatter{
		Kind: "type",
		Name: getTypeObject(sg.Name) + "UserErrors",
		NameTypes: []*NameTypeFormatter{
			{
				Name:    "code",
				Type:    getTypeObject(sg.Name) + "ErrorCode",
				Comment: "The error code.",
			},
			{
				Name:    "field",
				Type:    "[String]!",
				Comment: "The path to the input field that caused the error.",
			},
			{
				Name:    "message",
				Type:    "String!",
				Comment: "The error message.",
			},
		},
	}
	return tf.Format()
}

func (sg *SchemaGenerator) genUserErrorCodeEnum() string {

	tf := EnumFormatter{
		Name: getTypeObject(sg.Name) + "ErrorCode",
		ValueWithComment: map[string]string{
			"CODE": " ** 😄This is example code, PLEASE REPLEACE YOUR.**",
		},
	}
	return tf.Format()
}

func (sg *SchemaGenerator) genInput() string {
	tf := TypeFormatter{
		Kind: "input",
		Name: getTypeInputObject(sg.Name),
		NameTypes: asNameTypeFormatter(sg.RawColumns, func(c Column) bool {
			return c.Name != "id"
		}),
	}
	return tf.Format()
}

func (sg *SchemaGenerator) genQueries() string {
	qg := &QueryGenerator{Name: sg.Name, RawColumns: sg.RawColumns}
	return qg.Gen()
}

func (sg *SchemaGenerator) genMutations() string {
	mg := &MutationGenerator{Name: sg.Name, RawColumns: sg.RawColumns}
	return mg.Gen()
}

func (sg *SchemaGenerator) Gen() string {
	return sg.genQueries() + "\n" +
		sg.genMutations() + "\n" +
		sg.genObject() + "\n" +
		sg.genPayload() + "\n" +
		sg.genInput() + "\n" +
		sg.genUserErrors() + "\n" +
		sg.genUserErrorCodeEnum()
}

func getTypeObject(name string) string {
	return asCamStyle(name)
}
func getTypeInputObject(name string) string {
	return asCamStyle(name) + "Input"
}

func getTypePayloadObject(name string) string {
	return asCamStyle(name) + "Payload"
}
func getAPINameForUpdate(name string) string {
	return asLowCaseCamStyle(name) + "Update"
}
