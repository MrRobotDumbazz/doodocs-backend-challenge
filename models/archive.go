package models

type Archive struct {
	Filename     string  `json:"filename"`
	Archive_size float64 `json:"archive_size"`
	Total_size   float64 `json:"total_size"`
	Total_files  float64 `json:"total_files"`
	Files        []File  `json:"files"`
}
