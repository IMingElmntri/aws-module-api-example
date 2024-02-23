package aws_apis

import (
	"log"
	"net/http"
	"strings"
	"fmt"
	"os"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"

	"github.com/IMingElmntri/aws-module-api-example/shared"
)

func (a *APIs) uploadFile(c *gin.Context) {
	var result = shared.APIResponse{
		Code: 200,
		Message: "",
		Data: nil,
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		result.Message = err.Error()
		c.JSON(http.StatusBadRequest, result)
	}
	defer file.Close()

	// Check file extension
	ext := filepath.Ext(header.Filename)
	if !isValidExtension(ext) {
		result.Message = "Invalid file type"
		c.JSON(http.StatusBadRequest, result)
	}

	// Generate a UUID for temporary file name
	tempFileName := uuid.New().String() + ext
	tempFilePath := filepath.Join(os.TempDir(), tempFileName)

	// Create a temporary file to store the uploaded file
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		result.Message = err.Error()
		c.JSON(http.StatusInternalServerError, result)
	}
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		result.Message = err.Error()
		c.JSON(http.StatusInternalServerError, result)
	}
	tempFile.Seek(0, 0)

	awsBucket := viper.GetString("AWS_S3_BUCKET")
	awsRegion := viper.GetString("AWS_S3_REGION")

	var contentType = ""
	switch ext {
		case ".csv":
			contentType = "text/csv"
		case ".png":
			contentType = "image/png"
		case ".jpeg", ".jpg":
			contentType = "image/jpeg"
		default:
			contentType = "" // Fallback content type if extension is unknown
	}

	if len(contentType) == 0 {
		c.JSON(http.StatusInternalServerError, "Invalid Data")
	}

	err = a.params.BucketConnector.UploadFile(
		tempFileName, 
		awsBucket,
		tempFile,
		contentType,
	)
	if err != nil {
		result.Message = err.Error()
		c.JSON(http.StatusInternalServerError, result)
	}
	result.Data = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", awsBucket, awsRegion, tempFileName)
	c.JSON(http.StatusOK, result)
}

func isValidExtension(ext string) bool {
	allowedExts := []string{".csv", ".png", ".jpeg", ".jpg"}
	for _, allowedExt := range allowedExts {
		if strings.EqualFold(ext, allowedExt) {
			return true
		}
	}
	return false
}

func (a *APIs) listBuckets(c *gin.Context) {

	var result = shared.APIResponse{
		Code: 200,
		Message: "",
		Data: nil,
	}

	buckets, err := a.params.BucketConnector.ListBuckets()
	if err != nil {
		log.Fatal(err)
		result.Message = "Error"
		c.JSON(http.StatusInternalServerError, result)
	}

	result.Data = buckets
	c.JSON(http.StatusOK, result)
}