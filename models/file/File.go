package file

import (
	"encoding/json"

	"github.com/Owicca/controller/models/dir"
	"github.com/Owicca/controller/models/fileinfo"
)

type File struct {
	Info fileinfo.FileInfo `json:"info"`
}

//TODO: every element has an pointer to its parent,
//marshaling goes in a infinite loop if it follows pointers :(
func (f File) ToJson() (string, error) {
	byteArr, err := json.Marshal(f)

	return string(byteArr), err
}

func (f File) GetPath(childPath []byte) []byte {
	path := append([]byte(f.Info.Name), []byte{'/'}...)
	path = append(path, childPath...)

	parent := f.Info.Parent.(*dir.Dir)
	return parent.GetPath(path)
}
