package source

type ISource interface {
	Load(rootPath string) (map[string][]byte, error)
	Watch(rootPath string, callback func(namespaceKey string, data []byte))
}
