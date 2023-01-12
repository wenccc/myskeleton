package email

type Driver interface {
	Send(email Email, config map[string]string) bool
}

type Email struct {
	From    From
	To      []string
	Bcc     []string
	Cc      []string
	Subject string
	Text    []byte // Plaintext message (optional)
	HTML    []byte // Html message (optional)
}

type From struct {
	Address string
	Name    string
}
