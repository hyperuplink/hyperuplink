package search

import (
	"strings"

	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vsearchresult"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	searchrepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/search"
)

type Input struct {
	Query        string `json:"q" form:"q" validate:"required,min=1,max=256"`
	Author       string `json:"author" form:"author" validate:"omitempty,max=32"`
	FTitle       bool   `json:"f_title" form:"f_title"`
	FBody        bool   `json:"f_body" form:"f_body"`
	FReplies     bool   `json:"f_replies" form:"f_replies"`
	FAttachments bool   `json:"f_attachments" form:"f_attachments"`
	Page         int    `json:"page" form:"page"`
	PerPage      int    `json:"-" form:"-"`
}

type Results struct {
	Results *[]vsearchresult.VSearchResult `json:"results"`
	Total   int64                          `json:"total"`
	Pages   int                            `json:"pages"`
}

func (in *Input) Normalize() {
	in.Query = strings.TrimSpace(in.Query)
	in.Author = strings.TrimSpace(in.Author)

	if !in.FTitle && !in.FBody && !in.FReplies && !in.FAttachments {
		in.FTitle, in.FBody, in.FReplies = true, true, true
	}

	if in.Page < 1 {
		in.Page = 1
	}
}

func Query(
	rt *runtime.Runtime,
	perms *permission.Resolution,
	in *Input,
) (results *Results, err error) {
	res, total, err := rt.Repositories.Search.Query(
		in.Query,
		searchrepo.Options{
			Title:                in.FTitle,
			Body:                 in.FBody,
			Replies:              in.FReplies,
			Attachments:          in.FAttachments,
			Author:               in.Author,
			AllowedCategorySlugs: perms.AllowedReadSlugs(),
		},
		in.PerPage,
		in.Page,
	)
	if err != nil {
		return nil, err
	}

	return &Results{
		Results: res,
		Total:   total,
		Pages:   paging.Pages(total, in.PerPage),
	}, nil
}
