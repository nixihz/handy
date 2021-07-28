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
