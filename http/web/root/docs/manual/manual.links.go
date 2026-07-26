package manual

import (
	"strings"
	"text/template"

	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
)

func linkFuncs(
	req *request.Request,
	myRoute route.Route,
) template.FuncMap {
	return template.FuncMap{
		"hrefTo": req.Site.HrefTo,
		"manual": func(target string) string {
			target = strings.Trim(target, "/")
			if target == "" {
				return req.Site.HrefTo(myRoute.AsURL() + "/")
			}
			return req.Site.HrefTo(myRoute.AsURL() + "/" + target + "/")
		},
	}
}

func expandLinks(
	req *request.Request,
	myRoute route.Route,
	src string,
) (dst string, err error) {
	var tmpl *template.Template
	if tmpl, err = template.New("manual").
		Funcs(linkFuncs(req, myRoute)).
		Parse(src); err != nil {
		return dst, err
	}

	var buf strings.Builder
	if err = tmpl.Execute(&buf, nil); err != nil {
		return dst, err
	}

	return buf.String(), nil
}
