# putioSync
[![Build Status](https://travis-ci.org/jjdd12/putioSync.svg?branch=master)](https://travis-ci.org/jjdd12/putioSync)

A Go library to sync files from a remote Put.io account to a local directory\
It allows you to set a  ttl for local files and the age of the files you want to import

## Usage

```
package main

import (
	"github.com/mitchellh/go-homedir"
	"log"
)
import "github.com/jhony/putioSync"

func main() {
	home, _ := homedir.Dir()
	config,err := putioSync.LoadConfig( home +"/.config/putioSync.json")
	if err != nil {
		log.Fatal("There was an issue loading te config", err)
	}
	putioSync.Sync(config)
}
```
