/*
Copyright 2017 by GoSpider author.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License
*/
package jiandan

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/hunterhug/GoSpider/query"
	"github.com/hunterhug/GoSpider/util"
	"strings"
)

// 解析页面数量
func ParseIndexNum(data []byte) error {
	doc, e := query.QueryBytes(data)
	if e != nil {
		return e
	}
	s := doc.Find(".pages").Text()
	temp := strings.Split(s, "/")
	if len(temp) != 2 {
		return errors.New("index page not found")
	}
	result := strings.Replace(strings.TrimSpace(temp[1]), ",", "", -1)
	i, e := util.SI(result)
	if e != nil {
		return e
	}
	IndexPage = i
	return nil
}

// 提取信息
func ParseIndex(data []byte) []string {
	list := []string{}
	doc, e := query.QueryBytes(data)
	if e != nil {
		return list
	}
	doc.Find(".post").Each(func(num int, node *goquery.Selection) {
		//title := node.Find("h2").Text()
		//if title == "" {
		//	return
		//}
		url, _ := node.Find("h2").Find("a").Attr("href")
		//tag := node.Find(".time_s").Text()
		//if strings.Contains(tag, "·") {
		//	tag = strings.Split(tag, "·")[1]
		//}
		//fmt.Printf("%s,%s,%s\n", title, url, tag)
		list = append(list, url)
	})
	return list
}


func ParseDetail(data []byte){

}