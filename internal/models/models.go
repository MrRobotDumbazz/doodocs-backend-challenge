package models

type File struct {
	File_path string  `json:"file_path"`
	Size      float64 `json:"size"`
	Mimetype  string  `json:"mimetype"`
}

type Archive struct {
	Filename     string  `json:"filename"`
	Archive_size float64 `json:"archive_size"`
	Total_size   float64 `json:"total_size"`
	Total_files  float64 `json:"total_files"`
	Files        []File  `json:"files"`
}
