package walker

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/Owicca/controller/models/dir"
	"github.com/Owicca/controller/models/file"
	"github.com/Owicca/controller/models/fileinfo"
	"github.com/Owicca/controller/models/fsitem"
)

type Walker struct {
	Extensions map[string]bool
	FSTree     fsitem.FSItem
	TTL        int
}

func NewWalker() *Walker {
	walker := &Walker{
		Extensions: make(map[string]bool),
		FSTree:     nil,
		TTL:        10,
	}
	walker.setupExtensions()
	return walker
}

func (w Walker) setupExtensions() {
	file, err := os.Open("config/video_extensions.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w.Extensions[scanner.Text()] = true
	}

	file, err = os.Open("config/image_extensions.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		w.Extensions[scanner.Text()] = true
	}
}

func (w Walker) ParsePath(Dir *string) {
	log.Println("Started walking!")
	var root = dir.Dir{
		Children: make(map[string]fsitem.FSItem),
		Info: fileinfo.FileInfo{
			Name:   []byte(*Dir),
			Parent: nil,
		},
		IsRoot: true,
	}
	w.FSTree = &root
	var curParent = &root
	err := filepath.Walk(*Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			var item = w.NewFSItem(curParent, info)
			curParent.Children[info.Name()] = item
			curParent = item.(*dir.Dir)
		} else {
			ext := filepath.Ext(info.Name())
			_, validExt := w.Extensions[ext]
			if validExt {
				curParent.Children[info.Name()] = w.NewFSItem(curParent, info)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (w Walker) NewFSItem(parent fsitem.FSItem, info os.FileInfo) fsitem.FSItem {
	if !info.IsDir() {
		return &file.File{
			Info: fileinfo.FileInfo{
				Name:   []byte(info.Name()),
				Parent: parent,
			},
		}
	}
	return &dir.Dir{
		Children: make(map[string]fsitem.FSItem),
		Info: fileinfo.FileInfo{
			Name:   []byte(info.Name()),
			Parent: parent,
		},
		IsRoot: false,
	}
}
