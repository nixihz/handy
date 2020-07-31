/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"jaytaylor.com/html2text"
	"os"

	"github.com/spf13/cobra"
)

// mdviewCmd represents the mdview command
var mdviewCmd = &cobra.Command{
	Use:   "mdview",
	Short: "markdown viewer for command-line",
	Long:  `markdown viewer for command-line 命令行查看 markdown 文件`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("require markdown filepath")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			println(err.Error())
		}
		defer file.Close()

		content, err := ioutil.ReadAll(file)
		html := string(blackfriday.MarkdownCommon([]byte(content)))

		text, err := html2text.FromString(string(html), html2text.Options{PrettyTables: true})
		if err != nil {
			panic(err)
		}
		fmt.Println(text)
	},
}

func init() {
	rootCmd.AddCommand(mdviewCmd)
}
