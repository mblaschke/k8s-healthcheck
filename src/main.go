package main

import (
	"os"
	"fmt"
	"github.com/jessevdk/go-flags"
	"k8s-healthcheck/src/k8s"
	"k8s-healthcheck/src/logger"
)

const (
	Author  = "webdevops.io"
	Version = "0.1.0"
)

var (
	argparser   *flags.Parser
	args        []string
	Logger      *logger.DaemonLogger
	ErrorLogger *logger.DaemonLogger
	k8sService  *k8s.Kubernetes
)

var opts struct {
	KubeConfig  string `long:"kubeconfig"              env:"KUBECONFIG"              description:"Path to .kube/config"`
	KubeContext string `long:"kubecontext"             env:"KUBECONTEXT"             description:"Context of .kube/config"`
	ScrapeTime  int `   long:"scrape-time"             env:"SCRAPE_TIME"             description:"Scrape time in seconds"  default:"30"`
}

func main() {
	initArgparser()

	// Init logger
	Logger = logger.CreateDaemonLogger(0)
	ErrorLogger = logger.CreateDaemonErrorLogger(0)

	Logger.Messsage("init")
	k8sService = initK8sService()

	Logger.Messsage("starting metrics collection")
	initMetrics()

	Logger.Messsage("starting http server")
	startHttpServer()
}

func initArgparser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println(err)
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	if opts.KubeConfig == "" {
		kubeconfigPath := fmt.Sprintf("%s/.kube/config", UserHomeDir())
		if _, err := os.Stat(kubeconfigPath); err == nil {
			opts.KubeConfig = kubeconfigPath
		}
	}
}

func initK8sService() *k8s.Kubernetes {
	service := k8s.Kubernetes{}
	service.KubeConfig = opts.KubeConfig
	service.KubeContext = opts.KubeContext
	service.Logger = Logger

	return &service
}
