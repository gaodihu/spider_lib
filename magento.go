package spider_lib

// 基础包
import (
	//"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context/request" //必需
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
	//"strconv"
	"strings"

	// 其他包
	"fmt"
	// "math"
	// "time"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

func init() {
	Magento.Register()
}

type Site struct {
	Url          string `bson:"Url"`
	ParentUrl    string `bson:"ParentUrl"`
	DownloadTime string `bson:"DownloadTime"`
	Rank         string `bson:"rank"`
	Site         string `bson:"site"`
}

var Magento = &Spider{
	Name:        "Alexa网站前10万名，那些magento的网站",
	Description: "Alexa网站前10万名，那些magento的网站",

	EnableCookie: false,
	RuleTree: &RuleTree{

		Root: func(ctx *Context) {

			session, err := mgo.Dial("172.168.90.236:27017")
			if err != nil {
				panic(err)
			}
			defer session.Close()
			session.SetMode(mgo.Monotonic, true)
			db := session.DB("pholcus")
			mycollection := db.C("Alexa网站排名__获取网站排名")

			result := Site{}
			iter := mycollection.Find(nil).Sort("rank").Skip(20000).Limit(10000).Iter()
			for iter.Next(&result) {
				fmt.Printf("Result: %s\n", result.Site)
				geturl := "http://www." + result.Site
				//fmt.Println("Result: %s", geturl)
				ctx.AddQueue(&request.Request{Url: geturl, Rule: "ifMageto", Temp: map[string]interface{}{"Site": result.Site}})
			}

		},

		Trunk: map[string]*Rule{

			"ifMageto": {
				ItemFields: []string{
					"site",
					"if",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					src, _ := query.Html()
					if strings.Contains(src, "/skin/frontend/") {
						ctx.Output(map[int]interface{}{
							0: ctx.GetTemp("Site", ""),
							1: 1,
						})
					}

				},
			},
		},
	},
}
