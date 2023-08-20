package main

import (
	cloudstorage "cloud-lock-go/internal/gcp"
	lock "cloud-lock-go/pkg"
	"fmt"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go newStorageDoLock(&wg)
	}

	wg.Wait()
	fmt.Println("done")
}

func newStorageDoLock(wg *sync.WaitGroup) {
	defer wg.Done()
	gs := cloudstorage.NewStorage("tyler-lockett", "leader.txt")
	//gs := local.NewLocalStorage("leader.txt")

	if leader := lock.Lock(&gs); leader {
		fmt.Println("doing all se funky stuff")
		time.Sleep(5 * time.Second)
		lock.Unlock(&gs)
	}
}
