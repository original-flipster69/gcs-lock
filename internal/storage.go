package storage

type storage interface {
	LockFileExists() bool

	Lock()

	GetLockContent() (string, error)

	HasLock() bool

	DeleteLock()

	Unlock()
}
