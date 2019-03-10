package putioSync

import (
	"gotest.tools/assert"
	"testing"
)

func TestSync(t *testing.T) {
	//TODO
}

func TestDeleteOlderFiles(t *testing.T) {
	//TODO
}

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("something_bogus")
	assert.Error(t, err, "open something_bogus: no such file or directory")
}
