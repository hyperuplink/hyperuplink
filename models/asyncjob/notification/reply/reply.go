package reply

type Reply struct {
	Recipient  Recipient
	Subject    string
	ByUsername string
	Text       string
	HTML       string
	URL        string
	ReplyTo    string
	Category   Category
	Forum      Forum
	Topic      Topic
}

type Recipient struct {
	Username string
	Address  string
	Lang     string
}

type Topic struct {
	Title string
	URL   string
}

type Forum struct {
	Title string
	URL   string
}

type Category struct {
	Title string
	URL   string
}
