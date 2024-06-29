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

package ui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/dperfly/xmind2testcase/cmd/gui/ui/chineseTheme"
	"github.com/dperfly/xmind2testcase/pingcode"
	"github.com/dperfly/xmind2testcase/version"
	"os"
	"path/filepath"
	"strings"
)

func Run() {
	app := app.New()
	app.Settings().SetTheme(chineseTheme.MyTheme{})

	MainPage := app.NewWindow(fmt.Sprintf("xmind2testcase %s ", version.Version))

	MainPage.CenterOnScreen()
	MainPage.Resize(fyne.NewSize(600, 500))
	MainPage.SetOnClosed(func() {
		os.Exit(0)
	})

	// 用例维护者姓名
	CreateName := widget.NewEntry()
	// XMind 文件路径
	XMindPath := ""
	XMindFileName := widget.NewLabel("")
	// excel 存放文件夹路径
	ExcelAbPath := widget.NewLabel("")

	MainPage.SetContent(container.NewVBox(
		container.NewGridWithColumns(2, widget.NewLabel("用例维护人"), CreateName),
		container.NewGridWithColumns(3, widget.NewLabel("XMind文件上传"), XMindFileName, widget.NewButton("点击选择", func() {
			dl := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
				if err == nil && closer != nil {
					XMindPath = closer.URI().Path()

					XMindFileName.SetText(filepath.Base(XMindPath))
				}
			}, MainPage)
			dl.SetFilter(&storage.ExtensionFileFilter{
				Extensions: []string{".xmind"},
			})

			dl.Show()
		})),
		container.NewGridWithColumns(3, widget.NewLabel("Excel存放目录"), ExcelAbPath, widget.NewButton("点击选择", func() {
			folderOpen := dialog.NewFolderOpen(func(folderURI fyne.ListableURI, err error) {
				if err != nil {
					return
				}
				if folderURI == nil {
					return
				}
				ExcelAbPath.SetText(folderURI.Path())
			}, MainPage)
			folderOpen.Show()
		})),
		container.NewGridWithColumns(1, widget.NewButton("生成", func() {
			if CreateName.Text == "" {
				dialog.NewError(errors.New("请输入用例维护人"), MainPage).Show()
				return
			}
			_, err := os.Stat(XMindPath)
			if err != nil || XMindPath == "" {
				dialog.NewError(errors.New("Xmind路径错误"), MainPage).Show()
				return
			}
			ExcelFileInfo, err := os.Stat(ExcelAbPath.Text)
			if err != nil || !ExcelFileInfo.IsDir() || ExcelAbPath.Text == "" {
				dialog.NewError(errors.New("保存的文件路径错误"), MainPage).Show()
				return

			}
			excelFileName := strings.TrimSuffix(XMindFileName.Text, ".xmind")
			err = pingcode.WriteExcel(XMindPath, filepath.Join(ExcelAbPath.Text, fmt.Sprintf("%s.%s", excelFileName, "xlsx")), CreateName.Text)
			if err != nil {
				dialog.NewError(err, MainPage).Show()
				return
			}
			dialog.NewInformation("ok", "生成成功", MainPage).Show()
		})),
	))

	MainPage.ShowAndRun()
}
