package menu

import (
	"fmt"

	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/models/user"
)

type Menu struct {
	menuItems           []MenuItem
	footerMenuItems     []MenuItem
	role                user.Role
	i18n                func(msg string) string
	currentCategorySlug string
	currentForumSlug    string
	general             setting.General
}

type MenuItem struct {
	IsSeparator bool
	IsCheckbox  bool
	Checked     bool
	Disabled    bool
	Label       string
	Title       string
	Href        string
	SubItems    []MenuItem
}

func New() (m *Menu) {
	m = new(Menu)

	return m
}

func (m *Menu) SetRole(role user.Role) {
	m.role = role
	m.generate()
}

func (m *Menu) SetI18n(fn func(msg string) string) {
	m.i18n = fn
}

func (m *Menu) SetGeneral(general setting.General) {
	m.general = general
	m.generate()
}

func (m *Menu) SetCategoryForumSlugs(catSlug, forumSlug string) {
	m.currentCategorySlug = catSlug
	m.currentForumSlug = forumSlug
	m.generate()
}

func (m *Menu) T(msg string) string {
	return m.i18n(msg)
}

func (m *Menu) generate() {
	m.menuItems = []MenuItem{}

	m.menuItems = append(m.menuItems, m.FileMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.AccountMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.AdminMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.ViewMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.HelpMenu(m.role)...)

	m.footerMenuItems = []MenuItem{}
	m.footerMenuItems = append(m.footerMenuItems, m.FooterMenu(m.role)...)
}

func (m *Menu) FileMenu(forRole user.Role) []MenuItem {
	var subItems []MenuItem

	newpostUrl := route.For("New").AsURL()
	if m.currentCategorySlug != "" && m.currentForumSlug != "" {
		newpostUrl = fmt.Sprintf("%s?category=%s&forum=%s",
			newpostUrl, m.currentCategorySlug, m.currentForumSlug)
	}

	subItems = []MenuItem{
		{
			Disabled: (forRole == user.GuestRole),
			Label:    m.T("new"),
			Title:    m.T("new"),
			Href:     newpostUrl,
		},
		// TODO: Implement Open for recently viewed posts
		// {
		// 	Label: m.T("open"),
		// 	SubItems: []MenuItem{
		// 		{ // TODO: List of previously opened posts
		// 			Label: m.T("Hetzner experiences?"),
		// 			Title: m.T("Hetzner experiences?"),
		// 			Href:  route.For("CategoriesForumsTopics").AsURL(),
		// 		},
		// 	},
		// },
		{
			IsSeparator: true,
		},
		{
			Label: m.T(route.For("Search").AsTitle()),
			Title: m.T(route.For("Search").AsTitle()),
			Href:  route.For("Search").AsURL(),
		},
		{
			IsSeparator: true,
		},
		{
			Label: m.T("rss"),
			Title: m.T("rss"),
			Href:  route.For("current").AsURL() + "?format=rss", // TODO: Implement current route
		},
		{
			Label: m.T("print"),
			Title: m.T("print"),
			Href:  route.For("current").AsURL() + "?format=print", // TODO: Implement current route
		},
	}

	if m.general.EnableQuit {
		subItems = append(subItems,
			MenuItem{
				IsSeparator: true,
			},
			MenuItem{
				Label: m.T("quit"),
				Title: m.T("quit"),
				Href:  m.general.QuitURL,
			},
		)
	}

	return []MenuItem{
		{
			Label:    m.T("_f_file"),
			SubItems: subItems,
		},
	}
}

func (m *Menu) AccountMenu(forRole user.Role) []MenuItem {
	var subItems []MenuItem

	if forRole == user.GuestRole {
		subItems = []MenuItem{
			{
				Label: m.T(route.For("SessionSignin").AsTitle()),
				Title: m.T(route.For("SessionSignin").AsTitle()),
				Href:  route.For("SessionSignin").AsURL(),
			},
			{
				IsSeparator: true,
			},
			{
				Label: m.T(route.For("SessionSignup").AsTitle()),
				Title: m.T(route.For("SessionSignup").AsTitle()),
				Href:  route.For("SessionSignup").AsURL(),
			},
		}
	} else {
		subItems = []MenuItem{
			{
				Label: m.T(route.For("AccountProfile").AsTitle()),
				Title: m.T(route.For("AccountProfile").AsTitle()),
				Href:  route.For("AccountProfile").AsURL(),
			},
			{
				Label: m.T(route.For("AccountSettings").AsTitle()),
				Title: m.T(route.For("AccountSettings").AsTitle()),
				Href:  route.For("AccountSettings").AsURL(),
			},
			{
				IsSeparator: true,
			},
			{
				Label: m.T("security"),
				SubItems: []MenuItem{
					{
						Label: m.T(route.For("AccountPassword").AsTitle()),
						Title: m.T(route.For("AccountPassword").AsTitle()),
						Href:  route.For("AccountPassword").AsURL(),
					},
					{
						Label: m.T(route.For("AccountTwofactor").AsTitle()),
						Title: m.T(route.For("AccountTwofactor").AsTitle()),
						Href:  route.For("AccountTwofactor").AsURL(),
					},
				},
			},
			{
				IsSeparator: true,
			},
			{
				Label: m.T(route.For("SessionSignout").AsTitle()),
				Title: m.T(route.For("SessionSignout").AsTitle()),
				Href:  route.For("SessionSignout").AsURL(),
			},
		}
	}

	return []MenuItem{
		{
			Label:    m.T("_a_ccount"),
			SubItems: subItems,
		},
	}
}

