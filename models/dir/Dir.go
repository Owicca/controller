package dir

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Owicca/controller/models/fileinfo"
	"github.com/Owicca/controller/models/fsitem"
)

type Dir struct {
	Children map[string]fsitem.FSItem `json:"children"`
	Info     fileinfo.FileInfo        `json:"info"`
}

//TODO: every element has an pointer to its parent,
//marshaling goes in a infinite loop if it follows pointers :(
func (t Dir) ToJson() (string, error) {
	result := ""
	f, _ := os.Create("dump.json")
	jsn, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	f.Write(jsn)

	for idx, item := range t.Children {
		switch item.(type) {
		case Dir:
			break
		default:
			break
		}
		log.Printf("%s => %T", idx, item)
		continue
	}
	f.Close()
	os.Exit(-1)

	return result, nil
}

func (f Dir) GetPath(childPath []byte) []byte {
	path := append([]byte(f.Info.Name), []byte{'/'}...)
	path = append(path, childPath...)
	if f.Info.Parent == nil {
		return path
	}

	parent := f.Info.Parent.(*Dir)
	return parent.GetPath(path)
}
