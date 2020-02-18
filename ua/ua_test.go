package ua

import (
	"strings"
	"testing"
)

// Slice that contains all the tests. Each test is contained in a struct
// that groups the title of the test, the User-Agent string to be tested and the expected value.
var uastrings = []struct {
	title    string
	ua       string
	expected bool
}{
	{
		title:    "Curl",
		ua:       "curl/7.28.1",
		expected: false,
	},
	{
		title:    "ChromeLinux",
		ua:       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.97 Safari/537.11",
		expected: false,
	},
	// Bots
	{
		title:    "GoogleBot",
		ua:       "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		expected: false,
	},
	{
		title:    "GoogleBotSmartphone (iPhone)",
		ua:       "Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		expected: false,
	},
	{
		title:    "GoogleBotSmartphone (Android)",
		ua:       "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		expected: false,
	},
	{
		title:    "BingBot",
		ua:       "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		expected: true,
	},
	{
		title:    "BingPreview",
		ua:       "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534+ (KHTML, like Gecko) BingPreview/1.0b",
		expected: true,
	},
	{
		title:    "BaiduBot",
		ua:       "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
		expected: true,
	},
	{
		title:    "Baidu Image",
		ua:       "Baiduspider-image+(+http://www.baidu.com/search/spider.htm)",
		expected: true,
	},
	{
		title:    "Twitterbot",
		ua:       "Twitterbot",
		expected: true,
	},
	{
		title:    "YahooBot",
		ua:       "Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)",
		expected: true,
	},
	{
		title:    "FacebookExternalHit",
		ua:       "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php",
		expected: true,
	},
	{
		title:    "FacebookPlatform",
		ua:       "facebookplatform/1.0 (+http://developers.facebook.com)",
		expected: true,
	},
	{
		title:    "FaceBot",
		ua:       "Facebot",
		expected: true,
	},
	{
		title:    "iframely",
		ua:       "iframely/1.2.7 (+http://iframely.com/;)",
		expected: true,
	},
	{
		title:    "NutchCVS",
		ua:       "NutchCVS/0.8-dev (Nutch; http://lucene.apache.org/nutch/bot.html; nutch-agent@lucene.apache.org)",
		expected: true,
	},
	{
		title:    "MJ12bot",
		ua:       "Mozilla/5.0 (compatible; MJ12bot/v1.2.4; http://www.majestic12.co.uk/bot.php?+)",
		expected: true,
	},
	{
		title:    "MJ12bot",
		ua:       "MJ12bot/v1.0.8 (http://majestic12.co.uk/bot.php?+)",
		expected: true,
	},
	{
		title:    "AhrefsBot",
		ua:       "Mozilla/5.0 (compatible; AhrefsBot/4.0; +http://ahrefs.com/robot/)",
		expected: true,
	},
	{
		title:    "AdsBotGoogle",
		ua:       "AdsBot-Google (+http://www.google.com/adsbot.html)",
		expected: false,
	},
	{
		title:    "AdsBotGoogleMobile",
		ua:       "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		expected: false,
	},
	{
		title:    "WhatsApp",
		ua:       "WhatsApp/2.16.57 A",
		expected: true,
	},
	{
		title:    "TelegramBot",
		ua:       "TelegramBot (like TwitterBot)",
		expected: true,
	},
	{
		title:    "DuckDuckBot",
		ua:       "DuckDuckBot/1.1; (+http://duckduckgo.com/duckduckbot.html)",
		expected: true,
	},
	{
		title:    "pinterest",
		ua:       "Pinterest/0.2 (+http://www.pinterest.com/bot.html)",
		expected: true,
	},
	{
		title:    "LinkedInBot",
		ua:       "LinkedInBot/1.0 (compatible; Mozilla/5.0; Jakarta Commons-HttpClient/4.3 +http://www.linkedin.com)",
		expected: true,
	},
	{
		title:    "Quora",
		ua:       "Quora Link Preview/1.0 (http://www.quora.com)",
		expected: true,
	},
	{
		title:    "ia_archiver",
		ua:       "ia_archiver (+http://www.alexa.com/site/help/webmasters; crawler@alexa.com)",
		expected: true,
	},
	{
		title:    "HTTrack",
		ua:       "Mozilla/4.5 (compatible; HTTrack 3.0x; Windows 98)",
		expected: true,
	},
	{
		title:    "Sogou",
		ua:       "Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
		expected: true,
	},
	{
		title:    "Proximic",
		ua:       "Mozilla/5.0 (compatible; proximic; +https://www.comscore.com/Web-Crawler)",
		expected: true,
	},
	{
		title:    "SEMRush",
		ua:       "Mozilla/5.0 (compatible; SemrushBot/2~bl; +http://www.semrush.com/bot.html)",
		expected: true,
	},
	{
		title:    "Yandex Images",
		ua:       "Mozilla/5.0 (compatible; YandexImages/3.0; +http://yandex.com/bots)",
		expected: true,
	},
	{
		title:    "Ask Jeeves",
		ua:       "Mozilla/5.0 (compatible; Ask Jeeves/Teoma; +http://about.ask.com/en/docs/about/webmasters.shtml)",
		expected: true,
	},
	{
		title:    "Jobboerse",
		ua:       "Mozilla/5.0 (X11; U; Linux Core i7-4980HQ; de; rv:32.0; compatible; Jobboerse.com; http://www.xn--jobbrse-d1a.com) Gecko/20100401 Firefox/24.0",
		expected: true,
	},
	{
		title:    "Yahoo MMCrawler",
		ua:       "Yahoo-MMCrawler/3.x (mms dash mmcrawler dash support at yahoo dash inc dot com)",
		expected: true,
	},
	{
		title:    "Special Archiver",
		ua:       "Mozilla/5.0 (compatible; special_archiver/3.1.1 +http://www.archive.org/details/archive.org_bot)",
		expected: true,
	},
	{
		title:    "ZGrab",
		ua:       "Mozilla/5.0 zgrab/0.x",
		expected: true,
	},
	// {
	// 	title:    "Google Feedfetcher",
	// 	ua:       "Feedfetcher-Google; (+http://www.google.com/feedfetcher.html; 1 subscribers; feed-id=4296914164355380091)",
	// 	expected: true,
	// },
	{
		title:    "Accoona",
		ua:       "Accoona-AI-Agent/1.1.1 (crawler at accoona dot com)",
		expected: true,
	},
	{
		title:    "ADmantX",
		ua:       "ADmantX Platform Semantic Analyzer - ADmantX Inc. - www.admantx.com - support@admantx.com",
		expected: true,
	},
	{
		title:    "Netcraft",
		ua:       "Mozilla/4.0 (compatible; Netcraft Web Server Survey)",
		expected: true,
	},
	{
		title:    "Evaliant",
		ua:       "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36 evaliant Evaliant Impressions Bot 41 Windows Blink Common",
		expected: true,
	},
	{
		title:    "envolk",
		ua:       "envolk/1.7 (+http://www.envolk.com/envolkspiderinfo.html)",
		expected: true,
	},
	{
		title:    "aboundex",
		ua:       "Aboundex/0.2 (http://www.aboundex.com/crawler/)",
		expected: true,
	},
	{
		title:    "eright",
		ua:       "Mozilla/5.0 (compatible; eright/1.0; +bot@eright.com)",
		expected: true,
	},
}

// The test suite.
func TestUserAgent(t *testing.T) {
	for _, tt := range uastrings {
		ua := strings.ToLower(tt.ua)
		got := IsCrawler(ua)
		if tt.expected != got {
			t.Errorf("\nTest     %v\ngot:     %t\nexpected %t\n", tt.title, got, tt.expected)
		}
	}
}

// Benchmark: it parses each User-Agent string on the uastrings slice b.N times.
func BenchmarkUserAgent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, tt := range uastrings {
			b.StartTimer()
			ua := strings.ToLower(tt.ua)
			IsCrawler(ua)
		}
	}
}
