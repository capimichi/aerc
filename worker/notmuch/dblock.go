//go:build notmuch
// +build notmuch

package notmuch

import "sync"

// dbLocks holds one mutex per notmuch database path.
//
// Every notmuch account runs its own worker with its own connection to the
// database. When several accounts share the same database (e.g. one maildir
// root with maildir-account-path set per account), a worker may iterate
// query results while another one writes to the database (delete, move,
// flag). libnotmuch does not cope with the database being modified during an
// iteration and may return NULL objects or crash internally. Serializing the
// workers that share a database prevents this.
var dbLocks sync.Map // map[string]*sync.Mutex

func dbLockFor(path string) *sync.Mutex {
	l, _ := dbLocks.LoadOrStore(path, &sync.Mutex{})
	return l.(*sync.Mutex)
}

// withDBLock runs f while holding the lock of the worker's database, if the
// worker is configured. Configure itself runs unlocked, since the database
// path is only known once the database has been opened.
func (w *worker) withDBLock(f func() error) error {
	if l := w.dbLock; l != nil {
		l.Lock()
		defer l.Unlock()
	}
	return f()
}
