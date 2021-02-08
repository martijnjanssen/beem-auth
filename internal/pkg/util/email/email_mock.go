package email

import (
	"github.com/stretchr/testify/mock"
)

type EmailMock struct {
	mock.Mock
}

func NewEmailMock() *EmailMock {
	return &EmailMock{}
}

func (e *EmailMock) SendEmail(email Email) error {
	args := e.MethodCalled("SendEmail", email)
	return args.Error(0)
}
