package fileinfo

import (
	"github.com/Owicca/controller/models/fsitem"
)

type FileInfo struct {
	Name   []byte `json:"name"`
	Parent fsitem.FSItem
}
