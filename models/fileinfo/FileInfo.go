package fileinfo

type FileInfo struct {
	Name       string      `json:"name"`
	PseudoName string      `json:"pseudoname"`
	Parent     interface{} `json:"parent"`
}
