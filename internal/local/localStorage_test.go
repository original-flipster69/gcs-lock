package local

import (
	"testing"
)

func TestLockFileExists(t *testing.T) {
	ls := NewLocalStorage("leader.txt")
	res := ls.LockFileExists()
	if res {
		t.Errorf("LockFileExists should result in false, but returned true")
	}
}

func TestHasLock(t *testing.T) {
	testTable := []struct {
		name           string
		input          bool
		expectedResult bool
	}{
		{
			"true when lock acquired", true, true,
		},
		{
			"false when lock not acquired", false, false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ls := LocalStorage{tc.input, ""}
			if hasLock := ls.HasLock(); hasLock != tc.expectedResult {
				t.Errorf("HasLock should result in %v, but was %v", tc.expectedResult, hasLock)
			}
		})
	}
}

func TestUnlock(t *testing.T) {
	testTable := []struct {
		name           string
		input          bool
		expectedResult bool
	}{
		{
			"false after unlocking with lock", true, false,
		},
		{
			"false after unlocking without lock", false, false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ls := LocalStorage{tc.input, ""}
			ls.Unlock()
			if ls.lockAcquired != tc.expectedResult {
				t.Errorf("after Unlock lockAcquired should be %v, but was %v", tc.expectedResult, ls.lockAcquired)
			}
		})
	}
}
