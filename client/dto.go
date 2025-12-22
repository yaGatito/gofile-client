package gofile

import "fmt"

// Public Models
type CreateFolderResponseBody struct {
	Status string `json:"status"`
	Data   struct {
		Id             string `json:"id"`
		Owner          string `json:"owner"`
		Name           string `json:"name"`
		ParentFolderId string `json:"parentFolder"`
		CreateTime     int64  `json:"createTime"`
		Code           string `json:"code"`
	} `json:"data"`
}

func (c CreateFolderResponseBody) String() string {
	return fmt.Sprintf("Status: %s; Data.Id: %s; Data.Owner: %s; Data.Name: %s; Data.ParentFolderId: %s; Data.CreateTime: %d; Data.Code: %s",
		c.Status, c.Data.Id, c.Data.Owner, c.Data.Name, c.Data.ParentFolderId, c.Data.CreateTime, c.Data.Code)
}

type UploadFileResponseBody struct {
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

func (u UploadFileResponseBody) String() string {
	return fmt.Sprintf("Status: %s; Data.Id: %s; Data.Name: %s; Data.Md5: %s; Data.Size: %d; Data.Type: %s; Data.Mimetype: %s; Data.CreateTime: %d; Data.ParentFolderId: %s; Data.DownloadPage: %s",
		u.Status, u.Data.Id, u.Data.Name, u.Data.Md5, u.Data.Size, u.Data.Type, u.Data.Mimetype, u.Data.CreateTime, u.Data.ParentFolderId, u.Data.DownloadPage)
}

type GetFileInfoResponseBody struct {
	Status string `json:"status"`
	Data   struct {
		Id             string   `json:"id"`
		ParentFolderId string   `json:"parentFolder"`
		Type           string   `json:"type"`
		Name           string   `json:"name"`
		CreateTime     int64    `json:"createTime"`
		Size           int64    `json:"size"`
		Mimetype       string   `json:"mimetype"`
		Servers        []string `json:"servers"`
		ServerSelected string   `json:"serverSelected"`
		DownloadPage   string   `json:"link"`
		ThumbNailLink  string   `json:"thumbnail"`
		Md5            string   `json:"md5"`
	} `json:"data"`
}

func (u GetFileInfoResponseBody) String() string {
	return fmt.Sprintf("Status: %s; Data.Id: %s; Data.Name: %s; Data.Md5: %s; Data.Size: %d; Data.Type: %s; Data.Mimetype: %s; Data.CreateTime: %d; Data.ParentFolderId: %s; Data.DownloadPage: %s",
		u.Status, u.Data.Id, u.Data.Name, u.Data.Md5, u.Data.Size, u.Data.Type, u.Data.Mimetype, u.Data.CreateTime, u.Data.ParentFolderId, u.Data.DownloadPage)
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