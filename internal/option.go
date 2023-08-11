package internal

type Option struct {
	Key   string
	Value any
}

type Options []Option

func (o Options) Value(key string) (value any) {
	for i := len(o) - 1; i >= 0; i-- {
		option := o[i]
		if option.Key == key {
			return option.Value
		}
	}
	return nil
}
