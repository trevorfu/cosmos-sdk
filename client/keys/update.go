package keys

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	keys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"

	"github.com/spf13/cobra"
)

func updateKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <name>",
		Args:    cobra.ExactArgs(1),
		Short:   "Change the password used to protect private key",
		Long:    "Change the password used to protect private key",
		Example: "TODO",
		RunE:    runUpdateCmd,
	}
	return cmd
}

func runUpdateCmd(cmd *cobra.Command, args []string) error {
	name := args[0]

	buf := client.BufferStdin()
	kb, err := GetKeyBase()
	if err != nil {
		return err
	}
	oldpass, err := client.GetPassword(
		"Enter the current passphrase:", buf)
	if err != nil {
		return err
	}

	getNewpass := func() (string, error) {
		return client.GetCheckPassword(
			"Enter the new passphrase:",
			"Repeat the new passphrase:", buf)
	}

	err = kb.Update(name, oldpass, getNewpass)
	if err != nil {
		return err
	}
	fmt.Println("Password successfully updated!")
	return nil
}

///////////////////////
// REST

// update key request REST body
type UpdateKeyBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

// update key REST handler
func UpdateKeyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var kb keys.Keybase
	var m UpdateKeyBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	kb, err = GetKeyBase()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	getNewpass := func() (string, error) { return m.NewPassword, nil }

	// TODO check if account exists and if password is correct
	err = kb.Update(name, m.OldPassword, getNewpass)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
}
