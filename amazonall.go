package spider_lib

// 基础包
import (
	"bufio"

	"io"
	"net/url"
	"os"

	"strings"

	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/config"
	. "github.com/henrylee2cn/pholcus/logs"
	//. "github.com/henrylee2cn/pholcus/spider/common"          //选用
)

func init() {
	AmazonAll.Register()
}

var AmazonAll = &Spider{
	Name:         "Amazon All",
	Description:  "Amazon All商品数据 [Auto Page] [www.Amazon.com]",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			f, err := os.Open(config.WORK_ROOT + "/amazon_list.txt")
			if err != nil {
				Log.Debug("file err")
				panic(err)
			}
			defer f.Close()
			rd := bufio.NewReader(f)
			for {
				line, _, err := rd.ReadLine()
				if err != nil || io.EOF == err {
					Log.Debug("error")
					break
				}

				line_s := strings.Trim(string(line), " \r\n")
				Log.Debug(line_s)
				if line_s != "" {
					ctx.AddQueue(
						&request.Request{
							Url:  line_s,
							Rule: "list",
						},
					)
					Log.Debug("task add")
				}

			}

			/*
				ctx.AddQueue(
					&request.Request{
						Url:  "http://www.amazon.com/gp/site-directory/ref=nav_shopall_btn",
						Rule: "list",
					},
				)
			*/

		},

		Trunk: map[string]*Rule{

			"list": {
				ItemFields: []string{
					"ASIN",
					"pageurl",
					"html",
				},
				ParseFunc: func(ctx *Context) {

					query := ctx.GetDom()
					query.Find("header").Remove()
					query.Find("#footer").Remove()
					query.Find("script").Remove()
					query.Find("#navFooter").Remove()
					pageurl := ctx.GetUrl()

					is_product := false
					ASIN := ""
					ASIN_item := query.Find("#ASIN")
					if ASIN_item.Size() > 0 {
						ASIN_text, ok := ASIN_item.Attr("value")
						if ok {
							is_product = true
							ASIN = strings.Trim(ASIN_text, " ")
						}
					}

					if is_product {
						src, _ := query.Html()

						ctx.Output(map[int]interface{}{
							0: ASIN,
							1: pageurl,
							2: src,
						})

					} else {
						query.Find("a").Each(func(i int, s *goquery.Selection) {

							if next_url, ok := s.Attr("href"); ok {
								if strings.Contains(next_url, "/b/") || strings.Contains(next_url, "/s/") || strings.Contains(next_url, "/dp/") {
									url_arr, err := url.Parse(next_url)
									if err == nil {
										if url_arr.Host == "" {
											ctx.AddQueue(&request.Request{
												Url:  "http://www.amazon.com" + next_url,
												Rule: "list",
											},
											)

										} else if strings.Contains(url_arr.Host, "www.amazon.com") {
											ctx.AddQueue(&request.Request{
												Url:  next_url,
												Rule: "list",
											},
											)
										}
									}
								}
							}

						})
					}

				},
			},
		},
	},
}
