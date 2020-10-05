package imgurfetch

//Image - information about image.
type Image struct {
	Hash   string `json:"hash"`
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Ext    string `json:"ext"`
}
