package auth2provider

var registry = map[string]func() Provider{}

func Register(name string, provider func() Provider) {
	registry[name] = provider
}

func GetProvider(name string) (Provider, bool) {
	provider,ok := registry[name]

	if !ok {
		return nil,false
	}

	return provider(),true
}
