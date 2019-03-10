//A Go library to sync files from a remote Put.io account to a local directory
//It allows you to set a  ttl for local files and the age of the files you want to import
//
// Usage:
//	import "github.com/jhony/putioSync"
//
//	func main() {
//		config,err := putioSync.LoadConfig("user_home/dir/.config/putioSync.json")
//		if err != nil {
//			log.Fatal("There was an issue loading the config", err)
//		}
//		putioSync.Sync(config)
//	}
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

//Given a Configuration creates a downstream sync, that pulls newely added files
//and removes older files based on the config
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

// given a directory address checks if the directory is empty
func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer close(f)

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
		downloadNode(client, ctx, e, config.Path)
	}
}

func close(c io.Closer) {
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
