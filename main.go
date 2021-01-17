package main

import (
	"flag"
	"os"

	"github.com/halvfigur/mqtty/control"
	"github.com/halvfigur/mqtty/network"
)

type repeatedStringArg []string

func (a *repeatedStringArg) String() string {
	return ""
}

func (a *repeatedStringArg) Set(v string) error {
	*a = append(*a, v)
	return nil
}

func parseArgs() control.Config {
	var conf control.Config

	flag.StringVar(&conf.Server, "b", "tcp://localhost", `Broker endpoint on the form <scheme>://<host>:
where "scheme" is one of tcp, ssl or ws`)

	flag.IntVar(&conf.Port, "p", 1883, "Broker port")

	flag.Var((*repeatedStringArg)(&conf.Topics), "t", "Topic to subscribe to. May be repeated")

	if *flag.Bool("h", false, "Display this message") {
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()

	return conf
}

func main() {
	conf := parseArgs()
	c := network.NewMqttClient()
	control.NewMqttApp(c, conf).Start()
}
