package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/mengzhuo/cookiestxt"
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
	err := mainApp()
	if err != nil {
		log.Fatal(err)
	}
}
func mainApp() error {
	var cookies []*http.Cookie
	if _, err := os.Stat(cookiesFile); err == nil {
		if strings.HasSuffix(cookiesFile, ".txt") {
			log.Printf("getting cookies from cookies.txt format file")
			f, err := os.Open(cookiesFile)
			if err != nil {
				return fmt.Errorf("could not open cookies file: %w", err)
			}
			defer f.Close()
			cookies, err = cookiestxt.Parse(f)
			if err != nil {
				return fmt.Errorf("could not parse cookies.txt format file: %w", err)
			}
		} else {
			log.Printf("getting cookies from JSON")
			if buf, err := os.ReadFile(cookiesFile); err == nil {
				json.Unmarshal(buf, &cookies)
			}
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
	return nil
}
