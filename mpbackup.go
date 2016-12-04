package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func parseArgs()  {
	app := cli.NewApp()
	app.Run(os.Args)
}


func main()  {
	parseArgs()
}
