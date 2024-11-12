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
		Kind:      "type",
		Name:      sg.Name,
		NameTypes: asNameTypeFormatter(sg.RawColumns),
	}
	return tf.Format()
}

func asNameTypeFormatter(cols []Column) []*NameTypeFormatter {
	result := []*NameTypeFormatter{}
	if len(cols) <= 0 {
		return result
	}
	for _, v := range cols {
		result = append(result, &NameTypeFormatter{
			Name:    v.Name,
			Type:    mapMySQLTypeToGraphQL(v.Type, v.Null),
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
				Name:    asLowCaseCamStyle(sg.Name + "userErrors"),
				Type:    fmt.Sprintf("[%s!]!", asCamStyle(sg.Name+"userErrors")),
				Comment: "The list of errors that occurred from executing the mutation.",
			},
		},
	}
	return tf.Format()
}

func (sg *SchemaGenerator) genInput() string {
	tf := TypeFormatter{
		Kind:      "input",
		Name:      sg.Name,
		NameTypes: asNameTypeFormatter(sg.RawColumns),
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
		sg.genInput()
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
