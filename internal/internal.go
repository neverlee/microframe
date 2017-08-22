package internal

import (
	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/pluginer"

	"github.com/neverlee/microframe/internal/broker"
	"github.com/neverlee/microframe/internal/client"
	"github.com/neverlee/microframe/internal/logger"
	"github.com/neverlee/microframe/internal/registry"
	"github.com/neverlee/microframe/internal/selector"
	"github.com/neverlee/microframe/internal/server"
	"github.com/neverlee/microframe/internal/transport"
	"github.com/neverlee/microframe/internal/ulimit"
)

var Plugins = map[string]func(*config.SrvConf, *config.RawYaml) (pluginer.SrvPluginer, error){
	"selector":  selector.NewPlugin,
	"server":    server.NewPlugin,
	"transport": transport.NewPlugin,
	"logger":    logger.NewPlugin,
	"client":    client.NewPlugin,
	"broker":    broker.NewPlugin,
	"ulimit":    ulimit.NewPlugin,
	"registry":  registry.NewPlugin,
}
