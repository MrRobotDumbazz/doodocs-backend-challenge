package models

type File struct {
	File_path string  `json:"file_path"`
	Size      float64 `json:"size"`
	Mimetype  string  `json:"mimetype"`
}
