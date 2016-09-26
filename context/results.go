package context

type Result interface {
	Marshal() interface{}
}

type ResultList []Result
