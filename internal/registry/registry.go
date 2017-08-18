package registry

import (
	"errors"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry/mdns"

	// "github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

type registryConf struct {
	Registry string   `yaml:"registry"`
	Address  []string `yaml:"address"`
}

type registryPlugin struct {
	pluginer.SrvPluginBase
	conf registryConf
}

func NewPlugin(mconf *config.SrvConf, pconf *config.RawYaml) (pluginer.SrvPluginer, error) {
	p := &registryPlugin{
		SrvPluginBase: pluginer.SrvPluginBase{
			Phase: pluginer.BasePhase,
		},
		conf: registryConf{
			Registry: "consul",
			Address:  []string{"127.0.0.1:8500"},
		},
	}
	if err := pconf.Decode(&p.conf); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *registryPlugin) Preinit(s *service.Service) error {
	conf := p.conf

	// Set the registry
	switch conf.Registry {
	case "consul":
		s.Registry = consul.NewRegistry(registry.Addrs(conf.Address...))
	case "mdns":
		s.Registry = mdns.NewRegistry(registry.Addrs(conf.Address...))
	default:
		return errors.New("No such registry " + conf.Registry)
	}

	return nil
}
