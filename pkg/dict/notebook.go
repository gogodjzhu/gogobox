package dict

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type FileNotebook struct {
	filename string
}

func NewFileNotebook(filename string) (*FileNotebook, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			return nil, errors.New("[Err] create notebook filename failed")
		}
		f.Close()
	}
	return &FileNotebook{
		filename: filename,
	}, nil
}

func (f *FileNotebook) Mark(word string, action Action) error {
	notes, err := f.readNote()
	if err != nil {
		return err
	}
	// filter note
	var note *WordNote
	for _, n := range notes {
		if n.Word == word {
			note = n
			break
		}
	}
	if note == nil {
		note = &WordNote{
			Word: word,
		}
		notes = append(notes, note)
	}
	switch action {
	case Learning:
		note.LookupTimes++
		note.LastLookupTime = time.Now().Unix()
	case Learned:
		note.LookupTimes--
		note.LastLookupTime = time.Now().Unix()
	case Delete:
		// delete note
		var newNotes []*WordNote
		for _, n := range notes {
			if n.Word != word {
				newNotes = append(newNotes, n)
			}
		}
		return f.writeNote(newNotes)
	default:
		return errors.New("[Err] invalid action")
	}
	return f.writeNote(notes)
}

func (f *FileNotebook) Get(word string) (*WordNote, error) {
	notes, err := f.readNote()
	if err != nil {
		return nil, err
	}
	// filter note
	for _, note := range notes {
		if note.Word == word {
			return note, nil
		}
	}
	return nil, nil
}

func (f *FileNotebook) Review() (*WordNote, error) {
	notes, err := f.readNote()
	if err != nil {
		return nil, err
	}
	// select one word from notes map randomly
	for _, note := range notes {
		return note, nil
	}
	return nil, nil
}

func (f *FileNotebook) List() ([]*WordNote, error) {
	notes, err := f.readNote()
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (f *FileNotebook) readNote() ([]*WordNote, error) {
	bytes, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return nil, errors.New("[Err] read notebook file failed")
	}
	var notes []*WordNote
	err = yaml.Unmarshal(bytes, &notes)
	if err != nil {
		return nil, errors.New("[Err] unmarshal notebook file failed")
	}

	// sort notes by lookup times
	sort.SliceStable(notes, func(i, j int) bool {
		if notes[i].LookupTimes == notes[j].LookupTimes {
			return notes[i].LastLookupTime > notes[j].LastLookupTime
		} else {
			return notes[i].LookupTimes > notes[j].LookupTimes
		}
	})
	return notes, nil
}

func (f *FileNotebook) writeNote(notes []*WordNote) error {
	bytes, err := yaml.Marshal(notes)
	if err != nil {
		return errors.New("[Err] marshal notebook file failed")
	}
	err = ioutil.WriteFile(f.filename, bytes, 0666)
	if err != nil {
		return errors.New("[Err] write notebook file failed")
	}
	return nil
}
