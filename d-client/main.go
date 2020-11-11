package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {

	var name, host, simplify, lang string
	var port int
	required := []string{"n"}

	flag.StringVar(&name, "n", "ddz", "连接用户名")
	flag.StringVar(&simplify, "simplify", "true", "简化显示")
	flag.StringVar(&lang, "lang", "zh", "语言")
	flag.StringVar(&host, "host", "localhost", "服务器地址")
	flag.IntVar(&port, "port", 8080, "服务器端口")
	flag.Parse()

	seen := make(map[string]bool)
	valid := true
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			valid = false
			fmt.Printf("缺少必要参数 -%s\n", req)
		}
	}
	if !valid {
		fmt.Printf("使用 -h 查看帮助\n")
		os.Exit(-1)
	}

	simplifyBool, err := strconv.ParseBool(simplify)
	if err != nil {
		simplifyBool = false
	}

	dc := NewDdzClient(name, host, port, simplifyBool, lang)
	dc.Run()
}
