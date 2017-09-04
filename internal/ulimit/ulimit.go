package ulimit

import (
	"syscall"

	"github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

type ulimitConf struct {
	FileDescriptors uint64 `yaml:"fids"`
}

type ulimitPlugin struct {
	pluginer.SrvPluginBase
	conf ulimitConf
}

func NewPlugin(pconf *config.RawYaml) (pluginer.SrvPluginer, error) {
	p := &ulimitPlugin{
		SrvPluginBase: pluginer.SrvPluginBase{
			Phase: pluginer.PreaccessPhase,
		},
		conf: ulimitConf{
			FileDescriptors: 8192,
		},
	}
	if err := pconf.Decode(&p.conf); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *ulimitPlugin) Preinit(s *service.Service) error {
	if p.conf.FileDescriptors == 0 {
		return nil
	}

	rLimit := syscall.Rlimit{
		Max: p.conf.FileDescriptors,
		Cur: p.conf.FileDescriptors,
	}
	xclog.Infoln("Setting fid ulimit:", p.conf.FileDescriptors)
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		xclog.Errorln("Setting fid ulimit fail", err)
	}
	return nil
}
