package spider_lib

// 基础包
import (
	"bufio"
	//"fmt"
	"io"

	"os"
	//"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	. "github.com/henrylee2cn/pholcus/logs"
	//. "github.com/henrylee2cn/pholcus/spider/common"          //选用
	"github.com/henrylee2cn/pholcus/config"
)

func init() {
	Amazonproduct.Register()
}

var Amazonproduct = &Spider{
	Name:        "Amazonproduct",
	Description: "Amazonproduct [Auto Page] [www.Amazon.com]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Optional: &Optional{},
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
							Rule: "listpages",
						},
					)
					Log.Debug("task add")
				}

			}
			/*
				ctx.AddQueue(

					&request.Request{
						Url:  "http://www.amazon.com/IPOW/b/ref=bl_dp_s_web_8287399011?ie=UTF8&node=8287399011&field-lbr_brands_browse-bin=IPOW",
						Rule: "list",
					},
				)


					ctx.AddQueue(
						&request.Request{
							Url:  "http://www.amazon.com/Flowtron-BK-15D-Electronic-Insect-Coverage/dp/B00004R9VZ/ref=sr_1_2?s=lawn-garden&ie=UTF8&qid=1459232087&sr=1-2&keywords=bug+zappers+%7C+repellents+%7C+traps",
							Rule: "product",
						},
					)
			*/

		},

		Trunk: map[string]*Rule{
			"listpages": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					total := 0
					if total_item := query.Find(".pagnDisabled"); total_item.Size() > 0 {
						total_text := total_item.Text()
						if totaltemp, err := strconv.Atoi(total_text); err == nil {
							total = totaltemp
						}
					}
					next_url := ""
					if next_page := query.Find("#pagnNextLink"); next_page.Size() > 0 {
						if next_url_text, ok := next_page.Attr("href"); ok {
							if !strings.Contains(next_url_text, "www.amazon.com") {
								next_url = "http://www.amazon.com" + next_url_text
							} else {
								next_url = next_url_text
							}
						}
					}
					pagnCur := 0
					if pagnCur_item := query.Find(".pagnCur"); pagnCur_item.Size() > 0 {
						pagnCur_text := pagnCur_item.Text()
						if pagnCurtemp, err := strconv.Atoi(pagnCur_text); err == nil {
							pagnCur = pagnCurtemp
						}
					}
					Log.Debug("next_url:" + next_url)
					Log.Debug("total:" + strconv.Itoa(total))

					Log.Debug("pagecur:" + strconv.Itoa(pagnCur))

					if next_url != "" && total > 0 && pagnCur > 0 {
						next_page_num := pagnCur + 1
						next_page_num_text := strconv.Itoa(next_page_num)
						for ii := 1; ii <= total; ii++ {
							if pagnCur == ii {
								continue
							}
							ii_text := strconv.Itoa(ii)
							tempurl := strings.Replace(next_url, "page="+next_page_num_text, "page="+ii_text, -1)
							//ref_re, _ := regexp.Compile(`ref\=(.*)\_(\d+)\?`)
							//ref_arr := ref_re.FindAllStringSubmatch(tempurl, -1)
							//old_ref := "ref=" + ref_arr[0][1] + "_" + ref_arr[0][2]
							//new_ref := "ref=" + ref_arr[0][1] + "_" + ii_text
							//tempurl = strings.Replace(tempurl, old_ref, new_ref, -1)
							Log.Debug("tempurl:" + tempurl)
							ctx.AddQueue(&request.Request{
								Url:  tempurl,
								Rule: "list",
							},
							)
						}
					}
					ctx.Parse("list")

				},
			},

			"list": {

				ParseFunc: func(ctx *Context) {

					query := ctx.GetDom()
					//src, _ := query.Html()

					lis := query.Find("li[data-asin]")

					lis.Each(func(i int, s *goquery.Selection) {
						ASIN, _ := s.Attr("data-asin")
						FBA := 0
						if s.Find(".a-icon-prime").Size() != 0 {
							FBA = 1
						}
						item := s.Find(".s-access-detail-page").Eq(0)
						if url, ok := item.Attr("href"); ok {
							tit, _ := item.Attr("title")
							//fmt.Println(url)
							//fmt.Println(tit)
							if !strings.Contains(url, "www.amazon.com") {
								url = "http://www.amazon.com" + url
							}

							ctx.AddQueue(
								&request.Request{
									Url:  url,
									Rule: "product",
									Temp: map[string]interface{}{
										"ASIN":    ASIN,
										"type":    tit,
										"baseUrl": url,
										"FBA":     FBA,
									},
								},
							)

						}

					})

				},
			},

			"product": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"ASIN",
					"pageurl",
					"html",
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
