package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
)

const version = "3.10.1"
const protocURLTemplate = "https://github.com/protocolbuffers/protobuf/releases/download/v%s/protoc-%s-%s-x86_64.zip"
const protocZipPath = "bin/protoc"

var goosToProtocOS = map[string]string{
	"darwin": "osx",
	"linux":  "linux",
}

func main() {
	output := flag.String("output", "", "Path where we should write the protoc binary")
	flag.Parse()

	log.Printf("downloading protoc to local file %s ...", *output)
	outputFile, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	protocURL := fmt.Sprintf(protocURLTemplate, version, version, goosToProtocOS[runtime.GOOS])
	log.Printf("downloading protoc from %s ...", protocURL)
	resp, err := http.Get(protocURL)
	if err != nil {
		panic(err)
	}
	protocZipBytes, err := ioutil.ReadAll(resp.Body)
	err2 := resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(protocZipBytes), int64(len(protocZipBytes)))
	if err != nil {
		panic(err)
	}
	found := false
	for _, f := range zipReader.File {
		log.Println(f.Name)
		if f.Name == protocZipPath {
			fileReader, err := f.Open()
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(outputFile, fileReader)
			if err != nil {
				panic(err)
			}
			found = true
			break
		}
	}
	if !found {
		err = os.Remove(*output)
		if err != nil {
			panic(err)
		}
		panic("protoc not found")
	}
}
