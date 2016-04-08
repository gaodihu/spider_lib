package spider_lib

// 基础包
import (
	//"bufio"
	//"fmt"
	//"io"
	"io/ioutil"

	//"net/url"
	"os"
	//"regexp"
	//"strconv"
	//"strings"

	//"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//. "github.com/henrylee2cn/pholcus/logs"
	//. "github.com/henrylee2cn/pholcus/spider/common"          //选用
	//"github.com/henrylee2cn/pholcus/config"
)

func init() {
	Amazonstock.Register()
}

var Amazonstock = &Spider{
	Name:        "Amazonstock",
	Description: "Amazonstock [Auto Page] [www.Amazon.com]",

	//EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			ctx.AddQueue(
				&request.Request{
					Url:          "http://www.amazon.com/Flowtron-BK-15D-Electronic-Insect-Coverage/dp/B00004R9VZ/ref=sr_1_2?s=lawn-garden&ie=UTF8&qid=1459232087&sr=1-2&keywords=bug+zappers+%7C+repellents+%7C+traps",
					Rule:         "product",
					Method:       "POST",
					EnableCookie: true,
					PostData:     "session-id=179-5119383-1363315&ASIN=B00IYA2ZJW&offerListingID=WuR9NdrA86V1nlg9qrMRgpW4rDjOa%2FFpH1e7YOGPso%2BhCKBgAZn0CJOEeD%2BIgi9sHclkPWvP9QwEZ8X2XSIxtFzqZTxw7O9IEaJjhgmDFCJYypQIeKtX0A%3D%3D&isMerchantExclusive=0&merchantID=ATVPDKIKX0DER&isAddon=0&nodeID=172282&sellingCustomerID=A2R2RITDJNW1Q6&qid=&sr=&storeID=electronics&tagActionCode=172282&viewID=glance&rsid=179-5119383-1363315&sourceCustomerOrgListID=&sourceCustomerOrgListItemID=&wlPopCommand=&quantity=1&freeTrialBBOP=1&itemCount=4&quantity.1=1&asin.1=B0083J7G2A&quantity.2=1&asin.2=B0083J7EUO&submit.add-to-cart-prime-buy-box.x=Add to Cart&dropdown-selection=add-new&usedMerchantID=A2L77EE7U53NWQ&usedOfferListingID=2AsicI0CCs%2B7yo3TLWAM%2BjuotnUGdS%2FTNCKIHSkDA%2Bu6eaTs%2FjxnSDvUXbBTqn26vEy3oDOv6wl8Nfhlqj%2FvrUt1uxc8mRsM%2FHlBy%2F6tdbnnMkJhdKgLeftPKKkz49mV2kLfRthRsQ6iDm9PJYBgZwR7wSj55z8J&usedSellingCustomerID=A2L77EE7U53NWQ",
				},
			)

		},

		Trunk: map[string]*Rule{

			"product": {

				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					html, _ := query.Html()
					ioutil.WriteFile("t.txt", []byte(html), os.ModeAppend)
				},
			},
		},
	},
}
