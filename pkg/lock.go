package lock

import (
	storage "cloud-lock-go/internal"
	log "github.com/sirupsen/logrus"
	"time"
)

const lifetime = 5

func Lock(st storage.Storage) bool {
	if !st.LockFileExists() {
		return tryLock(st)
	}

	content, err := st.GetLockContent()
	if err != nil {
		log.Errorf("Lock(): %v", err)
		return false
	}
	lockTime, err := time.Parse(time.RFC3339, content)
	if err != nil {
		log.Errorf("Lock(), error parsing time string: %v", err)
		return false
	}
	diff := int(lockTime.Sub(time.Now()).Minutes())

	if diff < lifetime {
		return false
	}
	st.DeleteLock()
	return tryLock(st)
}

func tryLock(st storage.Storage) bool {
	st.Lock()
	if !st.HasLock() {
		log.Infof("failed acquiring lock")
		return false
	}
	log.Infof("successfully acquired lock")
	return true
}

func Unlock(st storage.Storage) {

	if !st.HasLock() {
		return
	}

	content, err := st.GetLockContent()
	if err != nil {
		log.Errorf("Unlock(): %v", err)
		return
	}
	lockTime, err := time.Parse(time.RFC3339, content)
	if err != nil {
		log.Errorf("Unlock(), error parsing time string: %v", err)
		return
	}
	diff := int(lockTime.Sub(time.Now()).Seconds())
	if diff < lifetime {
		log.Infof("Unlock(): waiting minimum lock duration")
		time.Sleep(lifetime * time.Second)
	}

	st.DeleteLock()
	log.Infof("lock released")
}
