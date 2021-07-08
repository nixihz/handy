package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/tomarkdown"

	"github.com/spf13/cobra"
)

// blogationCmd represents the blogation command
var blogationCmd = &cobra.Command{
	Use:   "blogation",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: publish,
}

func init() {
	rootCmd.AddCommand(blogationCmd)
}

var logf = log.Printf
var panicIf = log.Panic

func publish(cmd *cobra.Command, args []string) {
	fmt.Println("blogation called")
	authToken2 := "3fae99bf3ebd3573ad80a8cc13dad050cad2094d5c03c8a6f1cde859562be8ecc392723ae668859d3a5e777d6ae4bdcb8a2fa95441f84f7bac8da9c6111e2c2ad9fbf1d168ae43392ac94f4d2e32"
	pageID := "47d84abb6c0b4e0eb799eec100785fab"
	postsDir := "/Users/admin/Documents/fencex.github.io/content/posts/"
	imageDir := "/Users/admin/Documents/fencex.github.io/static/images"

	client := &notionapi.Client{AuthToken: authToken2}
	rootPage, _ := client.DownloadPage(pageID)

	levelOneIds := rootPage.GetSubPages()
	for _, levelOneId := range levelOneIds {
		// 第一级分类
		levelOnePage, _ := client.DownloadPage(levelOneId)
		levelOneTitle := levelOnePage.Root().GetTitle()[0].Text

		indexText := getIndexContent(levelOneTitle)
		writeNewFile(indexText, postsDir+levelOneTitle, "_index.md")

		levelTwoIds := levelOnePage.GetSubPages()
		for _, levelTwoId := range levelTwoIds {

			// 第二级文章
			levelTwoPage, _ := client.DownloadPage(levelTwoId)
			title := levelTwoPage.Root().GetTitle()[0].Text

			// 修改 markdown 文件
			levelTwoPageMdStr := string(tomarkdown.ToMarkdown(levelTwoPage))
			nn := strings.Index(levelTwoPageMdStr, "\n\n")
			levelTwoPageMdStr = levelTwoPageMdStr[nn+2:]

			writeNewFile(strings.Replace(levelTwoPageMdStr, "![](", "![](/images/", -1), postsDir+levelOneTitle, title+".md")

			levelThreeContent := levelTwoPage.Root().Content
			for _, block := range levelThreeContent {
				if block.IsImage() {
					// 下载文件
					filename := getImageName(block.Source)
					resp, _ := client.DownloadFile(block.Source, block)
					writeNewFile(string(resp.Data), imageDir, filename)
				}
			}
		}
	}
}

func getImageName(source string) string {
	//input:https://s3-us-west-2.amazonaws.com/secure.notion-static.com/820803a1-f8d4-4a35-95a8-e251bb961c09/wallpaper.png
	//output:wallpaper-820803a1-f8d4-4a35-95a8-e251bb961c09
	parts := strings.Split(source, "/")
	originFileName := parts[len(parts)-1]
	filenameParts := strings.Split(originFileName, ".")

	return filenameParts[0] + "-" + parts[len(parts)-2] + "." + filenameParts[1]
}

func getIndexContent(title string) string {
	return "---\ntitle: " + title + "\nbookToc: false\nbookCollapseSection: true\n\n---\nmd\n\n"
}

func writeNewFile(content string, path string, filename string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
		os.Chmod(path, 0755)
	}

	filepath := path + "/" + filename
	if err := ioutil.WriteFile(filepath, []byte(content), 0644); err != nil {
		log.Println(err.Error())
		os.Exit(128)
	}
}
