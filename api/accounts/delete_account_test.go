package accounts_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/keratin/authn-server/api/accounts"
	"github.com/keratin/authn-server/api/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteAccount(t *testing.T) {
	app := test.App()
	server := test.Server(app, accounts.Routes(app))
	defer server.Close()

	client := test.NewClient(server).Authenticated(app.Config)

	t.Run("unknown account", func(t *testing.T) {
		res, err := client.Delete("/accounts/999999")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("unarchived account", func(t *testing.T) {
		account, err := app.AccountStore.Create("unlocked@test.com", []byte("bar"))
		require.NoError(t, err)

		res, err := client.Delete(fmt.Sprintf("/accounts/%v", account.Id))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		account, err = app.AccountStore.Find(account.Id)
		require.NoError(t, err)
		assert.NotEmpty(t, account.DeletedAt)
	})

	t.Run("archived account", func(t *testing.T) {
		account, err := app.AccountStore.Create("locked@test.com", []byte("bar"))
		require.NoError(t, err)
		app.AccountStore.Archive(account.Id)

		res, err := client.Delete(fmt.Sprintf("/accounts/%v", account.Id))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		account, err = app.AccountStore.Find(account.Id)
		require.NoError(t, err)
		assert.NotEmpty(t, account.DeletedAt)
	})
}