package httpserv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *ServerHttp) CreateCloudGroup() error {
	group := s.server.Group("/cloud")

	group.GET("/buckets", s.GetBuckets)
	group.PUT("/bucket", s.CreateBucket)
	group.DELETE("/:bucket", s.RemoveBucket)

	group.POST("/:bucket/files", s.GetFiles)
	group.POST("/:bucket/file/copy", s.CopyFile)
	group.POST("/:bucket/file/move", s.MoveFile)
	group.PUT("/:bucket/file/upload", s.UploadFile)
	group.POST("/:bucket/file/download", s.DownloadFile)
	group.DELETE("/:bucket/file/remove", s.RemoveFile)

	group.POST("/:bucket/file/share", s.ShareFile)

	return nil
}

// GetBuckets
// @Summary Get watched bucket list
// @Description Get watched bucket list
// @ID get-buckets
// @Tags buckets
// @Produce  json
// @Success 200 {array} string "Ok"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/buckets [get]
func (s *ServerHttp) GetBuckets(c echo.Context) error {
	ctx := c.Request().Context()
	watcherDirs, err := s.cloud.Cloud.GetBuckets(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, watcherDirs)
}

// CreateBucket
// @Summary Create new bucket into cloud
// @Description Create new bucket into cloud
// @ID create-bucket
// @Tags buckets
// @Accept  json
// @Produce json
// @Param jsonQuery body CreateBucketForm true "Bucket name to create"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/bucket [put]
func (s *ServerHttp) CreateBucket(c echo.Context) error {
	jsonForm := &CreateBucketForm{}
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(jsonForm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = s.cloud.Cloud.CreateBucket(ctx, jsonForm.BucketName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// RemoveBucket
// @Summary Remove bucket from cloud
// @Description Remove bucket from cloud
// @ID remove-bucket
// @Tags buckets
// @Produce  json
// @Param bucket path string true "Bucket name to remove"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/{bucket} [delete]
func (s *ServerHttp) RemoveBucket(c echo.Context) error {
	bucket := c.Param("bucket")
	ctx := c.Request().Context()
	err := s.cloud.Cloud.RemoveBucket(ctx, bucket)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// CopyFile
// @Summary Copy file to another location into bucket
// @Description Copy file to another location into bucket
// @ID copy-file
// @Tags files
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name of src file"
// @Param jsonQuery body CopyFileForm true "Params to copy file"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/{bucket}/file/copy [post]
func (s *ServerHttp) CopyFile(c echo.Context) error {
	bucket := c.Param("bucket")

	jsonForm := &CopyFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(jsonForm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = s.cloud.Cloud.CopyFile(ctx, bucket, jsonForm.SrcPath, jsonForm.DstPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// MoveFile
// @Summary Move file to another location into bucket
// @Description Move file to another location into bucket
// @ID move-file
// @Tags files
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name of src file"
// @Param jsonQuery body CopyFileForm true "Params to move file"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/{bucket}/file/move [post]
func (s *ServerHttp) MoveFile(c echo.Context) error {
	bucket := c.Param("bucket")

	jsonForm := &CopyFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(jsonForm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = s.cloud.Cloud.CopyFile(ctx, bucket, jsonForm.SrcPath, jsonForm.DstPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// UploadFile
// @Summary Upload files to cloud
// @Description Upload files to cloud
// @ID upload-files
// @Tags files
// @Accept  multipart/form
// @Produce  json
// @Param bucket path string true "Bucket name to upload files"
// @Param expired query string false "File datetime expired like 2025-01-01T12:01:01Z"
// @Param files formData file true "Files multipart form"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/{bucket}/file/upload [put]
func (s *ServerHttp) UploadFile(c echo.Context) error {
	var fileData bytes.Buffer

	multipartForm, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	bucket := c.Param("bucket")
	if exist, err := s.cloud.Cloud.IsBucketExist(c.Request().Context(), bucket); err != nil || !exist {
		retErr := fmt.Errorf("specified bucket %s does not exist", bucket)
		return echo.NewHTTPError(http.StatusBadRequest, retErr.Error())
	}

	if multipartForm.File["files"] == nil {
		err = fmt.Errorf("there are no files into multipart form")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	expired := c.QueryParam("expired")
	timeVal, timeParseErr := time.Parse(time.RFC3339, expired)
	if timeParseErr != nil {
		log.Println("failed to parse expired time param: ", expired, timeParseErr)
	}

	ctx := c.Request().Context()
	for _, fileForm := range multipartForm.File["files"] {
		fileName := fileForm.Filename
		fileHandler, err := fileForm.Open()
		if err != nil {
			log.Println("failed to open file form", err)
			continue
		}
		defer func() {
			if err := fileHandler.Close(); err != nil {
				log.Println("failed to close file handler: ", fileName, err)
				return
			}
		}()

		fileData.Reset()
		_, err = fileData.ReadFrom(fileHandler)
		if err != nil {
			log.Println("failed to read file form", fileName, err)
			continue
		}

		if timeParseErr == nil {
			err = s.cloud.Cloud.UploadExpired(ctx, bucket, fileName, timeVal, fileData)
		} else {
			err = s.cloud.Cloud.UploadFile(ctx, bucket, fileName, fileData)
		}

		if err != nil {
			log.Println("failed to upload file to cloud: ", fileName, err)
			continue
		}
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// DownloadFile
// @Summary Download file from cloud
// @Description Download file from cloud
// @ID download-file
// @Tags files
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name to download file"
// @Param jsonQuery body DownloadFileForm true "Parameters to download file"
// @Success 200 {file} io.Writer "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/{bucket}/file/download [post]
func (s *ServerHttp) DownloadFile(c echo.Context) error {
	bucket := c.Param("bucket")

	jsonForm := &DownloadFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	fileData, err := s.cloud.Cloud.DownloadFile(ctx, bucket, jsonForm.FileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer fileData.Reset()

	return c.Blob(200, echo.MIMEMultipartForm, fileData.Bytes())
}

// RemoveFile
// @Summary Remove file from cloud
// @Description Remove file from cloud
// @ID remove-file
// @Tags files
// @Produce  json
// @Param bucket path string true "Bucket name to remove file"
// @Param jsonQuery body RemoveFileForm true "Parameters to remove file"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/{bucket}/file/remove [delete]
func (s *ServerHttp) RemoveFile(c echo.Context) error {
	bucket := c.Param("bucket")

	jsonForm := &RemoveFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if err := s.cloud.Cloud.RemoveFile(ctx, bucket, jsonForm.FileName); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// GetFiles
// @Summary Get files list into bucket
// @Description Get files list into bucket
// @ID get-list-files
// @Tags files
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name to get list files"
// @Param jsonQuery body GetFilesForm true "Parameters to get list files"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/{bucket}/files [post]
func (s *ServerHttp) GetFiles(c echo.Context) error {
	bucket := c.Param("bucket")

	jsonForm := &GetFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(jsonForm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	listObjects, err := s.cloud.Cloud.GetFiles(ctx, bucket, jsonForm.DirectoryName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, listObjects)
}

// ShareFile
// @Summary Get share URL for file
// @Description Get share URL for file
// @ID share-file
// @Tags share
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name to share file"
// @Param jsonQuery body ShareFileForm true "Parameters to share file"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /cloud/{bucket}/file/share [post]
func (s *ServerHttp) ShareFile(c echo.Context) error {
	bucket := c.Param("bucket")

	jsonForm := &ShareFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(jsonForm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	expired := time.Second * time.Duration(jsonForm.ExpiredSecs)

	ctx := c.Request().Context()
	url, err := s.cloud.Cloud.GetShareURL(ctx, bucket, jsonForm.FileName, expired)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, createStatusResponse(200, url))
}
