package spider_lib

// 基础包
import (
	
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//. "github.com/henrylee2cn/pholcus/logs"
	
	"regexp"
	"strings"
	"strconv"
	"fmt"

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

		},

		Trunk: map[string]*Rule{

			"product": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"sku",
					"name",
					"desc",
					"price",
					"shippingWeight",
					"cat1",
					"cat2",
					"cat3",
					"ShipFrom"
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					src, _ := query.Html()
					
					name := query.Find(".showInformation h1").Text()
					sku  := ""
					re,_ := regexp.Compile(`\<span\> \( Item#: (.*?) \)\<\/span\>`)
					skuArr  := re.FindAllStringSubmatch(name,"")
					if skuArr != nil && len(skuArr) > 0 {
						sku = skuArr[0][1]
					}
					
					price := query.Find("#detailPrice").Text()
					
					desc  := query.Find("#description").Html()
					
					cat1 := ""
					cat2 := ""
					cat3 := ""
					query.Find(".Bread_crumbs li a").Each(func(index int, s *goquery.Selection) {
						 cat := s.Text()
						if index == 1 {
							cat1 := cat
						}
						if index == 2 {
							cat1 := cat
						}
						if index == 3 {
							cat1 := cat
						}
					})
					shippingFrom := ""
					if shippingFromItem := query.Find(".shippingFrom .selectActive");shippingFromItem.Size()>0 {
						shippingFrom = shippingFromItem.Text()
					}
					desc_text := query.Find("#description").Text()
					weight := ""
					if strings.Contains(desc,"Package weight") {
						
						
					}else if strings.Contains(desc,"Item weight") {
						
					}
					
					
					

					

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0:  ctx.GetTemp("ASIN",""),
						1:  name,
						2:  price,
						3:  ctx.GetTemp("FBA",""),
						4:  brand,
						5:  reviews,
						6:  avgRating,
						7:  star5Rating,
						8:  star4Rating,
						9:  star3Rating,
						10: star2Rating,
						11: star1Rating,
						12: mainRank,
						13: mainRankCategory,
						14: addDate,
						15: shippingWeight,
						16: subRank,
						17: mainImg,
						18: store_name,
						19: productStatus,
					})


				},
			},
		},
	},
}
