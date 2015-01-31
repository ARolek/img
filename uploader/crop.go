package uploader

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oliamb/cutter"

	"img/config"
)

func Crop(w http.ResponseWriter, r *http.Request) {
	//	read hash, fetch from s3

	filePath := ""

	imgFile, err := os.Open(filePath)
	defer imgFile.Close()

	img, imgFormat, err := image.Decode(imgFile)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("image format: ", imgFormat)

	cImg, err := cutter.Crop(img, cutter.Config{
		Height: 600,
		Width:  600,
		Mode:   cutter.Centered,
	})
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//	temp file location
	filePath = filepath.Join(*config.FS_TEMP, RandomHash())

	toImg, err := os.Create(filePath)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch imgFormat {
	case "jpeg":
		if err = jpeg.Encode(toImg, cImg, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
			Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "png":
		if err = png.Encode(toImg, cImg); err != nil {
			Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		Error(w, "image format not supported", http.StatusInternalServerError)
		return
	}

	//	delete temp file

	//	delete downloaded file
}
