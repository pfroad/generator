/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"generator/parser"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	jPackage string
)

// jgenCmd represents the jgen command
var jgenCmd = &cobra.Command{
	Use:   "jgen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("jgen called")
		ts := strings.Split(tables, ",")
		if _, err := os.Stat(output); os.IsNotExist(err) {
			os.MkdirAll(output, os.ModePerm)
		}

		for _, t := range ts {
			err := parser.GenForJava(schema, t, jPackage)
			if err != nil {
				fmt.Printf("failed to gen schema[%s] table[%s] java files. %v\n", schema, t, err)
			}
		}
		fmt.Println("done to gen!")
	},
}
