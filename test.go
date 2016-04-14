package main

import (
	"fmt"
	//"os"
	//"os/exec"
	//"regexp"
	//"strconv"
	//"strings"

	"github.com/PuerkitoBio/goquery"
	//"math"
	"bytes"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	Id           bson.ObjectId `bson:"_id"`
	ASIN         string        `bson:"ASIN"`
	Pageurl      string        `bson:"Pageurl"`
	Html         string        `bson:"Html"`
	Url          string        `bson:"Url"`
	ParentUrl    string        `bson:"ParentUrl"`
	DownloadTime string        `bson:"DownlondTime"`
}

const URL = "172.168.90.236:27017"

var (
	mgoSession *mgo.Session
	dataBase   = "pholcus"
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic(err)
		}
	}
	return mgoSession.Clone()
}
func ParseProduct(html string) {
	f := bytes.NewBufferString(html)
	query, error := goquery.NewDocumentFromReader(f)
	if error != nil {
		panic(error)
	}

	src := html

	//productStatus
	productStatus := ""
	if availability := query.Find("#availability"); availability.Size() > 0 {
		availabilityText := availability.Text()
		availabilityText = strings.ToLower(availabilityText)
		productStatus = availabilityText
		productStatus = strings.Trim(productStatus, " \n\r")

	}
	//name
	name := query.Find("#productTitle").Text()

	//price
	price := ""
	if priceblock_ourprice := query.Find("#priceblock_ourprice"); priceblock_ourprice.Size() > 0 {
		price = priceblock_ourprice.Text()
	}
	if priceblock_saleprice := query.Find("#priceblock_saleprice"); priceblock_saleprice.Size() > 0 {
		price = priceblock_saleprice.Text()
	}
	if priceblock_dealprice := query.Find("#priceblock_dealprice"); priceblock_dealprice.Size() > 0 {
		price = priceblock_dealprice.Text()
	}
	//discount
	discount := ""
	if discount_item := query.Find("#regularprice_savings"); discount_item.Size() > 0 {
		discount_txt := discount_item.Text()
		re_dis, _ := regexp.Compile(`\((\d+%)\)`)
		discount_arr := re_dis.FindAllStringSubmatch(discount_txt, -1)
		discount = discount_arr[0][1]
	}

	//brand
	brand := query.Find("#brand").Text()

	//reviews
	reviews := ""
	if summaryStars := query.Find("#summaryStars"); summaryStars.Size() > 0 {
		reviews = summaryStars.Text()
		re, _ := regexp.Compile(`([\d\.]+)`)
		reviews_arr := re.FindAllStringSubmatch(reviews, -1)

		reviews = reviews_arr[1][0]
	}
	//avgRating
	avgRating := ""
	if avgRatingItem := query.Find("#avgRating"); avgRatingItem.Size() > 0 {
		avgRating, _ = avgRatingItem.Html()
		re, _ := regexp.Compile(`\<[\S\s]+?\>`)
		avgRating = re.ReplaceAllString(avgRating, "")
		avgRating = strings.Replace(avgRating, "out of 5 stars", "", -1)
		avgRating = strings.Trim(avgRating, " \n\r")
	}
	//start5 start4 start3 start2 start1
	star5Rating := ""
	star4Rating := ""
	star3Rating := ""
	star2Rating := ""
	star1Rating := ""
	if histogramTable := query.Find("#histogramTable"); histogramTable.Size() > 0 {

		if star5Rating_text, ok := histogramTable.Find("div[aria-label]").Eq(0).Attr("aria-label"); ok {
			star5Rating = star5Rating_text
		}
		if star4Rating_text, ok := histogramTable.Find("div[aria-label]").Eq(1).Attr("aria-label"); ok {
			star4Rating = star4Rating_text
		}
		if star3Rating_text, ok := histogramTable.Find("div[aria-label]").Eq(2).Attr("aria-label"); ok {
			star3Rating = star3Rating_text
		}
		if star2Rating_text, ok := histogramTable.Find("div[aria-label]").Eq(3).Attr("aria-label"); ok {
			star2Rating = star2Rating_text
		}
		if star1Rating_text, ok := histogramTable.Find("div[aria-label]").Eq(4).Attr("aria-label"); ok {
			star1Rating = star1Rating_text
		}
	}
	//mainRank subRank
	mainRank := ""
	mainRankCategory := ""
	subRank := map[string]string{}
	subRank = make(map[string]string, 5)
	if SalesRank := query.Find("#SalesRank"); SalesRank.Size() > 0 {
		re, _ := regexp.Compile(`#(\d+(,\d+){0,})\s*in\s*(.*)\s*\(\<a\s*href=\"(.*)\">[^<]*<\s*/\s*a\s*\>`)
		SalesRank_src, _ := SalesRank.Html()

		rank_arr := re.FindAllStringSubmatch(SalesRank_src, -1)
		if rank_arr != nil && len(rank_arr) > 0 {
			mainRank = rank_arr[0][1]
			mainRank = strings.Replace(mainRank, ",", "", -1)
			mainRankCategory = rank_arr[0][3]
		}

		SalesRank.Find("ul li").Each(func(i int, s *goquery.Selection) {
			tmp_sub_rank := s.Find(".zg_hrsr_rank").Text()
			tmp_sub_rank = strings.Replace(tmp_sub_rank, "#", "", -1)
			tmp_sub_rank = strings.Trim(tmp_sub_rank, " \n\r")
			tmp_sub_rank_category := ""
			s.Find(".zg_hrsr_ladder a").Each(func(ii int, s *goquery.Selection) {
				//tmp_url, _ := s.Attr("href")
				tmp_sub_rank_category = tmp_sub_rank_category + s.Text() + "->"
			})
			subRank[tmp_sub_rank_category] = tmp_sub_rank

		})

	} else if productDetails_item := query.Find("#productDetails_detailBullets_sections1 tbody tr"); productDetails_item.Size() > 0 {
		productDetails_item.Each(func(i int, s *goquery.Selection) {
			td_title := s.Find("th").Eq(0).Text()

			if strings.Contains(td_title, "Best Sellers Rank") {
				if tmptd := s.Find("td").Eq(0); tmptd != nil {
					rank_text := tmptd.Text()
					rank_arr2 := strings.Split(rank_text, "#")
					if len(rank_arr2) > 1 {
						mainrank_re, _ := regexp.Compile(`(\d+(,\d+){0,})\s*in\s*(.*)\s*\(`)
						mainrank_re_arr := mainrank_re.FindAllStringSubmatch(rank_arr2[1], -1)
						if mainrank_re_arr != nil && len(mainrank_re_arr) > 0 {
							mainRank = mainrank_re_arr[0][1]
							mainRank = strings.Replace(mainRank, ",", "", -1)
							mainRankCategory = mainrank_re_arr[0][3]
						}
						if len(rank_arr2) >= 3 {
							subrank_re, _ := regexp.Compile(`(\d+(,\d+){0,})\s*in\s*(.*)\s*`)
							for ii := 2; ii < len(rank_arr2); ii++ {
								subrank_re_arr := subrank_re.FindAllStringSubmatch(rank_arr2[ii], -1)
								if subrank_re_arr != nil && len(subrank_re_arr) > 0 {
									subRank_temp := subrank_re_arr[0][1]
									subRank_temp = strings.Replace(subRank_temp, ",", "", -1)
									subRankCategory := subrank_re_arr[0][3]

									subRank[subRankCategory] = subRank_temp
								}
							}
						}

					}
				}

			}
		})
	}

	/**解析产品数据***/
	text := query.Text()

	shippingWeight := ""
	shippingWeight_re, _ := regexp.Compile(`(?i)Shipping(\s*)Weight(\D*?)([\d\.]+)(\s*)`)
	if shippingWeight_arr := shippingWeight_re.FindAllStringSubmatch(text, -1); shippingWeight_arr != nil {
		shippingWeight = shippingWeight_arr[0][3]
	}
	//7.4 x 7 x 2 Length Width Height
	dimensions_length := ""
	dimensions_width := ""
	dimensions_height := ""
	dimensions_re, _ := regexp.Compile(`(?i)Dimensions(\D*?)([\d\.]+)(\s*)x(\s*)([\d\.]+)(\s*)x(\s*)([\d\.]+)`)
	if dimensions_arr := dimensions_re.FindAllStringSubmatch(text, -1); dimensions_arr != nil {
		dimensions_length = dimensions_arr[0][2]
		dimensions_width = dimensions_arr[0][5]
		dimensions_height = dimensions_arr[0][8]
	}

	//mainImg
	mainImg, _ := query.Find("#landingImage").Attr("src")

	//merchant-info
	store_url := ""
	store_name := ""
	store := query.Find("#merchant-info a").Eq(0)
	if store.Size() > 0 {
		store_url, _ = store.Attr("href")
		store_name = store.Text()
	}
	merchantID := ""
	if merchantID_item := query.Find("#merchantID"); merchantID_item.Size() > 0 {
		if merchantID_text, ok := merchantID_item.Attr("value"); ok {
			merchantID = merchantID_text
		}
	}
	storeID := ""
	if storeID_item := query.Find("#storeID"); storeID_item.Size() > 0 {
		if storeID_text, ok := storeID_item.Attr("value"); ok {
			storeID = storeID_text
		}
	}

	FBA := ctx.GetTemp("FBA", "")

	ASIN := ""

	ASIN_item := query.Find("#ASIN")
	if ASIN_item.Size() > 0 {
		ASIN_text, ok := ASIN_item.Attr("value")
		if ok {
			ASIN = ASIN_text
		}
	}
	//}
	addDate := ""
	// 结果存入Response中转
	ctx.Output(map[int]interface{}{
		0:  ASIN,
		1:  name,
		2:  price,
		3:  FBA,
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
		20: discount,
		21: merchantID,
		22: storeID,
		23: dimensions_length,
		24: dimensions_width,
		25: dimensions_height,
	})

	//reviews list

	if reviewsLinkItem := query.Find("#revSum #summaryStars a"); reviewsLinkItem.Size() > 0 {
		if reviewsLink, ok := reviewsLinkItem.Attr("href"); ok {
			if strings.Index(reviewsLink, "www.amazon.com") == -1 {
				reviewsLink = "http://www.amazon.com" + reviewsLink
			}
			ctx.AddQueue(&request.Request{
				Url:  reviewsLink,
				Rule: "reviewsList",
				Temp: map[string]interface{}{
					"ASIN": ASIN,
				},
			},
			)
		}
	}

	//跟买

	if moreBuyingitem := query.Find("#mbc-action-panel-wrapper #mbc .a-box .a-padding-base .a-size-small a"); moreBuyingitem.Size() > 0 {
		if tempMoreBuyingLink, ok := moreBuyingitem.Attr("href"); ok {
			if strings.Index(tempMoreBuyingLink, "www.amazon.com") == -1 {
				tempMoreBuyingLink = "http://www.amazon.com" + tempMoreBuyingLink
			}

			ctx.AddQueue(&request.Request{
				Url:  tempMoreBuyingLink,
				Rule: "buyList",
				Temp: map[string]interface{}{
					"ASIN": ASIN,
				},
			},
			)

		}
	}

	//store 数据
	if store_url != "" {
		if !strings.Contains(store_url, "www.amazon.com") {
			store_url = "http://www.amazon.com" + store_url
		}
		ctx.AddQueue(&request.Request{
			Url:  store_url,
			Rule: "store",
			Temp: map[string]interface{}{
				"ASIN": ASIN,
			},
		})
	}

	//存储源码
	//ctx.Parse("product_html")

	//product item
	query.Find("#twister_feature_div  li").Each(func(index int, s *goquery.Selection) {
		dataurl, ok := s.Attr("data-dp-url")
		_ASIN, _ := s.Attr("data-defaultasin")
		if ok {
			if dataurl != "" {
				if !strings.Contains(dataurl, "www.amazon.com") {
					dataurl = "http://www.amazon.com" + dataurl
				}
				ctx.AddQueue(
					&request.Request{
						Url:  dataurl,
						Rule: "product",
						Temp: map[string]interface{}{
							"ASIN":    _ASIN,
							"type":    "",
							"baseUrl": dataurl,
							"FBA":     "",
						},
					},
				)

			}
		}
	})

}

func main() {

	p := Product{}
	session := getSession()

	session.DB("pholcus").C("Amazonbest__product").Find(nil).One(&p)
	fmt.Println(p)

	s, _ := query.Html()
	fmt.Println(s)

}
