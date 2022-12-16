package network

import (
	"fmt"
	"net"
)

type Interface struct {
	Name      string
	IpAddress string
}

func (i *Interface) String() string {
	return fmt.Sprintf("%s: %s", i.Name, i.IpAddress)
}

func GetInterfaces() []Interface {
	ifaces, _ := net.Interfaces()

	result := make([]Interface, 0)

	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			inter := Interface{}
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			inter.IpAddress = ip.String()
			inter.Name = i.Name
			result = append(result, inter)
		}
	}
	return result
}
