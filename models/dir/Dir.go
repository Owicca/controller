package dir

import (
	"encoding/json"

	"github.com/Owicca/controller/models/fileinfo"
	"github.com/Owicca/controller/models/fsitem"
)

type Dir struct {
	Children map[string]fsitem.FSItem `json:"children"`
	Info     fileinfo.FileInfo        `json:"info"`
}

func (t Dir) ToJson() ([]byte, error) {
	jsn, err := json.Marshal(t)
	if err != nil {
		return []byte(""), err
	}

	return jsn, nil
}

func (f Dir) GetPath() []byte {
	path := append([]byte(f.Info.Path), []byte{'/'}...)
	path = append(path, []byte(f.Info.Name)...)

	return path
}
