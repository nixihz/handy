package cmd

import (
	"github.com/spf13/cobra"
	"net/http"
	"strconv"
)
import "github.com/gin-gonic/gin"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "http serve your web directory to users.",
	Long:  `http serve your web directory to users. 启动文件服务器(映射web文件夹)`,
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()
		router.StaticFS("/", http.Dir("./web"))
		var port, _ = cmd.Flags().GetInt("port")

		router.Run(":" + strconv.Itoa(port))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	var port int
	serveCmd.Flags().IntVarP(&port, "port", "t", 8080, "serve port")
}
