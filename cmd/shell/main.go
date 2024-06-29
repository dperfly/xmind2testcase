/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/dperfly/xmind2testcase/pingcode"
	"github.com/dperfly/xmind2testcase/version"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var CreateName string
var XMindAbPath string
var ExcelSaveFolderPath string
var Version bool

func init() {
	flag.BoolVar(&Version, "v", false, "版本")
	flag.StringVar(&CreateName, "n", "", "用例编辑者姓名")
	flag.StringVar(&XMindAbPath, "f", "", "XMind文件")
	flag.StringVar(&ExcelSaveFolderPath, "o", "", "保存路径，默认为XMind的目录")
	flag.Parse()
}

func main() {
	if Version {
		log.Println(version.Version)
		return
	}
	if CreateName == "" {
		log.Fatalln("请输入用例编辑者姓名 -n username ")
	}
	if XMindAbPath == "" {
		log.Fatalln("请输入XMind文件路径 -f filepath ")
	}
	// 默认为XMind的目录
	if ExcelSaveFolderPath == "" {
		fp, _ := filepath.Abs(XMindAbPath)
		ExcelSaveFolderPath = filepath.Dir(fp)
	}

	_, err := os.Stat(XMindAbPath)
	if err != nil {
		log.Fatalln("Xmind路径错误: ", XMindAbPath)

	}
	ExcelFileInfo, err := os.Stat(ExcelSaveFolderPath)
	if err != nil || !ExcelFileInfo.IsDir() {
		log.Fatalln("保存的文件路径错误: ", ExcelFileInfo)
	}

	excelFileName := strings.TrimSuffix(filepath.Base(XMindAbPath), ".xmind")
	err = pingcode.WriteExcel(XMindAbPath, filepath.Join(ExcelSaveFolderPath, fmt.Sprintf("%s.%s", excelFileName, "xlsx")), CreateName)
	if err != nil {
		log.Fatalln("生成失败: ", err)
	}
	log.Println(fmt.Sprintf("生成成功: %s", excelFileName))
}
