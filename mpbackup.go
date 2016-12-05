package main

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/BurntSushi/toml"
	"os"
	"io/ioutil"
	"log"
	"fmt"
	"github.com/yuya008/mpbackup/mcsqueue"
)

var mpb Mpbackup

type Mpbackup struct {
	Config string
}

func (mpb *Mpbackup) parseArgs() {
	app := cli.NewApp()
	app.Name = "mpbackup"
	app.Usage = "米拍图片资源备份工具"
	app.UsageText = "mpbackup [arguments...]"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:"config, c",
			Usage:"指定配置文件所在",
			EnvVar:"MPBACKUPCONF",
			Value:"mpbackup.toml",
			Destination:&mpb.Config,
		},
	}
	app.Action = func(c *cli.Context) error {
		if (c.IsSet("config")) {
			mpb.Config = c.String("config")
		}
		return nil
	}
	app.Run(os.Args)
}

func (mpb *Mpbackup) parseConfig() {
	confByte, err := ioutil.ReadFile(mpb.Config)
	if err != nil {
		log.Panicf("%s %s", mpb.Config, err.Error())
	}
	err = toml.Unmarshal(confByte, mpb)
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Println(confByte)
}

func main()  {
	mpb.parseArgs()
	mpb.parseConfig()
	mcs := mcsqueue.New()
	mcs.Put()
}
