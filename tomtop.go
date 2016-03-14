package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//. "github.com/henrylee2cn/pholcus/logs"

	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
	//. "github.com/henrylee2cn/pholcus/spider/common"          //选用
)

func init() {
	Tomtop.Register()
}

var Tomtop = &Spider{
	Name:        "Tomtop free shipping",
	Description: "Tomtop free shipping [Auto Page] [www.tomtop.com]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Optional: &Optional{},
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			/*
							f, err := os.Open("E:/go/src/github.com/henrylee2cn/pholcus/tomtop_list.txt")
							if err != nil {
								Log.Debug("file err")
								panic(err)
							}
							defer f.Close()
							rd := bufio.NewReader(f)
							for {
								line,_, err := rd.ReadLine()
								if err != nil || io.EOF == err {
									Log.Debug("error")
				            		break
				        		}

								line_s := strings.Trim(string(line), " \r\n")
								Log.Debug(line_s)
								if line_s != "" {
									ctx.AddQueue(
										&context.Request{
											Url:  line_s,
											Rule: "list",
										},
									)
									Log.Debug("task add")
								}

							}
			*/

			//base_url := "http://www.tomtop.com/product/freeshipping?limit=20&p=510&category="
			for i := 1; i < 2; i++ {
				url := "http://www.tomtop.com/product/freeshipping?limit=20&p=" + strconv.Itoa(i) + "&category="

				ctx.AddQueue(
					&context.Request{
						Url:  url,
						Rule: "list",
					},
				)
			}

		},

		Trunk: map[string]*Rule{

			"list": {
				ParseFunc: func(ctx *Context) {

					query := ctx.GetDom()
					//src, _ := query.Html()

					lis := query.Find(".arrangeLess li a.publicTitle")

					lis.Each(func(i int, s *goquery.Selection) {

						if url, ok := s.Attr("href"); ok {

							ctx.AddQueue(
								&context.Request{
									Url:         "http://www.tomtop.com" + url,
									Rule:        "product",
									Priority:    1,
									ConnTimeout: -1,
									Temp:        map[string]interface{}{"l": "en"},
								},
							)
							//de
							ctx.AddQueue(
								&context.Request{
									Url:         "http://de.tomtop.com" + url,
									Rule:        "product",
									Priority:    1,
									ConnTimeout: -1,
									Temp:        map[string]interface{}{"l": "de"},
								},
							)
							//fr
							ctx.AddQueue(
								&context.Request{
									Url:         "http://fr.tomtop.com" + url,
									Rule:        "product",
									Priority:    1,
									ConnTimeout: -1,
									Temp:        map[string]interface{}{"l": "fr"},
								},
							)
							//es
							ctx.AddQueue(
								&context.Request{
									Url:         "http://es.tomtop.com" + url,
									Rule:        "product",
									Priority:    1,
									ConnTimeout: -1,
									Temp:        map[string]interface{}{"l": "es"},
								},
							)

							//it
							ctx.AddQueue(
								&context.Request{
									Url:         "http://it.tomtop.com" + url,
									Rule:        "product",
									Priority:    1,
									ConnTimeout: -1,
									Temp:        map[string]interface{}{"l": "it"},
								},
							)
							//

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
					"price",
					"shippingWeight",
					"cat1",
					"cat2",
					"cat3",
					"ShipFrom",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//src, _ := query.Html()
					l := ctx.GetTemp("l", "")
					lang := l

					name := query.Find(".showInformation h1").Text()
					sku := ""
					re, _ := regexp.Compile(`Item\#\:\s*(.*?)\s`)
					skuArr := re.FindAllStringSubmatch(name, -1)
					if skuArr != nil && len(skuArr) > 0 {
						sku = skuArr[0][1]
					}
					re2, _ := regexp.Compile(`\(\s*Item\#\:.*`)
					name = re2.ReplaceAllString(name, "")
					name = strings.Trim(name, " ")

					fmt.Println(sku)
					price := query.Find("#detailPrice").Text()

					desc, _ := query.Find("#description").Html()

					cat1 := ""
					cat2 := ""
					cat3 := ""
					query.Find(".Bread_crumbs li a").Each(func(index int, s *goquery.Selection) {
						cat := s.Text()
						if index == 1 {
							cat1 = cat
						}
						if index == 2 {
							cat2 = cat
						}
						if index == 3 {
							cat3 = cat
						}
					})
					shippingFrom := ""
					if shippingFromItem := query.Find(".shippingFrom .selectActive"); shippingFromItem.Size() > 0 {
						shippingFrom = shippingFromItem.Text()
					}
					desc_text := query.Find("#description").Text()

					shippingWeight := ""
					//weight := ""
					re3, _ := regexp.Compile(`(?i)Package(\s*)weight(\D*)([\d\.]+)(\s*)g`)
					re4, _ := regexp.Compile(`(?i)Item(\s*)weight(\D*)([\d\.]+)(\s*)g`)
					re5, _ := regexp.Compile(`(?i)weight(\D*)([\d\.]+)(\s*)g`)
					if weightArr := re3.FindAllStringSubmatch(desc_text, -1); weightArr != nil {
						fmt.Println("w 1")
						fmt.Println(weightArr)
						shippingWeight = weightArr[0][3]
					} else if weightArr2 := re4.FindAllStringSubmatch(desc_text, -1); weightArr2 != nil {
						fmt.Println("w 2")
						fmt.Println(weightArr2)
						shippingWeight = weightArr2[0][3]
					} else if weightArr3 := re5.FindAllStringSubmatch(desc_text, -1); weightArr3 != nil {
						fmt.Println("w 3")
						fmt.Println(weightArr3)
						shippingWeight = weightArr3[0][2]
					}
					fmt.Println(shippingWeight)

					if lang == "en" {
						query.Find(".productSmallPic img").Each(func(index int, s *goquery.Selection) {
							image_url, _ := s.Attr("src")
							image_url = strings.Replace(image_url, "/60/60/", "/2000/2000/", -1)
							image_name := sku + "/" + sku + "_" + strconv.Itoa(index) + ".jpg"
							fmt.Println(image_name)
							ctx.AddQueue(&context.Request{
								Url:         image_url,
								Rule:        "images",
								Temp:        map[string]interface{}{"image": image_name},
								Priority:    10,
								ConnTimeout: -1,
							})
						})

						query.Find("#description img").Each(func(index int, s *goquery.Selection) {
							image_url, _ := s.Attr("src")
							fmt.Println(image_url)
							fmt.Println(index)
							image_name := "sku/" + path.Base(image_url)
							//image_name := sku + "/desc_" + strconv.Itoa(index) + ".jpg"
							fmt.Println(image_name)
							ctx.AddQueue(&context.Request{
								Url:         image_url,
								Rule:        "images",
								Temp:        map[string]interface{}{"image": image_name},
								Priority:    10,
								ConnTimeout: -1,
							})
						})
					}

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: sku,
						1: lang,
						2: name,
						3: desc,
						4: price,
						5: shippingWeight,
						6: cat1,
						7: cat2,
						8: cat3,
						9: shippingFrom,
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
