package getcookie

import (
	"log"
	"net/http"
	"strings"

	"github.com/browserutils/kooky"

	// comment out browsers you don't want
	// _ "github.com/browserutils/kooky/browser/chrome"
	_ "github.com/browserutils/kooky/browser/chromium"
	_ "github.com/browserutils/kooky/browser/edge"
	_ "github.com/browserutils/kooky/browser/firefox"
	_ "github.com/browserutils/kooky/browser/safari"
	//	_ "github.com/browserutils/kooky/browser/all"
)

func GetHoyoCookie() []*http.Cookie {
	var validCookies []*http.Cookie
	cookies := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix("hoyolab.com"))
	for _, c := range cookies {
		if strings.HasPrefix(c.Name, "_ga") {
			continue
		}
		log.Printf("found cookie: domain: %s, name: %s, val %s ", c.Domain, c.Name, c.Value)
		validCookies = append(validCookies, &c.Cookie)
	}
	return validCookies
}
