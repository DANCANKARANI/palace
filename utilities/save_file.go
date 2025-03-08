package utilities

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofiber/fiber/v2"
)

func SaveFile(c *fiber.Ctx, fieldName string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	accountName := os.Getenv("ACCOUNT_NAME")
	accountKey := os.Getenv("ACCOUNT_KEY")
	containerName := os.Getenv("CONTAINER_NAME")

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", errors.New("failed to create credentials")
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)
	url, _ := url.Parse(serviceURL)
	pipeline := azblob.NewPipeline(cred, azblob.PipelineOptions{})
	serviceURLObj := azblob.NewServiceURL(*url, pipeline)

	containerURL := serviceURLObj.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(file.Filename)

	_, err = azblob.UploadStreamToBlockBlob(context.Background(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		log.Println("error uploading file:",err.Error())
		return "", errors.New("failed to upload file")
	}

	return blobURL.String(), nil
}