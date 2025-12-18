package site

import "strconv"

type Pager struct {
	Pages      int
	PerPage    int
	ActivePage int
}

type PagerItem struct {
	Page     int
	Label    string
	Active   bool
	Disabled bool
}

type PagerView []PagerItem

func NewPager(
	pages int,
	perPage int,
	activePage int,
) (pg *Pager) {
	pg = new(Pager)
	pg.Pages = pages
	pg.PerPage = perPage
	pg.ActivePage = activePage

	return pg
}

func (s *Site) SetPager(pg *Pager) {
	s.pager = pg
}

func (s *Site) GetPager() *Pager {
	return s.pager
}

func (pg *Pager) GetView() (pgv PagerView) {
	maxPerSide := 2

	if pg.ActivePage < 1 {
		pg.ActivePage = 1
	}
	if pg.ActivePage > pg.Pages {
		pg.ActivePage = pg.Pages
	}

	lowerBoundary := (pg.ActivePage - maxPerSide)
	upperBoundary := (pg.ActivePage + maxPerSide)

	for i := lowerBoundary; i <= upperBoundary; i++ {
		if i == lowerBoundary {
			pgv = append(pgv, PagerItem{
				Page:     1,
				Label:    "|&lt;",
				Active:   false,
				Disabled: (i <= 1),
			})
			pgv = append(pgv, PagerItem{
				Page:     i,
				Label:    "&lt;-",
				Active:   false,
				Disabled: (i < 1),
			})
		} else if i == upperBoundary {
			pgv = append(pgv, PagerItem{
				Page:     i,
				Label:    "-&gt;",
				Active:   false,
				Disabled: (i > pg.Pages),
			})
			pgv = append(pgv, PagerItem{
				Page:     pg.Pages,
				Label:    "&gt;|",
				Active:   false,
				Disabled: (i >= pg.Pages),
			})
		} else {
			if i < 1 || i > pg.Pages {
				continue
			}
			pgv = append(pgv, PagerItem{
				Page:     i,
				Label:    strconv.Itoa(i),
				Active:   (i == pg.ActivePage),
				Disabled: false,
			})
		}
	}

	return pgv
}
