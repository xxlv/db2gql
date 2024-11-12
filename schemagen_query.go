package main

import (
	"fmt"
	"strings"
)

type QueryGenerator struct {
	Name       string
	RawColumns []Column
}

func (qg *QueryGenerator) genQueryByID() string {
	result := &NameTypeFormatter{
		Name:    asLowCaseCamStyle(qg.Name),
		Comment: fmt.Sprintf("Query %s", qg.Name),
		Type:    asCamStyle(qg.Name),
	}
	return result.Format()
}

func (qg *QueryGenerator) genQueryPage() string {
	name := asLowCaseCamStyle(qg.Name)
	if !strings.HasSuffix(name, "s") {
		name = name + "s"
	} else {
		name = name + "List"
	}
	result := &NameTypeFormatter{
		Name:    name,
		Comment: fmt.Sprintf("Query by page %s", qg.Name),
		Type:    asCamStyle(qg.Name),
	}
	return result.Format()
}

func (qg *QueryGenerator) Gen() string {
	return fmt.Sprintf("extend type Query {%s\n}", qg.genQueryByID()+
		qg.genQueryPage())
}
