package local

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type Storage struct {
	lockAcquired bool
	lockFilePath string
}

func (ls *Storage) LockFileExists() bool {
	_, err := os.Stat(ls.lockFilePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (ls *Storage) Lock() {
	f, err := os.OpenFile(ls.lockFilePath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("error while writing: %v", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf("Lock(): %v", err)
		}
	}()
	if _, err := f.Write([]byte("wow")); err != nil {
		log.Errorf("error while writing: %v", err)
		return
	}
	ls.lockAcquired = true
}

func (ls *Storage) GetLockContent() (string, error) {
	cont, err := os.ReadFile(ls.lockFilePath)
	if err != nil {
		return "", fmt.Errorf("GetLockContent(): %v", err)
	}
	return string(cont), nil
}

func (ls *Storage) HasLock() bool {
	return ls.lockAcquired
}

func (ls *Storage) DeleteLock() {
	err := os.Remove(ls.lockFilePath)
	if err != nil {
		log.Errorf("DeleteLock(): %v", err)
	}
	ls.lockAcquired = false
}

func NewLocalStorage(lockFilePath string) Storage {
	return Storage{false, lockFilePath}
}
