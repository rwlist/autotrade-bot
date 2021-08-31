package chatex

import (
	"os"
	"testing"

	chatexsdk "github.com/chatex-com/sdk-go"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func getClient(t *testing.T) *chatexsdk.Client {
	token := os.Getenv("CHATEX_TEST_TOKEN")
	if token == "" {
		t.Skip("no chatex token is provided")
	}
	cli := chatexsdk.NewClient("https://api.chatex.com/v1", token)
	return cli
}

func TestChatex_GetAllRates(t *testing.T) {
	cli := getClient(t)
	srv := NewChatex(cli, nil)
	rates, err := srv.GetAllRates(USDT)
	assert.NoError(t, err)

	spew.Dump(rates)
}

func TestChatex_AccountBalance(t *testing.T) {
	cli := getClient(t)
	srv := NewChatex(cli, nil)
	b, err := srv.AccountBalance()
	assert.NoError(t, err)

	spew.Dump(b)
}
