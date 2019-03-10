package putioSync

import (
	"context"
	"fmt"
	"github.com/igungor/go-putio/putio"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func getLatestFiles(client *putio.Client, ctx context.Context, findFrom time.Duration) []putio.File {
	files, _, _ := client.Files.List(ctx, 0)

	latest := Filter(files, func(file putio.File) bool {
		return file.CreatedAt.After(time.Now().Add(-findFrom))

	})
	return latest
}

func downloadNode(client *putio.Client, ctx context.Context, e putio.File, rootPath string) {
	ios, _ := client.Files.Download(ctx, e.ID, false, nil)
	path := fmt.Sprintf("%s%s", rootPath, e.Name)
	if e.IsDir() {
		CreateDirIfNotExist(path)
		children, _, _ := client.Files.List(ctx, e.ID)
		for _, ch := range children {
			downloadNode(client, ctx, ch, filepath.Clean(path)+"/"+e.Name)
		}
		return
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}
	newFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found %s", e.Name)
	defer newFile.Close()
	io.Copy(newFile, ios)
	if err != nil {
		log.Fatal(err)
		ctx.Err()
	}
	log.Printf("Downloaded %s byte file.\n", e.Name)
}

type File putio.File

func Filter(vs []putio.File, f func(putio.File) bool) []putio.File {
	vsf := make([]putio.File, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func CreateDirIfNotExist(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
		return true
	}
	return false
}
