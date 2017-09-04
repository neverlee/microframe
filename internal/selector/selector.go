package selector

import (
	"errors"

	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/selector/cache"

	// "github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

type selectorConf struct {
	Selector string `yaml:"selector"`
}

type selectorPlugin struct {
	pluginer.SrvPluginBase
	conf selectorConf
}

func NewPlugin(pconf *config.RawYaml) (pluginer.SrvPluginer, error) {
	p := &selectorPlugin{
		SrvPluginBase: pluginer.SrvPluginBase{
			Phase: pluginer.BasePhase,
		},
		conf: selectorConf{
			Selector: "",
		},
	}
	if err := pconf.Decode(&p.conf); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *selectorPlugin) Preinit(s *service.Service) error {
	conf := p.conf

	// Set the selector
	// "default": selector.NewSelector, "cache":   cache.NewSelector,
	switch conf.Selector {
	case "", "default":
		s.Selector = selector.NewSelector()
	case "cache":
		s.Selector = cache.NewSelector()
	default:
		return errors.New("No such selector " + conf.Selector)
	}

	return nil
}
