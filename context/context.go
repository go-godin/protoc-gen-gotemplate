package context

import (
	"fmt"
	"strings"
)

type FileContext struct {
	Options  []*FileOption
	Services []*Service
	Enums    []*Enum
}

type FileOption struct {
	GoPackage string
}

type Service struct {
	Name     string
	Package  string
	Methods  []*ServiceMethod
	Messages []*Message
	Enums    []*Enum
	Events   []*Event
}

type ServiceMethod struct {
	Name         string
	Input        []*FunctionInput
	Output       []*FunctionOutput
	RequestName  string
	ResponseName string
}

func (meth ServiceMethod) String() string {
	var params []string
	for _, input := range meth.Input {
		for _, field := range input.Fields {
			params = append(params, fmt.Sprintf("%s %s `json:\"%s\"", field.Name, field.Type, field.JsonName))
		}
	}

	return fmt.Sprintf("%s %s { %s }", meth.Name, meth.RequestName, strings.Join(params, "\n\t"))
}

type Field struct {
	Name     string
	JsonName string
	Type     FieldType
	Value    string
	Repeated bool
}

type FieldType struct {
	Name         string
	Enum         bool
	ProtoName    string
	DefaultValue string
}

type FunctionInput struct {
	Name   string
	Fields []*Field
}

type FunctionOutput struct {
	Name  string
	Field []*Field
}

type Event struct {
	Name     string
	JsonName string
	Fields   []*Field
}

type Message struct {
	Name      string
	Namespace string
	Fields    []*Field
}

func (msg Message) String() string {
	var fields []string
	for _, field := range msg.Fields {
		if field.Repeated && !field.Type.Enum {
			fields = append(fields, fmt.Sprintf("  %s []*%s `json:\"%s\" [default=%s]", field.Name, field.Type.Name, field.JsonName, field.Type.DefaultValue))
		} else if field.Repeated && field.Type.Enum {
			fields = append(fields, fmt.Sprintf("  %s []%s `json:\"%s\" [default=%s]", field.Name, field.Type.Name, field.JsonName, field.Type.DefaultValue))
		} else {
			fields = append(fields, fmt.Sprintf("  %s %s `json:\"%s\" [default=%s]", field.Name, field.Type.Name, field.JsonName, field.Type.DefaultValue))
		}
	}

	return fmt.Sprintf("\n%s.%s {\n%s\n}", msg.Namespace, msg.Name, strings.Join(fields, ",\n"))
}

type Enum struct {
	Name      string
	Namespace string
	Values    []*EnumValue
}

func (e Enum) String() string {
	var vals []string
	for _, val := range e.Values {
		vals = append(vals, fmt.Sprintf("  %s: %d", val.Name, val.Number))
	}
	return fmt.Sprintf("\n%s.%s {\n%s\n}", e.Namespace, e.Name, strings.Join(vals, ",\n"))
}

type EnumValue struct {
	Name   string
	Number int32
}
