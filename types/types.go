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

package types

import "archive/zip"

type TestCaseInterface interface {
	GetTestCase(file *zip.File) ([]TestCase, error)
}

// TestCase 测试用例Excel表头
type TestCase struct {
	Module       string   `json:"module"`         // 模块
	Number       string   `json:"number"`         // 编号
	Title        string   `json:"title"`          //标题
	CreateUser   string   `json:"create_user"`    // 维护人 传入
	TestCaseType string   `json:"test_case_type"` // 用例类型
	Priority     string   `json:"priority"`       // 优先级
	TestRunType  string   `json:"test_run_type"`  // 测试类型，默认手动
	Story        string   `json:"story"`          // 关联工作项
	Before       string   `json:"before"`         //前置条件
	TestSteps    TestStep `json:"test_steps"`     // 测试步骤
	Follower     string   `json:"follower"`       // 关注人
	Comment      string   `json:"comment"`        // 备注
	Version      string   `json:"version"`        // 迭代版本
}

type TestStep struct {
	Steps []*Step `json:"steps"`
}

type Step struct {
	Step     string `json:"step"`     // 测试步骤
	Expected string `json:"expected"` //期望结果
}
