package protoval

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func Reflect(validationTag string, message proto.Message, vc map[string]func() error, prefix string, index int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	reflector := message.ProtoReflect()
	fields := reflector.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		options := field.Options().(*descriptorpb.FieldOptions)
		fieldName := field.TextName()
		if field.JSONName() != "" {
			fieldName = field.JSONName()
		}
		name := ""
		if prefix == "" {
			name = fieldName
		} else {
			name = fmt.Sprintf("%s.%s", prefix, fieldName)
		}
		if index != -1 {
			name += fmt.Sprintf("[%d]", index)
		}
		proto.RangeExtensions(options, func(et protoreflect.ExtensionType, i interface{}) bool {
			fullName := et.TypeDescriptor().FullName()
			if fullName != protoreflect.FullName(validationTag) {
				return true
			}
			rule, ok := i.(string)
			if !ok {
				return true
			}
			value := reflector.Get(field).Interface()
			rules, err := ExprParser(rule)
			if err != nil {
				panic(err)
			}
			vc[name] = func() error {
				for ruleName, rule := range rules {
					validator, ok := _rules[ruleName]
					if !ok {
						return fmt.Errorf("%s: rule not found", ruleName)
					}
					err := validator(name, value, rule)
					if err != nil {
						return err
					}
				}
				return nil
			}
			return false
		})
		if field.Kind() == protoreflect.MessageKind {
			if field.IsList() {
				list := reflector.Get(field).List()
				for i := 0; i < list.Len(); i++ {
					err := Reflect(validationTag, list.Get(i).Message().Interface(), vc, name, i)
					if err != nil {
						return err
					}
				}
				continue
			}
			err := Reflect(validationTag, reflector.Get(field).Message().Interface(), vc, name, -1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
