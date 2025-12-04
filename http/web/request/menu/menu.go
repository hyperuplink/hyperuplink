package menu

import "github.com/mrusme/hyperuplink/models/user"

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
				Href:  "session/signin",
			},
			{
				IsSeparator: true,
			},
			{
				Label: "Sign up",
				Title: "Sign up",
				Href:  "session/signup",
			},
		}
	} else {
		subItems = []MenuItem{
			{
				Label: "Profile",
				Title: "Profile",
				Href:  "account/profile",
			},
			{
				Label: "Settings",
				Title: "Settings",
				Href:  "account/settings",
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
						Href:  "account/password",
					},
					{
						Label: "2-Factor Authentication",
						Title: "2-Factor Authentication",
						Href:  "account/2fa",
					},
				},
			},
			{
				IsSeparator: true,
			},
			{
				Label: "Sign out",
				Title: "Sign out",
				Href:  "session/signout",
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
						Href:  "admin/general",
					},
					{
						Label: "Authentication",
						Title: "Authentication",
						Href:  "admin/auth",
					},
					{
						Label: "Communication",
						SubItems: []MenuItem{
							{
								Label: "E-Mail",
								Title: "E-Mail",
								Href:  "admin/comms/email",
							},
							{
								Label: "XMPP",
								Title: "XMPP",
								Href:  "admin/comms/xmpp",
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
								Href:  "admin/board/categories",
							},
							{
								Label: "Forums",
								Title: "Forums",
								Href:  "admin/board/forums",
							},
							{
								IsSeparator: true,
							},
							{
								Label: "Posts",
								Title: "Posts",
								Href:  "admin/board/posts",
							},
							{
								Label: "Attachments",
								Title: "Attachments",
								Href:  "admin/board/attachments",
							},
							{
								Label: "Profiles",
								Title: "Profiles",
								Href:  "admin/board/profiles",
							},
							{
								Label: "Signatures",
								Title: "Signatures",
								Href:  "admin/board/signatures",
							},
							{
								IsSeparator: true,
							},
							{
								Label: "Look &amp; Feel",
								Title: "Look & Feel",
								Href:  "admin/board/theme",
							},
						},
					},
					{
						IsSeparator: true,
					},
					{
						Label: "Users",
						Title: "Users",
						Href:  "admin/users",
					},
					{
						IsSeparator: true,
					},
					{
						Label: "Adminlog",
						Title: "Adminlog",
						Href:  "admin/log",
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
							Href:       "account/settings?mode=light",
						},
						{
							IsCheckbox: true,
							Checked:    false,
							Label:      "Dark",
							Title:      "Dark",
							Href:       "account/settings?mode=dark",
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
					Href:       "account/settings?banner=false",
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Footer",
					Title:      "Footer",
					Href:       "account/settings?footer=false",
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Profile Pictures",
					Title:      "Profile Pictures",
					Href:       "account/settings?profile_pictures=false",
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
					Href:  "",
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Terms of Service",
					Title: "Terms of Service",
					Href:  "",
				},
				{
					Label: "Privacy Policy",
					Title: "Privacy Policy",
					Href:  "",
				},
				{
					Label: "Contact",
					Title: "Contact",
					Href:  "",
				},
				{
					IsSeparator: true,
				},
				{
					Label: "About",
					Title: "About",
					Href:  "version",
				},
			},
		},
	}
}

func (m *Menu) Get() []MenuItem {
	return m.menuItems
}
