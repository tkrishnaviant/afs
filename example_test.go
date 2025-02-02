package afs_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/asset"
	"github.com/viant/afs/file"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/option"
	"github.com/viant/afs/scp"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

//ExampleService_List reading location content
func ExampleService_List() {
	service := afs.New()
	ctx := context.Background()
	objects, err := service.List(ctx, "/tmp/folder")
	if err != nil {
		log.Fatal(err)
	}
	for _, object := range objects {
		fmt.Printf("%v %v\n", object.Name(), object.URL())
		if object.IsDir() {
			continue
		}
		reader, err := service.Download(ctx, object)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", data)
	}
}

//ExampleService_Upload uploading content
func ExampleService_Upload() {
	service := afs.New()
	ctx := context.Background()
	keyAuth, err := scp.LocalhostKeyAuth("")
	if err != nil {
		log.Fatal(err)
	}
	err = service.Init(ctx, "scp://127.0.0.1:22/", keyAuth)
	if err != nil {
		log.Fatal(err)
	}
	err = service.Upload(ctx, "scp://127.0.0.1:22/folder/asset.txt", 0644, strings.NewReader("test me"))
	if err != nil {
		log.Fatal(err)
	}
	ok, err := service.Exists(ctx, "scp://127.0.0.1:22/folder/asset.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("has file: %v\n", ok)
}

//ExampleService_Copy copy content
func ExampleService_Copy() {
	service := afs.New()
	ctx := context.Background()
	keyAuth, err := scp.LocalhostKeyAuth("")
	if err != nil {
		log.Fatal(err)
	}
	err = service.Copy(ctx, "s3://mybucket/myfolder", "scp://127.0.0.1/tmp", option.NewSource(), option.NewDest(keyAuth))
	if err != nil {
		log.Fatal(err)
	}
}

//ExampleService_Move moves content
func ExampleService_Move() {
	service := afs.New()
	ctx := context.Background()
	keyAuth, err := scp.LocalhostKeyAuth("")
	if err != nil {
		log.Fatal(err)
	}
	err = service.Move(ctx, "/tmp/transient/app", "scp://127.0.0.1/tmp", option.NewSource(), option.NewDest(keyAuth))
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleService_Uploader() {
	service := afs.New()
	ctx := context.Background()
	upload, closer, err := service.Uploader(ctx, "/tmp/clone")
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	assets := []*asset.Resource{
		asset.NewFile("asset1.txt", []byte("test 1"), 0644),
		asset.NewFile("asset2.txt", []byte("test 2"), 0644),
		asset.NewDir("folder1", file.DefaultDirOsMode),
		asset.NewFile("folder1/asset1.txt", []byte("test 1"), 0644),
		asset.NewFile("folder1/asset2.txt", []byte("test 2"), 0644),
	}
	for _, asset := range assets {
		relative := ""
		var reader io.Reader
		if strings.Contains(asset.Name, "/") {
			relative, _ = path.Split(asset.Name)
		}
		if !asset.Dir {
			reader = bytes.NewReader(asset.Data)
		}
		err = upload(ctx, relative, asset.Info(), reader)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//Example_DownloadError download error simulation example
func Example_DownloadError() {
	service := afs.NewFaker()
	ctx := context.Background()
	err := service.Upload(ctx, "gs://myBucket/folder/asset.txt", 0, strings.NewReader("some data"), option.NewUploadError(io.EOF))
	if err != nil {
		log.Fatalf("expect upload error: %v", err)
	}
}

//Example_DownloadError download error simulation example
func Example_UploadError() {
	service := afs.NewFaker()
	ctx := context.Background()
	err := service.Upload(ctx, "gs://myBucket/folder/asset.txt", 0, strings.NewReader("some data"), option.NewDownloadError(io.EOF))
	if err != nil {
		log.Fatal(err)
	}
	_, err = service.DownloadWithURL(ctx, "gs://myBucket/folder/asset.txt")
	if err != nil {
		log.Fatalf("expect download error: %v", err)
	}
}

//Example_DownloadError download error simulation example
func Example_ReaderError() {
	service := afs.NewFaker()
	ctx := context.Background()
	err := service.Upload(ctx, "gs://myBucket/folder/asset.txt", 0, strings.NewReader("some data"), option.NewReaderError(fmt.Errorf("test error")))
	if err != nil {
		log.Fatal(err)
	}
	reader, err := service.DownloadWithURL(ctx, "gs://myBucket/folder/asset.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	_, err = ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("expect download error: %v", err)
	}
}

func Example_IgnoreMatcher() {
	ignoreMatcher, err := matcher.NewIgnore([]string{"*.txt", ".ignore"})
	if err != nil {
		log.Fatal(err)
	}
	service := afs.New()
	ctx := context.Background()
	objects, err := service.List(ctx, "/tmp/folder", ignoreMatcher)
	if err != nil {
		log.Fatal(err)
	}
	for _, object := range objects {
		fmt.Printf("%v %v\n", object.Name(), object.URL())
		if object.IsDir() {
			continue
		}
	}

}
