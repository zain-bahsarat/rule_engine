package evaluator

import (
	"reflect"
)

func NewEnvironment(bindings map[string]interface{}) *Environment {
	s := make(map[string]Object)
	e := &Environment{store: s}
	buildEnv(e, bindings)

	return e
}

type Environment struct {
	store map[string]Object
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func buildEnv(env *Environment, bindings map[string]interface{}) {
	for k, v := range bindings {

		valPtr := reflect.ValueOf(v)
		switch valPtr.Kind() {
		case reflect.Slice:
			val, ok := v.([]string)
			if !ok {
				env.Set(k, &Error{Message: "Invalid value"})
				continue
			}
			env.Set(k, NewRegexList(val))

		case reflect.String:
			env.Set(k, &String{Value: v.(string)})
		case reflect.Bool:
			env.Set(k, &Boolean{Value: v.(bool)})
		case reflect.Float64:
			env.Set(k, &Number{Value: v.(float64)})
		case reflect.Int64:
			env.Set(k, &Number{Value: float64(v.(int64))})
		case reflect.Int:
			env.Set(k, &Number{Value: float64(v.(int))})

		default:
			env.Set(k, &Error{Message: "Invalid value"})
		}
	}
}
