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
							Href:  "account/password", // TODO: ...
						},
						{
							Label: "2-Factor Authentication",
							Title: "2-Factor Authentication",
							Href:  "account/2fa", // TODO: ...
						},
					},
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Sign out",
					Title: "Sign out",
					Href:  "session/signout", // TODO: ...
				},
			},
		},
		{
			Label: "Ad<u>m</u>inistration",
			SubItems: []MenuItem{
				{
					Label: "General",
					Title: "General",
					Href:  "admin/general", // TODO: ...
				},
				{
					Label: "Authentication",
					Title: "Authentication",
					Href:  "admin/auth", // TODO: ...
				},
				{
					Label: "Communication",
					SubItems: []MenuItem{
						{
							Label: "E-Mail",
							Title: "E-Mail",
							Href:  "admin/comms/email", // TODO: ...
						},
						{
							Label: "XMPP",
							Title: "XMPP",
							Href:  "admin/comms/xmpp", // TODO: ...
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
							Href:  "admin/board/categories", // TODO: ...
						},
						{
							Label: "Forums",
							Title: "Forums",
							Href:  "admin/board/forums", // TODO: ...
						},
						{
							IsSeparator: true,
						},
						{
							Label: "Posts",
							Title: "Posts",
							Href:  "admin/board/posts", // TODO: ...
						},
						{
							Label: "Attachments",
							Title: "Attachments",
							Href:  "admin/board/attachments", // TODO: ...
						},
						{
							Label: "Profiles",
							Title: "Profiles",
							Href:  "admin/board/profiles", // TODO: ...
						},
						{
							Label: "Signatures",
							Title: "Signatures",
							Href:  "admin/board/signatures", // TODO: ...
						},
						{
							IsSeparator: true,
						},
						{
							Label: "Look &amp; Feel",
							Title: "Look & Feel",
							Href:  "admin/board/theme", // TODO: ...
						},
					},
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Users",
					Title: "Users",
					Href:  "admin/users", // TODO: ...
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Adminlog",
					Title: "Adminlog",
					Href:  "admin/log", // TODO: ...
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
							Href:       "account/settings?mode=light", // TODO: ...
						},
						{
							IsCheckbox: true,
							Checked:    false,
							Label:      "Dark",
							Title:      "Dark",
							Href:       "account/settings?mode=dark", // TODO: ...
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
					Href:       "account/settings?banner=false", // TODO: ...
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Footer",
					Title:      "Footer",
					Href:       "account/settings?footer=false", // TODO: ...
				},
				{
					IsCheckbox: true,
					Checked:    true,
					Label:      "Profile Pictures",
					Title:      "Profile Pictures",
					Href:       "account/settings?profile_pictures=false", // TODO: ...
				},
			},
		},
		{
			Label: "<u>H</u>elp",
			SubItems: []MenuItem{
				{
					Label: "User Manual",
					Title: "User Manual",
					Href:  "", // TODO: ...
				},
				{
					IsSeparator: true,
				},
				{
					Label: "Terms of Service",
					Title: "Terms of Service",
					Href:  "", // TODO: ...
				},
				{
					Label: "Privacy Policy",
					Title: "Privacy Policy",
					Href:  "", // TODO: ...
				},
				{
					Label: "Contact",
					Title: "Contact",
					Href:  "", // TODO: ...
				},
				{
					IsSeparator: true,
				},
				{
					Label: "About",
					Title: "About",
					Href:  "version", // TODO: ...
				},
			},
		},
	}
}
