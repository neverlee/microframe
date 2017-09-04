package logger

import (
	"time"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"

	"github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

type loggerConf struct {
	Handle    bool `yaml:"handle"`
	Subscribe bool `yaml:"subscribe"`
	Call      bool `yaml:"call"`
}

type loggerPlugin struct {
	pluginer.SrvPluginBase
	conf loggerConf
}

func NewPlugin(pconf *config.RawYaml) (pluginer.SrvPluginer, error) {
	p := &loggerPlugin{
		SrvPluginBase: pluginer.SrvPluginBase{
			Phase: pluginer.LogPhase,
		},
		conf: loggerConf{
			Handle:    true,
			Subscribe: true,
			Call:      true,
		},
	}
	if err := pconf.Decode(&p.conf); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *loggerPlugin) handlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			t := time.Now()
			err := fn(ctx, req, rsp)
			xclog.Infof("[logger:han]: %s duration: %v", req.Method(), time.Since(t))
			return err
		}
	}
}

func (p *loggerPlugin) subscriberWrapper() server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, req server.Publication) error {
			t := time.Now()
			err := fn(ctx, req)
			xclog.Info("[logger:sub]: %s duration: %v", req.Topic(), time.Since(t))
			return err
		}
	}
}

func (p *loggerPlugin) callWrapper() client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) error {
			t := time.Now()
			err := cf(ctx, addr, req, rsp, opts)
			xclog.Infoln("[logger:cal]: %s %s.%s duration: %v", addr, req.Service(), req.Method(), time.Since(t))
			return err
		}
	}
}

//func (p *loggerPlugin) ClientWrapper() client.Wrapper {
//	return func(c client.Client) client.Client {
//		return c
//	}
//}
//

func (p *loggerPlugin) Preinit(s *service.Service) error {
	if p.conf.Handle {
		s.ServerOption(server.WrapHandler(p.handlerWrapper()))
	}
	if p.conf.Subscribe == false {
		s.ServerOption(server.WrapSubscriber(p.subscriberWrapper()))
	}
	if p.conf.Call == false {
		s.ClientOption(client.WrapCall(p.callWrapper()))
	}
	return nil
}
