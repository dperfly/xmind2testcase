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

package json

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/dperfly/xmind2testcase/types"
	"io"
	"strings"
)

// JsonXMindContent json XMind 10
type JsonXMindContent []struct {
	Title     string `json:"title,omitempty" `
	RootTopic Topic  `json:"rootTopic,omitempty" `
}

type Topic struct {
	Title    string `json:"title,omitempty"`
	Children struct {
		Attached []Topic `json:"attached,omitempty"`
	} `json:"children,omitempty" xml:"children"`
}

func (x JsonXMindContent) GetTestCase(file *zip.File) ([]types.TestCase, error) {
	testCases := make([]types.TestCase, 0)
	// 用于解析测试用例的后半部分
	var dfsTestStep func(i int, step *types.Step, topic Topic)
	dfsTestStep = func(i int, step *types.Step, topic Topic) {
		if i == 0 {
			step.Step = topic.Title
		}
		if i == 1 {
			step.Expected = topic.Title
		}
		if i > 1 {
			return
		}
		for _, c := range topic.Children.Attached {
			dfsTestStep(1, step, c)
		}
	}
	// 用于解析用例的title 和前置条件等信息
	var dfs func(c types.TestCase, index int, children Topic, isStep bool, isExp bool)
	dfs = func(c types.TestCase, index int, children Topic, isStep bool, isExp bool) {
		// index 用于标记xmind 的深度，深度为1时表示是模块+所属需求
		if index == 1 {
			if strings.Contains(children.Title, "(") {
				ss := strings.Split(children.Title, "(")
				// 没有所属需求
				if len(ss) > 0 {
					c.Module = c.Module + "/" + ss[0]
				}
				if len(ss) > 1 {
					// 需要去掉最后的")" 或者“）”
					c.Story = strings.Split(ss[1], ")")[0]
				}
			} else if strings.Contains(children.Title, "（") {
				ss := strings.Split(children.Title, "（")
				// 没有所属需求
				if len(ss) > 0 {
					c.Module = c.Module + "/" + ss[0]
				}
				if len(ss) > 1 {
					// 需要去掉最后的")" 或者“）”
					c.Story = strings.Split(ss[1], "）")[0]
				}
			} else {
				c.Module = c.Module + "/" + children.Title
			}

		} else {

			if strings.Contains(children.Title, "验证") {

				if strings.Contains(children.Title, "p1") || strings.Contains(children.Title, "P1") {
					c.Priority = "P1"
				}
				if strings.Contains(children.Title, "p2") || strings.Contains(children.Title, "P2") {
					c.Priority = "P2"
				}
				if strings.Contains(children.Title, "p3") || strings.Contains(children.Title, "P3") {
					c.Priority = "P3"
				}
				if strings.Contains(children.Title, "p4") || strings.Contains(children.Title, "P4") {
					c.Priority = "P4"
				}
				if strings.Contains(children.Title, "(") {
					c.Title += strings.Split(children.Title, "(")[0]
				} else if strings.Contains(children.Title, "（") {
					c.Title += strings.Split(children.Title, "（")[0]
				} else {
					c.Title += children.Title
				}
				// 进行测试步骤设置
				for _, children := range children.Children.Attached {
					s := &types.Step{}
					dfsTestStep(0, s, children)
					c.TestSteps.Steps = append(c.TestSteps.Steps, s)

				}
				testCases = append(testCases, c)

			} else {
				// 前置步骤的处理
				if strings.Contains(children.Title, "（") {
					before := strings.Split(children.Title, "（")
					c.Before += strings.Split(before[1], "）")[0] + "\n"
					c.Title += fmt.Sprint("【" + before[0] + "】")
				} else if strings.Contains(children.Title, "(") {
					before := strings.Split(children.Title, "(")
					c.Before += strings.Split(before[1], ")")[0] + "\n"
					c.Title += fmt.Sprint("【" + before[0] + "】")

				} else { // title 的 【】 部分
					c.Title += fmt.Sprint("【" + children.Title + "】")
				}
			}
		}
		for _, cld := range children.Children.Attached {
			dfs(c, index+1, cld, isStep, isExp)
		}

	}

	contentFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	contentByte, err := io.ReadAll(contentFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(contentByte, &x)
	if err != nil {
		return nil, err
	}

	for _, sheet := range x {
		for _, children := range sheet.RootTopic.Children.Attached {
			dfs(types.TestCase{
				TestCaseType: "功能测试",
				TestRunType:  "手动",
				Priority:     "P2",
				Module:       sheet.RootTopic.Title,
				TestSteps: types.TestStep{
					Steps: make([]*types.Step, 0),
				},
			}, 1, children, false, false)
		}
	}
	return testCases, nil
}