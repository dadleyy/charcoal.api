package context

type Result interface {
	Public() interface{}
}

type ResultList []Result

func (list ResultList) Apply() []interface{} {
	out := make([]interface{}, 0)

	for _, result := range list {
		out = append(out, result.Public())
	}

	return out
}

type ResultString struct {
	String string
}

func (str *ResultString) Public() interface{} {
	return str.String
}
