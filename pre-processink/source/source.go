package source

type Source interface {
	Fetch() (string, error)
}
