package ucsclient

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	utils "github.com/ContainerSolutions/go-utils"
)

type StubHTTPClient struct {
	StatusCode int
	Body       []byte
}

func (c StubHTTPClient) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	if c.StatusCode >= 500 {
		return nil, errors.New("Something went wrong")
	} else {
		res := &http.Response{
			StatusCode: c.StatusCode,
			Body:       ioutil.NopCloser(bytes.NewReader(c.Body)),
		}
		return res, nil
	}
}

type StubHTTPClientWithAssertion struct {
	StatusCode      int
	Body            []byte
	ExpectedPayload []byte
	t               *testing.T
}

func (c StubHTTPClientWithAssertion) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	payload, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(payload, c.ExpectedPayload) {
		c.t.Errorf("%s expected; got %s", c.ExpectedPayload, payload)
	}
	res := &http.Response{
		StatusCode: c.StatusCode,
		Body:       ioutil.NopCloser(bytes.NewReader(c.Body)),
	}
	return res, nil
}

func newTestConfig() *Config {
	return &Config{
		IpAddress:             "1.2.3.4",
		Username:              "john",
		Password:              "doe",
		TslInsecureSkipVerify: true,
		LogLevel:              0,
		LogFilename:           "foo.log",
		AppName:               "testapp",
	}
}

func TestServiceProfileDN(t *testing.T) {
	expected := "org-foo/ls-bar"
	sp := ServiceProfile{
		Name:      "bar",
		TargetOrg: "org-foo",
	}
	out := sp.DN()
	if out != expected {
		t.Errorf("%s expected; got %s", expected, out)
	}
}

func TestServiceProfileToJSON(t *testing.T) {
	expected := `{"Name":"blah","Template":"foo","TargetOrg":"root","Hierarchical":true,"VNICs":[{"Name":"eth0","Mac":"01:02:03:04:05:06","CIDR":"127.0.0.1/24","Ip":"127.0.0.1"}]}`
	sp := ServiceProfile{
		Name:         "blah",
		Template:     "foo",
		TargetOrg:    "root",
		Hierarchical: true,
		VNICs: []VNIC{
			VNIC{
				Name: "eth0",
				Mac:  "01:02:03:04:05:06",
				Ip:   net.ParseIP("127.0.0.1"),
				CIDR: "127.0.0.1/24",
			},
		},
	}

	actual, err := sp.ToJSON()
	utils.FailOnError(t, err)

	if actual != expected {
		t.Errorf("%s expected; got %s", expected, actual)
	}
}

func TestNewUCSClient(t *testing.T) {
	config := newTestConfig()
	ucsClient := NewUCSClient(config)

	if config.IpAddress != ucsClient.ipAddress {
		t.Errorf("%v expected but got %v", config.IpAddress, ucsClient.ipAddress)
	}

	if config.Username != ucsClient.username {
		t.Errorf("%v expected but got %v", config.Username, ucsClient.username)
	}

	if config.Password != ucsClient.password {
		t.Errorf("%v expected but got %v", config.Password, ucsClient.password)
	}

	if config.TslInsecureSkipVerify != ucsClient.tslInsecureSkipVerify {
		t.Errorf("%v expected but got %v", config.TslInsecureSkipVerify, ucsClient.tslInsecureSkipVerify)
	}

	if config.AppName != ucsClient.appName {
		t.Errorf("%v expected but got %v", config.AppName, ucsClient.appName)
	}

	if config.LogLevel != ucsClient.Logger.Level {
		t.Errorf("%v expected but got %v", config.LogLevel, ucsClient.Logger.Level)
	}

	if ucsClient.httpClient == nil {
		t.Errorf("*http.Client expected but got nil")
	}
}

func TestEndpointURL(t *testing.T) {
	var endpointURLTests = []struct {
		ip  string
		uex string // expected url
	}{
		{"1.2.3.4", "https://1.2.3.4/nuova/"},
		{"127.0.0.1", "https://127.0.0.1/nuova/"},
		{"192.168.1.1", "https://192.168.1.1/nuova/"},
	}
	for _, test := range endpointURLTests {
		config := newTestConfig()
		config.IpAddress = test.ip
		client := NewUCSClient(config)
		out := client.endpointURL()
		if out != test.uex {
			t.Errorf("%s expected; got %s", test.uex, out)
		}
	}
}

func TestPostWithError(t *testing.T) {
	payload := []byte("blah")
	config := newTestConfig()
	ucsClient := NewUCSClient(config)

	for errorCode := 500; errorCode <= 505; errorCode++ {
		ucsClient.httpClient = StubHTTPClient{
			StatusCode: errorCode,
		}
		res, err := ucsClient.Post(payload)
		if len(res) > 0 {
			t.Errorf("expected blank response but got %s", res)
		}

		if err == nil {
			t.Errorf("error expected but got nil")
		}
	}
}

