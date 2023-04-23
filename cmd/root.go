/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "summary.go",
	Short: "CLI Tool using OpenAI",
	Long:  `A simple command line interface to to use OpenAI API.`,
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
	printBanner()

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.summary.go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".summary.go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".summary.go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		viper.SetConfigFile(filepath.Join(home, ".summary.go.yaml"))
		err = viper.SafeWriteConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err = viper.WriteConfig()
				if err != nil {
					fmt.Println("Error creating config file:", err)
					os.Exit(1)
				}
			} else {
				fmt.Println("Error saving config file:", err)
				os.Exit(1)
			}
		}
	}

	checkOpenAIKey()
}

func checkOpenAIKey() {
	openAIKey := viper.GetString("openai-key")

	if openAIKey == "" {
		fmt.Println("It looks like you haven't set your OpenAI API key.")
		fmt.Print("Please enter your OpenAI API key: ")

		reader := bufio.NewReader(os.Stdin)
		openAIKey, _ = reader.ReadString('\n')
		openAIKey = strings.TrimSpace(openAIKey)

		// Save the OpenAI key to the config file
		viper.Set("openai-key", openAIKey)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Error saving OpenAI API key to config file:", err)
			os.Exit(1)
		}

		fmt.Println("Your OpenAI API key has been saved to the config file.")
	}
}
