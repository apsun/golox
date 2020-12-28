package lox

type Environment struct {
	values map[string]Value
}

func NewEnvironment() *Environment {
	return &Environment{
		values: map[string]Value{},
	}
}

func (e *Environment) Define(key string, value Value) {
	e.values[key] = value
}

func (e *Environment) Get(key string) *Value {
	value, ok := e.values[key]
	if !ok {
		return nil
	}
	return &value
}
