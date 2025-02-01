package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Game interface {
}

var Games = map[string]*HoyoGame{

	"genshin": &HoyoGame{
		Origin:  "https://act.hoyolab.com",
		Referer: "https://act.hoyolab.com/ys/event/signin-sea-v3/index.html?act_id=e202102251931481",

		ActionID:   "e202102251931481",
		InfoAPIURL: "https://sg-hk4e-api.hoyolab.com/event/sol/info?lang=en-us&act_id=e202102251931481",
		SignAPIURL: "https://sg-hk4e-api.hoyolab.com/event/sol/sign",
	},

	"honkai": &HoyoGame{
		Origin:  "https://act.hoyolab.com",
		Referer: "https://act.hoyolab.com/bbs/event/signin/hkrpg/index.html?act_id=e202303301540311",

		ActionID:   "e202303301540311",
		InfoAPIURL: "https://sg-public-api.hoyolab.com/event/luna/os/info?lang=en-us&act_id=e202303301540311",
		SignAPIURL: "https://sg-public-api.hoyolab.com/event/luna/os/sign",
	},

	"zzz": &HoyoGame{
		Origin:     "https://act.hoyolab.com",
		Referer:    "https://act.hoyolab.com/", //"https://act.hoyolab.com/bbs/event/signin/zzz/e202406031448091.html?act_id=e202406031448091",
		InfoAPIURL: "https://sg-public-api.hoyolab.com/event/luna/zzz/os/info?lang=en-us&act_id=e202406031448091",
		ActionID:   "e202406031448091",
		SignAPIURL: "https://sg-public-api.hoyolab.com/event/luna/zzz/os/sign",
		MoreHeaders: func(header http.Header) {
			header.Add("x-rpc-signgame", "zzz")
		},
	},
}

type HoyoGame struct {
	Origin  string // http Origin policy origin
	Referer string // http referer (page you're on)

	ActionID    string
	InfoAPIURL  string // URL to the "info" action
	SignAPIURL  string // URL to the "sign" action
	MoreHeaders func(http.Header)
}

func (g *HoyoGame) CanClaim(client *http.Client) (canClaim bool, err error) {
	// https://act.hoyolab.com/ys/event/signin-sea-v3/index.html?act_id=e202102251931481
	req, err := http.NewRequest("GET", g.InfoAPIURL, nil)
	if err != nil {
		return false, err
	}
	commonHeaders(g, req.Header)

	var infoResponse InfoResponse
	_, body, err := DoRequestDecode(client, req, &infoResponse)
	if err == nil {
		log.Printf("%s\n\n%v", body, infoResponse.Data)
		if infoResponse.Data == nil || infoResponse.Data.IsSign {
			log.Printf("Already claimed, quitting")
			return false, nil
		}
		return true, nil
	}
	return false, err
}

func (g *HoyoGame) Claim(client *http.Client) error {
	buf, err := json.Marshal(map[string]string{"act_id": g.ActionID})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", g.SignAPIURL, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	commonHeaders(g, req.Header)
	_, body, err := DoRequestGetBody(client, req)
	if err != nil {
		return err
	}
	log.Printf("sign body: %s", body)
	return nil
}

func commonHeaders(g *HoyoGame, header http.Header) {
	header.Set("Accept", "application/json, text/plain, */*")
	header.Set("Accept-Language", "en-US,en,q=0.5")
	header.Set("Origin", g.Origin)
	header.Set("Referer", g.Referer)
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	if g.MoreHeaders != nil {
		g.MoreHeaders(header)
	}
}

// {"retcode":0,"message":"OK","data":{"total_sign_day":12,"today":"2024-11-20","is_sign":true,"first_bind":false,"is_sub":false,"region":"","month_last_day":false}}
// {"data":null,"message":"网络出小差了，请稍后重试~","retcode":-500001} <-- some error when you don't have the game maybe??
type InfoResponse struct {
	RetCode int               `json:"retcode"`
	Data    *InfoResponseData `json:"data"`
}
type InfoResponseData struct {
	IsSign bool `json:"is_sign"` // true if already claimed
}

// Already checked in:
//
//	{"data":null,"message":"Traveler, you've already checked in today~","retcode":-5003}
//
// Good response:
//
//	{"retcode":0,"message":"OK","data":{"code":"","risk_code":0,"gt":"","challenge":"","success":0,"is_risk":false}}
type SignResponse struct {
	Message string `json:"message"`
	Retcode int    `json:"retcode"`
}
