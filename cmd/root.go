/*
Copyright © 2025 Zingui Fred Mike <mikezingui@yahoo.com>
*/
package cmd

import (
	"github.com/opencontainers/runc/libcontainer"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Init process for libcontainer (INTERNAL ONLY)",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Verrouillage du thread OS
		// Les namespaces Linux sont liés au thread, pas à la Goroutine.
		// On doit s'assurer que ce code reste sur le thread principal.
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()

		// 2. Initialisation du conteneur (Nouvelle API)
		// Plus de Factory, plus de New(), plus de Cgroupfs.
		// Init() récupère automatiquement la configuration envoyée par le processus parent via le pipe.
		libcontainer.Init()

	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Nexus",
	Short: "Cloud Storage & computation tool",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Nexus.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(initCmd)
}
