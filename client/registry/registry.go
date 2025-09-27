package registry

import (
	"google.golang.org/grpc"
)

type ServerDescriptor struct {
	Address string
	Port    int
}

type Registry struct {
	Descriptors map[string]*ServerDescriptor
	Connections []*grpc.ClientConn
}

func NewRegistry() *Registry {
	return &Registry{Descriptors: make(map[string]*ServerDescriptor)}
}

func (r *Registry) Register(serviceName string, descriptor *ServerDescriptor) {
	r.Descriptors[serviceName] = descriptor
}

func (r *Registry) Get(serviceName string) *ServerDescriptor {
	return r.Descriptors[serviceName]
}

func (r *Registry) ClearEntry(serviceName string) {
	delete(r.Descriptors, serviceName)
}
