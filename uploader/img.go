package uploader

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"

	"img/config"
)

type QueryDetails struct {
	Action string
	Width  int64
	Height int64
	Algo   string //	resizing / thumbnail algo
}

//	parse a URL extracting information necessary to render the file
//	pull uri params and ext hbvsdn.png?w=300&h=300&tag=test&bg=#000&ver=3
//		w = width in px
//		h = height in px
//		ver = version index. if not provided will be set to -1
//		bg = background color in hex format [not implemented]
//		tag = a tracking tag [not implemented]
func (this *QueryDetails) Parse(url *url.URL) (err error) {
	//	parse our request url
	vals := url.Query()

	//	check action
	if _, ok := vals["action"]; ok {
		this.Action = vals["action"][0]
	} else {
		return errors.New("no action provided")
	}

	//	check width
	if _, ok := vals["w"]; ok {
		//	parse width
		this.Width, err = strconv.ParseInt(vals["w"][0], 10, 0)
		if err != nil {
			return
		}
	} else {
		this.Width = 450
	}

	//	check for height
	if _, ok := vals["h"]; ok {
		//	parse height
		this.Height, err = strconv.ParseInt(vals["h"][0], 10, 0)
		if err != nil {
			return
		}
	} else {
		this.Height = 450
	}

	//	check for algo
	if _, ok := vals["algo"]; ok {
		this.Algo = vals["algo"][0]
	}

	return
}

func Img(w http.ResponseWriter, r *http.Request) {
	var err error

	//	extract our key
	key := r.URL.Path[len("/img/"):]
	//	confirm we have key
	if key == "" {
		Error(w, "no key provided", http.StatusBadRequest)
		return
	}

	//	download file from s3
	tmpFilePath, err := FetchFromS3(key, *config.S3_BUCKET)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//	parse the request params
	query := QueryDetails{}
	if err = query.Parse(r.URL); err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//	location of the temp file we will deliver to the client
	var modFilePath string

	switch query.Action {
	case "crop":
		modFilePath, err = Crop(query, tmpFilePath)
		if err != nil {
			Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "thumbnail":
		modFilePath, err = Thumbnail(query, tmpFilePath)
		if err != nil {
			Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		Error(w, "action not supported: "+query.Action, http.StatusInternalServerError)
		return
	}

	//	set cache control headers
	//	fastly.com CDN key used for purging
	w.Header().Set("Surrogate-Key", key)

	//	will cache the asset with fastly forever, but with the
	//	browser for only 1 hour. this way we don't have to
	//	re-render the asset every hour.
	//
	//	reference: https://docs.fastly.com/guides/tutorials/cache-control-tutorial
	w.Header().Set("Cache-Control", "max-age=3600")
	w.Header().Set("Surrogate-Control", "max-age=2592000")

	//	send our file in the response
	http.ServeFile(w, r, modFilePath)

	//	clean up the created file
	err = os.Remove(modFilePath)
	if err != nil {
		//	who do we tell?
		log.Println(err)
	}
}

func Crop(query QueryDetails, tmpFilePath string) (filePath string, err error) {
	tmpImgFile, err := os.Open(tmpFilePath)
	defer tmpImgFile.Close()
	if err != nil {
		return
	}

	img, imgFormat, err := image.Decode(tmpImgFile)
	if err != nil {
		return
	}

	cImg, err := cutter.Crop(img, cutter.Config{
		Height: int(query.Height),
		Width:  int(query.Width),
		Mode:   cutter.Centered,
	})
	if err != nil {
		return
	}

	//	cropped temp file location
	filePath = filepath.Join(*config.FS_TEMP, RandomHash())

	toImg, err := os.Create(filePath)
	if err != nil {
		return
	}

	//	encode our image
	if err = EncodeImage(toImg, cImg, imgFormat); err != nil {
		return
	}

	return
}

//	creats a thumbnail of an image maintaing aspect ratio
//	if the width and height are larger than the original image, the original image is preserved
func Thumbnail(query QueryDetails, tmpFilePath string) (filePath string, err error) {
	tmpImgFile, err := os.Open(tmpFilePath)
	defer tmpImgFile.Close()
	if err != nil {
		return
	}

	img, imgFormat, err := image.Decode(tmpImgFile)
	if err != nil {
		return
	}

	var algo resize.InterpolationFunction
	switch query.Algo {
	case "nearestNeighbor":
		algo = resize.NearestNeighbor
		break
	case "bilinear":
		algo = resize.Bilinear
		break
	case "mitchellNetravali":
		algo = resize.MitchellNetravali
		break
	case "lanczos2":
		algo = resize.Lanczos2
		break
	case "lanczos3":
		algo = resize.Lanczos3
		break
	default:
		algo = resize.NearestNeighbor
	}

	//	resize the image
	cImg := resize.Thumbnail(uint(query.Width), uint(query.Height), img, algo)

	//	cropped temp file location
	filePath = filepath.Join(*config.FS_TEMP, RandomHash())

	toImg, err := os.Create(filePath)
	if err != nil {
		return
	}

	//	encode our image
	if err = EncodeImage(toImg, cImg, imgFormat); err != nil {
		return
	}

	return
}

//	encodes an image reference to a file reference based on the passed format
func EncodeImage(file *os.File, img image.Image, format string) (err error) {
	switch format {
	case "jpeg":
		if err = jpeg.Encode(file, img, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
			return
		}
	case "png":
		if err = png.Encode(file, img); err != nil {
			return
		}
	default:
		err = errors.New("image format not supported")
		return
	}

	return
}
