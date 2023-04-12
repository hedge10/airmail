package mail

type (
	Mailgun struct {
		Domain      string
		PrivateKey  string
		UseEuDomain bool
	}
)

func (mg Mailgun) Send(e *Email) error {

	return nil
}
