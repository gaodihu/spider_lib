package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//"github.com/henrylee2cn/pholcus/logs"                   //信息输出
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用
	//"github.com/henrylee2cn/pholcus/app"

	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

const HOST = "http://www.geekbuying.com/"

func init() {
	GeekBuying.Register()
}

var GeekBuying = &Spider{
	Name:        "GeekBuying free shipping",
	Description: "GeekBuying free shipping [www.geekbuying.com]",

	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			f, err := os.Open("D:/go/geekbuying_list.txt")
			if err != nil {

				panic(err)
			}
			defer f.Close()
			rd := bufio.NewReader(f)
			for {
				line, _, err := rd.ReadLine()
				if err != nil || io.EOF == err {

					break
				}

				line_s := strings.Trim(string(line), " \r\n")

				if line_s != "" {
					ctx.AddQueue(
						&request.Request{
							Url:  line_s,
							Rule: "list",
							Temp: map[string]interface{}{"p": 1},
						},
					)

				}

			}
		},

		Trunk: map[string]*Rule{

			"list": {
				ParseFunc: func(ctx *Context) {
					var curr int
					ctx.GetTemp("p", &curr)
					current_page := ctx.GetDom().Find("#pagination .current").Text()
					current_page = strings.Trim(current_page, " \n\r")
					if current_page != strconv.Itoa(curr) {
						fmt.Println(current_page)
						fmt.Println(strconv.Itoa(curr))
						fmt.Println("p=" + strconv.Itoa(curr))
						return
					}
					prefixUrl := ctx.GetUrl()
					toUrl := ""
					if curr == 1 {
						toUrl = prefixUrl + strconv.Itoa(curr+1) + "-40-3-0-0-0-grid.html"
					} else {
						toUrl = strings.Replace(prefixUrl, strconv.Itoa(curr)+"-40-3-0-0-0-grid.html", strconv.Itoa(curr+1)+"-40-3-0-0-0-grid.html", -1)
					}

					ctx.AddQueue(&request.Request{
						Url:         toUrl,
						Rule:        "list",
						Temp:        map[string]interface{}{"p": curr + 1},
						ConnTimeout: -1,
					})

					ctx.GetDom().Find(".gridView .name a").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							if !strings.Contains(url, HOST) {
								url = HOST + url
							}

							ctx.AddQueue(
								&request.Request{
									Url:  url,
									Rule: "product",
									Temp: map[string]interface{}{"position": curr*40 + i},
								},
							)
						}
					})

				},
			},

			"product": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"sku",
					"lang",
					"name",
					"desc",
					"currency",
					"price",
					"cat1",
					"cat2",
					"cat3",
					"freeshipping",
					"position",
					"attrProperty",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					name := query.Find("#productName").Text()

					lang := "en"

					sku := query.Find("#iconCodeDiv1 b").Text()

					currency := query.Find("#currencyF").Text()

					price := query.Find("#saleprice").Text()

					desc, _ := query.Find("#DESCRIPTION_HTML").Html()

					cat1 := ""
					cat2 := ""
					cat3 := ""
					query.Find("#crumbs li").Each(func(index int, s *goquery.Selection) {
						if index == 1 {
							cat1 = s.Find(".jb_farther a").Text()
						}
						if index == 2 {
							cat2 = s.Find("a").Text()
						}
						if index == 3 {
							cat3 = s.Find("a").Text()
						}
					})

					freeshipping := 0

					if query.Find(".stockFree").Size() > 0 {
						freeshipping = 1
					}

					position := 0
					ctx.GetTemp("position", &position)

					attrProperty := make(map[string]string, 10)

					query.Find(".attrProperty li a").Each(func(index int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							attrProperty[url] = s.Text()
						}
					})
					fmt.Println(attrProperty)
					attrPropertyJson, _ := json.Marshal(attrProperty)

					//image list

					query.Find("#thumbnail img").Each(func(index int, s *goquery.Selection) {
						if src, ok := s.Attr("src"); ok {
							bigSrc := strings.Replace(src, "make_pic", "ggo_pic", -1)
							bigSrc = strings.Replace(src, "1.jpg", ".jpg", -1)
							imageName := sku + "/" + path.Base(bigSrc)
							ctx.AddQueue(&request.Request{
								Url:         bigSrc,
								Rule:        "images",
								Temp:        map[string]interface{}{"image": imageName},
								Priority:    1,
								ConnTimeout: -1,
							})
						}
					})
					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0:  sku,
						1:  lang,
						2:  name,
						3:  desc,
						4:  currency,
						5:  price,
						6:  cat1,
						7:  cat2,
						8:  cat3,
						9:  freeshipping,
						10: position,
						11: string(attrPropertyJson),
					})

				},
			},
			"images": {
				ParseFunc: func(ctx *Context) {
					// 文件输出方式一（推荐）
					ctx.FileOutput(ctx.GetTemp("image", "").(string))

				},
			},
		},
	},
}
