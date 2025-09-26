package registry

type ServerDescriptor struct {
	Address string
	Port    int
}

type Registry struct {
	Descriptors map[string][]*ServerDescriptor
}

func NewRegistry() *Registry {
	return &Registry{Descriptors: make(map[string][]*ServerDescriptor)}
}

func (r *Registry) Register(serviceName string, descriptor *ServerDescriptor) {
	r.Descriptors[serviceName] = append(r.Descriptors[serviceName], descriptor)
}

func (r *Registry) Get(serviceName string) []*ServerDescriptor {
	return r.Descriptors[serviceName]
}

func (r *Registry) ClearEntry(serviceName string) {
	delete(r.Descriptors, serviceName)
}
