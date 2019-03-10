package putioSync

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	_ "golang.org/x/oauth2"
	"io"
	"log"
	"os"
	"path/filepath"
	_ "path/filepath"
	"time"
)
import "github.com/igungor/go-putio/putio"

func Sync(config Configuration) {
	ctx, client := buildClient(config)
	cutOutTime := time.Now().AddDate(0, 0, -config.FilesTTLInDays)
	deleteOlderFiles(config.Path, config, cutOutTime)
	syncLatest(config, client, ctx)
}

func deleteOlderFiles(path string, config Configuration, cutOutTime time.Time) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if empty, _ := IsEmpty(path); empty && info.IsDir() {
			log.Printf("Removing empty dir %s", path)
			err := os.Remove(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		if info.Name() != filepath.Base(filepath.Dir(path)) && info.ModTime().Before(cutOutTime) {
			log.Printf("removing %s\n", info.Name())
			err = os.Remove(path)
			if err != nil && os.IsExist(err) {
				deleteOlderFiles(path+"/", config, cutOutTime)
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer Close(f)

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func syncLatest(config Configuration, client *putio.Client, ctx context.Context) {
	findFrom, _ := time.ParseDuration(fmt.Sprintf("%dh", config.FromTimeInHours))
	latest := getLatestFiles(client, ctx, findFrom)
	for _, e := range latest {
		DownloadNode(client, ctx, e, config.Path)
	}
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func buildClient(config Configuration) (context.Context, *putio.Client) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	ctx := context.Background()
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	client := putio.NewClient(oauthClient)
	return ctx, client
}
