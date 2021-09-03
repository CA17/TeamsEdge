package main

import (
	"flag"
	"fmt"
	"github.com/teamsweHere/teamsedge-hy/common/installer"
	"github.com/teamsweHere/teamsedge-hy/config"
	"github.com/teamsweHere/teamsedge-hy/jobs"
	"github.com/teamsweHere/teamsedge-hy/service/app"
	"log"
	"os"
	"runtime"
	_ "time/tzdata"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group

	BuildVersion   string
	ReleaseVersion string
	BuildTime      string
	BuildName      string
	CommitID       string
	CommitDate     string
	CommitUser     string
	CommitSubject  string
)

// 命令行定义
var (
	h         = flag.Bool("h", false, "help usage")
	showVer   = flag.Bool("v", false, "show version")
	conffile  = flag.String("c", "/etc/teamsacs.yaml", "config yaml/json file")
	install   = flag.Bool("install", false, "run install")
	uninstall = flag.Bool("uninstall", false, "run uninstall")
	initcfg   = flag.Bool("initcfg", false, "write default config > /etc/teamsedge.yaml")
)

// PrintVersion Print version information
func PrintVersion() {
	fmt.Fprintf(os.Stdout, "build name:\t%s\n", BuildName)
	fmt.Fprintf(os.Stdout, "build version:\t%s\n", BuildVersion)
	fmt.Fprintf(os.Stdout, "build time:\t%s\n", BuildTime)
	fmt.Fprintf(os.Stdout, "release version:\t%s\n", ReleaseVersion)
	fmt.Fprintf(os.Stdout, "Commit ID:\t%s\n", CommitID)
	fmt.Fprintf(os.Stdout, "Commit Date:\t%s\n", CommitDate)
	fmt.Fprintf(os.Stdout, "Commit Username:\t%s\n", CommitUser)
	fmt.Fprintf(os.Stdout, "Commit Subject:\t%s\n", CommitSubject)
}

func printHelp() {
	if *h {
		ustr := fmt.Sprintf("%s version: %s, Usage:%s -h\nOptions:", BuildName, BuildVersion, BuildName)
		fmt.Fprintf(os.Stderr, ustr)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *showVer {
		PrintVersion()
		os.Exit(0)
	}

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
