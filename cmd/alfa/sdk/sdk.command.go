package sdk

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/postlink/alfa/exporter"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Command = &cobra.Command{
	Use:   "sdk",
	Short: "生成 SDK",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd)
	},
}

func init() {
	Command.Flags().StringP("address", "a", "", "指定服务地址，如：http://localhost:8090")
	Command.Flags().StringP("lang", "l", "", "指定 SDK 语言，可选值：go、angular、axios")
	Command.Flags().StringP("target", "t", "", "指定 SDK 存放目录")
	Command.Flags().StringP("package", "p", "", "指定 SDK 包名称")
	Command.Flags().BoolP("yes", "y", false, "如果指定 target 目录不存在，是否自动创建")
}

func run(cmd *cobra.Command) (err error) {
	address, err := cmd.Flags().GetString("address")
	if err != nil {
		return
	}
	if address == "" {
		err = fmt.Errorf("请通过 --address 选项指定服务地址, ，如：--address http://localhost:8090")
		return
	}
	lang, err := cmd.Flags().GetString("lang")
	if err != nil {
		return
	}
	if lang == "" {
		err = fmt.Errorf("请通过 --lang 选项指定 SDK 语言, , 如 --lang go")
		return
	}
	target, err := cmd.Flags().GetString("target")
	if err != nil {
		return
	}
	if target == "" {
		err = fmt.Errorf("请通过 --target 选项指定 SDK 存放目录, 如 --target ./sdk")
		return
	}
	pkg, err := cmd.Flags().GetString("package")
	if err != nil {
		return
	}
	if pkg == "" {
		err = fmt.Errorf("请通过 --package 选项指定 SDK 包名称, 如 --package foo-sdk")
		return
	}
	yes, err := cmd.Flags().GetBool("yes")
	if err != nil {
		return
	}
	files, err := request(address, lang, pkg)
	if err != nil {
		return
	}
	err = askMakeTarget(target, yes)
	if err != nil {
		return
	}
	for _, v := range files {
		path := filepath.Join(target, v.Name)
		err = writeFile(path, []byte(v.Content))
		if err != nil {
			err = fmt.Errorf("文件 '%s' 写入错误: %s", v.Name, err)
			return
		}
		fmt.Printf("文件 '%s' 写入成功\n", path)
	}
	return
}

func request(address, lang, pkg string) (files []*exporter.File, err error) {
	url := fmt.Sprintf("%s/sdk?lang=%s&package=%s", address, lang, pkg)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("SDK 下载构建错误: %s", err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("SDK 下载请求错误: %s", err)
		return
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("SDK 下载读取错误: %s", err)
		return
	}
	err = json.Unmarshal(body, &files)
	if err != nil {
		err = fmt.Errorf("SDK 下载解码错误: %s", err)
		return
	}
	return
}

func askMakeTarget(target string, yes bool) (err error) {
	if dirExist(target) {
		return
	}
	if !yes {
		fmt.Printf("SDK 存放目录'%s' 不存在，是否自动创建？[Y/n]", target)
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			line := input.Text()
			fmt.Println("输入:", line)
			if strings.ToLower(line) == "y" {
				yes = true
			}
			break
		}
	}
	if yes {
		err = makeTarget(target)
		if err != nil {
			return
		}
		fmt.Printf("创建 SDK 存放目录 '%s' 成功\n", target)
	} else {
		err = fmt.Errorf("SDK 存放目录 '%s' 不存在", target)
	}
	return
}

func dirExist(path string) (exists bool) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		return
	}
	exists = true
	return
}

func makeTarget(path string) (err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			cmd := exec.Command("bash", "-c", fmt.Sprintf("mkdir -p %s", path))
			err = cmd.Run()
			if err != nil {
				err = fmt.Errorf("make dir error: %s", err)
				return
			}
		} else {
			return
		}
	}
	return
}

func writeFile(path string, content []byte) (err error) {
	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return
	}
	return
}
