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
	"fmt"
	"github.com/hunterhug/GoSpider/util"
	"path/filepath"
	//"strings"
	"github.com/hunterhug/GoSpider/query"
	"testing"
)

func TestInit(t *testing.T) {
	e := RedisClient.Set("test", "1", 0)
	if e != nil {
		t.Error(e.Error())
	}
}
func TestParseIndexNum(t *testing.T) {
	data, e := util.ReadfromFile(filepath.Join(util.CurDir(), "data", "index.html"))
	if e != nil {
		t.Error(e.Error())
	} else {
		e = ParseIndexNum(data)
		if e != nil {
			t.Error(e.Error())
		} else {
			t.Log(IndexPage)
		}
	}

}

func TestIndexStep(t *testing.T) {

	IndexStep()
	urllist := []string{}
	for i := 2; i <= IndexPage; i++ {
		urllist = append(urllist, fmt.Sprintf("%s/page/%d", Url, i))
	}
	// 分配任务
	tasks, e := util.DevideStringList(urllist, IndexSpiderNum)
	if e != nil {
		t.Error(e.Error())
	}
	t.Logf("%#v", tasks)
}

func TestIndexSpiderRun(t *testing.T) {
	IndexSpiderRun()
}

func TestParseIndex(t *testing.T) {
	data, e := util.ReadfromFile(filepath.Join(RootDir, "data", "index.html"))
	//data, e := util.ReadfromFile(filepath.Join(RootDir, "data", "http###jandan.net#page#2.html"))
	if e != nil {
		t.Error(e.Error())
	}
	result := ParseIndex(data)
	t.Logf("%#v", result)
}

func TestRedis(t *testing.T) {
	RedisClient.Hset(RedisListDone, "ddd", "ddd")
	s, e := RedisClient.Hget(RedisListDone, "dsdd")
	if e.Error() == "redis: nil" {
		t.Error(e.Error())
	}
	t.Log(s)

	ok, e := RedisClient.Hexists(RedisListDone+"dd", "dddddd")
	t.Logf("%v,%v", ok, e)
}

func TestParseDetail(t *testing.T) {
	data, e := util.ReadfromFile(filepath.Join(RootDir, "data", "detail", "http###jandan.net#2011#02#03#skin-gun.html"))
	if e != nil {
		t.Error(e.Error())
	}
	doc, e := query.QueryBytes(data)
	if e != nil {
		t.Error(e.Error())
		return
	}
	// 标题
	title :=doc.Find("title").Text()
	// 标签
	tag :=doc.Find("#content").Find("h3 a").Text()

	result:= doc.Find("#content").Find(".post").Nodes
	t.Logf("%v,%v,%v", title,tag,result)
}
