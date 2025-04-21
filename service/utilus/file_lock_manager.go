package utilus

import (
	"sync"
)

// FileLockManager 管理针对特定会议ID的读写锁
type FileLockManager struct {
	locks map[string]*sync.RWMutex
	mu    sync.Mutex
}

var fileLockManager *FileLockManager
var fileLockOnce sync.Once

// GetFileLockManager 获取 FileLockManager 单例
func GetFileLockManager() *FileLockManager {
	fileLockOnce.Do(func() {
		fileLockManager = &FileLockManager{
			locks: make(map[string]*sync.RWMutex),
		}
	})
	return fileLockManager
}

// GetLock 获取指定会议ID的锁
func (flm *FileLockManager) GetLock(meetingID string) *sync.RWMutex {
	flm.mu.Lock()
	defer flm.mu.Unlock()

	if lock, exists := flm.locks[meetingID]; exists {
		return lock
	}

	lock := &sync.RWMutex{}
	flm.locks[meetingID] = lock
	return lock
}
