package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	router      *mux.Router
	nameGetter  NameGetter
	notesGetter NotesGetter
	tasksGetter TasksGetter
}

//go:generate counterfeiter -o mock/name_getter.go -fake-name NameGetter . NameGetter
type NameGetter interface {
	GetName(int) (string, error)
}

//go:generate counterfeiter -o mock/notes_getter.go -fake-name NotesGetter . NotesGetter
type NotesGetter interface {
	GetNotes(int) []string
}

//go:generate counterfeiter -o mock/tasks_getter.go -fake-name TasksGetter . TasksGetter
type TasksGetter interface {
	GetTasks(int) []string
}

func NewHTTPHandler(nameGetter NameGetter, notesGetter NotesGetter, tasksGetter TasksGetter) *HTTPHandler {
	handler := &HTTPHandler{
		router:      mux.NewRouter(),
		nameGetter:  nameGetter,
		notesGetter: notesGetter,
		tasksGetter: tasksGetter,
	}

	handler.router.HandleFunc("/users/{userid}", handler.serveGetUserData).Methods(http.MethodGet)

	return handler
}

func (h *HTTPHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(resp, req)
}

func (h *HTTPHandler) serveGetUserData(resp http.ResponseWriter, req *http.Request) {
	var (
		userData *UserDataResponse
		err      error
	)
	doneChan := make(chan struct{})

	go func() {
		userData, err = h.getUserData(resp, req)
		close(doneChan)
	}()

	select {
	case <-req.Context().Done():
		// client closed the connection.
	case <-doneChan:
		// user data successfully retrieved
		if err != nil {
			h.sendResponseError(resp, http.StatusInternalServerError, err)
			return
		}
		h.sendResponseOK(resp, userData)
	}
}

func (h *HTTPHandler) getUserData(resp http.ResponseWriter, req *http.Request) (*UserDataResponse, error) {
	userID, err := h.extractUserID(req, resp)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(3)

	var name string
	go func() {
		defer wg.Done()
		var err error
		name, err = h.nameGetter.GetName(userID)
		if err != nil {
			cancel()
		}
	}()

	var notes []string
	go func() {
		defer wg.Done()
		for {
			select {
			case <-req.Context().Done():
				// client closed the connection.
				return
			case <-ctx.Done():
				// one of the other user data retrieval
				// goroutines failed
				return
			default:
				notes = h.notesGetter.GetNotes(userID)
				notes, err = trimSlice(notes, 10)
				if err == nil {
					return
				}
			}
		}
	}()

	var tasks []string
	go func() {
		defer wg.Done()
		for {
			select {
			case <-req.Context().Done():
				// client closed the connection.
				return
			case <-ctx.Done():
				// one of the other user data retrieval
				// goroutines failed
				return
			default:
				tasks = h.tasksGetter.GetTasks(userID)
				tasks, err = trimSlice(tasks, 10)
				if err == nil {
					return
				}
			}
		}
	}()
	wg.Wait()

	return &UserDataResponse{
		UserID:   userID,
		UserName: name,
		Notes:    notes,
		Tasks:    tasks,
	}, nil
}

func trimSlice(slice []string, trimLength int) ([]string, error) {
	index := len(slice) - trimLength
	if index < 0 {
		return nil, errors.New("slice is not long enough to trim")
	}

	slice = slice[index:]

	return slice, nil
}

func (h *HTTPHandler) extractUserID(req *http.Request, resp http.ResponseWriter) (int, error) {
	userID, ok := mux.Vars(req)["userid"]
	if !ok {
		err := errors.New("no userID specified in request")
		return 0, err
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid non-integer userid: %s", userID)
	}
	return id, nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *HTTPHandler) sendResponseError(resp http.ResponseWriter, code int, err error) {
	encoder := json.NewEncoder(resp)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(code)
	if err := encoder.Encode(&ErrorResponse{Error: err.Error()}); err != nil {
		log.Printf("failed to encode error. err=%v", err)
	}
}

type UserDataResponse struct {
	UserID   int      `json:"user_id"`
	UserName string   `json:"user_name"`
	Tasks    []string `json:"tasks"`
	Notes    []string `json:"notes"`
}

func (h *HTTPHandler) sendResponseOK(resp http.ResponseWriter, content interface{}) {
	encoder := json.NewEncoder(resp)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	if err := encoder.Encode(content); err != nil {
		log.Printf("failed to encode content. err=%v", err)
	}
}
