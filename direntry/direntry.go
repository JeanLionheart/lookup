package direntry

import (
	"os"
)

type subFile struct {
	name string
	size int64
}

type DirEntry struct {
	path     string
	size     int64
	subDirs  []*DirEntry
	subFiles []subFile
}

func New() *DirEntry {
	d := DirEntry{}
	d.size = 0
	d.subDirs = make([]*DirEntry, 0)
	d.subFiles = make([]subFile, 0)
	return &d
}

func (d *DirEntry) MakePath(path string) {
	d.path = path
}

func (d *DirEntry) Size() int64 {
	return d.size
}

func (d *DirEntry) Collect(info os.FileInfo, nowPath string) *DirEntry {
	if info.IsDir() {
		subd := New()
		subd.path = nowPath + "\\" + info.Name()
		d.subDirs = append(d.subDirs, subd)
		return subd
	} else {
		subf := subFile{info.Name(), info.Size()}
		d.subFiles = append(d.subFiles, subf)
		return nil
	}
}

func (d *DirEntry) Path() string {
	return d.path
}

func (d *DirEntry) SubDirs() []*DirEntry {
	return d.subDirs
}

func (d *DirEntry) Sum() {
	for _, sf := range d.subFiles {
		d.size += sf.size
	}
	for _, sd := range d.subDirs {
		sd.Sum()
		d.size += sd.size
	}
}
