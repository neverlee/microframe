package service

import (
	// "os" "os/signal" "syscall"

	"time"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"

	"golang.org/x/net/context"
)

type Service struct {
	Broker    broker.Broker
	Client    client.Client
	Server    server.Server
	Registry  registry.Registry
	Selector  selector.Selector
	Transport transport.Transport

	// Register loop interval
	RegisterInterval time.Duration

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// registryOpts  []registry.Option
	// transportOpts []transport.Option
	brokerOpts   []broker.Option
	serverOpts   []server.Option
	clientOpts   []client.Option
	selectorOpts []selector.Option

	exit chan bool
}

func NewService() *Service {
	srv := Service{
		Broker:    broker.DefaultBroker,
		Client:    client.DefaultClient,
		Server:    server.DefaultServer,
		Registry:  registry.DefaultRegistry,
		Transport: transport.DefaultTransport,
		Selector:  selector.DefaultSelector,

		Context: context.Background(),

		exit: make(chan bool),
	}

	return &srv
	// mylient = &clientWrapper{
	// 	opt.Client,
	// 	metadata.Metadata{
	// 		HeaderPrefix + "From-Service": opt.Server.Options().Name,
	// 	},
	// }
}

// Init initialises service.
func (s *Service) Init() {
	s.SelectorOption(
		selector.Registry(s.Registry),
	)

	s.BrokerOption(
		broker.Registry(s.Registry),
	)

	s.ClientOption(
		client.Broker(s.Broker),
		client.Registry(s.Registry),
		client.Transport(s.Transport),
		client.Selector(s.Selector),
	)

	s.ServerOption(
		server.Broker(s.Broker),
		server.Registry(s.Registry),
		server.Transport(s.Transport),
	)

	// s.Registry.Init(s.registryOpts)
	// s.Transport.Init(s.transportOpts)
	s.Selector.Init(s.selectorOpts...)
	s.Broker.Init(s.brokerOpts...)
	s.Client.Init(s.clientOpts...)
	s.Server.Init(s.serverOpts...)

	s.cleanOption()
}

func (s *Service) cleanOption() {
	// s.registryOpts = nil
	// s.transportOpts = nil
	s.brokerOpts = nil
	s.serverOpts = nil
	s.clientOpts = nil
	s.selectorOpts = nil
}

// func (s *Service) RegistryOption(o ...registry.Option) {
// 	s.registryOpts = append(s.registryOpts, o...)
// }
// func (s *Service) TransportOption(o ...transport.Option) {
// 	s.transportOpts = append(s.transportOpts, o...)
// }
func (s *Service) BrokerOption(o ...broker.Option) {
	s.brokerOpts = append(s.brokerOpts, o...)
}
func (s *Service) ServerOption(o ...server.Option) {
	s.serverOpts = append(s.serverOpts, o...)
}
func (s *Service) ClientOption(o ...client.Option) {
	s.clientOpts = append(s.clientOpts, o...)
}
func (s *Service) SelectorOption(o ...selector.Option) {
	s.selectorOpts = append(s.selectorOpts, o...)
}

func (s *Service) Run() {
	if s.RegisterInterval <= time.Duration(0) {
		return
	}

	t := time.NewTicker(s.RegisterInterval)

	for {
		select {
		case <-t.C:
			s.Server.Register()
		case <-s.exit:
			t.Stop()
			return
		}
	}
}

func (s *Service) String() string {
	return "microframe-service"
}

func (s *Service) Start() error {
	if err := s.Server.Start(); err != nil {
		return err
	}

	if err := s.Server.Register(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Stop() error {
	var gerr error

	if err := s.Server.Deregister(); err != nil {
		return err
	}

	if err := s.Server.Stop(); err != nil {
		return err
	}

	return gerr
}

func (s *Service) CloseExit() {
	close(s.exit)
}

// func (s *Service) Run() error {
// 	if err := s.Start(); err != nil {
// 		return err
// 	}
//
// 	// start reg loop
// 	go s.run()
//
// 	ch := make(chan os.Signal, 1)
// 	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
//
// 	select {
// 	// wait on kill signal
// 	case <-ch:
// 	// wait on context cancel
// 	case <-s.Context.Done():
// 	}
//
// 	// exit reg loop
// 	s.CloseExit()
//
// 	if err := s.Stop(); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
