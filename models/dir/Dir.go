package dir

import (
	"github.com/Owicca/controller/models/fileinfo"
	"github.com/Owicca/controller/models/fsitem"
)

type Dir struct {
	Children map[string]fsitem.FSItem
	Info     fileinfo.FileInfo
	IsRoot   bool
}

func (t Dir) ToJson() (string, error) {
	result := ""
	for _, item := range t.Children {
		jsonString, err := item.ToJson()
		if err != nil {
			return result, err
		}
		result = result + jsonString
	}

	return result, nil
}

//func (f Dir) ToJson() (string, error) {
//	byteArr, err := json.Marshal(f)
//
//	return string(byteArr), err
//}

func (f Dir) GetPath(childPath []byte) []byte {
	path := append(f.Info.Name, []byte{'/'}...)
	path = append(path, childPath...)
	if f.IsRoot {
		return path
	}

	parent := f.Info.Parent
	return parent.GetPath(path)
}
