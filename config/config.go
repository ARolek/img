package config

import (
	"flag"
	"log"
	"runtime"

	"github.com/stvp/go-toml-config"
)

var (
	Config = flag.String("conf", "config.toml", "reference to config file")

	//	how many CPUs do we have
	CPUcount = runtime.NumCPU()

	//	http
	HTTP_PORT        = config.String("http.port", ":8080")
	HTTP_STATIC_ROOT = config.String("http.staticRoot", "http/static")

	//	filesystem
	FS_TEMP = config.String("fs.temp", "/tmp")

	//	aws
	AWS_ACCESS_KEY_ID     = config.String("aws.accessKeyId", "")
	AWS_SECRET_ACCESS_KEY = config.String("aws.secretAccessKey", "")
	S3_BUCKET             = config.String("aws.s3.bucket", "ogol")

	// cdn base
	CDN = config.String("cdn.url", "")

	//	fastly
	FASTLY_SERVICEID = config.String("fastly.serviceId", "")
	FASTLY_APIKEY    = config.String("fastly.apiKey", "")
)

func init() {
	//	parse command line flags
	flag.Parse()

	if err := config.Parse(*Config); err != nil {
		log.Println("error parsing ", *Config, ": ", err, " using defaults.")
	}
}
