package main

import (
	"flag"
	"fmt"
	gohttp "net/http"

	"github.com/github/freno/go/config"
	"github.com/github/freno/go/group"
	"github.com/github/freno/go/http"
	"github.com/outbrain/golib/log"
)

// AppVersion has to be filled by ldflags:
var AppVersion string

func main() {
	if AppVersion == "" {
		AppVersion = "local-build"
	}
	log.Infof("starting freno %s", AppVersion)

	configFile := flag.String("config", "", "config file name")
	http := flag.Bool("http", false, "spawn the HTTP API server")
	help := flag.Bool("help", false, "show the help")
	flag.Parse()

	loadConfiguration(*configFile)

	flag.Parse()

	switch {
	case *http:
		err := httpServe()
		log.Errore(err)
	case *help:
		printHelp()
	default:
		printUsage()
	}
}

func loadConfiguration(configFile string) {
	var err error
	if configFile != "" {
		err = config.Instance().Read(configFile)
	} else {
		err = config.Instance().Read("/etc/freno.conf.json", "conf/freno.conf.json")
	}

	if err != nil {
		log.Fatalf("Error reading configuration, please check the logs. Error was: %s", err.Error())
	}
}

func httpServe() error {
	log.Infof("Starting raft")
	if err := group.Setup(); err != nil {
		return err
	}
	go group.Monitor()

	api := http.NewAPIImpl()
	router := http.ConfigureRoutes(api)
	port := config.Settings().ListenPort
	log.Infof("Starting server in port %d", port)
	return gohttp.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func printHelp() {
	panic("not yet implemented")
}

func printUsage() {
	fmt.Println(`Usage: freno [OPTIONS]

	For more help options use: freno -help.
	`)
}
