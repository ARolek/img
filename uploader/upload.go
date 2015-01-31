package uploader

import (
	"crypto/md5"
	//	"errors"
	"fmt"
	"io"
	//	"log"
	"net/http"
	"os"
	"path/filepath"
	//	"strconv"
	//	"strings"

	"github.com/arolek/tools"

	"img/config"
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

/*
type RequestDetails struct {
	LogoKey string
	Ext     string //	file extension
	Width   float64
	Height  float64
	Version int64
}

//	parse a URL extracting information necessary to render the file
//	pull uri params and ext hbvsdn.png?w=300&h=300&tag=test&bg=#000&ver=3
//		w = width in px
//		h = height in px
//		ver = version index. if not provided will be set to -1
//		bg = background color in hex format [not implemented]
//		tag = a tracking tag [not implemented]
func (this *RequestDetails) Parse(r *http.Request, params map[string]string) error {
	var err error

	file := strings.Split(params["logo"], ".")
	if len(file) != 2 {
		return errors.New("missing file extension. try adding .png, .jpg, or .pdf to the url")
	}

	//	file details
	this.LogoKey = file[0]
	this.Ext = file[1]

	//	parse our request url
	vals := r.URL.Query()

	//	check width
	if _, ok := vals["w"]; ok {
		//do something here
		this.Width, err = strconv.ParseFloat(vals["w"][0], 64)
		if err != nil {
			return err
		}
	} else {
		this.Width = 450
	}

	//	check for height
	if _, ok := vals["h"]; ok {
		this.Height, err = strconv.ParseFloat(vals["h"][0], 64)
		if err != nil {
			return err
		}
	} else {
		this.Height = 450
	}

	//	check for version. if no version, set to -1
	if _, ok := vals["ver"]; ok {
		this.Version, err = strconv.ParseInt(vals["ver"][0], 10, 64)
		if err != nil {
			return err
		}
	} else {
		this.Version = -1
	}

	return nil
}

func Img(w http.ResponseWriter, r *http.Request) {
	var err error

	//	file hash
	key := r.URL.Path[len("/img/")]

	log.Println(key)

	rasterTmp, err := DownloadAndConvert(w, r)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//	send image to response
	http.ServeFile(w, r, rasterTmp)

	//	clean up rendered and downloaded files
	os.Remove(rasterTmp)
}

func DownloadAndConvert(w http.ResponseWriter, r *http.Request) {
	var err error

	req := RequestDetails{}
	err = req.Parse(r, params)
	if err != nil {
		return "", err
	}

	downloadTmp := filepath.Join("/tmp", tools.GetRand())

	//	our downloadURL includes the version identifier
	downloadURL := fmt.Sprintf("http://ogol.s3.amazonaws.com/")

	//	TODO: check if image exists at the given size and format - if so deliver
	//	if not, resize image. we will only do this if we need to optimize our flow
	//	for now we're going to render on demand and leverage the CDN

	//	download the file from S3
	if err = tools.HTTPdownload(downloadURL, downloadTmp); err != nil {
		return "", err
	}

	//	raster temp file with the file extension
	deliverTmp := fmt.Sprintf("%v.%v", downloadTmp, req.Ext)

	//	raserize the EPS. default is PNG
	switch req.Ext {
	case "jpg":

	}

	//	fastly.com CDN key used for purging
	w.Header().Set("Surrogate-Key", "hash")

	//	will cache the asset with fastly forever, but with the
	//	browser for only 1 hour. this way we don't have to
	//	re-render the asset every hour.
	//
	//	reference: https://docs.fastly.com/guides/tutorials/cache-control-tutorial
	w.Header().Set("Cache-Control", "max-age=3600")
	w.Header().Set("Surrogate-Control", "max-age=2592000")

	//	clean up the download temp
	go os.Remove(downloadTmp)

	return deliverTmp, nil
}
*/
