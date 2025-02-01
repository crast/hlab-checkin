package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/publicsuffix"

	"hlab-checkin/app"
	"hlab-checkin/internal/getcookie"
)

var cookiesFile = "cookies.json"
var chosenGame string

func main() {
	flag.StringVar(&chosenGame, "game", "genshin", "Which game to do (genshin/zzz/honkai)")
	flag.StringVar(&cookiesFile, "file", "cookies.json", "file to store/load cookies")

	flag.Parse()

	var cookies []*http.Cookie
	if _, err := os.Stat(cookiesFile); err == nil {
		log.Printf("getting cookies from file")
		if buf, err := os.ReadFile(cookiesFile); err == nil {
			json.Unmarshal(buf, &cookies)
		}
	}
	if len(cookies) == 0 {
		log.Printf("reading cookies")
		cookies = getcookie.GetHoyoCookie()
		buf, err := json.MarshalIndent(cookies, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		os.WriteFile("cookies.json", buf, 0666)
	}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	u, _ := url.Parse("https://sg-hk4e-api.hoyolab.com/")
	jar.SetCookies(u, cookies)

	client := &http.Client{
		Jar: jar,
	}

	gameKey := strings.ToLower(chosenGame)
	canClaim, err := app.Games[gameKey].CanClaim(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("can claim: %v", canClaim)
	if canClaim {
		app.Games[gameKey].Claim(client)
	}
}
