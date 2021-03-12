package internal

import (
	"reflect"
	"sync"
	"testing"
)

// AssertMutexUnlocked checks if a mutex is locked
func AssertMutexUnlocked(t *testing.T, m *sync.Mutex) {
	var state = reflect.ValueOf(m).Elem().FieldByName("state")
	if state.Int() == 1 {
		t.Errorf("mutex still locked")
	}
}
