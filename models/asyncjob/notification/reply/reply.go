package reply

type Reply struct {
	Recipient  Recipient
	Subject    string
	ByUsername string
	Text       string
	HTML       string
	URL        string
	ReplyTo    string
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
	Forum Forum
}

type Forum struct {
	Title    string
	URL      string
	Category Category
}

type Category struct {
	Title string
	URL   string
}
