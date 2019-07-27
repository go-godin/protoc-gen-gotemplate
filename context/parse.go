package context

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func FindMessageType(file *descriptor.FileDescriptorProto, search *descriptor.FieldDescriptorProto) *descriptor.DescriptorProto {
	for _, message := range file.MessageType {
		m := strings.Replace(search.GetTypeName(), fmt.Sprintf(".%s.", file.GetPackage()), "", 1)

		if message.GetName() == m {
			return message
		}
	}

	return nil
}

func RemoveNamespace(namespacedName string, packageName string) string {
	return strings.Replace(namespacedName, fmt.Sprintf(".%s.", packageName), "", 1)
}

func ParseMessages(file *descriptor.FileDescriptorProto) []*Message {
	var messages []*Message
	for _, protoMessage := range file.MessageType {

		var fields []*Field
		for _, messageField := range protoMessage.Field {
			field := ParseField(file, messageField)
			fields = append(fields, field)
		}

		msg := &Message{
			Name:      RemoveNamespace(protoMessage.GetName(), file.GetPackage()),
			Namespace: file.GetPackage(),
			Fields:    fields,
		}
		messages = append(messages, msg)
	}

	return messages
}

func ParseEnums(file *descriptor.FileDescriptorProto) []*Enum {
	var enums []*Enum
	for _, protoEnum := range file.EnumType {

		var enumValues []*EnumValue
		for _, enumValue := range protoEnum.Value {
			val := &EnumValue{
				Name: enumValue.GetName(),
				Number: enumValue.GetNumber(),
			}
			enumValues = append(enumValues, val)
		}

		enum := &Enum{
			Name: protoEnum.GetName(),
			Namespace: file.GetPackage(),
			Values: enumValues,
		}

		enums = append(enums, enum)
	}

	return enums
}

func ParseField(file *descriptor.FileDescriptorProto, protoField *descriptor.FieldDescriptorProto) *Field {
	field := &Field{}
	field.Name = protoField.GetName()
	field.JsonName = protoField.GetJsonName()
	field.Type.Enum = false

	// TYPES
	switch protoField.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		typ := FindMessageType(file, protoField)
		if typ == nil {
			break
		}
		field.Type.ProtoName = descriptor.FieldDescriptorProto_Type_name[int32(descriptor.FieldDescriptorProto_TYPE_MESSAGE)]
		field.Type.Name = typ.GetName()
		field.Type.DefaultValue = "nil"

		break

	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		field.Type.ProtoName = descriptor.FieldDescriptorProto_Type_name[int32(descriptor.FieldDescriptorProto_TYPE_BOOL)]
		field.Type.Name = "bool"
		field.Type.DefaultValue = "false"

		break

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		field.Type.ProtoName = descriptor.FieldDescriptorProto_Type_name[int32(descriptor.FieldDescriptorProto_TYPE_STRING)]
		field.Type.Name = "string"
		field.Type.DefaultValue = "\"\""

		break

	case descriptor.FieldDescriptorProto_TYPE_INT64:
		field.Type.ProtoName = descriptor.FieldDescriptorProto_Type_name[int32(descriptor.FieldDescriptorProto_TYPE_INT64)]
		field.Type.Name = "int64"
		field.Type.DefaultValue = "0"

		break

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		field.Type.ProtoName = descriptor.FieldDescriptorProto_Type_name[int32(descriptor.FieldDescriptorProto_TYPE_ENUM)]
		field.Type.Name = RemoveNamespace(protoField.GetTypeName(), file.GetPackage())
		field.Type.DefaultValue = protoField.GetDefaultValue()
		field.Type.Enum = true

		break
	}

	// LABELS
	switch *protoField.Label {
	case descriptor.FieldDescriptorProto_LABEL_REPEATED:
		field.Repeated = true
	}

	return field
}
