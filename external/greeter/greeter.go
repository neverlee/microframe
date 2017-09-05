package main

import (
	"github.com/neverlee/xclog/go"
	hello "greeter/proto"
	"golang.org/x/net/context"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

// 配置结构体
type sayConf struct {
	Salutation string `yaml:"salutation"`
}

type sayPlugin struct {
	pluginer.SrvPluginBase
	conf sayConf
}

// NewPlugin 生成插件ctx
func NewPlugin(pconf *config.RawYaml) (pluginer.SrvPluginer, error) {
	s := &sayPlugin{
		SrvPluginBase: pluginer.SrvPluginBase{
			Phase: pluginer.ContentPhase,
		},
		conf: sayConf{
			Salutation: "Hi",
		},
	}
	if err := pconf.Decode(&s.conf); err != nil {
		return nil, err
	}
	return s, nil
}

// Hello 对外的Hello rpc函数
func (s *sayPlugin) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	xclog.Infoln("Received Say.Hello request")
	rsp.Msg = s.conf.Salutation + " " + req.Name
	return nil
}

func (s *sayPlugin) Init(srv *service.Service) error {
	// 向框架注册handler
	hello.RegisterSayHandler(srv.Server, s)
	return nil
}

// 此函数必须要有
func main() {
}
