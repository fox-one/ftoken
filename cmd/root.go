package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/structs"
	"github.com/fox-one/ftoken/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	cfg       config.Config
	debugMode bool

	initialized bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ftoken",
	Short: "ftoken is a tool for generating new tokens",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	rootCmd.Version = ver
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogging, initDone)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ftoken.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "toggle debug mode")
}

func initConfig() {
	if initialized {
		return
	}

	if cfgFile == "" {
		dir, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		filename := path.Join(dir, ".ftoken.yaml")
		info, err := os.Stat(filename)
		if !os.IsNotExist(err) && !info.IsDir() {
			cfgFile = filename
		}
	}

	if cfgFile == "" {
		filename := "config.yaml"
		if info, err := os.Stat(filename); !os.IsNotExist(err) && !info.IsDir() {
			cfgFile = filename
		}
	}

	if cfgFile != "" {
		logrus.Debugln("use config file", cfgFile)
	}

	if err := config.Load(cfgFile, &cfg); err != nil {
		panic(err)
	}
}

func initLogging() {
	if initialized {
		return
	}

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	structs.DefaultTagName = "json"
}

func initDone() {
	initialized = true
}
