package email

type Email struct {
	Recipient string
	Subject   string
	Content   string
}

type Mailer interface {
	SendEmail(Email) error
}
