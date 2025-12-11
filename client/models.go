package gofile

// Public Models
type CreateFolderResponseData struct {
	Status string `json:"status"`
	Data   struct {
		Id             string `json:"folderId"`
		Owner          string `json:"owner"`
		Name           string `json:"name"`
		ParentFolderId string `json:"parentFolder"`
		CreateTime     int64  `json:"createTime"`
		Code           string `json:"code"`
	} `json:"data"`
}

type UploadFileResponseData struct {
	Status string `json:"status"`
	Data   struct {
		CreateTime       int64    `json:"createTime"`
		DownloadPage     string   `json:"downloadPage"`
		Id               string   `json:"id"`
		Md5              string   `json:"md5"`
		Mimetype         string   `json:"mimetype"`
		Name             string   `json:"name"`
		ParentFolderId   string   `json:"parentFolder"`
		ParentFolderCode string   `json:"parentFolderCode"`
		Servers          []string `json:"servers"`
		Size             int64    `json:"size"`
		Type             string   `json:"type"`
	} `json:"data"`
}

// Private Models
type createFolderRequestBody struct {
	ParentFolderId string `json:"parentFolderId"`
	FolderName     string `json:"folderName,omitempty"`
}

type getAccountInfoResponseData struct {
	Status string `json:"status"`
	Data   struct {
		RootFolder string `json:"rootFolder"`
		Stats      struct {
			FolderCount int `json:"folderCount"`
			FileCount   int `json:"fileCount"`
			Storage     int `json:"storage"`
		} `json:"statsCurrent"`
		Email string `json:"email"`
	} `json:"data"`
}

type getIdResponseData struct {
	Status string `json:"status"`
	Data   struct {
		Id    string `json:"id"`
		Tier  string `json:"tier"`
		Email string `json:"email"`
	} `json:"data"`
}

