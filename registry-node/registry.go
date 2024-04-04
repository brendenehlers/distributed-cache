package registry

type RegistryEntry struct {
	Url string
}

type Registry interface {
	Register(url string) error
	Unregister(url string) error
	GetNode() (*RegistryEntry, error)
}