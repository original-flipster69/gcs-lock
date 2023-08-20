package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

type Storage struct {
	lockAcquired bool
	bucket       string
	lockFile    string
}

func (gs *Storage) LockFileExists() bool {
	ctx := context.Background()
	obj := gs.objHandle(ctx)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return false
	}
	defer r.Close()
	return true
}

func (gs *Storage) objHandle(ctx context.Context) *storage.ObjectHandle {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return client.Bucket(gs.bucket).Object(gs.lockFile)
}

func (gs *Storage) Lock() {
	ctx := context.Background()
	obj := gs.objHandle(ctx)
	w := obj.If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
	if _, err := fmt.Fprint(w, time.Now().Format(time.RFC3339)); err != nil {
		log.Errorf("error while writing: %v", err)
		return
	}
	if err := w.Close(); err != nil {
		log.Debugf("could not acquire lock: %v", err)
		return
	}
	gs.lockAcquired = true
}

func (gs *Storage) GetLockContent() (string, error) {
	ctx := context.Background()
	obj := gs.objHandle(ctx)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("getLockContent(): %v", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Errorf("getLockContent(): %v", err)
		}
	}()
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r); err != nil {
		return "", fmt.Errorf("getLockContent(): %v", err)
	}
	return buf.String(), nil
}

func (gs *Storage) HasLock() bool {
	return gs.lockAcquired
}

func (gs *Storage) DeleteLock() {
	ctx := context.Background()
	obj := gs.objHandle(ctx)
	if err := obj.Delete(ctx); err != nil {
		log.Errorf("DeleteLock(): %v", err)
	}
	gs.lockAcquired = false
}

func NewStorage(bucket string, lockFile string) Storage {
	return Storage{false, bucket, lockFile}
}
