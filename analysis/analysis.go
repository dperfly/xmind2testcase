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

package analysis

import (
	"archive/zip"
	"errors"
	"github.com/dperfly/xmind2testcase/analysis/json"
	"github.com/dperfly/xmind2testcase/analysis/xml"
	"github.com/dperfly/xmind2testcase/types"
	"log"
)

func GetXMindTestCase(XMindPath string) ([]types.TestCase, error) {
	zipFp, err := zip.OpenReader(XMindPath)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range zipFp.File {
		// content.json 和content.xml 两个版本
		if file.FileInfo().Name() == "content.json" {
			return json.JsonXMindContent{}.GetTestCase(file)
		}
		if file.FileInfo().Name() == "content.xml" {
			return xml.XmlXMindContent{}.GetTestCase(file)
		}
	}
	return nil, errors.New("xmind中未发现content.json 或 content.xml文件,请将xmind后缀改为zip,然后查看")
}
