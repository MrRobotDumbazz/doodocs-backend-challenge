package models

type File struct {
	File_path string
	Size      float64
	Mimetype  string
}

type Archive struct {
	Filename     string
	Archive_size float64
	Total_size   float64
	Files        []File
}