func TestPost(t *testing.T) {
	tests := []struct {
		Payload          []byte
		ExpectedResponse []byte
	}{
		{[]byte(`<someRequest></someRequest>`), []byte(`<foo><bar><baz /></bar></foo>`)},
		{[]byte(`<foo></bar>`), []byte(`<lol></lol>`)},
		{[]byte(`<blah></blah>`), []byte(`<katz></katz>`)},
	}
	config := newTestConfig()
	ucsClient := NewUCSClient(config)

	for _, test := range tests {
		ucsClient.httpClient = StubHTTPClient{
			StatusCode: 200,
			Body:       test.ExpectedResponse,
		}

		out, err := ucsClient.Post([]byte(test.Payload))
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(out, test.ExpectedResponse) {
			t.Errorf("%s expected; got %s", test.ExpectedResponse, out)
		}
	}
}

func TestLoginWithServerError(t *testing.T) {
	config := newTestConfig()
	ucsClient := NewUCSClient(config)
	ucsClient.httpClient = StubHTTPClient{
		StatusCode: 500,
	}
	err := ucsClient.Login()
	if err == nil {
		t.Errorf("error expected but got nil")
	}
}

func TestLogin(t *testing.T) {
	body := []byte(`<aaaLogin cookie="some-cookie" response="yes" outCookie="chipsahoy!" outRefreshPeriod="123" outPriv="admin,read-only" outDomains="" outChannel="noencssl" outEvtChannel="noencssl" outSessionId="session-123" outVersion="2.2" outName="blahblah"></aaaLogin>`)
	cex := "chipsahoy!"
	config := newTestConfig()
	ucsClient := NewUCSClient(config)
	ucsClient.httpClient = StubHTTPClient{
		StatusCode: 200,
		Body:       body,
	}
	err := ucsClient.Login()
	if err != nil {
		t.Error(err)
	}

	if ucsClient.cookie != cex {
		t.Errorf("%s expected; got %s", cex, ucsClient.cookie)
	}

}

func TestLoginSendsProperPayload(t *testing.T) {
	tests := []struct {
		Username string
		Password string
		Payload  []byte
	}{
		{"john", "doe", []byte(`<aaaLogin inName="john" inPassword="doe"></aaaLogin>`)},
		{"anna", "pascal", []byte(`<aaaLogin inName="anna" inPassword="pascal"></aaaLogin>`)},
		{"dmitry", "ovtcharov", []byte(`<aaaLogin inName="dmitry" inPassword="ovtcharov"></aaaLogin>`)},
	}
	config := newTestConfig()

	for _, test := range tests {
		config.Username = test.Username
		config.Password = test.Password
		ucsClient := NewUCSClient(config)
		ucsClient.httpClient = StubHTTPClientWithAssertion{
			Body:            make([]byte, 0),
			ExpectedPayload: test.Payload,
			t:               t,
		}
		ucsClient.Login()
	}
}

func TestIsLoggedIn(t *testing.T) {
	config := newTestConfig()
	ucsClient := NewUCSClient(config)

	if ucsClient.IsLoggedIn() {
		t.Errorf("expected false; got true")
	}

	ucsClient.cookie = "foo"
	if !ucsClient.IsLoggedIn() {
		t.Errorf("expected true; got false")
	}
}

func TestCreateServiceProfile(t *testing.T) {
	pex := []byte(`<lsInstantiateNNamedTemplate cookie="chipsahoy!" dn="org-root/ls-test-template" inTargetOrg="org-root" inHierarchical="false" inErrorOnExisting="true"><inNameSet><dn value="deathstar"></dn></inNameSet></lsInstantiateNNamedTemplate>`)
	body, err := ioutil.ReadFile("testdata/service-profile.xml")
	utils.FailOnError(t, err)

	if err != nil {
		t.Fatalf("could not fetch ServiceProfileXML fixture:\n%s", err)
	}

	config := newTestConfig()
	ucsClient := NewUCSClient(config)
	ucsClient.cookie = "chipsahoy!"
	ucsClient.httpClient = StubHTTPClientWithAssertion{
		StatusCode:      200,
		Body:            body,
		ExpectedPayload: pex,
		t:               t,
	}
	sp := &ServiceProfile{
		Name:         "deathstar",
		Template:     "test-template",
		TargetOrg:    "org-root",
		Hierarchical: false,
	}
	created, err := ucsClient.CreateServiceProfile(sp)
	if err != nil {
		t.Error(err)
	}

	if !created {
		t.Error("expected true but got false")
	}
}

