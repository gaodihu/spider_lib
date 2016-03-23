package spider_lib

// 基础包
import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"                                //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"                   //必需
	. "github.com/henrylee2cn/pholcus/logs"
	//. "github.com/henrylee2cn/pholcus/spider/common"          //选用
)

func init() {
	Amazon.Register()
}

var Amazon = &Spider{
	Name:        "Amazon",
	Description: "Amazon商品数据 [Auto Page] [www.Amazon.com]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Optional: &Optional{},
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			f, err := os.Open("E:/go/src/github.com/henrylee2cn/pholcus/amazon_list.txt")
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
							Rule: "list",
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
						Url:  "http://www.amazon.com/iWG-Dual-slot-Flashlight-Battery-Standard/dp/B016Y4RYU4/ref=sr_1_2?ie=UTF8&qid=1449729556&sr=8-2&keywords=myled",
						Rule: "product",
					},
				)
			*/
		},

		Trunk: map[string]*Rule{

			"list": {
				ParseFunc: func(ctx *Context) {

					query := ctx.GetDom()
					//src, _ := query.Html()

					if next_page := query.Find("#pagnNextLink"); next_page.Size() > 0 {
						if next_url, ok := next_page.Attr("href"); ok {
							ctx.AddQueue(&request.Request{
								Url:  "http://www.amazon.com" + next_url,
								Rule: "list",
							},
							)
						}
					}

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
					"name",
					"price",
					"FBA",
					"brand",
					"reviews",
					"avgRating",
					"star5Rating",
					"star4Rating",
					"star3Rating",
					"star2Rating",
					"star1Rating",
					"mainRank",
					"mainRankCategory",
					"addDate",
					"shippingWeight",
					"subRank",
					"mainImg",
					"store",
					"productStatus",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					src, _ := query.Html()
					productStatus := ""
					if availability := query.Find("#availability"); availability.Size() > 0 {
						availabilityText := availability.Text()
						availabilityText = strings.ToLower(availabilityText)
						productStatus = availabilityText
						productStatus = strings.Trim(productStatus, " \n\r")
						/*
							if !strings.Contains(availabilityText, "in stock") {
								return
							}
						*/
					}
					name := query.Find("#productTitle").Text()

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

					brand := query.Find("#brand").Text()
					reviews := ""
					if summaryStars := query.Find("#summaryStars"); summaryStars.Size() > 0 {
						reviews, _ = summaryStars.Html()
						re, _ := regexp.Compile(`\<[\S\s]+?\>`)
						reviews = re.ReplaceAllString(reviews, "")
						reviews = strings.Replace(reviews, ",", "", -1)
						reviews = strings.Trim(reviews, " \n\r")
					}
					avgRating := ""
					if avgRatingItem := query.Find("#avgRating"); avgRatingItem.Size() > 0 {
						avgRating, _ = avgRatingItem.Html()
						re, _ := regexp.Compile(`\<[\S\s]+?\>`)
						avgRating = re.ReplaceAllString(avgRating, "")
						avgRating = strings.Replace(avgRating, "out of 5 stars", "", -1)
						avgRating = strings.Trim(avgRating, " \n\r")
					}
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

					}

					addDate := ""
					if addDateItem := query.Find(".date-first-available"); addDateItem.Size() > 0 {
						addDate = addDateItem.Find(".value").Text()
					}
					//<li><b> Date first available at Amazon.com:</b> October 8, 2014</li>
					re, _ := regexp.Compile(`\<b\> Date first available at Amazon.com:\<\/b\> (.*)\<\/li\>`)
					//src, _ := query.Html()
					addDateArr := re.FindAllStringSubmatch(src, -1)
					if addDateArr != nil && len(addDateArr) > 0 {
						addDate = addDateArr[0][1]
					}

					shippingWeight := ""
					if shippingWeightItem := query.Find(".shipping-weight"); shippingWeightItem.Size() > 0 {
						shippingWeight = shippingWeightItem.Find(".value").Text()
						shippingWeight = strings.Replace(shippingWeight, "ounces (View shipping rates and policies)", "", -1)
						shippingWeight = strings.Trim(shippingWeight, " \n\r")
					}

					mainImg, _ := query.Find("#landingImage").Attr("src")

					//merchant-info
					store_url := ""
					store_name := ""
					store := query.Find("#merchant-info a").Eq(0)
					if store.Size() > 0 {
						store_url, _ = store.Attr("href")
						store_name = store.Text()
					}

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0:  ctx.GetTemp("ASIN", ""),
						1:  name,
						2:  price,
						3:  ctx.GetTemp("FBA", ""),
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

					ASIN := ctx.GetTemp("ASIN", "")
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
					if store_url != "" {
						ctx.AddQueue(&request.Request{
							Url:  "http://www.amazon.com" + store_url,
							Rule: "store",
							Temp: map[string]interface{}{
								"ASIN": ASIN,
							},
						})
					}

				},
			},
			"buyList": {
				ItemFields: []string{
					"ASIN",
					"offeringID",
					"price",
					"shippingFee",
					"FBA",
					"buyer",
					"condition",
				},
				ParseFunc: func(ctx *Context) {

					query := ctx.GetDom()
					query.Find("#olpTabContent .olpOffer").Each(func(index int, s *goquery.Selection) {
						price := s.Find(".olpOfferPrice").Text()
						price = strings.Trim(price, " ")

						shippingFee := ""
						if olpShippingInfo := s.Find(".olpShippingInfo"); olpShippingInfo.Size() > 0 {
							if shippingFeeItem := olpShippingInfo.Find(".olpShippingPrice"); shippingFeeItem.Size() > 0 {
								shippingFee = shippingFeeItem.Text()
							} else {
							}
							olpShippingInfoText := olpShippingInfo.Text()
							if strings.Contains(olpShippingInfoText, "on orders over") {
								shippingFee = olpShippingInfoText
							} else {
								if strings.Contains(olpShippingInfoText, "FREE Shipping") {
									shippingFee = "0"
								}
							}
						}
						condition := s.Find(".olpCondition").Text()
						buyer := s.Find(".olpSellerColumn .olpSellerName span a").Text()

						FBA := 0
						if prime := s.Find(".a-icon-prime"); prime.Size() > 0 {
							FBA = 1
						}
						//offeringID := ""
						offeringID, _ := s.Find("input[name^=offeringID]").Attr("value")

						ASIN := ctx.GetTemp("ASIN", "")
						// 结果存入Response中转
						ctx.Output(map[int]interface{}{
							0: ASIN,
							1: offeringID,
							2: price,
							3: shippingFee,
							4: FBA,
							5: buyer,
							6: condition,
						})

					})
					if nextUrlItem := query.Find(".a-pagination .a-last a"); nextUrlItem.Size() > 0 {
						nextUrl, _ := nextUrlItem.Attr("href")
						if !strings.Contains(nextUrl, "www.amazon.com") {
							nextUrl = "http://www.amazon.com" + nextUrl
						}
						ctx.AddQueue(&request.Request{
							Url:  nextUrl,
							Rule: "buyList",
							Temp: map[string]interface{}{
								"ASIN": ctx.GetTemp("ASIN", ""),
							},
						},
						)
					}
				},
			},
			"reviewsList": {
				ItemFields: []string{
					"ASIN",
					"addDate",
					"star",
					"verifiedPurchase",
					"author",
					"content",
					"joinNumber",
					"helpful",
					"reviewId",
				},
				ParseFunc: func(ctx *Context) {

					query := ctx.GetDom()
					if reviewsList := query.Find("#cm_cr-review_list div.review"); reviewsList.Size() > 0 {
						reviewsList.Each(func(index int, s *goquery.Selection) {
							reviewId, _ := s.Attr("id")

							helpful := ""
							joinNumber := ""
							helpfulText := s.Find(".a-row helpful-votes-count span").Text()
							re, _ := regexp.Compile(`(\d+)\s*of\s*(\d+)`)
							helpfulArr := re.FindAllStringSubmatch(helpfulText, -1)
							if helpfulArr != nil && len(helpfulArr) > 0 {
								helpful = helpfulArr[0][1]
								joinNumber = helpfulArr[0][2]
							}

							star := s.Find(".review-rating span").Text()
							star = strings.Replace(star, "out of 5 stars", "", -1)
							star = strings.Trim(star, " \n\r")

							addDate := s.Find(".review-date").Text()
							addDate = strings.Replace(addDate, "on", "", -1)
							addDate = strings.Trim(addDate, " ")

							author := s.Find(".review-byline .author").Text()

							verifiedPurchase := 0
							if strings.Contains(s.Text(), "Verified Purchase") {
								verifiedPurchase = 1
							}

							content, _ := s.Find(".review-text").Html()

							// 结果存入Response中转
							ctx.Output(map[int]interface{}{
								0: ctx.GetTemp("ASIN", ""),
								1: addDate,
								2: star,
								3: verifiedPurchase,
								4: author,
								5: content,
								6: joinNumber,
								7: helpful,
								8: reviewId,
							})

						})
					}

					if nextUrlItem := query.Find("#cm_cr-pagination_bar .a-pagination .a-last a"); nextUrlItem.Size() > 0 {
						nextUrl, _ := nextUrlItem.Attr("href")
						if !strings.Contains(nextUrl, "www.amazon.com") {
							nextUrl = "http://www.amazon.com" + nextUrl
						}
						ctx.AddQueue(&request.Request{
							Url:  nextUrl,
							Rule: "buyList",
							Temp: map[string]interface{}{
								"ASIN": ctx.GetTemp("ASIN", ""),
							},
						},
						)
					}
				},
			},
			"store": {
				ItemFields: []string{
					"storeName",
					"Positive_30days",
					"Neutral_30days",
					"Negative_30days",
					"Count_30days",
					"Positive_90days",
					"Neutral_90days",
					"Negative_90days",
					"Count_90days",
					"Positive_12months",
					"Neutral_12months",
					"Negative_12months",
					"Count_12months",
					"Positive_lifetime",
					"Neutral_lifetime",
					"Negative_lifetime",
					"Count_lifetime",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					table := query.Find("#feedback-summary-table")
					if table.Size() > 0 {
						positive_30days := table.Find("tbody tr").Eq(1).Find("td").Eq(1).Text()
						positive_30days = strings.Trim(positive_30days, " \r\n")

						positive_90days := table.Find("tbody tr").Eq(1).Find("td").Eq(2).Text()
						positive_90days = strings.Trim(positive_90days, " \r\n")

						positive_12months := table.Find("tbody tr").Eq(1).Find("td").Eq(3).Text()
						positive_12months = strings.Trim(positive_12months, " \r\n")

						positive_lifetime := table.Find("tbody tr").Eq(1).Find("td").Eq(4).Text()
						positive_lifetime = strings.Trim(positive_lifetime, " \r\n")

						neutral_30days := table.Find("tbody tr").Eq(2).Find("td").Eq(1).Text()
						neutral_30days = strings.Trim(neutral_30days, " \r\n")

						neutral_90days := table.Find("tbody tr").Eq(2).Find("td").Eq(2).Text()
						neutral_90days = strings.Trim(neutral_90days, " \r\n")

						neutral_12months := table.Find("tbody tr").Eq(2).Find("td").Eq(3).Text()
						neutral_12months = strings.Trim(neutral_12months, " \r\n")

						neutral_lifetime := table.Find("tbody tr").Eq(2).Find("td").Eq(4).Text()
						neutral_lifetime = strings.Trim(neutral_lifetime, " \r\n")

						negative_30days := table.Find("tbody tr").Eq(3).Find("td").Eq(1).Text()
						negative_30days = strings.Trim(negative_30days, " \r\n")

						negative_90days := table.Find("tbody tr").Eq(3).Find("td").Eq(2).Text()
						negative_90days = strings.Trim(negative_90days, " \r\n")

						negative_12months := table.Find("tbody tr").Eq(3).Find("td").Eq(3).Text()
						negative_12months = strings.Trim(negative_12months, " \r\n")

						negative_lifetime := table.Find("tbody tr").Eq(3).Find("td").Eq(4).Text()
						negative_lifetime = strings.Trim(negative_lifetime, " \r\n")

						count_30days := table.Find("tbody tr").Eq(3).Find("td").Eq(1).Text()
						count_30days = strings.Trim(count_30days, " \r\n")

						count_90days := table.Find("tbody tr").Eq(3).Find("td").Eq(2).Text()
						count_90days = strings.Trim(count_90days, " \r\n")

						count_12months := table.Find("tbody tr").Eq(3).Find("td").Eq(3).Text()
						count_12months = strings.Trim(count_12months, " \r\n")

						count_lifetime := table.Find("tbody tr").Eq(3).Find("td").Eq(4).Text()
						count_lifetime = strings.Trim(count_lifetime, " \r\n")

						storeName := query.Find("#sellerName").Text()
						storeName = strings.Trim(storeName, " \r\n")

						// 结果存入Response中转
						ctx.Output(map[int]interface{}{
							0:  storeName,
							1:  positive_30days,
							2:  neutral_30days,
							3:  negative_30days,
							4:  count_30days,
							5:  positive_90days,
							6:  neutral_90days,
							7:  negative_90days,
							8:  count_90days,
							9:  positive_12months,
							10: neutral_12months,
							11: negative_12months,
							12: count_12months,
							13: positive_lifetime,
							14: neutral_lifetime,
							15: negative_lifetime,
							16: count_lifetime,
						})

					}

				},
			},
		},
	},
}
