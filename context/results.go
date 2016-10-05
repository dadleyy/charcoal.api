package context

type Result interface {
	Public() interface{}
}

type ResultList []Result

func (list *ResultList) Apply() []interface{} {
	result := make([]interface{}, len(*list))

	for i, r := range *list {
		result[i] = r.Public()
	}

	return result
}

type ResultString struct {
	String string
}

func (str *ResultString) Public() interface{} {
	return str.String
}
