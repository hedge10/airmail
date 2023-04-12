package mail

import (
	"errors"
	"testing"

	"github.com/hedge10/airmail/pkg/conf"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransfert(t *testing.T) {
	type test struct {
		input string
		want  error
	}

	tests := []test{
		{input: "mailgun", want: nil},
		{input: "smtp", want: nil},
		{input: "unknown", want: errors.New("cannot create transfer")},
	}

	for _, tc := range tests {
		c := &conf.Config{
			MailService: tc.input,
		}
		_, got := CreateTransfer(c)
		if tc.want != nil {
			assert.Equal(t, tc.want.Error(), got.Error())
		} else {
			assert.Nil(t, got)
		}
	}

}