func TestMarshalServiceProfile(t *testing.T) {
	pex := []byte(`<lsInstantiateNNamedTemplate cookie="chipsahoy!" dn="org-root/ls-test-template" inTargetOrg="org-root" inHierarchical="false" inErrorOnExisting="true"><inNameSet><dn value="deathstar"></dn></inNameSet></lsInstantiateNNamedTemplate>`)
	sp := &ServiceProfile{
		Name:         "deathstar",
		Template:     "test-template",
		TargetOrg:    "org-root",
		Hierarchical: false,
	}
	out, err := sp.Marshal("chipsahoy!")
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(pex, out) {
		t.Errorf("%s expected; got %s", pex, out)
	}
}

func TestDestroy(t *testing.T) {
	pex := []byte(`<configConfMos cookie="chipsahoy!" inHierarchical="true"><inConfigs><pair key="org-root/ls-deathstar"><lsServer dn="org-root/ls-deathstar" status="deleted"></lsServer></pair></inConfigs></configConfMos>`)
	body := []byte(`<configConfMos cookie="chipsahoy!" response="yes"><outConfigs></outConfigs></configConfMos>`)
	config := newTestConfig()
	ucsClient := NewUCSClient(config)
	ucsClient.cookie = "chipsahoy!"
	ucsClient.httpClient = StubHTTPClientWithAssertion{
		StatusCode:      200,
		Body:            body,
		ExpectedPayload: pex,
		t:               t,
	}
	ucsClient.Destroy("deathstar", "org-root", true)
}

func TestLogout(t *testing.T) {
	pex := []byte(`<aaaLogout inCookie="chipsahoy!"></aaaLogout>`)
	config := newTestConfig()
	ucsClient := NewUCSClient(config)
	ucsClient.cookie = "chipsahoy!"
	ucsClient.httpClient = StubHTTPClientWithAssertion{
		StatusCode:      200,
		Body:            make([]byte, 0),
		ExpectedPayload: pex,
		t:               t,
	}
	ucsClient.Logout()

	if ucsClient.cookie != "" {
		t.Errorf("expected cookie to be unset; got %s", ucsClient.cookie)
	}

}

func TestConfigResolveDN(t *testing.T) {
	req, err := utils.Fixture("config-resolve-dn-req.xml")
	utils.FailOnError(t, err)

	res, err := utils.Fixture("config-resolve-dn-res.xml")
	utils.FailOnError(t, err)

	dn := "org-root/ls-foobar"
	config := newTestConfig()
	ucsClient := NewUCSClient(config)
	ucsClient.cookie = "chipsahoy!"
	ucsClient.httpClient = StubHTTPClientWithAssertion{
		StatusCode:      200,
		Body:            res,
		ExpectedPayload: req,
		t:               t,
	}

	expectedSP := ServiceProfile{
		Name:      "foobar",
		Template:  "mamamia",
		TargetOrg: "org-root",
		VNICs: []VNIC{
			VNIC{
				Name: "eth1",
				Mac:  "00:25:B5:00:00:8F",
			},
			VNIC{
				Name: "eth0",
				Mac:  "00:25:B5:00:00:9F",
			},
		},
	}

	sp, err := ucsClient.ConfigResolveDN(dn)
	utils.FailOnError(t, err)

	if sp.Name != expectedSP.Name {
		t.Errorf("%s expected; got %s", expectedSP.Name, sp.Name)
	}

	if sp.Template != expectedSP.Template {
		t.Errorf("%s expected; got %s", expectedSP.Template, sp.Template)
	}

	if sp.TargetOrg != expectedSP.TargetOrg {
		t.Errorf("%s expected; got %s", expectedSP.TargetOrg, sp.TargetOrg)
	}

	if sp.Hierarchical != expectedSP.Hierarchical {
		t.Errorf("%s expected; got %s", expectedSP.Hierarchical, sp.Hierarchical)
	}

	if x, y := len(sp.VNICs), len(expectedSP.VNICs); x != y {
		t.Errorf("number of VNICs don't match: %d = %d", x, y)
	}

	// Because we know how the data is organised in the fixture it's safe to assume the order of
	// the vNICs.
	for i := 0; i < len(sp.VNICs); i++ {
		vnic := sp.VNICs[i]
		evnic := expectedSP.VNICs[i]

		if vnic.Name != evnic.Name {
			t.Errorf("%s expected; got %s", evnic.Name, vnic.Name)
		}

		if vnic.Mac != evnic.Mac {
			t.Errorf("%s expected; got %s", evnic.Mac, vnic.Mac)
		}

		if vnic.CIDR != evnic.CIDR {
			t.Errorf("%s expected; got %s", evnic.CIDR, vnic.CIDR)
		}
	}
}
