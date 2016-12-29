package main

import (
	"github.com/mkideal/cli"
	"github.com/yuya008/mpbackup/mpbackup"
)

func main() {
	mpbp := new(mpbackup.Mpbackup)
	cli.Run(&mpbackup.Cfg{}, func(ctx *cli.Context) error {
		cfg := ctx.Argv().(*mpbackup.Cfg)
		mpbp.Run(cfg)
		return nil
	}, "米拍图片备份程序")
}
