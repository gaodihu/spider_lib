package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
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
	// "fmt"
	// "math"
	// "time"
)

func init() {
	Proxyip.Register()
}

var Proxyip = &Spider{
	Name:        "每日代理ip地址",
	Description: "每日代理ip地址",

	EnableCookie: false,
	RuleTree: &RuleTree{

		Root: func(ctx *Context) {

			for i := 0; i < 10; i++ {
				ctx.AddQueue(&context.Request{Url: "http://www.ip181.com/daili/" + strconv.Itoa(i) + ".html", Rule: "每日代理ip地址"})
			}

		},

		Trunk: map[string]*Rule{

			"每日代理ip地址": {
				ItemFields: []string{
					"Ip",
					"Port",
					"Anonymous",
					"Type",
					"Response",
					"Address",
					"Url",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					lis := query.Find(".col-md-12 table tbody tr")
					lis.Each(func(i int, s *goquery.Selection) {
						ip := s.Find("td").Eq(0).Text()
						ip = strings.Trim(ip, " \n\r")

						port := s.Find("td").Eq(1).Text()
						port = strings.Trim(port, " \n\r")

						anonymous := s.Find("td").Eq(2).Text()
						anonymous = strings.Trim(anonymous, " \n\r")

						proxy_type := s.Find("td").Eq(3).Text()
						proxy_type = strings.Trim(proxy_type, " \n\r")

						response := s.Find("td").Eq(4).Text()
						response = strings.Trim(response, " \n\r")

						address := s.Find("td").Eq(5).Text()
						address = strings.Trim(address, " \n\r")

						url := "http://" + ip + ":" + port

						ctx.Output(map[int]interface{}{
							0: ip,
							1: port,
							2: anonymous,
							3: proxy_type,
							4: response,
							5: address,
							6: url,
						})

					})
				},
			},
		},
	},
}
