package main

import "fmt"

type MutationGenerator struct {
	Name       string
	RawColumns []Column // keep origin col ref
}

func (mg *MutationGenerator) genDelete() string {
	args := []*ArgsFormatter{}
	args = append(args, &ArgsFormatter{
		Name:     "id",
		Type:     "ID",
		Required: true,
	})

	result := &NameTypeFormatter{
		Name:    getAPINameForDelete(mg.Name),
		Comment: fmt.Sprintf("Delete %s", mg.Name),
		Args:    args,
		Type:    asCamStyle(mg.Name) + "DeletePayload",
	}
	return result.Format()
}

func (mg *MutationGenerator) genUpdate() string {
	args := []*ArgsFormatter{}
	inputArg := &ArgsFormatter{
		Name:     "input",
		Type:     getTypeInputObject(mg.Name),
		Required: true,
	}

	args = append(args, &ArgsFormatter{
		Name:     "id",
		Type:     "ID",
		Required: true,
	})
	args = append(args, inputArg)

	result := &NameTypeFormatter{
		Name:    getAPINameForUpdate(mg.Name),
		Comment: fmt.Sprintf("Update %s", mg.Name),
		Args:    args,
		Type:    asCamStyle(mg.Name) + "Payload",
	}
	return result.Format()
}

func (mg *MutationGenerator) Gen() string {

	return fmt.Sprintf("type Mutation {%s\n}", mg.genUpdate()+mg.genDelete())
}
