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

package pingcode

import (
	"errors"
	"fmt"
	"github.com/dperfly/xmind2testcase/analysis"
	"github.com/tealeg/xlsx"
	"strings"
)

func WriteExcel(XMindPath string, SaveFileName string, CreateUserName string) error {
	sheetName := "test_case"
	f := xlsx.NewFile()

	sheet, err := f.AddSheet(sheetName)
	if err != nil {
		return err
	}
	// 需要加个空行，才能满足pingCode的格式
	sheet.AddRow()
	row1 := sheet.AddRow()
	// 标题行
	firstData := []string{"模块", "编号", "*标题", "维护人", "用例类型", "重要程度", "测试类型", "预估工时", "剩余工时", "关联工作项", "前置条件", "步骤描述", "预期结果", "关注人", "备注", "迭代版本"}
	for _, v := range firstData {
		cell := row1.AddCell()
		cell.Value = v
	}
	cases, err := analysis.GetXMindTestCase(XMindPath)
	if err != nil {
		return err
	}
	for _, c := range cases {
		steps := strings.Builder{}
		exp := strings.Builder{}
		for i, v := range c.TestSteps.Steps {
			steps.WriteString(fmt.Sprintf("%d.%s\n", i+1, v.Step))
			exp.WriteString(fmt.Sprintf("%d.%s\n", i+1, v.Expected))
		}

		thisCase := []string{c.Module, c.Number, c.Title, CreateUserName, c.TestCaseType, c.Priority, c.TestRunType, "", "", c.Story, c.Before, steps.String(), exp.String(), "", "", ""}
		row := sheet.AddRow()
		for _, v := range thisCase {
			cell := row.AddCell()
			cell.Value = v
		}
	}

	if err = f.Save(SaveFileName); err != nil {
		return errors.New("请关掉之前生成的excel文件，然后重新运行")
	}

	return nil
}
