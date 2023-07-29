package cloudstorage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"strings"
)

type GoogleCloudStorage struct {
	lockAquired bool
	bucket      string
	lockFile    string
}

func (gs *GoogleCloudStorage) LockFileExists() bool {
	ctx := context.Background()
	obj := gs.objHandle(ctx)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return false
	}
	defer r.Close()
	return true
}

func (gs *GoogleCloudStorage) objHandle(ctx context.Context) *storage.ObjectHandle {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return client.Bucket(gs.bucket).Object(gs.lockFile)
}

func (gs *GoogleCloudStorage) Lock() {
	ctx := context.Background()
	obj := gs.objHandle(ctx)
	w := obj.If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
	if _, err := fmt.Fprint(w, "yui yui yui"); err != nil {
		fmt.Printf("error while writing: %v", err)
		return
	}
	if err := w.Close(); err != nil {
		fmt.Printf("could not acquire lock: %v", err)
		return
	}
	gs.lockAquired = true
	fmt.Println("successfully acquired lock")

}

func (gs *GoogleCloudStorage) GetLockContent() (string, error) {
	ctx := context.Background()
	obj := gs.objHandle(ctx)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("getLockContent(): %v", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			fmt.Printf("getLockContent(): %v", err)
		}
	}()
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r); err != nil {
		return "", fmt.Errorf("getLockContent(): %v", err)
	}
	return buf.String(), nil
}

func (gs *GoogleCloudStorage) HasLock() bool {
	return gs.lockAquired
}

func (gs *GoogleCloudStorage) DeleteLock() {
	ctx := context.Background()
	obj := gs.objHandle(ctx)
	if err := obj.Delete(ctx); err != nil {
		fmt.Printf("DeleteLock(): %v", err)
	}
}

func (gs *GoogleCloudStorage) Unlock() {
	gs.lockAquired = false
}

func New(bucket string, lockFile string) GoogleCloudStorage {
	return GoogleCloudStorage{false, bucket, lockFile}
}
