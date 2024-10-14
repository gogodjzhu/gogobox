package dict

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gogobox/internal/config"
	"gogobox/pkg/dict/entity"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Action string

const (
	Learning Action = "learning"
	Learned  Action = "learned"
	Delete   Action = "delete"
)

type Notebooks interface {
	Mark(word string, action Action) (*entity.WordNote, error)
	ListNotes() ([]*entity.WordNote, error)
	ListChapter() ([]string, error)
}

func OpenNotebook(conf *config.NotebookConfig) (Notebooks, error) {
	currentNotebook := conf.CurrentChapter
	/* fileNotebook */
	if conf.FileNotebookConfig != nil {
		return openOrCreateFileNotebook(conf.FileNotebookConfig, currentNotebook)
	}
	return nil, errors.New("[Err] notebook named " + currentNotebook + " not found")
}

type fileNotebook struct {
	directory    string
	notebookName string
	filename     string
}

func openOrCreateFileNotebook(conf *config.NotebookFileNotebookConfig, notebookName string) (*fileNotebook, error) {
	filename := filepath.Join(conf.Directory, notebookName+".yaml")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := os.MkdirAll(conf.Directory, 0755)
		if err != nil {
			return nil, errors.New("[Err] create notebook directory failed")
		}
		f, err := os.Create(filename)
		if err != nil {
			return nil, errors.New("[Err] create notebook chapter file failed")
		}
		_ = f.Close()
	} else if err != nil {
		return nil, errors.New("[Err] open notebook chapter file failed")
	}
	return &fileNotebook{
		directory:    conf.Directory,
		notebookName: notebookName,
		filename:     filename,
	}, nil
}

func (f *fileNotebook) Mark(word string, action Action) (*entity.WordNote, error) {
	notes, err := f.readNote()
	if err != nil {
		return nil, err
	}
	note := &entity.WordNote{
		WordItemId:     entity.WordId(word),
		Word:           word,
		LookupTimes:    0,
		CreateTime:     time.Now().Unix(),
		LastLookupTime: time.Now().Unix(),
	}
	isOld := false
	for _, n := range notes {
		if n.WordItemId == entity.WordId(word) {
			note, isOld = n, true
			break
		}
	}
	if !isOld {
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
		var newNotes []*entity.WordNote
		for _, n := range notes {
			if n.WordItemId != entity.WordId(word) {
				newNotes = append(newNotes, n)
			}
		}
		notes = newNotes
	default:
		return nil, errors.New("[Err] invalid action:" + string(action))
	}
	return note, f.writeNote(notes)
}

func (f *fileNotebook) ListNotes() ([]*entity.WordNote, error) {
	notes, err := f.readNote()
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (f *fileNotebook) ListChapter() ([]string, error) {
	files, err := os.ReadDir(f.directory)
	if err != nil {
		return nil, errors.New("[Err] list notebook directory failed")
	}
	var chapters []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if !strings.HasSuffix(filename, ".yaml") {
			log.Warnf("ignore invalid notebook file: %s", filename)
			continue
		}
		filename = strings.TrimSuffix(filename, ".yaml")
		chapters = append(chapters, filename)
	}
	return chapters, nil
}

func (f *fileNotebook) readNote() ([]*entity.WordNote, error) {
	bytes, err := os.ReadFile(f.filename)
	if err != nil {
		return nil, errors.New("[Err] read notebook file failed")
	}
	var notes []*entity.WordNote
	err = yaml.Unmarshal(bytes, &notes)
	if err != nil {
		return nil, errors.New("[Err] unmarshal notebook file failed")
	}

	// sort notes by lookup times
	sort.SliceStable(notes, func(i, j int) bool {
		return notes[i].CreateTime > notes[j].CreateTime
	})
	return notes, nil
}

func (f *fileNotebook) writeNote(notes []*entity.WordNote) error {
	bytes, err := yaml.Marshal(notes)
	if err != nil {
		return errors.New("[Err] marshal notebook file failed")
	}
	tmpFilename := f.filename + ".tmp"
	err = os.WriteFile(tmpFilename, bytes, 0666)
	if err != nil {
		return errors.New("[Err] write notebook file failed")
	}
	err = os.Rename(tmpFilename, f.filename)
	if err != nil {
		return err
	}
	return nil
}

type sqlNotebook struct {
	db           *gorm.DB
	notebookName string
}

type SQLNotebookWordNote struct {
	WordId         string `gorm:"word_id"`
	Chapter        string `gorm:"chapter"`
	Word           string `gorm:"word"`
	LookupTimes    int    `gorm:"lookup_times"`
	CreateTime     int64  `gorm:"create_time"`
	LastLookupTime int64  `gorm:"last_lookup_time"`
}

func (s *SQLNotebookWordNote) TableName() string {
	return "word_note"
}

func (s *SQLNotebookWordNote) toWordNote() *entity.WordNote {
	return &entity.WordNote{
		WordItemId:     s.WordId,
		Word:           s.Word,
		LookupTimes:    s.LookupTimes,
		CreateTime:     s.CreateTime,
		LastLookupTime: s.LastLookupTime,
	}
}

func (s *sqlNotebook) Mark(word string, action Action) (*entity.WordNote, error) {
	var sql string
	var params []interface{}
	now := time.Now().Unix()
	switch action {
	case Learning:
		sql = "INSERT INTO word_note (word_id, chapter, word, lookup_times, create_time, last_lookup_time) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE lookup_times = lookup_times + 1, last_lookup_time = ?"
		params = []interface{}{entity.WordId(word), s.notebookName, word, 1, now, now, now}
	case Learned:
		sql = "INSERT INTO word_note (word_id, chapter, word, lookup_times, create_time, last_lookup_time) VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE lookup_times = lookup_times - 1, last_lookup_time = ?"
		params = []interface{}{entity.WordId(word), s.notebookName, word, 1, now, now, now}
	case Delete:
		sql = "DELETE FROM word_note WHERE chapter = ? AND word_id = ?"
		params = []interface{}{s.notebookName, entity.WordId(word)}
	}
	tx := s.db.Exec(sql, params...)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "[Err] mark word failed")
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("[Err] mark word failed")
	}
	if action == Delete {
		return nil, nil
	}
	var updateWordNote SQLNotebookWordNote
	tx = s.db.First(&updateWordNote, "chapter = ? AND word_id = ?", s.notebookName, entity.WordId(word))
	tx.Order("create_time DESC")
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "[Err] get word note failed")
	}
	return updateWordNote.toWordNote(), nil
}

func (s *sqlNotebook) ListNotes() ([]*entity.WordNote, error) {
	var notes []*SQLNotebookWordNote
	tx := s.db.Find(&notes, "chapter = ?", s.notebookName)
	tx.Order("create_time DESC")
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "[Err] list word note failed")
	}
	var wordNotes []*entity.WordNote
	for _, note := range notes {
		wordNotes = append(wordNotes, note.toWordNote())
	}
	return wordNotes, nil
}

func (s *sqlNotebook) ListChapter() ([]string, error) {
	var chapters []string
	tx := s.db.Model(&SQLNotebookWordNote{}).Distinct("chapter").Find(&chapters)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "[Err] list chapter failed")
	}
	return chapters, nil
}
