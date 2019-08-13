package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/smartystreets/s3"
)

func main() {
	url, err := s3.NewPresignedGet(
		s3.EnvironmentCredentials(),
		s3.ExpiresIn(time.Hour),
		s3.Region("us-west-1"),
		s3.Bucket("smartystreets-downloads"),
		s3.Key("us-street-api/linux-amd64/latest.tar.gz"),
	)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("GET:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Panicln(err)
	}
	defer func() { _ = response.Body.Close() }()

	fmt.Println("HTTP", response.Status)
	if response.StatusCode != http.StatusOK {
		_, _ = io.Copy(os.Stderr, response.Body)
	}
}
