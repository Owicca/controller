package models

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type FSItem interface {
	ToJson() (string, error)
	GetPath() []byte
}

type FileInfo struct {
	Name   []byte `json:"name"`
	Parent *FSItem
}

//===========
//===========
//===========
type Dir struct {
	Children []*FSItem
	Info     FileInfo
}

func (f Dir) ToJson() (string, error) {
	byteArr, err := json.Marshal(f)

	return string(byteArr), err
}

func (f Dir) GetPath(childPath []byte, isRoot bool) []byte {
	path := append(f.Info.Name, []byte{'/'}...)
	path = append(path, childPath...)
	if isRoot {
		return path
	}

	parent := *f.Info.Parent
	return parent.GetPath(path)
}

//===========
//===========
//===========
type File struct {
	Info FileInfo
}

func (f File) ToJson() (string, error) {
	byteArr, err := json.Marshal(f)

	return string(byteArr), err
}

func (f File) GetPath(childPath []byte, isRoot bool) []byte {
	path := append(f.Info.Name, []byte{'/'}...)
	path = append(path, childPath...)
	if isRoot {
		return path
	}

	parent := *f.Info.Parent
	return parent.GetPath(path)
}

//============
var Tree = &Dir{}

//func (t FSItem) ToJson() (string, error) {
//	result := ""
//	for _, item := range t {
//		jsonString, err := (*item).ToJson()
//		if err != nil {
//			return result, err
//		}
//		result = result + jsonString
//	}
//
//	return result, nil
//}

//============
type Walker struct {
	Extensions map[string]bool
	FSTree     FSItem
	TTL        int
}

func NewWalker() *Walker {
	walker := &Walker{
		Extensions: make(map[string]bool),
		FSTree:     nil, //Dir{}, #@#TODO
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
}

func (w Walker) ParsePath(Dir *string) {
	log.Println("Started walking!")
	err := filepath.Walk(*Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(info.Name())
		_, validExt := w.Extensions[ext]

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (w Walker) NewFSItem(parent *FSItem, info os.FileInfo) *FSItem {
	if !info.IsDir() {
		return &File{
			Info: FileInfo{
				Name:   []byte(info.Name()),
				Parent: parent,
			},
		}
	}
	return &Dir{
		Children: nil, //#@#TODO
		Info: FileInfo{
			Name:   []byte{"root"},
			Parent: parent,
		},
	}
}
