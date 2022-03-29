package main

import (
	"fmt"
	"github.com/DataWorkbench/deploy/internal/installer"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const DefaultConfigFile = "/etc/dataomnis/configuration.yaml"

func main() {
	var configFile *string
	var debug, dryRun *bool
	var services *[]string

	app := kingpin.New("ds", "The command to operate dataomnis services on k8s")

	initCmd := app.Command("init", "Initialize configuration file of dataomnis")

	// installer
	serviceHelper := fmt.Sprintf("The services%s to install, default: all services of dataomnis.", installer.AllServices)
	installCmd := app.Command("install", "Install dataomnis service")
	debug = installCmd.Flag("debug", "Enable debug mode").Bool()
	dryRun = installCmd.Flag("dry-run", "if enable dry run install release for helm").Bool()
	configFile = installCmd.Flag("file", "The configuration file with full-path of dataomnis").Short('f').Default(DefaultConfigFile).String()
	services = installCmd.Flag("services", serviceHelper).Short('s').Default(installer.AllServices...).Strings()

	app.HelpFlag.Short('h')

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case initCmd.FullCommand():
		installer.InitConfiguration()
	case installCmd.FullCommand():
		_ = installer.Install(*configFile, services, *debug, *dryRun)
	}
}
