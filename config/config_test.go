package config_test

import (
	"chpassh/config"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	jsonFile := `
{
  "User": "USERNAME",
  "Oldpassword": "OLDPASSWORD",
  "Newpassword": "NEWPASSWORD",
  "Hosts": [
    {
      "Name": "mail.server",
      "Port": 22
    },
    {
      "Name": "web.server",
      "Port": 2022
    }
  ]
}`
	tempFile, err := ioutil.TempFile(os.TempDir(), "config_tmp.json")
	if err != nil {
		t.Fatalf("Can not create json tempfile %v", err)
		t.Fail()
	}
	_, err = tempFile.WriteString(jsonFile)
	if err != nil {
		t.Fatalf("Can not write to tempfile %s: %v", tempFile.Name(), err)
	}
	tempFile.Close()
	conf := config.Config{}
	err = conf.LoadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("Can not load config file %s %v", tempFile.Name(), err)
	}
	hosts := conf.Hosts
	assert.Equal(t, conf.User, "USERNAME", "Username password not loaded correctly")
	assert.Equal(t, conf.Oldpassword, "OLDPASSWORD", "Old password not loaded correctly")
	assert.Equal(t, conf.Newpassword, "NEWPASSWORD", "New password not loaded correctly")

	assert.Equal(t, hosts[0].Name, "mail.server", "Hostname not loaded correctly")
	assert.Equal(t, hosts[1].Name, "web.server", "Hostname not loaded correctly")

	assert.Equal(t, hosts[0].Port, 22, "Port not loaded correctly")
	assert.Equal(t, hosts[1].Port, 2022, "Port not loaded correctly")

	os.Remove(tempFile.Name())
}
