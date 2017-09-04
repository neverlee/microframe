package transport

import (
	"github.com/micro/go-micro/transport"

	// "github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

type transportConf struct {
	Address []string `yaml:"address"`
}

type transportPlugin struct {
	pluginer.SrvPluginBase
	conf transportConf
}

func NewPlugin(pconf *config.RawYaml) (pluginer.SrvPluginer, error) {
	p := &transportPlugin{
		SrvPluginBase: pluginer.SrvPluginBase{
			Phase: pluginer.BasePhase,
		},
		conf: transportConf{
			Address: nil,
		},
	}
	if err := pconf.Decode(&p.conf); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *transportPlugin) Preinit(s *service.Service) error {
	conf := p.conf

	// Set the transport
	// "default": transport.Newtransport, "cache":   cache.Newtransport,
	if len(conf.Address) > 0 {
		s.Transport = transport.NewTransport(transport.Addrs(conf.Address...))
	}

	return nil
}
