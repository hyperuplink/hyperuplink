package signup

type Signup struct {
	Recipient Recipient
	Subject   string
	Token     string
	URL       string
}

type Recipient struct {
	Username string
	Address  string
}
