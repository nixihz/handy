package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// mdviewCmd represents the mdview command
var wordCloudCmd = &cobra.Command{
	Use:   "wordcloud",
	Short: "wordcloud for analyse code word frequency",
	Long:  `wordcloud 检查代码单词词云`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("require project filepath and export filename")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectPath := args[0]
		files, err := GetAllFiles(projectPath)
		if err != nil {
			println(err.Error())
		}
		wordCount := make(map[string]int)
		for _, filePath := range files {
			file, err := os.Open(filePath)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			content, err := ioutil.ReadAll(file)
			if err != nil {
				panic(err)
			}

			fileType := "php"
			if strings.HasSuffix(file.Name(), ".go") {
				fileType = "go"
			}
			analyse(fileType, content, wordCount)
		}

		f, err := os.Create(args[1] + ".csv")
		f.WriteString("weight;\"word\";\"color\";\"url\"\n")

		for s, i := range wordCount {
			f.WriteString(strconv.Itoa(i) + ";\"" + s + "\";\"\";\"\"" + "\n")
		}
		f.Sync()

		fmt.Println("ok")
	},
}

func init() {
	rootCmd.AddCommand(wordCloudCmd)
}

func analyse(fileType string, content []byte, wordCount map[string]int) {
	keywords :=
		"null,bool,false,true,int,float,string,integer,php,define,this,abstract,and,as,break,callable,case,catch,class,clone,const,continue,declare,default,do,echo,else,elseif,empty,enddeclare,endfor,endforeach,endif,endswitch,endwhile,extends,final,finally,fn,for,foreach,function,global,goto,if,implements,include,include_once,instanceof,insteadof,interface,isset,list,namespace,new,or,print,private,protected,public,require,require_once,return,static,switch,throw,trait,try,unset,use,var,while,xor,yield,"
	if fileType == "go" {
		keywords = "string,int,break,default,func,interface,select,case,defer,go,map,struct,chan,else,goto,package,switch,const,fallthrough,if,range,type,continue,for,import,return,var,append,bool,byte,cap,close,complex,complex64,complex128,uint16,copy,false,float32,float64,imag,int,int8,int16,uint32,int32,int64,iota,len,make,new,nil,panic,uint64,print,println,real,recover,string,true,uint,uint8,uintp,"
	}

	//解析正则表达式，如果成功返回解释器
	reg1 := regexp.MustCompile(`([a-zA-Z]+)`)
	if reg1 == nil {
		fmt.Println("regexp err")
		return
	}
	//根据规则提取关键信息
	result1 := reg1.FindAll(content, -1)
	for _, w := range result1 {
		word := string(w)
		//word = strings.ToLower(word)
		if strings.Contains(keywords, word+",") {
			continue
		}
		wordCount[word]++
	}
}

//获取指定目录下的所有文件和目录
func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".php") || strings.HasSuffix(fi.Name(), ".go")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	return files, dirs, nil
}

//获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			if fi.Name() == "vendor" || fi.Name() == ".git" || fi.Name() == "tests" {
				continue
			}
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".php") || strings.HasSuffix(fi.Name(), ".go")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}
