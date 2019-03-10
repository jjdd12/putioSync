package putioSync

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Token           string
	Path            string
	FromTimeInHours int
	FilesTTLInDays  int
}

func LoadConfig(filename string) (Configuration, error) {
	conf := Configuration{}
	file, err := os.Open(filename)
	if err != nil {
		return Configuration{}, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	return conf, err
}
