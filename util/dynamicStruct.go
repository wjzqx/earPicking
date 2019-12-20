package util

import (
	"reflect"
)

var typeRegistry = make(map[string]reflect.Type)

// RegistryType ...
func RegistryType(elem interface{}){
	t := reflect.TypeOf(elem).Elem()
	typeRegistry[t.Name()] = t
}
// NewStruct
func NewStruct(name string) (interface{}, bool){
	elem, ok := typeRegistry[name]
	if !ok {
		return nil, false
	}

	return reflect.New(elem), true
}

type Operator interface {
	Apply(int, int) int
}

type Operation struct {
	Operator Operator
}


func (this *Operation) Operate(l, r int) int{
	return this.Operator.Apply(l,r)
}
