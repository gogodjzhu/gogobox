package dict

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

type FileNotebook struct {
	filename string
}

func NewFileNotebook(filename string) (*FileNotebook, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			return nil, errors.New("create notebook filename failed")
		}
		f.Close()
	}
	return &FileNotebook{
		filename: filename,
	}, nil
}

func (f *FileNotebook) Mark(word string, action Action) error {
	note, err := f.Get(word)
	if err != nil {
		return err
	}
	if note == nil {
		note = &WordNote{
			Word: word,
		}
	}
	switch action {
	case Learning:
		note.LookupTimes++
	case Learned:
		note.LookupTimes--
	default:
		return errors.New("invalid action")
	}
	note.LastLookupTime = time.Now().Unix()
	return f.writeNote(note)
}

func (f *FileNotebook) Get(word string) (*WordNote, error) {
	notes, err := f.readNote()
	if err != nil {
		return nil, err
	}
	return notes[word], nil
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

func (f *FileNotebook) readNote() (map[string]*WordNote, error) {
	file, err := os.OpenFile(f.filename, os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.New("open notebook file failed")
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("read notebook file failed")
	}
	notes := make(map[string]*WordNote)
	if len(bytes) == 0 {
		return notes, nil
	}
	err = json.Unmarshal(bytes, &notes)
	if err != nil {
		return nil, errors.New("unmarshal notebook file failed")
	}
	return notes, nil
}

func (f *FileNotebook) writeNote(note *WordNote) error {
	file, err := os.OpenFile(f.filename, os.O_RDWR, 0666)
	if err != nil {
		return errors.New("open notebook file failed")
	}
	defer file.Close()
	notes, err := f.readNote()
	if err != nil {
		return err
	}
	notes[note.Word] = note
	bytes, err := json.MarshalIndent(notes, "", " ")
	if err != nil {
		return errors.New("marshal notebook file failed")
	}
	_, err = file.Write(bytes)
	if err != nil {
		return errors.New("write notebook file failed")
	}
	return nil
}
