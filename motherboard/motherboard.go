package motherboard

import (
	"path/filepath"
	"plugin"
	"sort"

	"github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/internal"
	"github.com/neverlee/microframe/pluginer"
	"github.com/neverlee/microframe/service"
)

type newPluginFunc func(*config.SrvConf, *config.RawYaml) (pluginer.SrvPluginer, error)

type pluginCard struct {
	new   newPluginFunc
	conf  *config.PluginConfig
	ctx   pluginer.SrvPluginer
	index uint32
}

type motherBoard struct {
	conf  *config.SrvConf
	cards []*pluginCard
}

func NewMotherBoard(conf *config.SrvConf, plupath string) (*motherBoard, error) {
	board := &motherBoard{
		conf:  conf,
		cards: make([]*pluginCard, len(conf.Plugins)),
	}
	if plupath == "" {
		plupath = conf.PluginPath
	}
	for i, pconf := range conf.Plugins {
		sopath := filepath.Join(plupath, pconf.Name) + ".so"

		xclog.Infoln("Load plugin libary:", sopath)
		solib, err := plugin.Open(sopath)
		rnew := internal.Plugins[pconf.Name]
		if err == nil {
			if sby, serr := solib.Lookup("NewPlugin"); serr == nil {
				rnew = sby.(func(*config.SrvConf, *config.RawYaml) (pluginer.SrvPluginer, error))
			} else {
				xclog.Fatalln("Bad plugin:", pconf.Name, serr)
			}
		}
		if rnew == nil {
			xclog.Fatalln("Open plugin libary", sopath, "Fail!")
			return nil, err
		}

		board.cards[i] = &pluginCard{
			new:   newPluginFunc(rnew),
			conf:  pconf,
			index: uint32(i),
		}
	}
	return board, nil
}

func (mb *motherBoard) sort() {
	sort.Slice(mb.cards, func(i, j int) bool {
		li := (uint64(mb.cards[i].ctx.GetPhase()) << 32) + uint64(mb.cards[i].index)
		lj := (uint64(mb.cards[j].ctx.GetPhase()) << 32) + uint64(mb.cards[j].index)
		return li < lj
	})
}

func (mb *motherBoard) PluginsNew() error {
	for _, card := range mb.cards {
		var err error
		card.ctx, err = card.new(mb.conf, &card.conf.Config)
		if err != nil {
			xclog.Errorln("PluginNew error", card.conf.Name, err)
			return err
		}
	}

	mb.sort()
	xclog.Infoln("Plugin set order:")
	for i, card := range mb.cards {
		xclog.Infof("+ plugin_%03d plugin:%s", i, card.conf.Name)
	}
	return nil
}

func (mb *motherBoard) PluginsPreinit(service *service.Service) {
	for _, card := range mb.cards {
		if err := card.ctx.Preinit(service); err != nil {
			xclog.Fatalf("Init plugin:%s error: %v", card.conf.Name, err)
		}
	}
}

func (mb *motherBoard) PluginsInit(service *service.Service) {
	for _, card := range mb.cards {
		if err := card.ctx.Init(service); err != nil {
			xclog.Fatalf("Init plugin:%s error: %v", card.conf.Name, err)
		}
	}
}

func (mb *motherBoard) PluginsStart(service *service.Service) {
	for _, card := range mb.cards {
		if err := card.ctx.Start(service); err != nil {
			xclog.Fatalf("Start plugin:%s error: %v", card.conf.Name, err)
		}
	}
}

func (mb *motherBoard) PluginsStop(service *service.Service) {
	for i := len(mb.cards) - 1; i >= 0; i-- {
		card := mb.cards[i]
		card.ctx.Stop(service)
	}
}

func (mb *motherBoard) PluginsUninit(service *service.Service) {
	for i := len(mb.cards) - 1; i >= 0; i-- {
		card := mb.cards[i]
		card.ctx.Uninit(service)
	}
}
