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
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hunterhug/GoSpider/query"
	"github.com/hunterhug/GoSpider/spider"
	"github.com/hunterhug/GoSpider/util"
	"strings"
)

func init() {
	// 设置日志和超时时间
	spider.SetLogLevel("info")
	// 有些图片好大！
	spider.SetGlobalTimeout(100)
}
func main() {
	// 单只爬虫，请耐心爬取好吗
	url := "http://jandan.net/ooxx"
	urlpage := "http://jandan.net/ooxx/page-%d"

	// 保存的地方
	// 最后干脆不要了，图片太多，全部分散到hash文件中
	rootdir := "D:\\jiandan"

	// hash图片，不然图片太大了
	hashdir := "D:\\jiandanhash"
	
	// 图片太大硬盘会爆！

	// 图片集中地
	util.MakeDir(rootdir)

	// 16份
	dirs := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for _, i := range dirs {
		util.MakeDir(hashdir + "/" + i)

		// gif文件夹
		util.MakeDir(hashdir + "/gif-" + i)
	}

	// 初始化爬虫
	client, _ := spider.NewSpider(nil)
	// 随机UA
	client.SetUa(spider.RandomUa())

	// 开始抓取
	client.SetUrl(url)
	data, e := client.Go()

	// 首页都抓出错，直接结束
	if e != nil {
		spider.Log().Panic(e.Error())
	}

	// 保存在本地看看
	//util.SaveToFile(util.CurDir()+"/index.html", data)

	// 解析查看页数
	doc, _ := query.QueryBytes(data)
	temp := doc.Find(".current-comment-page").Text()
	pagenum := strings.Replace(strings.Split(temp, "]")[0], "[", "", -1)
	spider.Log().Info(pagenum)

	num, e := util.SI(pagenum)
	if e != nil {
		spider.Log().Panic(e.Error())
	}

	// 循环抓取开始
	for i := num; i > 2; i-- {
		index := fmt.Sprintf(urlpage, i)
		client.SetUrl(index)
		data, e = client.Go()
		if e != nil {
			spider.Log().Errorf("page %s error:%s", index, e.Error())
			continue
		}
		spider.Log().Infof("index %s done!", index)
		//util.SaveToFile(rootdir+"/"+util.ValidFileName(index)+".html", data)
		doc, _ = query.QueryBytes(data)
		doc.Find(".view_img_link").Each(func(num int, node *goquery.Selection) {
			imgurl, ok := node.Attr("href")
			if !ok {
				return
			}
			//spider.Log().Infof("img:%s", imgurl)

			// 去重 处理
			temp := strings.Split(imgurl, ".")
			tempnum := len(temp)
			if tempnum <= 1 {
				return
			}

			// 后缀
			houzui := temp[tempnum-1]

			// 文件名
			filename := util.Md5(imgurl) + "." + houzui

			// 图片首字母
			firstchar := string(filename[0])

			// 大本营文件路径
			filedir := rootdir + "/" + filename

			// hash分图
			var hashfiledir string
			if houzui == "gif" {
				hashfiledir = hashdir + "/gif-" + firstchar + "/" + filename
			} else {
				hashfiledir = hashdir + "/" + firstchar + "/" + filename
			}
			// 大本营存在？
			exist := util.FileExist(filedir)
			
			// hash存在？
			exist2 := util.FileExist(hashfiledir)
			
			// hash不存在，大本营存在
			if !exist2 && exist {
				if !exist2 {
					// 读出来
					temp, e := util.ReadfromFile(filedir)
					// 出错不管
					if e != nil {
						return
					}
					// 写，出错不管
					util.SaveToFile(hashfiledir, temp)
					return
				}
				// spider.Log().Infof("image file %s exist", filedir)
				return
			}
			
			// hash存在
			if exist2 {
				return
			}
			
			// 补充img url
			if strings.HasPrefix(imgurl, "//") {
				imgurl = "http:" + imgurl
			}

			// 抓取开始
			client.SetUrl(imgurl).SetRefer(index)
			data, e = client.Go()
			if e != nil {
				spider.Log().Errorf("image page %s error:%s", imgurl, e.Error())
				return
			}

			//spider.Log().Infof("image page %s done!", imgurl)

			// 大本营不再保存
			// e = util.SaveToFile(filedir, data)
			e = util.SaveToFile(hashfiledir, data)
			if e != nil {
				spider.Log().Errorf("image keep %s error:%s", hashfiledir, e.Error())
			} else {
				spider.Log().Infof("image save %s", hashfiledir)
			}

		})
	}
}
