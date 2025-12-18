package bcn

type BreadcrumbNavigation struct {
	breadcrumbs []Breadcrumb
}

type Breadcrumb struct {
	IsActive bool
	Label    string
	Title    string
	Href     string
}

func New() *BreadcrumbNavigation {
	b := new(BreadcrumbNavigation)

	return b
}

func NewBreadcrumb(
	isActive bool,
	label string,
	title string,
	href string,
) *Breadcrumb {
	b := new(Breadcrumb)

	b.IsActive = isActive
	b.Label = label
	b.Title = title
	b.Href = href

	return b
}

func (b *BreadcrumbNavigation) Get() []Breadcrumb {
	return b.breadcrumbs
}

func (b *BreadcrumbNavigation) Set(bcs []Breadcrumb) {
	b.breadcrumbs = bcs
}

func (b *BreadcrumbNavigation) Append(bc Breadcrumb) {
	b.breadcrumbs = append(b.breadcrumbs, bc)
}

func (b *BreadcrumbNavigation) Prepend(bc Breadcrumb) {
	b.breadcrumbs = append([]Breadcrumb{bc}, b.breadcrumbs...)
}

func (b *BreadcrumbNavigation) UpdateLabel(revIdx int, label string) {
	bl := len(b.breadcrumbs)
	if bl >= revIdx {
		b.breadcrumbs[bl-revIdx].Label = label
		b.breadcrumbs[bl-revIdx].Title = label
	}
}

func (b *BreadcrumbNavigation) UpdateHref(revIdx int, href string) {
	bl := len(b.breadcrumbs)
	if bl >= revIdx {
		b.breadcrumbs[bl-revIdx].Href = href
	}
}
