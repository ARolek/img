package main

import (
	"log"
	"net/http"

	"img/config"
	"img/uploader"
)

func main() {
	http.HandleFunc("/upload", uploader.UploadHandler)
	http.HandleFunc("/img", uploader.Crop)

	log.Println("listening on port ", *config.HTTP_PORT)

	log.Fatal(http.ListenAndServe(":"+*config.HTTP_PORT, nil))
}
