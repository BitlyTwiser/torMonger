package database

type LinkReference struct {
	Id       string `json:"id,omitempty"`
	LinkHash string `json:"link_hash,omitempty"`
	Link     string `json:"link,omitempty"`
}

type SubdirctoryReference struct {
	Id               string `json:"id,omitempty"`
	TormongerDataId  string `json:"tormonger_data_id,omitempty"`
	SubdirectoryPath string `json:"subdirectory_path,omitempty"`
}

type HtmlDataReference struct {
	Id                            string `json:"id,omitempty"`
	TormongerDataId               string `json:"tormonger_data_id,omitempty"`
	TormongerDataSubDirectoriesId string `json:"tormonger_data_sub_directories_id,omitempty"`
	HtmlData                      string `json:"html_data,omitempty"`
	FoundValues                   bool
}
