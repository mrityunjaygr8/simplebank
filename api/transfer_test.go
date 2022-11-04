package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/mrityunjaygr8/simplebank/db/mock"
	db "github.com/mrityunjaygr8/simplebank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferApi(t *testing.T) {
	account1 := randomAccount()
	account1.Currency = "USD"
	account2 := randomAccount()
	account2.Currency = "USD"
	amount := int64(10)

	transferResponse := db.TransferTxResult{
		Transfer: db.Transfer{
			ID:            1,
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        amount,
		},
		FromAccount: db.Account{
			ID:        account1.ID,
			Owner:     account1.Owner,
			Currency:  account1.Currency,
			CreatedAt: account1.CreatedAt,
			Balance:   account1.Balance - amount,
		},
		ToAccount: db.Account{
			ID:        account2.ID,
			Owner:     account2.Owner,
			Currency:  account2.Currency,
			CreatedAt: account2.CreatedAt,
			Balance:   account2.Balance + amount,
		},
		FromEntry: db.Entry{
			ID:        1,
			AccountID: account1.ID,
			Amount:    -amount,
		},
		ToEntry: db.Entry{
			ID:        2,
			AccountID: account2.ID,
			Amount:    amount,
		},
	}

	testCases := []struct {
		name             string
		account1         db.Account
		account2         db.Account
		amount           int64
		currency         string
		transferParams   db.TransferTxParams
		transferResponse db.TransferTxResult
		buildStubs       func(store *mockdb.MockStore)
		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:             "ok",
			account1:         account1,
			account2:         account2,
			amount:           amount,
			currency:         "USD",
			transferResponse: transferResponse,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				})).Times(1).Return(transferResponse, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchCreateTransfer(t, recorder.Body, transferResponse)
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:             "Mismatch-currency",
			account1:         account1,
			account2:         account2,
			amount:           amount,
			currency:         "CAD",
			transferResponse: transferResponse,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatchError(t, recorder.Body, fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account1.ID, account1.Currency, "CAD"))
			},
		},
		{
			name:             "TransferTX-Internal-Error",
			account1:         account1,
			account2:         account2,
			amount:           amount,
			currency:         "USD",
			transferResponse: transferResponse,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				})).Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				// requireBodyMatchError(t, recorder.Body, fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account1.ID, account1.Currency, "CAD"))
			},
		},
		{
			name:             "crossed-currency",
			account1:         account1,
			account2:         db.Account{ID: account2.ID, Owner: account2.Owner, Balance: account2.Balance, CreatedAt: account2.CreatedAt, Currency: "CAD"},
			amount:           amount,
			currency:         "USD",
			transferResponse: transferResponse,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{ID: account2.ID, Owner: account2.Owner, Balance: account2.Balance, CreatedAt: account2.CreatedAt, Currency: "CAD"}, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatchError(t, recorder.Body, fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account2.ID, "CAD", "USD"))
			},
		},
		{
			name:             "account-missing",
			account1:         account1,
			account2:         account2,
			amount:           amount,
			currency:         "USD",
			transferResponse: transferResponse,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:             "internal-error-valid-account",
			account1:         account1,
			account2:         account2,
			amount:           amount,
			currency:         "USD",
			transferResponse: transferResponse,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:             "Invalid-Currency",
			account1:         account1,
			account2:         account2,
			amount:           amount,
			currency:         "INR",
			transferResponse: transferResponse,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatchError(t, recorder.Body, fmt.Sprintf("Key: 'transferRequestParams.Currency' Error:Field validation for 'Currency' failed on the 'currency' tag"))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewServer(store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/transfers")
			jsonStr := []byte(fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "currency": "%s", "amount": %d}`, tc.account1.ID, tc.account2.ID, tc.currency, tc.amount))
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
func requireBodyMatchCreateTransfer(t *testing.T, body *bytes.Buffer, transferResult db.TransferTxResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTransfer db.TransferTxResult

	err = json.Unmarshal(data, &gotTransfer)
	require.NoError(t, err)

	require.Equal(t, transferResult, gotTransfer)
}
func requireBodyMatchError(t *testing.T, body *bytes.Buffer, expected string) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotError struct {
		Error string `json:"error"`
	}

	err = json.Unmarshal(data, &gotError)
	require.NoError(t, err)

	require.Equal(t, expected, gotError.Error)
}
