package main

type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
}

func (l *LoxClass) String() string {
	return l.Name
}

func (l *LoxClass) Call(i *Interpreter, args []interface{}) (interface{}, *RuntimeError) {
	return &LoxInstance{
		Class: l,
	}, nil
}

func (l *LoxClass) Arity() int {
	return 0
}

func (l *LoxClass) findMethod(name string) (*LoxFunction, bool) {
	if method, ok := l.Methods[name]; ok {
		return method, true
	}

	return nil, false
}
