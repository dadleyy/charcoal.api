package context

type SessionStore interface {
	SetInt(string, int) error
	SetString(string, string) error
	GetInt(string) (int, error)
	GetString(string) (string, error)
}
