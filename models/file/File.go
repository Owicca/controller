package file

import (
	"encoding/json"

	"github.com/Owicca/controller/models/fileinfo"
)

type File struct {
	Info fileinfo.FileInfo `json:"info"`
}

func (f File) ToJson() ([]byte, error) {
	jsn, err := json.Marshal(f)

	return jsn, err
}

func (f File) GetPath() []byte {
	path := append([]byte(f.Info.Path), []byte{'/'}...)
	path = append(path, []byte(f.Info.Name)...)

	return path
}
