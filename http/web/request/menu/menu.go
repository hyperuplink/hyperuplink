package menu

import (
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/models/user"
)

type Menu struct {
	menuItems       []MenuItem
	footerMenuItems []MenuItem
	role            user.Role
	i18n            func(msg string) string
}

type MenuItem struct {
	IsSeparator bool
	IsCheckbox  bool
	Checked     bool
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

func (m *Menu) T(msg string) string {
	return m.i18n(msg)
}

func (m *Menu) generate() {
	m.menuItems = []MenuItem{}

	m.menuItems = append(m.menuItems, m.AccountMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.AdminMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.ViewMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.HelpMenu(m.role)...)

	m.footerMenuItems = []MenuItem{}
	m.footerMenuItems = append(m.footerMenuItems, m.FooterMenu(m.role)...)
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
					{
						IsSeparator: true,
					},
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
								Label: m.T(route.For("AdminBoardSignatures").AsTitle()),
								Title: m.T(route.For("AdminBoardSignatures").AsTitle()),
								Href:  route.For("AdminBoardSignatures").AsURL(),
							},
							{
								IsSeparator: true,
							},
							{
								Label: m.T(route.For("AdminBoardTheme").AsTitle()),
								Title: m.T(route.For("AdminBoardTheme").AsTitle()),
								Href:  route.For("AdminBoardTheme").AsURL(),
							},
						},
					},
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
						Label: m.T(route.For("AdminLog").AsTitle()),
						Title: m.T(route.For("AdminLog").AsTitle()),
						Href:  route.For("AdminLog").AsURL(),
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
							Href:       route.For("SessionSettings").AsURL() + "?mode=light",
						},
						{
							IsCheckbox: true,
							Checked:    false,
							Label:      m.T("dark"),
							Title:      m.T("dark"),
							Href:       route.For("SessionSettings").AsURL() + "?mode=dark",
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
					Href:       route.For("SessionSettings").AsURL() + "?banner=false",
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      m.T("footer"),
					Title:      m.T("footer"),
					Href:       route.For("SessionSettings").AsURL() + "?footer=false",
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      m.T("profile_pictures"),
					Title:      m.T("profile_pictures"),
					Href:       route.For("SessionSettings").AsURL() + "?profile_pictures=false",
				},
			},
		},
	}
}

func (m *Menu) HelpMenu(forRole user.Role) []MenuItem {
	return []MenuItem{
		{
			Label: m.T("_h_elp"),
			SubItems: []MenuItem{
				{
					Label: m.T(route.For("DocsManual").AsTitle()),
					Title: m.T(route.For("DocsManual").AsTitle()),
					Href:  route.For("DocsManual").AsURL(),
				},
				{
					IsSeparator: true,
				},
				{
					Label: m.T(route.For("DocsTerms").AsTitle()),
					Title: m.T(route.For("DocsTerms").AsTitle()),
					Href:  route.For("DocsTerms").AsURL(),
				},
				{
					Label: m.T(route.For("DocsPrivacy").AsTitle()),
					Title: m.T(route.For("DocsPrivacy").AsTitle()),
					Href:  route.For("DocsPrivacy").AsURL(),
				},
				{
					Label: m.T(route.For("DocsContact").AsTitle()),
					Title: m.T(route.For("DocsContact").AsTitle()),
					Href:  route.For("DocsContact").AsURL(),
				},
				{
					IsSeparator: true,
				},
				{
					Label: m.T(route.For("DocsAbout").AsTitle()),
					Title: m.T(route.For("DocsAbout").AsTitle()),
					Href:  route.For("DocsAbout").AsURL(),
				},
			},
		},
	}
}

func (m *Menu) FooterMenu(forRole user.Role) []MenuItem {
	return []MenuItem{
		{
			Label: m.T(route.For("DocsTerms").AsTitle()),
			Title: m.T(route.For("DocsTerms").AsTitle()),
			Href:  route.For("DocsTerms").AsURL(),
		},
		{
			Label: m.T(route.For("DocsPrivacy").AsTitle()),
			Title: m.T(route.For("DocsPrivacy").AsTitle()),
			Href:  route.For("DocsPrivacy").AsURL(),
		},
		{
			Label: m.T(route.For("DocsContact").AsTitle()),
			Title: m.T(route.For("DocsContact").AsTitle()),
			Href:  route.For("DocsContact").AsURL(),
		},
	}
}

func (m *Menu) Get() []MenuItem {
	return m.menuItems
}

func (m *Menu) GetFooter() []MenuItem {
	return m.footerMenuItems
}
