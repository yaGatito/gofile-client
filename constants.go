package gofile

// RootFolder used to specify the root folder ID that is behind the scene.
const RootFolder = "root"

const (
	postFolderEndpoint = "https://api.gofile.io/contents/createFolder"
	contentsBaseURL    = "https://api.gofile.io/contents/"
	accountsBaseURL    = "https://api.gofile.io/accounts/"
	getFileEndpoint    = "https://%s.gofile.io/download/web/%s/%s"
	postFileEndpoint   = "https://upload.gofile.io/uploadfile"
)
