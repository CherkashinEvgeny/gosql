package internal

type Option struct {
	Key   string
	Value any
}

func FindOption(key string, options ...[]Option) (value any) {
	for i := len(options) - 1; i >= 0; i-- {
		o := options[i]
		for j := len(o) - 1; j >= 0; j-- {
			option := o[j]
			if option.Key == key {
				return option.Value
			}
		}
	}
	return nil
}
