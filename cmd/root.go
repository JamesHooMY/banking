/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "banking",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Println("go version: ", runtime.Version())
		log.Println("os: ", runtime.GOOS)
		log.Println("arch: ", runtime.GOARCH)
	},
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
	cobra.OnInitialize(initConfig)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.banking.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// log.Println("Using config file:", viper.ConfigFileUsed())
	log.Println("Server mode:", viper.Get("server.runMode"))

	// set apm server env
	if err := os.Setenv("ELASTIC_APM_SERVICE_NAME", viper.GetString("apm.serviceName")); err != nil {
		panic(fmt.Errorf("Fatal error set apm service name: %s \n", err))
	}
	if err := os.Setenv("ELASTIC_APM_SERVER_URL", viper.GetString("apm.serverUrl")); err != nil {
		panic(fmt.Errorf("Fatal error set apm server url: %s \n", err))
	}
	if err := os.Setenv("ELASTIC_APM_SECRET_TOKEN", viper.GetString("apm.secretToken")); err != nil {
		panic(fmt.Errorf("Fatal error set apm secret token: %s \n", err))
	}

	// set elasticsearch env
	if err := os.Setenv("ELASTICSEARCH_URL", viper.GetString("elasticsearch.url")); err != nil {
		panic(fmt.Errorf("Fatal error set elasticsearch url: %s \n", err))
	}
	if err := os.Setenv("ELASTICSEARCH_USERNAME", viper.GetString("elasticsearch.username")); err != nil {
		panic(fmt.Errorf("Fatal error set elasticsearch username: %s \n", err))
	}
	if err := os.Setenv("ELASTICSEARCH_PASSWORD", viper.GetString("elasticsearch.password")); err != nil {
		panic(fmt.Errorf("Fatal error set elasticsearch password: %s \n", err))
	}
	if err := os.Setenv("ELASTICSEARCH_INDEX", viper.GetString("elasticsearch.index")); err != nil {
		panic(fmt.Errorf("Fatal error set elasticsearch index: %s \n", err))
	}
}
