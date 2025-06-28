package shopify

import (
	"net/url"
	"strings"
	"testing"
)

var (
	testShops = []string{
		"https://test.myshopify.com",
		"https://test.myshopify.com/admin",
		"test.myshopify.com",
	}
)

func TestShopNameParse(t *testing.T) {
	for _, dest := range testShops {
		t.Log(dest)
		if !strings.HasPrefix(dest, "https://") {
			dest = "https://" + dest
			t.Log(dest)
		}
		parsedURL, err := url.Parse(dest)
		if err != nil {
			t.Fatal(err)
		}
		host := parsedURL.Hostname()

		if strings.HasSuffix(host, ".myshopify.com") {
			shopName := strings.TrimSuffix(host, ".myshopify.com")
			t.Log(shopName)
		} else {
			t.Fatal(host)
		}
	}

}
