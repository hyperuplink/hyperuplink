package site

type MenuItem struct {
	IsSeparator bool
	IsCheckbox  bool
	Checked     bool
	Label       string
	Title       string
	Href        string
	SubItems    []MenuItem
}

func (s *Site) AccountMenu() []MenuItem {
	var subItems []MenuItem

	if false { // TODO: Check for user role
		subItems = []MenuItem{
			{
				Label: "Profile",
				Title: "Profile",
				Href:  s.HrefTo("account/profile"),
			},
			{
				Label: "Settings",
				Title: "Settings",
				Href:  s.HrefTo("account/settings"),
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
						Href:  s.HrefTo("account/password"),
					},
					{
						Label: "2-Factor Authentication",
						Title: "2-Factor Authentication",
						Href:  s.HrefTo("account/2fa"),
					},
				},
			},
			{
				IsSeparator: true,
			},
			{
				Label: "Sign out",
				Title: "Sign out",
				Href:  s.HrefTo("sessions/signout"),
			},
		}
	} else {
		subItems = []MenuItem{
			{
				Label: "Sign in",
				Title: "Sign in",
				Href:  s.HrefTo("sessions/signin"),
			},
			{
				IsSeparator: true,
			},
			{
				Label: "Sign up",
				Title: "Sign up",
				Href:  s.HrefTo("sessions/signup"),
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

func (s *Site) AdminMenu() []MenuItem {
	// TODO: Check user role
	return []MenuItem{}

	return []MenuItem{
		{
			Label: "Ad<u>m</u>inistration",
			SubItems: []MenuItem{
				{
					Label: "General",
					Title: "General",
					Href:  s.HrefTo("admin/general"),
				},
				{
					Label: "Authentication",
					Title: "Authentication",
					Href:  s.HrefTo("admin/auth"),
				},
				{
					Label: "Communication",
					SubItems: []MenuItem{
						{
							Label: "E-Mail",
							Title: "E-Mail",
							Href:  s.HrefTo("admin/comms/email"),
						},
						{
							Label: "XMPP",
							Title: "XMPP",
							Href:  s.HrefTo("admin/comms/xmpp"),
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
							Href:  s.HrefTo("admin/board/categories"),
						},
						{
							Label: "Forums",
							Title: "Forums",
							Href:  s.HrefTo("admin/board/forums"),
						},
						{
							IsSeparator: true,
						},
						{
							Label: "Posts",
							Title: "Posts",
							Href:  s.HrefTo("admin/board/posts"),
						},
						{
							Label: "Attachments",
							Title: "Attachments",
							Href:  s.HrefTo("admin/board/attachments"),
						},
						{
							Label: "Profiles",
							Title: "Profiles",
							Href:  s.HrefTo("admin/board/profiles"),
						},
						{
							Label: "Signatures",
							Title: "Signatures",
							Href:  s.HrefTo("admin/board/signatures"),
						},
						{
							IsSeparator: true,
						},
						{
							Label: "Look &amp; Feel",
							Title: "Look & Feel",
							Href:  s.HrefTo("admin/board/theme"),
						},
					},
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Users",
					Title: "Users",
					Href:  s.HrefTo("admin/users"),
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Adminlog",
					Title: "Adminlog",
					Href:  s.HrefTo("admin/log"),
				},
			},
		},
	}
}

func (s *Site) ViewMenu() []MenuItem {
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
							Href:       s.HrefTo("account/settings?mode=light"),
						},
						{
							IsCheckbox: true,
							Checked:    false,
							Label:      "Dark",
							Title:      "Dark",
							Href:       s.HrefTo("account/settings?mode=dark"),
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
					Href:       s.HrefTo("account/settings?banner=false"),
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Footer",
					Title:      "Footer",
					Href:       s.HrefTo("account/settings?footer=false"),
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Profile Pictures",
					Title:      "Profile Pictures",
					Href:       s.HrefTo("account/settings?profile_pictures=false"),
				},
			},
		},
	}
}

func (s *Site) HelpMenu() []MenuItem {
	return []MenuItem{
		{
			Label: "<u>H</u>elp",
			SubItems: []MenuItem{
				{
					Label: "User Manual",
					Title: "User Manual",
					Href:  s.HrefTo(""),
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Terms of Service",
					Title: "Terms of Service",
					Href:  s.HrefTo(""),
				},
				{
					Label: "Privacy Policy",
					Title: "Privacy Policy",
					Href:  s.HrefTo(""),
				},
				{
					Label: "Contact",
					Title: "Contact",
					Href:  s.HrefTo(""),
				},
				{
					IsSeparator: true,
				},
				{
					Label: "About",
					Title: "About",
					Href:  s.HrefTo("version"),
				},
			},
		},
	}
}

func (s *Site) Menu() []MenuItem {
	var menu []MenuItem

	menu = append(menu, s.AccountMenu()...)
	menu = append(menu, s.AdminMenu()...)
	menu = append(menu, s.ViewMenu()...)
	menu = append(menu, s.HelpMenu()...)

	return menu
}
