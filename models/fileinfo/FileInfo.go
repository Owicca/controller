package fileinfo

type FileInfo struct {
	Name       string `json:"name"`
	PseudoName string `json:"pseudoname"`
	Path       string `json:"-"`
}
