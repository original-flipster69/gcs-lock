package storage

type Storage interface {
	LockFileExists() bool

	Lock()

	GetLockContent() (string, error)

	HasLock() bool

	DeleteLock()
}
