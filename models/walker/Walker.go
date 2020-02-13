package walker

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Owicca/controller/models/dir"
	fl "github.com/Owicca/controller/models/file"
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

// load supported video/image extensions
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

func (w *Walker) ParsePath(Dir *string) (bool, error) {
	log.Println("Started walking!")
	tree, err := ParseDir(*Dir, *Dir)
	if err != nil {
		pErr := errors.New(fmt.Sprintf("Couldn't parse dir %s (%s)", *Dir, err))
		log.Println(pErr.Error())
		return false, pErr
	}
	w.FSTree = tree
	return true, nil
}

/*
* open dir and gather info
* parse files
* recursivelly parse sub directories
* return built hierarchy
 */
func ParseDir(Dir string, workingDir string) (dir.Dir, error) {
	var dr = dir.Dir{
		Children: make(map[string]fsitem.FSItem),
		Info: fileinfo.FileInfo{
			Name:       "",
			PseudoName: "",
			Path:       workingDir,
		},
	}

	subDirs := make([]string, 0)
	var (
		file *os.File
		fErr error
	)
	if Dir == workingDir {
		file, fErr = os.Open(Dir)
	} else {
		path := filepath.Join(workingDir, Dir)
		file, fErr = os.Open(path)
		workingDir = path
	}
	info, iErr := file.Stat()
	if fErr != nil || iErr != nil {
		var err string
		if fErr != nil {
			err = fErr.Error()
		} else {
			err = iErr.Error()
		}
		cErr := errors.New(fmt.Sprintf("Could not access %s (%s)\n", Dir, err))
		return dr, cErr
	}

	dr.Info.Name = info.Name()
	dr.Info.PseudoName, _, _ = GetPseudo(info.Name())
	log.Printf("Dir: %v\n", dr.Info)

	children, err := file.Readdir(0)
	if err != nil {
		log.Printf("Could not read children of %s (%s)\n", Dir, err.Error())
		cErr := errors.New(fmt.Sprintf("Could not read children of %s (%s)\n", Dir, err.Error()))
		return dr, cErr
	}

	for _, ch := range children {
		if ch.IsDir() {
			subDirs = append(subDirs, ch.Name())
		} else {
			ch_pseudo, _, _ := GetPseudo(ch.Name())
			dr.Children[ch_pseudo] = fl.File{
				Info: fileinfo.FileInfo{
					Name:       ch.Name(),
					PseudoName: ch_pseudo,
					Path:       workingDir,
				},
			}
			log.Printf("File: %s %s %s\n", ch.Name(), ch_pseudo, workingDir)
		}
	}

	if len(subDirs) > 0 {
		for _, subDir := range subDirs {
			subDirTree, err := ParseDir(subDir, workingDir)
			if err != nil {
				log.Println("In subdirs range\n", err)
			}
			dr.Children[subDirTree.Info.PseudoName] = subDirTree
		}
	}

	return dr, nil
}

// base64 encoded sha1 of string
func GetPseudo(str string) (string, []byte, error) {
	sh := sha1.New()
	sh.Write([]byte(str))
	hs := sh.Sum(nil)
	asString := base64.URLEncoding.EncodeToString(hs)

	return asString, hs, nil
}
