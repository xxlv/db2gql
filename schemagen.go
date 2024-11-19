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
				Name:    asCamStyleWithoutUnderline(sg.Name) + "UserErrors",
				Type:    fmt.Sprintf("[%s!]!", asCamStyle(sg.Name)+"UserErrors"),
				Comment: "The list of errors that occurred from executing the mutation.",
			},
		},
	}
	return tf.Format()
}

func (sg *SchemaGenerator) genDeletePayload() string {
	tf := &TypeFormatter{
		Kind: "type",
		Name: getDeleteTypePayloadObject(sg.Name),
		NameTypes: []*NameTypeFormatter{
			{
				Name:    "deleteId",
				Type:    "ID",
				Comment: "The globally-unique ID for the deleted cart transform.",
			},
			{
				Name:    asCamStyleWithoutUnderline(sg.Name) + "DeleteUserErrors",
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

func (sg *SchemaGenerator) genDeleteUserErrors() string {
	tf := TypeFormatter{
		Kind: "type",
		Name: getTypeObject(sg.Name) + "DeleteUserErrors",
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

func (sg *SchemaGenerator) genConnection() string {
	typeName := getTypeObject(sg.Name)
	tf := TypeFormatter{
		Kind: "type",
		Name: typeName + "Connection",
		NameTypes: []*NameTypeFormatter{
			{
				Name:    "edges",
				Type:    fmt.Sprintf("[%sEdge!]!", typeName),
				Comment: "A list of edges.",
			},
			{
				Name:    "pageInfo",
				Type:    "PageInfo!",
				Comment: "An object that's used to retrieve cursor information about the current page.",
			},
		},
	}
	return tf.Format()
}

func (sg *SchemaGenerator) genEdge() string {
	typeName := getTypeObject(sg.Name)

	tf := TypeFormatter{
		Kind: "type",
		Name: fmt.Sprintf("%sEdge", typeName),
		NameTypes: []*NameTypeFormatter{
			{
				Name:    "cursor",
				Type:    "String",
				Comment: "The position of each node in an array, used in pagination.",
			},
			{
				Name:    "node",
				Type:    typeName + "!",
				Comment: "Whether there are more pages to fetch following the current page.",
			},
		},
	}
	return tf.Format()
}
func (sg *SchemaGenerator) genPageInfo() string {
	tf := TypeFormatter{
		Kind: "type",
		Name: "PageInfo",
		NameTypes: []*NameTypeFormatter{
			{
				Name:    "endCursor",
				Type:    "String",
				Comment: "The cursor corresponding to the last node in edges.",
			},
			{
				Name:    "hasNextPage",
				Type:    "Boolean!",
				Comment: "Whether there are more pages to fetch following the current page.",
			}, {
				Name:    "hasPreviousPage",
				Type:    "Boolean!",
				Comment: "Whether there are any pages prior to the current page.",
			}, {
				Name:    "startCursor",
				Type:    "String",
				Comment: "The cursor corresponding to the first node in edges.",
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
		sg.genDeletePayload() + "\n" +
		sg.genInput() + "\n" +
		sg.genUserErrors() + "\n" +
		sg.genDeleteUserErrors() + "\n" +
		sg.genUserErrorCodeEnum() + "\n" +
		sg.genConnection() + "\n" +
		sg.genPageInfo() + "\n" +
		sg.genEdge()
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
func getDeleteTypePayloadObject(name string) string {
	return asCamStyle(name) + "DeletePayload"
}
func getAPINameForUpdate(name string) string {
	return asLowCaseCamStyle(name) + "Update"
}

func getAPINameForDelete(name string) string {
	return asLowCaseCamStyle(name) + "Delete"
}
