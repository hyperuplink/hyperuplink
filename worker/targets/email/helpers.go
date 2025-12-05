package email

import (
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/mrusme/hyperuplink/models/asyncjob"
)

type TmplCacheItem struct {
	TextTmpl *texttemplate.Template
	HtmlTmpl *htmltemplate.Template
}

type TmplCache map[string]TmplCacheItem

func (t Email) TemplatesFor(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	lang string,
) (
	textTmpl *texttemplate.Template,
	htmlTmpl *htmltemplate.Template,
	err error,
) {
	if cache, ok := t.tmplCache[lang]; ok {
		textTmpl = cache.TextTmpl
		htmlTmpl = cache.HtmlTmpl
	} else {
		textTmpl, htmlTmpl, err = t.LoadTemplates(
			jobType, jobSubType, lang)
		if err != nil {
			return nil, nil, err
		}

		t.tmplCache[lang] = TmplCacheItem{
			TextTmpl: textTmpl,
			HtmlTmpl: htmlTmpl,
		}
	}

	return textTmpl, htmlTmpl, nil
}

func (t *Email) LoadTemplates(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	lang string,
) (
	textTmpl *texttemplate.Template,
	htmlTmpl *htmltemplate.Template,
	err error,
) {
	var tmpl string = fmt.Sprintf(
		"templates/email/%s/%s.%s.tmpl",
		string(jobType), string(jobSubType), lang,
	)

	if textTmpl, err = texttemplate.ParseFS(t.rt.Embeds["templates"],
		tmpl+".txt"); err != nil {
		return nil, nil, err
	}
	if textTmpl, err = texttemplate.ParseFS(t.rt.Embeds["templates"],
		tmpl+".html"); err != nil {
		return nil, nil, err
	}

	return textTmpl, htmlTmpl, nil
}
