package server

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type Note struct {
	Timestamp time.Time
	Text      string
}

type Notes []Note

type NoteCache struct {
	fileName string
	lock     sync.RWMutex
	// map of userid -> notes
	notes map[int]Notes
}

func NewNoteCache(fileName string) *NoteCache {
	return &NoteCache{
		fileName: fileName,
		notes:    map[int]Notes{},
	}
}

// readNotes is run by the server in the background
// to continuously read the notes.csv file for new entries
func (n *NoteCache) readNotes() error {
	csvFile, err := os.Open(n.fileName)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)

	for {
		record, err := r.Read()
		if err == io.EOF {
			return errors.New("end of file")
		}
		if err != nil {
			return err
		}

		// timestamp is long, microseconds
		// timestamp := record[0]

		// userid is int
		userID := record[1]
		id, err := strconv.Atoi(userID)
		if err != nil {
			return fmt.Errorf("invalid user_id: %s", userID)
		}

		// note is string
		note := record[2]

		n.lock.Lock()
		n.notes[id] = append(n.notes[id],
			Note{
				// Timestamp: timestamp,
				Text: note,
			})
		n.lock.Unlock()
	}
}

func (n *NoteCache) GetNotes(id int) []string {
	n.lock.RLock()
	defer n.lock.RUnlock()

	notes := make([]string, len(n.notes[id]))
	for i, note := range n.notes[id] {
		notes[i] = note.Text
	}

	return notes
}
