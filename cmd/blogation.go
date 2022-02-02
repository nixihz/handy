package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/kjk/notionapi"
	"github.com/kjk/notionapi/tomarkdown"

	"github.com/spf13/cobra"
)

var blogationCmd = &cobra.Command{
	Use:   "blogation",
	Short: "sync notion page to hugo blog",
	Long:  `sync notion page to hugo blog`,
	Run:   publish,
}

func init() {
	rootCmd.AddCommand(blogationCmd)
}

func publish(cmd *cobra.Command, args []string) {
	fmt.Println("blogation generating")

	authToken2 := viper.GetString("authToken2")
	pageID := viper.GetString("pageID")
	postsDir := viper.GetString("postsDir")
	imageDir := viper.GetString("imageDir")

	client := &notionapi.Client{AuthToken: authToken2}
	rootPage, _ := client.DownloadPage(pageID)

	levelOneIds := rootPage.GetSubPages()

	dir, _ := ioutil.ReadDir(postsDir)
	for _, d := range dir {
		// todo 删除
		if d.IsDir() {
			os.RemoveAll(path.Join([]string{postsDir, d.Name()}...))
		}
	}

	for _, levelOneId := range levelOneIds {
		// 第一级分类
		levelOnePage, _ := client.DownloadPage(levelOneId.DashID)
		levelOneTitle := levelOnePage.Root().GetTitle()[0].Text
		println(levelOneTitle)

		indexText := getIndexContent(levelOneTitle)
		if levelOneTitle != "menu" {
			writeNewFile(indexText, postsDir+levelOneTitle, "_index.md")
		}

		levelTwoIds := levelOnePage.GetSubPages()
		for _, levelTwoId := range levelTwoIds {

			// 第二级文章
			levelTwoPage, _ := client.DownloadPage(levelTwoId.DashID)
			title := levelTwoPage.Root().GetTitle()[0].Text
			println("\t" + title)

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

	fmt.Println("blogation generated")
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
	return "---\ntitle: " + title + "\nbookToc: false\nbookCollapseSection: false\n\n---\nmd\n\n"
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
