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

func (s *Site) Menu() []MenuItem {
	return []MenuItem{
		{
			Label: "<u>A</u>ccount",
			SubItems: []MenuItem{
				{
					Label: "Settings",
					Title: "Settings",
					Href:  s.HrefTo("account/settings"), // TODO: Method for retrieving the relUrl
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
							Href:  s.HrefTo("account/password"), // TODO: ...
						},
						{
							Label: "2-Factor Authentication",
							Title: "2-Factor Authentication",
							Href:  s.HrefTo("account/2fa"), // TODO: ...
						},
					},
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Sign out",
					Title: "Sign out",
					Href:  s.HrefTo("session/signout"), // TODO: ...
				},
			},
		},
		{
			Label: "Ad<u>m</u>inistration",
			SubItems: []MenuItem{
				{
					Label: "General",
					Title: "General",
					Href:  s.HrefTo("admin/general"), // TODO: ...
				},
				{
					Label: "Authentication",
					Title: "Authentication",
					Href:  s.HrefTo("admin/auth"), // TODO: ...
				},
				{
					Label: "Communication",
					SubItems: []MenuItem{
						{
							Label: "E-Mail",
							Title: "E-Mail",
							Href:  s.HrefTo("admin/comms/email"), // TODO: ...
						},
						{
							Label: "XMPP",
							Title: "XMPP",
							Href:  s.HrefTo("admin/comms/xmpp"), // TODO: ...
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
							Href:  s.HrefTo("admin/board/categories"), // TODO: ...
						},
						{
							Label: "Forums",
							Title: "Forums",
							Href:  s.HrefTo("admin/board/forums"), // TODO: ...
						},
						{
							IsSeparator: true,
						},
						{
							Label: "Posts",
							Title: "Posts",
							Href:  s.HrefTo("admin/board/posts"), // TODO: ...
						},
						{
							Label: "Attachments",
							Title: "Attachments",
							Href:  s.HrefTo("admin/board/attachments"), // TODO: ...
						},
						{
							Label: "Profiles",
							Title: "Profiles",
							Href:  s.HrefTo("admin/board/profiles"), // TODO: ...
						},
						{
							Label: "Signatures",
							Title: "Signatures",
							Href:  s.HrefTo("admin/board/signatures"), // TODO: ...
						},
						{
							IsSeparator: true,
						},
						{
							Label: "Look &amp; Feel",
							Title: "Look & Feel",
							Href:  s.HrefTo("admin/board/theme"), // TODO: ...
						},
					},
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Users",
					Title: "Users",
					Href:  s.HrefTo("admin/users"), // TODO: ...
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Adminlog",
					Title: "Adminlog",
					Href:  s.HrefTo("admin/log"), // TODO: ...
				},
			},
		},
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
							Href:       s.HrefTo("account/settings?mode=light"), // TODO: ...
						},
						{
							IsCheckbox: true,
							Checked:    false,
							Label:      "Dark",
							Title:      "Dark",
							Href:       s.HrefTo("account/settings?mode=dark"), // TODO: ...
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
					Href:       s.HrefTo("account/settings?banner=false"), // TODO: ...
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Footer",
					Title:      "Footer",
					Href:       s.HrefTo("account/settings?footer=false"), // TODO: ...
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Profile Pictures",
					Title:      "Profile Pictures",
					Href:       s.HrefTo("account/settings?profile_pictures=false"), // TODO: ...
				},
			},
		},
		{
			Label: "<u>H</u>elp",
			SubItems: []MenuItem{
				{
					Label: "User Manual",
					Title: "User Manual",
					Href:  s.HrefTo(""), // TODO: ...
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Terms of Service",
					Title: "Terms of Service",
					Href:  s.HrefTo(""), // TODO: ...
				},
				{
					Label: "Privacy Policy",
					Title: "Privacy Policy",
					Href:  s.HrefTo(""), // TODO: ...
				},
				{
					Label: "Contact",
					Title: "Contact",
					Href:  s.HrefTo(""), // TODO: ...
				},
				{
					IsSeparator: true,
				},
				{
					Label: "About",
					Title: "About",
					Href:  s.HrefTo("version"), // TODO: ...
				},
			},
		},
	}
}
