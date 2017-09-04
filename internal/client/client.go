package client

import (
	"time"

	"github.com/micro/go-micro/client"

	// "github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

type clientConf struct {
	// Type: rpc # client for go-micro; rpc
	RequestTimeout config.Duration `yaml:"request_timeout"`
	Retries        int             `yaml:"retries"`
	PoolSize       int             `yaml:"pool_size"`
	PoolTTL        config.Duration `yaml:"pool_ttl"`
}

type clientPlugin struct {
	pluginer.SrvPluginBase
	conf clientConf
}

func NewPlugin(pconf *config.RawYaml) (pluginer.SrvPluginer, error) {
	p := &clientPlugin{
		SrvPluginBase: pluginer.SrvPluginBase{
			Phase: pluginer.BasePhase,
		},
		conf: clientConf{
			RequestTimeout: config.Duration(time.Second * 5),
			Retries:        1,
			PoolSize:       0,
			PoolTTL:        config.Duration(time.Second * 60),
		},
	}
	if err := pconf.Decode(&p.conf); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *clientPlugin) Preinit(s *service.Service) error {
	var clientOpts []client.Option

	conf := p.conf

	// client opts
	if conf.Retries > 0 {
		clientOpts = append(clientOpts, client.Retries(conf.Retries))
	}

	if conf.RequestTimeout > 0 { //  len(str) > 0
		clientOpts = append(clientOpts, client.RequestTimeout(conf.RequestTimeout.Value()))
	}

	if conf.PoolSize > 0 {
		clientOpts = append(clientOpts, client.PoolSize(conf.PoolSize))
	}

	if conf.PoolTTL > 0 {
		clientOpts = append(clientOpts, client.PoolTTL(conf.PoolTTL.Value()))
	}

	s.ClientOption(clientOpts...)

	return nil
}
