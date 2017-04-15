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
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/store/myredis"
	"github.com/hunterhug/GoSpider/store/mysql"
	"github.com/hunterhug/GoSpider/util"
)

func init() {
	// 新建Redis池，方便爬虫们插和抽！！
	client, err := myredis.NewRedisPool(RedisConfig, DetailSpiderNum+IndexSpiderNum+2)
	if err != nil {
		spider.Log().Error(err.Error())
	}
	RedisClient = client

	// 新建数据库
	e := mysqlconfig.CreateDb()
	if e != nil {
		spider.Log().Error(e.Error())
	}
	// a new db connection
	MysqlClient = mysql.New(mysqlconfig)

	// open connection
	MysqlClient.Open(500, 300)

	// create sql
	sql := `
  CREATE TABLE IF NOT EXISTS pages (
  id varchar(255) NOT NULL,
  url varchar(255) NOT NULL,
  title varchar(255) NOT NULL,
  shortcontent varchar(255) NOT NULL DEFAULT '',
  tags varchar(255) NOT NULL DEFAULT '',
  content longtext NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='煎蛋文章';`

	// create
	_, err = MysqlClient.Create(sql)
	if err != nil {
		spider.Log().Error(err.Error())
	}

}

func SentRedis(urls []string) {
	var interfaceSlice []interface{} = make([]interface{}, len(urls))
	for i, d := range urls {
		interfaceSlice[i] = d
	}
	_, e := RedisClient.Lpush(RedisListTodo, interfaceSlice...)
	if e != nil {
		spider.Log().Errorf("sent redis error:%s", e.Error())
	}
}

func SaveToMysql(url string, m map[string]string) {
	if m["title"] == "" {
		return
	}
	_, e := MysqlClient.Insert("INSERT INTO `jiandan`.`pages`(`id`,`url`,`title`,`shortcontent`,`tags`,`content`)VALUES(?,?,?,?,?,?)", util.Md5(url), url, m["title"], m["shortcontent"], m["tags"], m["content"])
	if e != nil {
		spider.Log().Error("save mysql error:" + e.Error())
	}
}
