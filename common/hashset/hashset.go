package hashset

import "sync"

type HashSet struct {
	set  map[string]bool
	lock sync.Mutex
}

func NewHashSet() *HashSet {
	return &HashSet{set: map[string]bool{}, lock: sync.Mutex{}}
}

func (set *HashSet) Add(i string) bool {
	set.lock.Lock()
	defer set.lock.Unlock()
	if set.set == nil {
		set.set = make(map[string]bool)
	}
	_, found := set.set[i]
	set.set[i] = true
	return !found //False if it existed already
}

func (set *HashSet) Get(i string) bool {
	set.lock.Lock()
	defer set.lock.Unlock()
	if set.set == nil {
		return false
	}
	_, found := set.set[i]
	return found //true if it existed already
}

func (set *HashSet) Remove(i string) {
	set.lock.Lock()
	defer set.lock.Unlock()
	if set.set == nil {
		return
	}
	delete(set.set, i)
}

func (set *HashSet) Values() []string {
	set.lock.Lock()
	defer set.lock.Unlock()
	if set.set == nil {
		return nil
	}
	values := make([]string, len(set.set))
	var idx = 0
	for k, _ := range set.set {
		values[idx] = k
		idx++
	}
	return values
}

func (set *HashSet) Len() int {
	set.lock.Lock()
	defer set.lock.Unlock()
	if set.set == nil {
		return 0
	}
	return len(set.set)
}
