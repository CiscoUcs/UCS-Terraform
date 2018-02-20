package ucsinternal

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/ContainerSolutions/go-utils"
)

func TestMarshalLoginRequest(t *testing.T) {
	ex := []byte(`<aaaLogin inName="john" inPassword="doesecret"></aaaLogin>`)
	req := LoginRequest{
		Username: "john",
		Password: "doesecret",
	}

	out, err := req.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out, ex) {
		t.Errorf("%s expected; got %s", ex, out)
	}
}

func TestNewLoginResponse(t *testing.T) {
	data := []byte(`<aaaLogin cookie="some-cookie" response="yes" outCookie="this-is-the-out-cookie" outRefreshPeriod="123" outPriv="admin,read-only" outDomains="org-blah" outChannel="noencssl" outEvtChannel="noencssl" outSessionId="session-123" outVersion="2.2" outName="blahblah"></aaaLogin>`)
	cex := "some-cookie"
	rex := "yes"
	ocex := "this-is-the-out-cookie"
	odex := "org-blah"

	res, err := NewLoginResponse(data)
	if err != nil {
		t.Error(err)
	}

	if res.Cookie != cex {
		t.Errorf("%s expected; got %s", cex, res.Cookie)
	}

	if res.Response != rex {
		t.Errorf("%s expected; got %s", cex, res.Response)
	}

	if res.OutCookie != ocex {
		t.Errorf("%s expected; got %s", ocex, res.OutCookie)
	}

	if res.OutDomains != odex {
		t.Errorf("%s expected; got %s", odex, res.OutDomains)
	}
}

func TestMarshalLogoutRequest(t *testing.T) {
	ex := []byte(`<aaaLogout inCookie="chipsahoy!"></aaaLogout>`)
	req := LogoutRequest{
		Cookie: "chipsahoy!",
	}
	out, err := req.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out, ex) {
		t.Errorf("%s expected; got %s", ex, out)
	}
}

func TestMarshalDestroyRequest(t *testing.T) {
	pex := []byte(`<configConfMos cookie="chipsahoy!" inHierarchical="true"><inConfigs><pair key="katz/ls-foobar"><lsServer dn="katz/ls-foobar" status="deleted"></lsServer></pair></inConfigs></configConfMos>`)
	req := DestroyRequest{
		Name:         "foobar",
		TargetOrg:    "katz",
		Hierarchical: true,
	}
	out, err := req.Marshal("chipsahoy!")
	if err != nil {
		t.Fatalf("could not marshalize DestroyRequest:\n%s", err)
	}

	if !bytes.Equal(out, pex) {
		t.Errorf("%s expected; got\n\t\t%s", pex, out)
	}
}

func TestNewServiceProfileResponse(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/service-profile.xml")
	utils.FailOnError(t, err)

	res, err := NewServiceProfileResponse(data)
	if err != nil {
		t.Fatal(err)
	}

	dn := "org-root/ls-deathstar"
	if res.OutConfigs.ServerConfig.Dn != dn {
		t.Errorf("%s expected; got %s", dn, res.OutConfigs.ServerConfig.Dn)
	}

	name := "deathstar"
	if res.OutConfigs.ServerConfig.Name != name {
		t.Errorf("%s expected; got %s", name, res.OutConfigs.ServerConfig.Name)
	}

	srcTempl := "test-template"
	if res.OutConfigs.ServerConfig.SrcTempl != srcTempl {
		t.Errorf("%s expected; got %s", srcTempl, res.OutConfigs.ServerConfig.SrcTempl)
	}

	status := "created"
	if res.OutConfigs.ServerConfig.Status != status {
		t.Errorf("%s expected; got %s", status, res.OutConfigs.ServerConfig.Status)
	}
}
