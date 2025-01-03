package tokenring

import (
	"log"
	"net"
)

func GetLocalAddr() (string, string) {
	// Get a list of all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the interfaces and get their addresses
	for _, iface := range interfaces {
		// Skip loopback interfaces (127.0.0.1)
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// Get the addresses for the interface
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		// Loop through addresses and print the first non-loopback IP address
		for _, addr := range addrs {
			// We check if the address is an IP address (not a MAC address)
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() {
				// Print the first non-loopback IP address
				return ipnet.IP.String(), ""
			}
		}
	}
	return "", "No addresses found"
}
