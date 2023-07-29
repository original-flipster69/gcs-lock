package local

import (
	"fmt"
	"os"
)

type LocalStorage struct {
	lockAcquired bool
	lockFilePath string
}

func (ls *LocalStorage) LockFileExists() bool {
	_, err := os.Stat(ls.lockFilePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (ls *LocalStorage) Lock() {
	f, err := os.OpenFile(ls.lockFilePath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error while writing: %v", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Lock(): %v", err)
		}
	}()
	if _, err := f.Write([]byte("wow")); err != nil {
		fmt.Printf("error while writing: %v", err)
		return
	}
	ls.lockAcquired = true
}

func (ls *LocalStorage) GetLockContent() (string, error) {
	cont, err := os.ReadFile(ls.lockFilePath)
	if err != nil {
		return "", fmt.Errorf("GetLockContent(): %v", err)
	}
	return string(cont), nil
}

func (ls *LocalStorage) HasLock() bool {
	return ls.lockAcquired
}

func (ls *LocalStorage) DeleteLock() {
	err := os.Remove(ls.lockFilePath)
	if err != nil {
		fmt.Printf("DeleteLock(): %v", err)
	}
}

func (ls *LocalStorage) Unlock() {
	ls.lockAcquired = false
}

func NewLocalStorage(lockFilePath string) LocalStorage {
	return LocalStorage{false, lockFilePath}
}
