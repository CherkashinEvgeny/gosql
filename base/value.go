package base

type Valuer interface {
	Value(key string) (value any)
}

var emptyInstance Valuer = empty{}

func Empty() Valuer {
	return emptyInstance
}

type empty struct {
}

func (e empty) Value(_ string) (value any) {
	return nil
}

func WithValue(parent Valuer, key string, value any) Valuer {
	return keyValue{parent, key, value}
}

type keyValue struct {
	parent Valuer
	key    string
	value  any
}

func (v keyValue) Value(key string) (value any) {
	if v.key == key {
		return v.value
	}
	if v.parent != nil {
		return v.parent.Value(key)
	}
	return nil
}
