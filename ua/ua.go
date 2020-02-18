package ua

import (
	//"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	rxMaybeCrawlerPattern    = regexp.MustCompile(`(?i)(?:bot|crawler|sp(i|y)der|search|worm|fetch|scaper|nutch)(?:[-_ ./;@()]|$)`)
	rxMaybeFeedParserPattern = regexp.MustCompile(`(?i)(?:feed|web) ?parser`)
	rxMaybeWatchdogPattern   = regexp.MustCompile(`(?i)watch ?dog`)
)

var CrawlerPrefixPatterns = []string{
	"Rome Client ",
	"UnwindFetchor/",
	"ia_archiver ",
	"Summify ",
	"PostRank/",
	// "Python-urllib",
}

var crawlerUserAgents = [...]string{
	"admantx",
	"ahrefssiteaudit",
	"alexa site audit",
	"ask jeeves",
	"biglotron",
	"bingpreview",
	"bitrix link preview",
	"charlotte",
	"chrome-lighthouse",
	"cloudflare-alwaysonline",
	"developers.google.com/+/web/snippet",
	"disqus",
	"docomo",
	"duckduckgo",
	"embedly",
	"envolk",
	"ezooms",
	"facebookexternalhit",
	"facebookplatform",
	"findlinks",
	"flipboard",
	"genieo",
	"gigablast",
	// "go-http-client",
	"google page speed",
	// "grub-client",
	"heritrix",
	"httrack",
	"ichiro",
	"iframely",
	"instagram",
	"ips-agent",
	"jobboerse",
	// "libwww",
	"linkdex",
	"mappy",
	"megaindex",
	"netcraft web server survey",
	"nuzzel",
	"outbrain",
	"page2rss",
	"panscient",
	"pinterest",
	"proximic",
	"quora link preview",
	"qwantify",
	"safesearch microdata crawler",
	"siteexplorer",
	"skypeuripreview",
	"slurp",
	"telegram",
	"tumblr",
	"viber",
	"vkShare",
	"whatsapp",
	"www.google.com/webmasters/tools/richsnippets",
	"xing-contenttabreceiver",
	"yahoo",
	"yandex",
	"yeti",
	"zgrab",
}

var skippedTypes = [...]string{
	".js",
	".css",
	".xml",
	".less",
	".png",
	".jpg",
	".jpeg",
	".gif",
	".pdf",
	".doc",
	".txt",
	".ico",
	".rss",
	".zip",
	".mp3",
	".rar",
	".exe",
	".wmv",
	".doc",
	".avi",
	".ppt",
	".mpg",
	".mpeg",
	".tif",
	".wav",
	".mov",
	".psd",
	".ai",
	".xls",
	".mp4",
	".m4a",
	".swf",
	".dat",
	".dmg",
	".iso",
	".flv",
	".m4v",
	".torrent",
	".ttf",
	".woff",
	".woff2",
	".svg",
	".eot",
	".webmanifest",
}

func ShouldPrerender(or *http.Request) bool {
	userAgent := strings.ToLower(or.Header.Get("User-Agent"))
	bufferAgent := or.Header.Get("X-Bufferbot")
	reqURL := strings.ToLower(or.URL.String())

	// No user agent, don't prerender
	if userAgent == "" {
		return false
	}

	//fmt.Printf("Prerender UserAgent [%s]\n", userAgent)

	// Not a GET or HEAD request, don't prerender
	if or.Method != "GET" && or.Method != "HEAD" {
		return false
	}

	// Static resource, don't prerender
	for _, extension := range skippedTypes {
		if strings.HasSuffix(reqURL, extension) {
			return false
		}
	}

	// Buffer Agent or requesting an excaped fragment, request prerender
	if _, ok := or.URL.Query()["_escaped_fragment_"]; bufferAgent != "" || ok {
		return true
	}

	// check is crawler
	if IsCrawler(userAgent) {
		return true
	}

	return false
}

func IsCrawler(agent string) bool {
	// No user agent, don't prerender
	if agent == "" {
		return false
	}

	// skip googlebot
	if ChallengeGoogle(agent) {
		return false
	}

	// Cralwer, request prerender
	for _, crawlerAgent := range crawlerUserAgents {
		if strings.Contains(agent, crawlerAgent) {
			return true
		}
	}

	if rxMaybeCrawlerPattern.MatchString(agent) || hasCrawlerPrefix(agent) || rxMaybeFeedParserPattern.MatchString(agent) || rxMaybeWatchdogPattern.MatchString(agent) {
		return true
	}

	return false
}

func hasCrawlerPrefix(agent string) bool {
	for _, pattern := range CrawlerPrefixPatterns {
		if strings.HasPrefix(agent, pattern) {
			return true
		}
	}

	if strings.Contains(agent, "ASP-Ranker Feed Crawler") {
		return true
	}

	return false
}

func ChallengeGoogle(agent string) bool {
	if !strings.Contains(agent, "google") {
		return false
	}

	if strings.Contains(agent, "compatible; googlebot") {
		return true
	}

	if strings.Contains(agent, "googlebot-image/") {
		return true
	}

	if strings.Contains(agent, "adsbot-google") {
		return true
	}

	if strings.Contains(agent, "mediapartners-moogle") {
		return true
	}

	if strings.Contains(agent, "feedfetcher-google;") {
		return true
	}

	if strings.Contains(agent, "appengine-google") {
		return true
	}

	return false
}
