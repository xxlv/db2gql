package main

import (
	"fmt"
)

// keep all formatters

// Args formatter like `name:String!`
type ArgsFormatter struct {
	Name     string
	Type     string
	Required bool
}

func (f *ArgsFormatter) Format() string {
	result := fmt.Sprintf("%s:%s", f.Name, f.Type)
	if f.Required {
		result += "!"
	}
	return result
}

type EnumFormatter struct {
	Name             string
	ValueWithComment map[string]string
}

func (f *EnumFormatter) Format() string {
	body := ""
	for enumValue, comment := range f.ValueWithComment {
		cf := CommentFormatter{Content: comment}
		body += "    " + cf.Format() + "\n"
		body += "    " + enumValue + "\n"
	}
	return fmt.Sprintf("%s %s {\n%s\n}", "enum", f.Name, body)
}

// Type formatter as `type Object {}`
type TypeFormatter struct {
	Kind      string
	Name      string
	NameTypes []*NameTypeFormatter
}

func (f *TypeFormatter) Format() string {
	body := ""
	for _, nameType := range f.NameTypes {
		body += nameType.Format()
	}
	return fmt.Sprintf("%s %s {\n%s\n}", f.Kind, f.Name, body)
}

// Name type formatter
type NameTypeFormatter struct {
	Name    string
	Type    string
	Comment string
	Args    []*ArgsFormatter
}

func (cf *NameTypeFormatter) Format() string {
	cf.Name = AsName(cf.Name)
	comment := &CommentFormatter{Content: cf.Comment}
	args := ""
	if len(cf.Args) > 0 {
		args = "("
		for i, arg := range cf.Args {
			if i < len(cf.Args)-1 {
				args += arg.Format() + ","
			} else {
				args += arg.Format() + ""
			}
		}
		args += ")"
	}
	result := fmt.Sprintf("\n    %s\n    %s%s:%s", comment.Format(), cf.Name, args, cf.Type)
	return result
}

// Comment formatter use `"""` comment style
type CommentFormatter struct {
	Content string
}

func (cf *CommentFormatter) Format() string {
	comment := fmt.Sprintf("%s %s %s", `"""`, cf.Content, `"""`)
	return comment
}
