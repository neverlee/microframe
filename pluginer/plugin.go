package pluginer

import (
	"github.com/neverlee/microframe/service"
)

const (
	BasePhase = 100 * (iota + 0)
	PreaccessPhase
	LogPhase
	AccessPhase
	ContentPhase
)

type SrvPluginer interface {
	GetPhase() uint32

	Preinit(service *service.Service) error
	Init(service *service.Service) error
	Start(service *service.Service) error
	Stop(service *service.Service) error
	Uninit(service *service.Service) error

	// UpdateConfig(service *service.Service) error

	// ClientWrapper() client.Wrapper
	// CallWrapper() client.CallWrapper
	// HandlerWrapper() server.HandlerWrapper
	// SubscriberWrapper() server.SubscriberWrapper
}

type SrvPluginBase struct {
	Phase uint32
}

func (pb *SrvPluginBase) GetPhase() uint32 {
	return uint32(pb.Phase)
}

func (pb *SrvPluginBase) Preinit(service *service.Service) error { return nil }
func (pb *SrvPluginBase) Init(service *service.Service) error    { return nil }
func (pb *SrvPluginBase) Start(service *service.Service) error   { return nil }
func (pb *SrvPluginBase) Stop(service *service.Service) error    { return nil }
func (pb *SrvPluginBase) Uninit(service *service.Service) error  { return nil }

// func (pb *SrvPluginBase) UpdateConfig(service *service.Service) error { return nil }
