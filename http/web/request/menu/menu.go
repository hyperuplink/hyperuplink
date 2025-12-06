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
				Label: m.T("sign_in_verb"),
				Title: m.T("sign_in_verb"),
				Href:  route.For("SessionSignin").AsURL(),
			},
			{
				IsSeparator: true,
			},
			{
				Label: m.T("sign_up_verb"),
				Title: m.T("sign_up_verb"),
				Href:  route.For("SessionSignup").AsURL(),
			},
		}
	} else {
		subItems = []MenuItem{
			{
				Label: m.T("profile"),
				Title: m.T("profile"),
				Href:  route.For("AccountProfile").AsURL(),
			},
			{
				Label: m.T("settings"),
				Title: m.T("settings"),
				Href:  route.For("AccountSettings").AsURL(),
			},
			{
				IsSeparator: true,
			},
			{
				Label: m.T("security"),
				SubItems: []MenuItem{
					{
						Label: m.T("password"),
						Title: m.T("password"),
						Href:  route.For("AccountPassword").AsURL(),
					},
					{
						Label: m.T("2_factor_authentication"),
						Title: m.T("2_factor_authentication"),
						Href:  route.For("AccountTwofactor").AsURL(),
					},
				},
			},
			{
				IsSeparator: true,
			},
			{
				Label: m.T("sign_out_verb"),
				Title: m.T("sign_out_verb"),
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
						Label: m.T("general"),
						Title: m.T("general"),
						Href:  route.For("AdminGeneral").AsURL(),
					},
					{
						Label: m.T("authentication"),
						Title: m.T("authentication"),
						Href:  route.For("AdminAuth").AsURL(),
					},
					{
						Label: m.T("communication"),
						SubItems: []MenuItem{
							{
								Label: m.T("email"),
								Title: m.T("email"),
								Href:  route.For("AdminCommsEmail").AsURL(),
							},
							{
								Label: m.T("xmpp"),
								Title: m.T("xmpp"),
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
								Label: m.T("categories"),
								Title: m.T("categories"),
								Href:  route.For("AdminBoardCategories").AsURL(),
							},
							{
								Label: m.T("forums"),
								Title: m.T("forums"),
								Href:  route.For("AdminBoardForums").AsURL(),
							},
							{
								IsSeparator: true,
							},
							{
								Label: m.T("posts"),
								Title: m.T("posts"),
								Href:  route.For("AdminBoardPosts").AsURL(),
							},
							{
								Label: m.T("attachments"),
								Title: m.T("attachments"),
								Href:  route.For("AdminBoardAttachments").AsURL(),
							},
							{
								Label: m.T("profiles"),
								Title: m.T("profiles"),
								Href:  route.For("AdminBoardProfiles").AsURL(),
							},
							{
								Label: m.T("signatures"),
								Title: m.T("signatures"),
								Href:  route.For("AdminBoardSignatures").AsURL(),
							},
							{
								IsSeparator: true,
							},
							{
								Label: m.T("look_and_feel"),
								Title: m.T("look_and_feel"),
								Href:  route.For("AdminBoardTheme").AsURL(),
							},
						},
					},
					{
						IsSeparator: true,
					},
					{
						Label: m.T("users"),
						Title: m.T("users"),
						Href:  route.For("AdminUsers").AsURL(),
					},
					{
						IsSeparator: true,
					},
					{
						Label: m.T("adminlog"),
						Title: m.T("adminlog"),
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
					Label: m.T("user_manual"),
					Title: m.T("user_manual"),
					Href:  route.For("DocsManual").AsURL(),
				},
				{
					IsSeparator: true,
				},
				{
					Label: m.T("terms_of_service"),
					Title: m.T("terms_of_service"),
					Href:  route.For("DocsTerms").AsURL(),
				},
				{
					Label: m.T("privacy_policy"),
					Title: m.T("privacy_policy"),
					Href:  route.For("DocsPrivacy").AsURL(),
				},
				{
					Label: m.T("contact"),
					Title: m.T("contact"),
					Href:  route.For("DocsContact").AsURL(),
				},
				{
					IsSeparator: true,
				},
				{
					Label: m.T("about"),
					Title: m.T("about"),
					Href:  route.For("DocsAbout").AsURL(),
				},
			},
		},
	}
}

func (m *Menu) FooterMenu(forRole user.Role) []MenuItem {
	return []MenuItem{
		{
			Label: m.T("terms_of_service"),
			Title: m.T("terms_of_service"),
			Href:  route.For("DocsTerms").AsURL(),
		},
		{
			Label: m.T("privacy_policy"),
			Title: m.T("privacy_policy"),
			Href:  route.For("DocsPrivacy").AsURL(),
		},
		{
			Label: m.T("contact"),
			Title: m.T("contact"),
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
