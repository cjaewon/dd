package server

import (
	"bytes"
	"errors"
	"html/template"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.CJK,
		meta.Meta,
		highlighting.NewHighlighting(
			highlighting.WithStyle("github"),
			highlighting.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
				chromahtml.WithClasses(true),
			),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func (a *ArticleTmplContext) MarkDown(content []byte) error {
	var buf bytes.Buffer
	ctx := parser.NewContext()

	if err := md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return err
	}

	metaData := meta.Get(ctx)

	title, ok := metaData["title"]
	if !ok {
		return errors.New("title is not available")
	}

	strTitle, ok := title.(string)
	if !ok {
		return errors.New("title can not assert as string type")
	}

	date, ok := metaData["date"]
	if !ok {
		return errors.New("date is not available")
	}

	strDate, ok := date.(string)
	if !ok {
		return errors.New("date can not assert as string type")
	}

	d, err := time.Parse(time.RFC3339, strDate)
	if err != nil {
		return err
	}

	a.Title = strTitle
	a.Date = d
	a.HTML = template.HTML(buf.String())

	return nil
}
