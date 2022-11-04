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
	"github.com/mrityunjaygr8/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func randomEntry() db.Entry {
	return db.Entry{
		ID:        utils.RandomInt(1, 1000),
		AccountID: utils.RandomInt(1, 1000),
		Amount:    utils.RandomMoney(),
	}

}

func TestListEntries(t *testing.T) {
	var entries []db.Entry
	for x := 0; x < 10; x++ {
		entries = append(entries, randomEntry())

	}

	testCases := []struct {
		name          string
		page_size     int32
		page_id       int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			page_size: 5,
			page_id:   1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntries(gomock.Any(), gomock.Eq(db.ListEntriesParams{
					Limit:  5,
					Offset: 0,
				})).Times(1).Return(entries[0:5], nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchListEntry(t, recorder.Body, entries[0:5])
			},
		},
		{
			name:      "Zero-Page",
			page_size: 5,
			page_id:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntries(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Big-Page-Size",
			page_size: 50,
			page_id:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntries(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Internal-error",
			page_size: 5,
			page_id:   1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntries(gomock.Any(), gomock.Eq(db.ListEntriesParams{
					Limit:  5,
					Offset: 0,
				})).Times(1).Return([]db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/entries?page_size=%d&page_id=%d", tc.page_size, tc.page_id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchListEntry(t *testing.T, body *bytes.Buffer, entries []db.Entry) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotEntries []db.Entry

	err = json.Unmarshal(data, &gotEntries)
	require.NoError(t, err)

	require.Equal(t, entries, gotEntries)
}
func TestListEntriesForAccount(t *testing.T) {
	var entries []db.Entry
	for x := 0; x < 10; x++ {
		entry := randomEntry()
		entry.AccountID = 1
		entries = append(entries, entry)
	}

	testCases := []struct {
		name          string
		page_size     int32
		page_id       int32
		account_id    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			page_size:  5,
			page_id:    1,
			account_id: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntriesForAccount(gomock.Any(), gomock.Eq(db.ListEntriesForAccountParams{
					Limit:     5,
					Offset:    0,
					AccountID: 1,
				})).Times(1).Return(entries[0:5], nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchListEntry(t, recorder.Body, entries[0:5])
			},
		},
		{
			name:       "Zero-Page",
			page_size:  5,
			page_id:    0,
			account_id: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntriesForAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "Big-Page-Size",
			page_size:  50,
			page_id:    0,
			account_id: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntriesForAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "Internal-error",
			page_size:  5,
			page_id:    1,
			account_id: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntriesForAccount(gomock.Any(), gomock.Eq(db.ListEntriesForAccountParams{
					Limit:     5,
					Offset:    0,
					AccountID: 1,
				})).Times(1).Return([]db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "Internal-error",
			page_size:  5,
			page_id:    1,
			account_id: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListEntriesForAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/entries/%d?page_size=%d&page_id=%d", tc.account_id, tc.page_size, tc.page_id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
