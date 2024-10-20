package files

import (
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return e.Wrap("cant save", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return e.Wrap("cant save", err)
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return e.Wrap("cant create file", err)
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap("cant serialize", err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	filePath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, e.Wrap("cant pick random", err)
	}

	if len(files) == 0 {
		return nil, e.Wrap("no saved page", err)
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(filePath, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return e.Wrap("cant remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(path); err != nil {
		return e.Wrap("cant remove file", err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("cant remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, err
	case err != nil:
		return false, e.Wrap("cant check exists", err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("cant decode page", err)
	}
	defer f.Close()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("cant decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
