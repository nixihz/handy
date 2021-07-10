# handy 

🧰  一个计划中的趁手工具箱 
## Contents

## Installation
```

export HANDY_WORK_DIR=.
go build handy.go
mv handy /usr/local/bin/
cp configs/prod-example.yaml configs/prod.yaml

```

### 生成数据库字典

生成markdown文件存放在 `/web/markdown`, html 文件存放`/web/html` 

```
./handy gendict --host 127.0.0.1 -P 3306 -u root -p test_pwd -d test
```

### 数据库字典在线查看服务

```
./handy servedict -p 8080

```

### 在 shell 中查看 markdown

```
handy mdview ./README.md

```

