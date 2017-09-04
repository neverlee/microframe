package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/micro/cli"

	"github.com/neverlee/xclog/go"

	"github.com/neverlee/microframe/config"
	"github.com/neverlee/microframe/motherboard"
	"github.com/neverlee/microframe/service"
)

func srvService(configPath, pluginPath string) {
	conf, err := config.NewBoardConf(configPath)
	if err != nil {
		xclog.Fatalln(err)
	}

	board, err := motherboard.NewMotherBoard(conf, pluginPath)
	if err != nil {
		xclog.Fatalln("NewMotherBoard", err)
	}

	if err = board.PluginsNew(); err != nil {
		xclog.Fatalln("PluginsNew", err)
	}

	mservice := service.NewService()

	board.PluginsPreinit(mservice)
	mservice.Init()

	board.PluginsInit(mservice)
	board.PluginsStart(mservice)

	err = func() error {
		s := mservice

		if err := s.Start(); err != nil {
			return err
		}

		// start reg loop
		go s.Run()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

		select {
		// wait on kill signal
		case <-ch:
		// wait on context cancel
		case <-s.Context.Done():
		}

		// exit reg loop
		s.CloseExit()

		if err := s.Stop(); err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		xclog.Fatalln("service.Run", err)
	}

	board.PluginsStop(mservice)
	board.PluginsUninit(mservice)

}

func main() {
	app := cli.NewApp()
	app.Name = "MicroFramework"
	app.Usage = "Micro Service Framework"
	app.Version = "1.0.0"
	app.Authors = []cli.Author{cli.Author{Name: "Never Lee", Email: "listarmb@gmail.com"}}

	app.Commands = []cli.Command{
		{
			Name:    "srv",
			Aliases: []string{"s"},
			Usage:   "start a micro rpc service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Value: "./config.yal",
					Usage: "srv yaml configure file",
					// Destination: &configPath,
				},
				cli.StringFlag{
					Name:  "plugin, p",
					Value: "./plugins",
					Usage: "srv plugin path",
					// Destination: &configPath,
				},
			},
			Action: func(c *cli.Context) {
				srvService(c.String("config"), c.String("plugin"))
			},
		},
		{
			Name:    "web",
			Aliases: []string{"w"},
			Usage:   "start a micro web service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Value: "./config.yal",
					Usage: "web yaml configure file",
					// Destination: &configPath,
				},
			},
			Action: func(c *cli.Context) {
				xclog.Infoln("web service: ", c.String("config"))
			},
		},
	}

	app.Run(os.Args)
}
