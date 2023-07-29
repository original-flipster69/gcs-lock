package main

import (
	"cloud-lock-go/internal/local"
	"fmt"
)

func main() {

	go newStorageDoLock()
	newStorageDoLock()
}

func newStorageDoLock() {
	//gs := cloudstorage.New("tyler-lockett", "leader.txt")
	gs := local.NewLocalStorage("leader.txt")

	if gs.LockFileExists() {
		fmt.Println("already locked... nothing to do here")
		return
	}
	gs.Lock()

	content, err := gs.GetLockContent()
	if err != nil {
		fmt.Printf("error getting lock content: %v", err)
		return
	}
	fmt.Println(content)

	gs.DeleteLock()
	gs.Unlock()
}
