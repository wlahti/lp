package server

import (
	"database/sql"
	"log"
	"net/http"
)

type Server struct {
	noteCache  *NoteCache
	tasksCache *TasksCache
	userDB     *UserDB
}

func (s *Server) Start() {
	s.userDB = NewUserDB(sql.Open)
	err := s.userDB.connect("mysql", "user", "pass", "userdb")
	if err != nil {
		// TODO
		panic(err)
	}
	defer s.userDB.db.Close()

	// initialize and continously update the note cache
	s.noteCache = NewNoteCache("notes.csv")
	go s.noteCache.readNotes()

	// initialize and continously update the tasks cache
	s.tasksCache = NewTasksCache()
	go s.tasksCache.readTasks()

	h := NewHTTPHandler(s.userDB, s.noteCache, s.tasksCache)
	log.Fatal(http.ListenAndServe(":8090", h))
}
