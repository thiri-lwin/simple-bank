package api

import (
	"context"
	"encoding/json"
	"net/http"

	db "github.com/thiri-lwin/thiri-bank/db/sqlc"
)

func (server *MuxServer) createAccount(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var arg createAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&arg); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(respondJson(err))
		return
	}

	acc, err := server.store.CreateAccount(context.Background(), db.CreateAccountParams{
		Owner:    arg.Owner,
		Balance:  0,
		Currency: arg.Currency,
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(errorResponse(err))
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(&acc)
}
