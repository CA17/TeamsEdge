package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	_ "time/tzdata"

	"github.com/ca17/teamsedge/assets"
	"github.com/ca17/teamsedge/common/installer"
	"github.com/ca17/teamsedge/config"
	"github.com/ca17/teamsedge/jobs"
	"github.com/ca17/teamsedge/service/app"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

// 命令行定义
var (
	h         = flag.Bool("h", false, "help usage")
	conffile  = flag.String("c", "/etc/teamsacs.yaml", "config yaml/json file")
	install   = flag.Bool("install", false, "run install")
	uninstall = flag.Bool("uninstall", false, "run uninstall")
	initcfg   = flag.Bool("initcfg", false, "write default config > /etc/teamsedge.yaml")
)

// PrintVersion Print version information
func PrintVersion() {
	fmt.Println(assets.BuildInfo)
}

func printHelp() {
	if *h {
		ustr := fmt.Sprintf("TeamsEdge Usage: teamsacs -h\n")
		fmt.Fprintf(os.Stderr, ustr)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	PrintVersion()

	printHelp()

	_config := config.LoadConfig(*conffile)

	if *initcfg {
		err := installer.InitConfig(_config)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 安装为系统服务
	if *install {
		err := installer.Install(_config)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// 卸载
	if *uninstall {
		installer.Uninstall()
		return
	}

	app.Init(_config)
	jobs.Init()

	g.Go(func() error {
		log.Println("Start Task Queue reciver ...")
		return app.StartSubscribe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
