package fixtures

import (
	"io/ioutil"
)

const (
	SERVICE_PROFILE_XML = "./fixtures/service-profile.xml"
)

func loadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func ServiceProfileXML() ([]byte, error) {
	return loadFile(SERVICE_PROFILE_XML)
}
