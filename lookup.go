package lookup

import (
	"container/list"
	"fmt"
	"lookup/direntry"
	"os"
	"path/filepath"
	"sync"
)

func dirWalkSearch(root string) *direntry.DirEntry {
	entry := direntry.New()
	p, _ := filepath.Abs(".")
	entry.MakePath(p)
	que := list.List{}
	que.Init()
	que.PushBack(entry)

	wideIteration := func(e *direntry.DirEntry, nowPath string, que *list.List, lock *sync.Mutex) {
		subs, err := os.ReadDir(e.Path())
		if err != nil {
			fmt.Println(e.Path())
			panic(err)
		}
		for _, sub := range subs {
			subinfo, _ := sub.Info()
			subaddr := e.Collect(subinfo, nowPath)
			if subaddr != nil {
				lock.Lock()
				que.PushBack(subaddr)
				lock.Unlock()
			}
		}
	}

	for que.Len() != 0 {
		cache := list.List{}
		cache.Init()
		cacheLock := sync.Mutex{}
		var wg sync.WaitGroup

		for que.Len() != 0 {
			wg.Add(1)
			e := que.Front().Value.(*direntry.DirEntry)
			que.Remove(que.Front())
			go func() {
				wideIteration(e, e.Path(), &cache, &cacheLock)
				wg.Done()
			}()
		}
		wg.Wait()

		for cache.Len() != 0 {
			que.PushBack(cache.Front().Value)
			cache.Remove(cache.Front())
		}
	}

	entry.Sum()

	return entry
}

func Run(root string, depth int) {
	fmt.Println()
	r := dirWalkSearch(root)

	que := list.List{}
	que.Init()
	que.PushBack(r)
	for depth > 0 && que.Len() != 0 {
		depth--
		loopQue := list.List{}
		loopQue.Init()
		for que.Len() != 0 {
			loopQue.PushBack(que.Front().Value)
			que.Remove(que.Front())
		}

		for loopQue.Len() != 0 {
			d := loopQue.Front().Value.(*direntry.DirEntry)
			loopQue.Remove(loopQue.Front())

			out := fmt.Sprintf("%s: %.2f MBytes", d.Path(), float64(d.Size())/1024/1024)
			fmt.Println(out)

			subdirs := d.SubDirs()
			for _, sd := range subdirs {
				que.PushBack(sd)
			}
		}

		fmt.Println()
	}
}
