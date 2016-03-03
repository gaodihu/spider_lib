package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	// "github.com/henrylee2cn/pholcus/logs"              //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	"strings"

	// 其他包
	//"fmt"
	// "math"
	// "time"
)

func init() {
	Alexa.Register()
}


var Alexa = &Spider{
	Name:        "Alexa网站排名",
	Description: "Alexa网站排名,前10万名",

	EnableCookie: false,
	RuleTree: &RuleTree{
		
		Root: func(ctx *Context) {
			for i:=0;i<200;i++ {
				//fmt.Println("http://ip-173-201-142-193.ip.secureserver.net/Alexa/Alexa_"+strconv.Itoa(i)+".html")
				ctx.AddQueue(&context.Request{Url: "http://ip-173-201-142-193.ip.secureserver.net/Alexa/Alexa_"+strconv.Itoa(i)+".html", Rule: "获取网站排名"})
			}
		},

		Trunk: map[string]*Rule{

			"获取网站排名": {
				ItemFields: []string{
					"rank",
					"site",
					
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					lis := query.Find("#Table-Uno tbody tr")
					lis.Each(func(i int, s *goquery.Selection) {
						rank := s.Find("td").Eq(0).Text()
						rank = strings.Trim(rank," \n\r")
						site := s.Find("td").Eq(1).Text()
						site = strings.Trim(site," \n\r")
						
						ctx.Output(map[int]interface{}{
							0: rank,
							1: site,
						})
						
					})
				},
			},

		},
	},
}
