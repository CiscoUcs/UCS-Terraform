package ipman

import (
	"bytes"
	"net"
	"os"
)

// Generates a new IP. The criteria for generating a new IP with this method is
// the following:
// Fetches the list of IPs from the given inventory file. If there is at least
// one IP in that list it will use it as its base to generate the next one.
// If there are no IPs in the inventory file then it will generate a new IP
// based on the given CIDR.
// Independent of the method used for generating the new IP, this method will
// save the whole list of IPs held in memory into the given inventory file.
func GenerateIP(inventoryFile, cidr string) (net.IP, error) {
	var ip net.IP

	inventory, err := Inventory(inventoryFile)
	if err != nil {
		return nil, err
	}

	// If there is an ip already in the inventory then generate a new one based on that
	if len(inventory) > 0 {
		lastIP := inventory[len(inventory)-1]
		ip = NextIP(lastIP)
	} else {
		x, _, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		ip = NextIP(x)
	}

	// Now that the IP has been generated, let's save it to the inventory file.
	err = SaveIP(inventoryFile, ip)
	if err != nil {
		return nil, err
	}

	return ip, nil
}

// NOTE:
// Unfortunatelly this implementation only supports up to 254 IPs in a single subnetwork.
// For ALPHA release this is ok but we must figure out a way to make this function to return
// IPs that do not end on .255 nor .0.
// I haven't been able yet to figure out a way to make it work. I'll leave this here for a bit
// until my brain comes back to me with a solution.
func NextIP(ip net.IP) net.IP {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
	return ip
}

// Returns an array of IPs (net.IP) from the given inventory file path.
// If the given inventory file does not exist it'll return an empty array.
func Inventory(inventoryFile string) ([]net.IP, error) {
	// Early return an empty array if the inventory file does not exist.
	file, err := os.Open(inventoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make([]net.IP, 0), nil
		} else {
			return nil, err
		}
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	b := make([]byte, fileInfo.Size())
	_, err = file.Read(b)
	if err != nil {
		return nil, err
	}

	// If the inventory file is empty the bytes array will contain only one element (EOL).
	if len(b) <= 1 {
		return nil, nil
	}

	// Trim the EOL character which gets appended by the `ioutil.ReadFile` function
	if b[len(b)-1] == 10 {
		b = b[0 : len(b)-1]
	}

	// 10 is unicode for \r.
	// IPs are expected to be stored in the inventory file one per line
	// and that is why the separation character expected is return (\r).
	unparsedIPs := bytes.Split(b, []byte{10})
	ips := make([]net.IP, len(unparsedIPs))

	// Loop over the IPs which at this stage are still strings and
	// convert them into actual net.IP objects.
	for i, unparsedIP := range unparsedIPs {
		ips[i] = net.ParseIP(string(unparsedIP))
	}
	return ips, nil
}

// Takes an IP (net.IP) and dumps ip into the given
// inventory file path.
func SaveIP(inventoryFile string, ip net.IP) error {
	file, err := os.OpenFile(inventoryFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := []byte(ip.String() + "\n")
	_, err = file.Write(data)
	return err
}
