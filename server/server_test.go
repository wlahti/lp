package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wlahti/lp/server"
	"github.com/wlahti/lp/server/mock"
)

func TestGetUserData(t *testing.T) {
	fakeNameGetter := &mock.NameGetter{}
	fakeNotesGetter := &mock.NotesGetter{}
	fakeTasksGetter := &mock.TasksGetter{}

	h := server.NewHTTPHandler(fakeNameGetter, fakeNotesGetter, fakeTasksGetter)

	t.Run("green path", func(t *testing.T) {
		fakeNameGetter.GetNameReturns("avocado", nil)
		fakeNotesGetter.GetNotesReturns(
			[]string{
				"test note 1",
				"test note 2",
				"test note 3",
				"test note 4",
				"test note 5",
				"test note 6",
				"test note 7",
				"test note 8",
				"test note 9",
				"test note 10",
				"test note 11",
				"test note 12",
			},
		)
		fakeTasksGetter.GetTasksReturns(
			[]string{
				"test task 1",
				"test task 2",
				"test task 3",
				"test task 4",
				"test task 5",
				"test task 6",
				"test task 7",
				"test task 8",
				"test task 9",
				"test task 10",
				"test task 11",
				"test task 12",
			},
		)

		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/1234", nil)
		h.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Result().StatusCode)
		require.Equal(t, "application/json", resp.Result().Header.Get("Content-Type"))

		userData := server.UserDataResponse{}
		err := json.Unmarshal(resp.Body.Bytes(), &userData)
		require.NoError(t, err, "cannot be unmarshaled")
		require.Equal(t, server.UserDataResponse{
			UserID:   1234,
			UserName: "avocado",
			Notes: []string{
				"test note 3",
				"test note 4",
				"test note 5",
				"test note 6",
				"test note 7",
				"test note 8",
				"test note 9",
				"test note 10",
				"test note 11",
				"test note 12",
			},
			Tasks: []string{
				"test task 3",
				"test task 4",
				"test task 5",
				"test task 6",
				"test task 7",
				"test task 8",
				"test task 9",
				"test task 10",
				"test task 11",
				"test task 12",
			},
		}, userData)
	})
}
