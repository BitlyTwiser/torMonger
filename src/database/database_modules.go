package database

type LinkReference struct {
	Id       string `json:"id,omitempty"`
	LinkHash string `json:"link_hash,omitempty"`
	Link     string `json:"link,omitempty"`
}

type SubdirctoryReference struct {
	TormongerDataId  string `json:"tormonger_data_id,omitempty"`
	HtmlDataId       string `json:"html_data_id,omitempty"`
	SubdirectoryPath string `json:"subdirectory_path,omitempty"`
}
