package ucsclient

import (
	"testing"
)

var config = &Config{
	IpAddress:             "1.2.3.4",
	Username:              "john",
	Password:              "doe",
	TslInsecureSkipVerify: true,
	LogLevel:              0,
	LogFilename:           "foo.log",
	AppName:               "testapp",
}

func TestClient(t *testing.T) {
	client := config.Client()
	if client == nil {
		t.Errorf("config.Client() = nil; expected an instance of a UCSClient")
	}
}
