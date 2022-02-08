package cmd

import (
	"fmt"
	"github.com/fencex/handy/pkg/notionservice/tomarkdown"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/kjk/notionapi"

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

	var removeAll string
	blogationCmd.Flags().StringVarP(&removeAll, "removeAll", "r", "0", "0 false|1 true")

}

func publish(cmd *cobra.Command, args []string) {
	removeAll, _ := cmd.Flags().GetString("removeAll")

	fmt.Println("blogation generating")

	authToken2 := viper.GetString("authToken2")
	pageID := viper.GetString("pageID")
	postsDir := viper.GetString("postsDir")
	imageDir := viper.GetString("imageDir")

	client := &notionapi.Client{AuthToken: authToken2}
	rootPage, _ := client.DownloadPage(pageID)

	levelOneIds := rootPage.GetSubPages()

	if removeAll == "1" {
		println("清空数据重新生成!")

		dir, _ := ioutil.ReadDir(postsDir)
		for _, d := range dir {
			// todo 删除
			if d.IsDir() {
				os.RemoveAll(path.Join([]string{postsDir, d.Name()}...))
			}
		}
	}

	// 不清理所有，则只是更新最近1天的文章
	compareTime := time.Now().AddDate(0, 0, -1).Unix() * 1000
	if removeAll == "1" {
		// 清理所有
		compareTime = 0
	}
	for _, levelOneId := range levelOneIds {
		// 第一级分类
		levelOnePage, _ := client.DownloadPage(levelOneId.DashID)
		levelOneTitle := levelOnePage.Root().GetTitle()[0].Text
		if levelOnePage.Root().LastEditedTime < compareTime {
			println(levelOneTitle + " not modified!")
		} else {
			println(levelOneTitle)

			indexText := getIndexContent(levelOneTitle)
			if levelOneTitle != "menu" {
				writeNewFile(indexText, postsDir+levelOneTitle, "_index.md")
			}
		}

		levelTwoIds := levelOnePage.GetSubPages()
		for _, levelTwoId := range levelTwoIds {

			// 第二级文章
			levelTwoPage, _ := client.DownloadPage(levelTwoId.DashID)
			title := levelTwoPage.Root().GetTitle()[0].Text
			if levelTwoPage.Root().LastEditedTime < compareTime {
				println("\t" + title + " not modified!")
				continue
			} else {
				println("\t" + title)
			}

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

/**
 * 获取文件名
 */
func getImageName(source string) string {
	//input:https://s3-us-west-2.amazonaws.com/secure.notion-static.com/820803a1-f8d4-4a35-95a8-e251bb961c09/wallpaper.png
	//output:wallpaper-820803a1-f8d4-4a35-95a8-e251bb961c09
	parts := strings.Split(source, "/")
	originFileName := parts[len(parts)-1]
	filenameParts := strings.Split(originFileName, ".")

	return filenameParts[0] + "-" + parts[len(parts)-2] + "." + filenameParts[1]
}

/**
 * 统一的 index 的内容
 */
func getIndexContent(title string) string {
	return "---\ntitle: " + title + "\nbookToc: false\nbookCollapseSection: false\n\n---\nmd\n\n"
}

/**
 * 写入文件
 */
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
