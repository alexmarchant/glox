package main

type LoxClass struct {
	Name       string
	Superclass *LoxClass
	Methods    map[string]*LoxFunction
}

func (l *LoxClass) String() string {
	return l.Name
}

func (l *LoxClass) Call(i *Interpreter, args []interface{}) (interface{}, *RuntimeError) {
	instance := &LoxInstance{
		Class:  l,
		Fields: map[string]interface{}{},
	}

	if initializer, ok := l.findMethod("init"); ok {
		initializer.bind(instance).Call(i, args)
	}

	return instance, nil
}

func (l *LoxClass) Arity() int {
	if initializer, ok := l.findMethod("init"); ok {
		return initializer.Arity()
	}

	return 0
}

func (l *LoxClass) findMethod(name string) (*LoxFunction, bool) {
	if method, ok := l.Methods[name]; ok {
		return method, true
	}

	if l.Superclass != nil {
		return l.Superclass.findMethod(name)
	}

	return nil, false
}
