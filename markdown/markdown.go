package markdown

import (
	"context"
	"regexp"

	"github.com/gleez/pkg/markdown/markdown"
	"github.com/microcosm-cc/bluemonday"
)

// Options represents option for rendering Markdown content.
type Options struct {
	// TODO(slimsag): add option for controlling relative links here.
}

// DefaultOptions is the default options for rendering Markdown content.
var DefaultOptions = Options{}

// Render renders Markdown content into sanitized HTML that is safe to render anywhere.
//
// When nil, options will default to DefaultOptions.
func Render(content string, options *Options) (string, error) {
	if options == nil {
		options = &DefaultOptions
	}

	// https://github.com/sourcegraph/docsite/blob/master/site.go

	//unsafeHTML := gfm.Markdown([]byte(content))
	// doc, err := markdown.Run(ctx, data, markdown.Options{
	// 	Base:                      base.ResolveReference(&url.URL{Path: urlPathPrefix}),
	// 	ContentFilePathToLinkPath: contentFilePathToPath,
	// 	Funcs:                     createMarkdownFuncs(s),
	// 	FuncInfo:                  markdown.FuncInfo{Version: contentVersion},
	// })

	doc, err := markdown.Run(context.Background(), []byte(content), markdown.Options{})
	if err != nil {
		return "", err
	}

	p := bluemonday.UGCPolicy()
	p.AllowAttrs("name").Matching(bluemonday.SpaceSeparatedTokens).OnElements("a")
	p.AllowAttrs("rel").Matching(regexp.MustCompile(`^nofollow$`)).OnElements("a")
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^anchor$`)).OnElements("a")
	p.AllowAttrs("aria-hidden").Matching(regexp.MustCompile(`^true$`)).OnElements("a")
	p.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
	p.AllowAttrs("checked", "disabled").Matching(regexp.MustCompile(`^$`)).OnElements("input")
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return string(p.SanitizeBytes(doc.HTML)), nil
}
