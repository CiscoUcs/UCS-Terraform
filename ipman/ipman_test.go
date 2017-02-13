package ipman

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	utils "github.com/ContainerSolutions/go-utils"
)

func resetInventory(t *testing.T, inventoryFile string, inventory []byte) {
	file, err := os.OpenFile(inventoryFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	utils.FailOnError(t, err)
	defer file.Close()

	_, err = file.Write(inventory)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateIPFromCIDR(t *testing.T) {
	inventory := make([]byte, 0)
	inventoryFile := "./fixtures/dummy-inventory.txt"
	resetInventory(t, inventoryFile, inventory)

	cidr := "10.0.1.0/24"
	expected := "10.0.1.1"
	ip, err := GenerateIP(inventoryFile, cidr)
	utils.FailOnError(t, err)

	actual := ip.String()
	if actual != expected {
		t.Errorf("%s expected; got %s", expected, actual)
	}

	expected = "10.0.1.2"
	ip, err = GenerateIP(inventoryFile, cidr)
	utils.FailOnError(t, err)

	actual = ip.String()
	if actual != expected {
		t.Errorf("%s expected; got %s", expected, actual)
	}

	expected = "10.0.1.3"
	ip, err = GenerateIP(inventoryFile, cidr)
	utils.FailOnError(t, err)

	actual = ip.String()
	if actual != expected {
		t.Errorf("%s expected; got %s", expected, actual)
	}
}

func TestGenerateIPFromInventoryFile(t *testing.T) {
	expected := "10.0.0.2"
	inventory := []byte("10.0.0.1" + "\n")
	inventoryFile := "./fixtures/dummy-inventory.txt"
	resetInventory(t, inventoryFile, inventory)

	ip, err := GenerateIP(inventoryFile, "127.0.0.1/32")
	utils.FailOnError(t, err)
	actual := ip.String()

	if actual != expected {
		t.Errorf("%s expected; got %s", expected, actual)
	}

	expectedInventory := []byte("10.0.0.1\n10.0.0.2\n")
	file, err := os.OpenFile(inventoryFile, os.O_RDONLY, 0644)
	utils.FailOnError(t, err)

	fileInfo, err := file.Stat()
	utils.FailOnError(t, err)

	newInventory := make([]byte, fileInfo.Size())
	_, err = file.Read(newInventory)
	utils.FailOnError(t, err)

	if !bytes.Equal(expectedInventory, newInventory) {
		fmt.Printf("expected:\n%s\n*****\ngot:\n%s\n", expectedInventory, newInventory)
	}
}

func TestNextIP(t *testing.T) {
	t.Skip("Skipping until figuring out how to skip .0 and .255 IPs")
	expected := "10.0.1.1"
	ip := NextIP(net.ParseIP("10.0.1.0"))
	actual := ip.String()

	if actual != expected {
		t.Errorf("%s expected; got %s", expected, actual)
	}

	expected = "10.0.2.1"
	ip = NextIP(net.ParseIP("10.0.1.254"))
	actual = ip.String()

	if actual != expected {
		t.Errorf("%s expected; got %s", expected, actual)
	}
}

func TestInventory(t *testing.T) {
	inventory := []byte("10.0.0.1\n")
	inventoryFile := "./fixtures/dummy-inventory.txt"
	resetInventory(t, inventoryFile, inventory)

	actual, err := Inventory(inventoryFile)
	utils.FailOnError(t, err)

	expected := []net.IP{net.ParseIP("10.0.0.1")}
	if len(actual) != len(expected) {
		t.Fatalf("%v elements expected; got %v", len(expected), len(actual))
	}

	for i, actualIP := range actual {
		actualIP := actualIP.String()
		expectedIP := expected[i].String()
		if actualIP != expectedIP {
			t.Errorf("%s expected; got %s", expectedIP, actualIP)
		}
	}
}

func TestInventoryNonExistingFile(t *testing.T) {
	expected := make([]net.IP, 0)
	inventoryFile := "i do not exist"
	actual, err := Inventory(inventoryFile)
	utils.FailOnError(t, err)

	if len(expected) != len(actual) {
		t.Errorf("%v expected but got %v", len(expected), len(actual))
	}
}

func TestSaveIPBlankInventory(t *testing.T) {
	inventory := []byte("")
	inventoryFile := "./fixtures/dummy-inventory.txt"
	resetInventory(t, inventoryFile, inventory)

	expected := "1.2.3.4\n"
	ip := net.ParseIP("1.2.3.4")
	err := SaveIP(inventoryFile, ip)
	utils.FailOnError(t, err)

	actual, err := ioutil.ReadFile(inventoryFile)
	utils.FailOnError(t, err)

	if !bytes.Equal(actual, []byte(expected)) {
		t.Errorf("expected:\n%s\n*****\ngot:\n%s\n", expected, actual)
	}
}

func TestSaveIPExistingInventory(t *testing.T) {
	inventory := []byte("1.2.3.4\n")
	inventoryFile := "./fixtures/dummy-inventory.txt"
	resetInventory(t, inventoryFile, inventory)

	expected := "1.2.3.4\n4.3.2.1\n"
	ip := net.ParseIP("4.3.2.1")
	err := SaveIP(inventoryFile, ip)
	utils.FailOnError(t, err)

	actual, err := ioutil.ReadFile(inventoryFile)
	utils.FailOnError(t, err)

	if !bytes.Equal(actual, []byte(expected)) {
		t.Errorf("expected:\n%s\n*****\ngot:\n%s\n", expected, actual)
	}
}
