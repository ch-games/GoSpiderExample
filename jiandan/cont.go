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
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/store/myredis"
	"github.com/hunterhug/GoSpider/store/mysql"
	"github.com/hunterhug/GoSpider/util"
	"path/filepath"
)

// 可抽离到配置文件中
const (
	// 网站
	Url  = "http://jandan.net"
	Host = "jandan.net"

	// 详情页爬虫数量
	DetailSpiderNum        = 30
	DetailSpiderNamePrefix = "detail"
	// 首页爬虫数量
	IndexSpiderNum        = 3
	IndexSpiderNamePrefix = "index"

	// 爬虫暂停时间
	StopTime = 1
	// 日志级别
	LogLevel = "info"
)

var (
	// 首页页数
	IndexPage int

	// 根目录
	//RootDir = util.CurDir()
	RootDir = "E:\\jiandan"

	// Redis配置
	RedisConfig = myredis.RedisConfig{
		DB:       0,
		Host:     "127.0.0.1:6379",
		Password: "smart2016",
	}

	RedisClient myredis.MyRedis

	RedisListTodo  = "jiandantodo"
	RedisListDoing = "jiandandoing"
	RedisListDone  = "jiandandone"

	// mysql config
	mysqlconfig = mysql.MysqlConfig{
		Username: "root",
		Password: "smart2016",
		Ip:       "127.0.0.1",
		Port:     "3306",
		Dbname:   "jiandan",
	}

	MysqlClient mysql.Mysql
)

// 设置全局
func init() {
	e := util.MakeDir(filepath.Join(RootDir, "data", "detail"))
	if e != nil {
		spider.Log().Panic(e.Error())
	}
	spider.SetGlobalTimeout(StopTime)
	spider.SetLogLevel(LogLevel)
	indexstopchan = make(chan bool, 1)

	// 初始化爬虫们，一种多爬虫方式，设置到全局MAP中
	for i := 0; i <= IndexSpiderNum; i++ {
		s, e := spider.New(nil)
		if e != nil {
			spider.Log().Panicf("index spider %d new error: %s", i, e.Error())
		}
		// 设置随机UA
		s.SetUa(spider.RandomUa())
		spider.Pool.Set(fmt.Sprintf("%s-%d", IndexSpiderNamePrefix, i), s)
	}
	for i := 0; i <= DetailSpiderNum; i++ {
		s, e := spider.New(nil)
		if e != nil {
			spider.Log().Panicf("detail spider %d new error: %s", i, e.Error())
		}
		s.SetUa(spider.RandomUa())
		spider.Pool.Set(fmt.Sprintf("%s-%d", DetailSpiderNamePrefix, i), s)
	}
}
