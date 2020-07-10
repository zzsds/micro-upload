package handler

import (
	"bytes"
	"math/rand"
	"path"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"

	. "github.com/zzsds/micro-upload/conf"
)

// Upload ...
type Upload struct{}

// UploadResponse ....
type UploadResponse struct {
	Original string `json:"original"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	URL      string `json:"url"`
}

// Aliyun ...
// @Tags Api.Upload
// @Summary Upload file
// @Description Upload file
// @ID file.upload
// @Accept  multipart/form-data
// @Produce  json
// @Param   image formData file true  "this is a test file"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} msg.Error
// @Router /upload/aliyun [post]
func (e *Upload) Aliyun(c *gin.Context) {
	log.Info("Received Say.Anything API request")
	var (
		config     = Conf.Aliyun
		now        = time.Now()
		objectName bytes.Buffer
	)
	file, _ := c.FormFile("image")
	if file.Size <= 0 {
		log.Errorf("上传小于 0")
		c.JSON(500, errors.InternalServerError("go.micro.api.upload.aliyun", "upload file not bull").Error)
		return
	}
	baseFile := path.Base(file.Filename)
	fileName := time.Now().Format("20060102150405") + strconv.Itoa(rand.Int())[0:5] + path.Ext(baseFile)
	objectName.WriteString(c.DefaultQuery("object", "default"))
	objectName.WriteString("/")
	objectName.WriteString(now.Format("2006-01-02"))
	objectName.WriteString("/")
	objectName.WriteString(fileName)
	log.Info(objectName.String())
	// 创建OSSClient实例。
	client, err := oss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		c.JSON(500, errors.InternalServerError("go.micro.api.upload.aliyun", err.Error()).Error)
		return
	}
	// 获取存储空间。
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		c.JSON(500, errors.InternalServerError("go.micro.api.upload.aliyun", err.Error()).Error)
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(500, errors.InternalServerError("go.micro.api.upload.aliyun", err.Error()).Error)
		return
	}
	defer src.Close()
	// 上传文件。
	err = bucket.PutObject(objectName.String(), src)
	if err != nil {
		log.Errorf("上传请求失败: %s", err.Error())
		c.JSON(500, errors.InternalServerError("go.micro.api.upload.aliyun", err.Error()))
	}
	var buffer bytes.Buffer
	buffer.WriteString(config.BucketHost)
	buffer.WriteString("/")
	buffer.Write(objectName.Bytes())

	c.JSON(200, map[string]string{
		"original": file.Filename,
		"name":     fileName,
		"path":     "/" + objectName.String(),
		"url":      buffer.String(),
	})
}
