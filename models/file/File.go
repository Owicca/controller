package file

import (
	"encoding/json"

	"github.com/Owicca/controller/models/fileinfo"
)

type File struct {
	Info fileinfo.FileInfo
}

func (f File) ToJson() (string, error) {
	byteArr, err := json.Marshal(f)

	return string(byteArr), err
}

func (f File) GetPath(childPath []byte) []byte {
	path := append(f.Info.Name, []byte{'/'}...)
	path = append(path, childPath...)

	parent := f.Info.Parent
	return parent.GetPath(path)
}
