package uploader

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arolek/tools"

	"github.com/arolek/img/config"
)

type File struct {
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

//	handles uploading of a file and generating a random file hash for it
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// the FormFile function takes in the POST input id file
	postData, _, err := r.FormFile("file")
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer postData.Close()

	//	hash for filename
	hash := RandomHash()

	//	new file path
	tmpPath := filepath.Join(*config.FS_TEMP, hash)

	//	temp file location
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
		return
	}

	// write the content from POST to the file
	_, err = io.Copy(tmpFile, postData)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpFile.Close()

	//	move file to S3 bucket
	if err = MoveToS3(tmpPath, hash); err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//	clean up temp file
	if err = os.Remove(tmpPath); err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//	prepare response
	fileData := File{
		Hash: hash,
		URL:  *config.CDN + "/" + hash,
	}

	//	send response to client
	Success(w, fileData)
}

//	generates a random file hash
func RandomHash() string {
	//	generate a random string
	rand := tools.GetRand()
	//	generate a md5 has from the random string
	return fmt.Sprintf("%x", md5.Sum([]byte(rand)))
}