func (m *Menu) AdminBoardMenu(forRole user.Role) []MenuItem {
	if forRole == user.AdminRole {
		return []MenuItem{
			{
				Label: m.T("board_settings"),
				SubItems: []MenuItem{
					{
						Label: m.T(route.For("AdminBoardCategories").AsTitle()),
						Title: m.T(route.For("AdminBoardCategories").AsTitle()),
						Href:  route.For("AdminBoardCategories").AsURL(),
					},
					{
						Label: m.T(route.For("AdminBoardForums").AsTitle()),
						Title: m.T(route.For("AdminBoardForums").AsTitle()),
						Href:  route.For("AdminBoardForums").AsURL(),
					},
					{
						IsSeparator: true,
					},
					{
						Label: m.T(route.For("AdminBoardTopics").AsTitle()),
						Title: m.T(route.For("AdminBoardTopics").AsTitle()),
						Href:  route.For("AdminBoardTopics").AsURL(),
					},
					{
						Label: m.T(route.For("AdminBoardAttachments").AsTitle()),
						Title: m.T(route.For("AdminBoardAttachments").AsTitle()),
						Href:  route.For("AdminBoardAttachments").AsURL(),
					},
					{
						Label: m.T(route.For("AdminBoardProfiles").AsTitle()),
						Title: m.T(route.For("AdminBoardProfiles").AsTitle()),
						Href:  route.For("AdminBoardProfiles").AsURL(),
					},
					{
						IsSeparator: true,
					},
					{
						Label: m.T(route.For("AdminBoardThemes").AsTitle()),
						Title: m.T(route.For("AdminBoardThemes").AsTitle()),
						Href:  route.For("AdminBoardThemes").AsURL(),
					},
				},
			},
		}
	}
	return []MenuItem{}
}

func (m *Menu) AdminCommsMenu(forRole user.Role) []MenuItem {
	if forRole == user.AdminRole {
		return []MenuItem{
			{
				Label: m.T("communication"),
				SubItems: []MenuItem{
					{
						Label: m.T(route.For("AdminCommsEmail").AsTitle()),
						Title: m.T(route.For("AdminCommsEmail").AsTitle()),
						Href:  route.For("AdminCommsEmail").AsURL(),
					},
					{
						Label: m.T(route.For("AdminCommsXmpp").AsTitle()),
						Title: m.T(route.For("AdminCommsXmpp").AsTitle()),
						Href:  route.For("AdminCommsXmpp").AsURL(),
					},
				},
			},
		}
	}
	return []MenuItem{}
}

func (m *Menu) AdminMenu(forRole user.Role) []MenuItem {
	if forRole == user.AdminRole {
		return []MenuItem{
			{
				Label: m.T("ad_m_inistration"),
				SubItems: []MenuItem{
					{
						Label: m.T(route.For("AdminGeneral").AsTitle()),
						Title: m.T(route.For("AdminGeneral").AsTitle()),
						Href:  route.For("AdminGeneral").AsURL(),
					},
					{
						Label: m.T(route.For("AdminAuth").AsTitle()),
						Title: m.T(route.For("AdminAuth").AsTitle()),
						Href:  route.For("AdminAuth").AsURL(),
					},
					m.AdminCommsMenu(forRole)[0],
					{
						IsSeparator: true,
					},
					m.AdminBoardMenu(forRole)[0],
					{
						IsSeparator: true,
					},
					{
						Label: m.T(route.For("AdminUsers").AsTitle()),
						Title: m.T(route.For("AdminUsers").AsTitle()),
						Href:  route.For("AdminUsers").AsURL(),
					},
					{
						IsSeparator: true,
					},
					{
						Label: m.T(route.For("AdminLogs").AsTitle()),
						Title: m.T(route.For("AdminLogs").AsTitle()),
						Href:  route.For("AdminLogs").AsURL(),
					},
				},
			},
		}
	}

	return []MenuItem{}
}

