package setting

type AddressType int

const (
	EmailOnly   AddressType = 0
	JID         AddressType = 1
	EmailAndJID AddressType = 2
)

type Auth struct {
	AddressType AddressType `json:"allowed_address_type"` // 0: EmailOnly, 1: JID, 2: Email and JID
}
