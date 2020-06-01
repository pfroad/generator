// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

var (
	cfgFile string
	schema  string
	tables  string
	output  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generator",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("Run rootCmd")
	//},
}

//var startCmd = &cobra.Command{
//	Use:   "start [COMMAND] [OPTIONS]",
//	Short: "start generator",
//	Run: func(cmd *cobra.Command, args []string) {
//		fmt.Println("Run start")
//	},
//}

//var helpCmd = &cobra.Command{
//	Use:   "help",
//	Short: "Print generator help",
//	Run: func(cmd *cobra.Command, args []string) {
//		fmt.Println("help called")
//	},
//}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "generate 1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file (default is $HOME/.config.yaml)")

	// Local flag only for this command
	gogenCmd.Flags().StringVar(&schema, "schema", "", "database schema")
	gogenCmd.Flags().StringVar(&tables, "tables", "", "model database tables, use \",\" split. a,b,c")
	gogenCmd.Flags().StringVar(&output, "output", "", "output directory")

	jgenCmd.Flags().StringVar(&schema, "schema", "", "database schema")
	jgenCmd.Flags().StringVar(&tables, "tables", "", "model database tables, use \",\" split. a,b,c")
	jgenCmd.Flags().StringVar(&jPackage, "jPackage", "", "java package name")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//rootCmd.AddCommand(startCmd)
	//rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(gogenCmd)
	rootCmd.AddCommand(jgenCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".generator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".generator")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
	// 	"apuser", "airparking", "10.35.22.61:3306", "octopus")
	//conn := "apuser:airparking@tcp(10.35.22.61:3306)/airparking"
	//viper.GetString("db.user"),
	//viper.GetString("db.password"),
	//viper.GetString("addr"),
	//viper.GetString("schema"))
	conn := viper.GetString("db.conn")

	var err error
	db, err = sql.Open("mysql", conn)

	if err != nil {
		os.Exit(0)
	}
}
