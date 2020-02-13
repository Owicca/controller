package walker

import (
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
	"github.com/Owicca/controller/config"
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
	array := make([]string, 0)
	array = append(array, config.NewVideoExtensions()...)
	array = append(array, config.NewImageExtensions()...)

	for _, ext := range array {
		w.Extensions[ext] = true
	}
}

func (w *Walker) ParsePath(Dir *string) (bool, error) {
	log.Println("Started walking!")
	tree, err := ParseDir(*Dir, *Dir, []byte(""))
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
func ParseDir(Dir string, workingDir string, DirPathPseudo []byte) (dir.Dir, error) {
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
	// log.Printf("Dir: %v\n", dr.Info)

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
			chPseudo, _, _ := GetPseudo(ch.Name())
			finalPseudo := GetPathPseudo(string(DirPathPseudo), chPseudo)
			dr.Children[string(finalPseudo)] = fl.File{
				Info: fileinfo.FileInfo{
					Name:       ch.Name(),
					PseudoName: chPseudo,
					Path:       workingDir,
				},
			}
			// log.Printf("=====File: %s (%s %s %s) %s\n", ch.Name(), DirPathPseudo, chPseudo, string(finalPseudo), workingDir)
		}
	}

	if len(subDirs) > 0 {
		for _, subDir := range subDirs {
			psd , _ , _ := GetPseudo(subDir)
			subDirTree, err := ParseDir(subDir, workingDir, []byte(psd))
			if err != nil {
				log.Println("In subdirs range\n", err)
			}

			finalPseudo := GetPathPseudo(string(DirPathPseudo), psd)
			// log.Println(subDir, psd, string(finalPseudo))
			dr.Children[string(finalPseudo)] = subDirTree
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

func GetPathPseudo(parentPseudo string, childPseudo string) []byte {
	pathPseudo := []byte(childPseudo)[:5]
	parentByteArr := make([]byte, 0)
	if len(parentPseudo) > 0 {
		parentByteArr = []byte(parentPseudo)[:5]
	}
	return append(parentByteArr, pathPseudo...)
}