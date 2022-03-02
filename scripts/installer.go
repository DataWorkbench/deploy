package main

import (
	"github.com/DataWorkbench/enginemanager/internal/installer"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const DefaultConfigFile = "/etc/dataomnis/configuration.yaml"

func main() {
	var configFile string

	app := kingpin.New("ds", "The command to operate dataomnis services on k8s.")

	initCmd := app.Command("init", "Initialize configuration file of dataomnis.")

	installCmd := app.Command("install", "Install dataomnis service.")
	configFile = installCmd.Flag("file", "The configuration file with full-path of dataomnis.").Short('f').Default(DefaultConfigFile).String()

	app.HelpFlag.Short('h')
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case initCmd.FullCommand():
		installer.InitConfiguration()
	case installCmd.FullCommand():
		installer.Install(configFile)
	}
}
