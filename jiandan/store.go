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
)

func init() {
	// 新建Redis池，方便爬虫们插和抽！！
	client, err := myredis.NewRedisPool(RedisConfig, DetailSpiderNum+IndexSpiderNum+2)
	if err != nil {
		spider.Log().Error(err.Error())
	}
	RedisClient = client
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
