package adapter

type Persister interface {
	Store(key, arg interface{}) error
	Load(key, arg interface{}) error
}
