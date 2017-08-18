package server

import (
	"time"

	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"

	// "github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

type serverConf struct {
	// Type: rpc # Server for go-micro; rpc
	Name             string            `yaml:"name"`
	Version          string            `yaml:"version"`
	Metadata         metadata.Metadata `yaml:"metadata"`
	ID               string            `yaml:"id"`
	Address          string            `yaml:"address"`
	Advertise        string            `yaml:"advertise"`
	RegisterInterval config.Duration   `yaml:"register_interval"`
	RegisterTTL      config.Duration   `yaml:"register_ttl"`
	Wait             bool              `yaml:"wait"`
}

type serverPlugin struct {
	pluginer.SrvPluginBase
	conf serverConf
}

func NewPlugin(mconf *config.SrvConf, pconf *config.RawYaml) (pluginer.SrvPluginer, error) {
	p := &serverPlugin{
		SrvPluginBase: pluginer.SrvPluginBase{
			Phase: pluginer.BasePhase,
		},
		conf: serverConf{
			Name:             "micro.frame.srv.service",
			Version:          "v0.0.1dev",
			Metadata:         nil,
			ID:               "",
			Address:          "",
			Advertise:        "",
			RegisterInterval: config.Duration(time.Second * 10),
			RegisterTTL:      config.Duration(time.Second * 5),
			Wait:             false,
		},
	}
	if err := pconf.Decode(&p.conf); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *serverPlugin) Preinit(s *service.Service) error {
	var serverOpts []server.Option

	conf := p.conf

	// Set the server
	if len(conf.Name) > 0 {
		serverOpts = append(serverOpts, server.Name(conf.Name))
	}

	if len(conf.Version) > 0 {
		serverOpts = append(serverOpts, server.Version(conf.Version))
	}

	if len(conf.Metadata) > 0 {
		serverOpts = append(serverOpts, server.Metadata(conf.Metadata))
	}

	if len(conf.ID) > 0 {
		serverOpts = append(serverOpts, server.Id(conf.ID))
	}

	if len(conf.Address) > 0 {
		serverOpts = append(serverOpts, server.Address(conf.Address))
	}

	if len(conf.Advertise) > 0 {
		serverOpts = append(serverOpts, server.Advertise(conf.Advertise))
	}

	s.RegisterInterval = conf.RegisterInterval.Value()

	if conf.RegisterTTL > 0 {
		serverOpts = append(serverOpts, server.RegisterTTL(conf.RegisterTTL.Value()))
	}

	serverOpts = append(serverOpts, server.Wait(conf.Wait))

	s.ServerOption(serverOpts...)

	return nil
}
