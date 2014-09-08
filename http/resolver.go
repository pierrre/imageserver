package http

// Resolver resolves a parameter to a HTTP parameter
//
// It returns the resolved HTTP parameter, or an empty string
type Resolver interface {
	Resolve(parameter string) (httpParameter string)
}

// ResolverFunc is a Resolver func
type ResolverFunc func(parameter string) string

// Resolve calls the func
func (f ResolverFunc) Resolve(parameter string) string {
	return f(parameter)
}

// Resolve resolves a parameter with a potential Resolver
//
// If it's not a resolver, it returns an empty string
func Resolve(i interface{}, parameter string) string {
	if resolver, ok := i.(Resolver); ok {
		return resolver.Resolve(parameter)
	}
	return ""
}