func (m *Menu) ViewMenu(forRole user.Role) []MenuItem {
	return []MenuItem{
		{
			Label: m.T("_v_iew"),
			SubItems: []MenuItem{
				{
					Label: m.T("mode"),
					SubItems: []MenuItem{
						{
							IsCheckbox: true,
							Checked:    true,
							Label:      m.T("light"),
							Title:      m.T("light"),
							Href:       route.For("Session").AsURL() + "?mode=light",
						},
						{
							IsCheckbox: true,
							Checked:    false,
							Label:      m.T("dark"),
							Title:      m.T("dark"),
							Href:       route.For("Session").AsURL() + "?mode=dark",
						},
					},
				},
				{
					IsSeparator: true,
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      m.T("banner"),
					Title:      m.T("banner"),
					Href:       route.For("Session").AsURL() + "?banner=false",
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      m.T("footer"),
					Title:      m.T("footer"),
					Href:       route.For("Session").AsURL() + "?footer=false",
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      m.T("profile_pictures"),
					Title:      m.T("profile_pictures"),
					Href:       route.For("Session").AsURL() + "?profile_pictures=false",
				},
			},
		},
	}
}

func (m *Menu) HelpMenu(forRole user.Role) []MenuItem {
	subItems := []MenuItem{
		{
			Label: m.T(route.For("DocsManual").AsTitle()),
			Title: m.T(route.For("DocsManual").AsTitle()),
			Href:  route.For("DocsManual").AsURL() + "/",
		},
	}

	if m.general.EnableTerms || m.general.EnablePrivacyPolicy || m.general.EnableContact {
		subItems = append(subItems, MenuItem{IsSeparator: true})
		if m.general.EnableTerms {
			subItems = append(subItems, MenuItem{
				Label: m.T(route.For("DocsTerms").AsTitle()),
				Title: m.T(route.For("DocsTerms").AsTitle()),
				Href:  route.For("DocsTerms").AsURL(),
			})
		}
		if m.general.EnablePrivacyPolicy {
			subItems = append(subItems, MenuItem{
				Label: m.T(route.For("DocsPrivacy").AsTitle()),
				Title: m.T(route.For("DocsPrivacy").AsTitle()),
				Href:  route.For("DocsPrivacy").AsURL(),
			})
		}
		if m.general.EnableContact {
			subItems = append(subItems, MenuItem{
				Label: m.T(route.For("DocsContact").AsTitle()),
				Title: m.T(route.For("DocsContact").AsTitle()),
				Href:  route.For("DocsContact").AsURL(),
			})
		}
	}

	if m.general.EnableAbout {
		subItems = append(subItems, MenuItem{IsSeparator: true})
		subItems = append(subItems, MenuItem{
			Label: m.T(route.For("DocsAbout").AsTitle()),
			Title: m.T(route.For("DocsAbout").AsTitle()),
			Href:  route.For("DocsAbout").AsURL(),
		})
	}

	return []MenuItem{
		{
			Label:    m.T("_h_elp"),
			SubItems: subItems,
		},
	}
}

func (m *Menu) FooterMenu(forRole user.Role) []MenuItem {
	var items []MenuItem

	if m.general.EnableTerms {
		items = append(items, MenuItem{
			Label: m.T(route.For("DocsTerms").AsTitle()),
			Title: m.T(route.For("DocsTerms").AsTitle()),
			Href:  route.For("DocsTerms").AsURL(),
		})
	}
	if m.general.EnablePrivacyPolicy {
		items = append(items, MenuItem{
			Label: m.T(route.For("DocsPrivacy").AsTitle()),
			Title: m.T(route.For("DocsPrivacy").AsTitle()),
			Href:  route.For("DocsPrivacy").AsURL(),
		})
	}
	if m.general.EnableContact {
		items = append(items, MenuItem{
			Label: m.T(route.For("DocsContact").AsTitle()),
			Title: m.T(route.For("DocsContact").AsTitle()),
			Href:  route.For("DocsContact").AsURL(),
		})
	}

	return items
}

func (m *Menu) Get() []MenuItem {
	return m.menuItems
}

func (m *Menu) GetFooter() []MenuItem {
	return m.footerMenuItems
}
