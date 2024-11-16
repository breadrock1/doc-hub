package httpserv

func createStatusResponse(status int, msg string) *ResponseForm {
	return &ResponseForm{Status: status, Message: msg}
}

// ResponseForm example
type ResponseForm struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Done"`
}

// BadRequestForm example
type BadRequestForm struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad Request message"`
}

// ServerErrorForm example
type ServerErrorForm struct {
	Status  int    `json:"status" example:"503"`
	Message string `json:"message" example:"Server Error message"`
}

// CreateBucketForm example
type CreateBucketForm struct {
	BucketName string `json:"bucket_name" example:"test-bucket"`
}

// MoveFilesForm example
type MoveFilesForm struct {
	TargetDirectory string   `json:"location" example:"common-folder"`
	SourceDirectory string   `json:"src_folder_id" example:"unrecognized"`
	DocumentPaths   []string `json:"document_ids" example:"./indexer/watcher/test.txt"`
}

// RemoveFileForm example
type RemoveFileForm struct {
	FileName string `json:"file_name" example:"test-file.docx"`
}

// DownloadFileForm example
type DownloadFileForm struct {
	FileName string `json:"file_name" example:"test-file.docx"`
}

// ShareFileForm example
type ShareFileForm struct {
	FileName    string `json:"file_name" example:"test-file.docx"`
	FileDirPath string `json:"dir_path" example:"test-folder/"`
	ExpiredSecs int32  `json:"expired_secs" example:"3600"`
}

// GetFilesForm example
type GetFilesForm struct {
	DirectoryName string `json:"directory" example:"test-folder/"`
}

// CopyFileForm example
type CopyFileForm struct {
	SrcPath string `json:"src_path" example:"old-test-document.docx"`
	DstPath string `json:"dst_path" example:"test-document.docx"`
}
