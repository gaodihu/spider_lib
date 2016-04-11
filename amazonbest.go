package spider_lib

// 基础包
import (
	"net/url"

	"strings"

	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//. "github.com/henrylee2cn/pholcus/logs"
	//. "github.com/henrylee2cn/pholcus/spider/common"          //选用
	//"github.com/henrylee2cn/pholcus/config"
)

func init() {
	Amazonbest.Register()
}

var Amazonbest = &Spider{
	Name:         "Amazonbest",
	Description:  "Amazonbest [Auto Page] [www.Amazon.com]",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			ctx.AddQueue(
				&request.Request{
					Url:  "http://www.amazon.com/Best-Sellers/zgbs/ref=zg_bs_unv_cps_0_cps_1",
					Rule: "listpages",
				},
			)
		},

		Trunk: map[string]*Rule{
			"listpages": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					//分类页
					c := query.Find("#zg_browseRoot  .zg_browseUp a")
					c.Remove()

					query.Find(" #zg_browseRoot  a").Each(func(i int, s *goquery.Selection) {

						if next_url, ok := s.Attr("href"); ok {

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

					})

				},
			},

			"list": {

				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					lis := query.Find("#zg_centerListWrapper .zg_itemImmersion .zg_title a")

					lis.Each(func(i int, s *goquery.Selection) {

						if productlink, ok := s.Attr("href"); ok {
							productlink = strings.Trim(productlink, " \n\r")
							if !strings.Contains(productlink, "www.amazon.com") {
								productlink = "http://www.amazon.com" + productlink
							}

							ctx.AddQueue(
								&request.Request{
									Url:  productlink,
									Rule: "product",
								},
							)

						}

					})

					pages := query.Find(".zg_pagination a")
					pages.Each(func(i int, s *goquery.Selection) {
						if next_url, ok := s.Attr("href"); ok {
							if css, ok2 := s.Attr("class"); ok2 {
								next_url = strings.Trim(next_url, " \n\r")
								if !strings.Contains(css, "zg_selected") {
									if !strings.Contains(next_url, "www.amazon.com") {
										next_url = "http://www.amazon.com" + next_url
									}
									ctx.AddQueue(
										&request.Request{
											Url:  next_url,
											Rule: "list",
										},
									)
								}
							}

						}
					})

				},
			},

			"product": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"ASIN",
					"Pageurl",
					"Html",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

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

					}

				},
			},
		},
	},
}
