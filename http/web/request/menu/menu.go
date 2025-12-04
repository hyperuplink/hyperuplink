package menu

import (
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/models/user"
)

type Menu struct {
	menuItems []MenuItem
	role      user.Role
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

func (m *Menu) generate() {
	m.menuItems = []MenuItem{}

	m.menuItems = append(m.menuItems, m.AccountMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.AdminMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.ViewMenu(m.role)...)
	m.menuItems = append(m.menuItems, m.HelpMenu(m.role)...)
}

func (m *Menu) AccountMenu(forRole user.Role) []MenuItem {
	var subItems []MenuItem

	if forRole == user.GuestRole {
		subItems = []MenuItem{
			{
				Label: "Sign in",
				Title: "Sign in",
				Href:  route.For("SessionSignin").AsURL(),
			},
			{
				IsSeparator: true,
			},
			{
				Label: "Sign up",
				Title: "Sign up",
				Href:  route.For("SessionSignup").AsURL(),
			},
		}
	} else {
		subItems = []MenuItem{
			{
				Label: "Profile",
				Title: "Profile",
				Href:  route.For("AccountProfile").AsURL(),
			},
			{
				Label: "Settings",
				Title: "Settings",
				Href:  route.For("AccountSettings").AsURL(),
			},
			{
				IsSeparator: true,
			},
			{
				Label: "Security",
				SubItems: []MenuItem{
					{
						Label: "Password",
						Title: "Password",
						Href:  route.For("AccountPassword").AsURL(),
					},
					{
						Label: "2-Factor Authentication",
						Title: "2-Factor Authentication",
						Href:  route.For("AccountTwofactor").AsURL(),
					},
				},
			},
			{
				IsSeparator: true,
			},
			{
				Label: "Sign out",
				Title: "Sign out",
				Href:  route.For("SessionSignout").AsURL(),
			},
		}
	}

	return []MenuItem{
		{
			Label:    "<u>A</u>ccount",
			SubItems: subItems,
		},
	}
}

func (m *Menu) AdminMenu(forRole user.Role) []MenuItem {
	if forRole == user.AdminRole {
		return []MenuItem{
			{
				Label: "Ad<u>m</u>inistration",
				SubItems: []MenuItem{
					{
						Label: "General",
						Title: "General",
						Href:  route.For("AdminGeneral").AsURL(),
					},
					{
						Label: "Authentication",
						Title: "Authentication",
						Href:  route.For("AdminAuth").AsURL(),
					},
					{
						Label: "Communication",
						SubItems: []MenuItem{
							{
								Label: "E-Mail",
								Title: "E-Mail",
								Href:  route.For("AdminCommsEmail").AsURL(),
							},
							{
								Label: "XMPP",
								Title: "XMPP",
								Href:  route.For("AdminCommsXmpp").AsURL(),
							},
						},
					},
					{
						IsSeparator: true,
					},
					{
						Label: "Board Settings",
						SubItems: []MenuItem{
							{
								Label: "Categories",
								Title: "Categories",
								Href:  route.For("AdminBoardCategories").AsURL(),
							},
							{
								Label: "Forums",
								Title: "Forums",
								Href:  route.For("AdminBoardForums").AsURL(),
							},
							{
								IsSeparator: true,
							},
							{
								Label: "Posts",
								Title: "Posts",
								Href:  route.For("AdminBoardPosts").AsURL(),
							},
							{
								Label: "Attachments",
								Title: "Attachments",
								Href:  route.For("AdminBoardAttachments").AsURL(),
							},
							{
								Label: "Profiles",
								Title: "Profiles",
								Href:  route.For("AdminBoardProfiles").AsURL(),
							},
							{
								Label: "Signatures",
								Title: "Signatures",
								Href:  route.For("AdminBoardSignatures").AsURL(),
							},
							{
								IsSeparator: true,
							},
							{
								Label: "Look &amp; Feel",
								Title: "Look & Feel",
								Href:  route.For("AdminBoardTheme").AsURL(),
							},
						},
					},
					{
						IsSeparator: true,
					},
					{
						Label: "Users",
						Title: "Users",
						Href:  route.For("AdminUsers").AsURL(),
					},
					{
						IsSeparator: true,
					},
					{
						Label: "Adminlog",
						Title: "Adminlog",
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
			Label: "<u>V</u>iew",
			SubItems: []MenuItem{
				{
					Label: "Mode",
					SubItems: []MenuItem{
						{
							IsCheckbox: true,
							Checked:    true,
							Label:      "Light",
							Title:      "Light",
							Href:       route.For("SessionSettings").AsURL() + "?mode=light",
						},
						{
							IsCheckbox: true,
							Checked:    false,
							Label:      "Dark",
							Title:      "Dark",
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
					Label:      "Banner",
					Title:      "Banner",
					Href:       route.For("SessionSettings").AsURL() + "?banner=false",
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Footer",
					Title:      "Footer",
					Href:       route.For("SessionSettings").AsURL() + "?footer=false",
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Profile Pictures",
					Title:      "Profile Pictures",
					Href:       route.For("SessionSettings").AsURL() + "?profile_pictures=false",
				},
			},
		},
	}
}

func (m *Menu) HelpMenu(forRole user.Role) []MenuItem {
	return []MenuItem{
		{
			Label: "<u>H</u>elp",
			SubItems: []MenuItem{
				{
					Label: "User Manual",
					Title: "User Manual",
					Href:  route.For("DocsManual").AsURL(),
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Terms of Service",
					Title: "Terms of Service",
					Href:  route.For("DocsTerms").AsURL(),
				},
				{
					Label: "Privacy Policy",
					Title: "Privacy Policy",
					Href:  route.For("DocsPrivacy").AsURL(),
				},
				{
					Label: "Contact",
					Title: "Contact",
					Href:  route.For("DocsContact").AsURL(),
				},
				{
					IsSeparator: true,
				},
				{
					Label: "About",
					Title: "About",
					Href:  route.For("DocsAbout").AsURL(),
				},
			},
		},
	}
}

func (m *Menu) Get() []MenuItem {
	return m.menuItems
}
