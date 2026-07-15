package tmpl

import (
	"fmt"
	htmltemplate "html/template"
	"io/fs"
	"sync"
	texttemplate "text/template"

	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
)

type Spec struct {
	Dir     string
	TextExt string
	HtmlExt string
}

var (
	EmailSpec = Spec{Dir: "email", TextExt: ".eml", HtmlExt: ".html"}
	XMPPSpec  = Spec{Dir: "xmpp", TextExt: ".md"}
)

type Item struct {
	TextTmpl *texttemplate.Template
	HtmlTmpl *htmltemplate.Template
}

type Cache struct {
	fsys fs.FS
	spec Spec

	mu    sync.RWMutex
	items map[string]Item
}

func NewCache(fsys fs.FS, spec Spec) (c *Cache) {
	c = new(Cache)

	c.fsys = fsys
	c.spec = spec
	c.items = make(map[string]Item)

	return c
}

func (c *Cache) base(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	lang string,
) string {
	return fmt.Sprintf(
		"templates/%s/%s/%s.%s.tmpl",
		c.spec.Dir, string(jobType), string(jobSubType), lang,
	)
}

func (c *Cache) TemplatesFor(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	lang string,
) (item Item, err error) {
	base := c.base(jobType, jobSubType, lang)

	c.mu.RLock()
	item, ok := c.items[base]
	c.mu.RUnlock()

	if ok {
		return item, nil
	}

	if item, err = c.load(base); err != nil {
		return Item{}, err
	}

	c.mu.Lock()
	c.items[base] = item
	c.mu.Unlock()

	return item, nil
}

func (c *Cache) load(base string) (item Item, err error) {
	if item.TextTmpl, err = texttemplate.ParseFS(
		c.fsys, base+c.spec.TextExt,
	); err != nil {
		return Item{}, err
	}

	if c.spec.HtmlExt == "" {
		return item, nil
	}

	if item.HtmlTmpl, err = htmltemplate.ParseFS(
		c.fsys, base+c.spec.HtmlExt,
	); err != nil {
		return Item{}, err
	}

	return item, nil
}
