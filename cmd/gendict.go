package cmd

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/russross/blackfriday"
	"github.com/spf13/cobra"
	"handy/models"
	"io/ioutil"
	"log"
	"os"
)

// gendictCmd represents the gendict command
var gendictCmd = &cobra.Command{
	Use:   "gendict",
	Short: "generate mysql table dictionary",
	Long:  `generate mysql table dictionary 生成数据库字典`,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(gendictCmd)

	var host string
	var port string
	var username string
	var password string
	var dbname string

	gendictCmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "host")
	gendictCmd.Flags().StringVarP(&port, "port", "P", "3306", "port")
	gendictCmd.Flags().StringVarP(&username, "username", "u", "root", "username")
	gendictCmd.Flags().StringVarP(&password, "password", "p", "", "password")
	gendictCmd.Flags().StringVarP(&dbname, "dbname", "d", "", "dbname")

}

func run(cmd *cobra.Command, args []string) {
	host, err := cmd.Flags().GetString("host")
	port, err := cmd.Flags().GetString("port")
	username, err := cmd.Flags().GetString("username")
	password, err := cmd.Flags().GetString("password")
	dbname, err := cmd.Flags().GetString("dbname")

	connectString := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		print(err.Error())
		return
	}
	var tableList []models.Table

	db.Exec("SET NAMES UTF8")
	db.Raw("SELECT TABLE_NAME as table_name,TABLE_COMMENT as table_comment FROM INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA = '" + dbname + "';").Scan(&tableList)

	var out = "[TOC]\n"
	for i := range tableList {
		tableComment := tableList[i].TableComment
		tableName := tableList[i].TableName
		var columns []models.Column

		out += "## 表名：" + tableName + " " + tableComment + "\n"
		out += "|字段名|数据类型|默认值|允许非空|自动递增|备注|\n"
		out += "| ---- | ---- | ---- | ---- | ---- | ---- |\n"

		db.Raw("SELECT " +
			"COLUMN_NAME column_name" +
			",COLUMN_TYPE column_type" +
			",COLUMN_DEFAULT column_default" +
			",IS_NULLABLE is_nullable" +
			",EXTRA extra" +
			",COLUMN_COMMENT column_comment" +
			" FROM INFORMATION_SCHEMA.COLUMNS" +
			" where table_name = '" + tableName + "'" +
			" and table_schema = '" + dbname + "'" +
			"").Scan(&columns)

		for i2 := range columns {
			column := columns[i2]

			out += "|" + column.ColumnName
			out += "|" + column.ColumnType
			out += "|" + column.ColumnDefault
			out += "|" + column.IsNullable
			if column.Extra == "auto_increment" {
				out += "|" + "是"
			} else {
				out += "|" + " "
			}
			out += "|" + column.ColumnComment + "|\n"
		}
	}

	html := string(blackfriday.MarkdownCommon([]byte(out)))
	html = "<link rel=\"stylesheet\" href=\"/assets/style.css\"/>" + html

	writeFile(html, "./web/html/"+dbname+".html")
	writeFile(out, "./web/markdown/"+dbname+".md")
	println("generate md file: ./web/markdown/" + dbname + ".md")
	println("generate html file: ./web/html/" + dbname + ".html")
}

func writeFile(content string, filePath string) {
	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Println(err.Error())
		os.Exit(128)
	}
}
