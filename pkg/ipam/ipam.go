package ipam

import (
	"fmt"
	"net"
	"strings"

	"k8s.io/klog"
)

// Manager - handles the addresses for each namespace/vip
var Manager []ipManager

// ipManager defines the mapping to a namespace and address pool
type ipManager struct {
	namespace      string
	cidr           string
	ipRange        string
	addressManager map[string]bool
	hosts          []string
}

// FindAvailableHostFromRange - will look through the cidr and the address Manager and find a free address (if possible)
func FindAvailableHostFromRange(namespace, ipRange string) (string, error) {

	// Look through namespaces and update one if it exists
	for x := range Manager {
		if Manager[x].namespace == namespace {
			// Check that the address range is the same
			if Manager[x].ipRange != ipRange {
				// If not rebuild the available hosts
				ah, err := buildHostsFromRange(ipRange)
				if err != nil {
					return "", err
				}
				Manager[x].hosts = ah
			}
			// TODO - currently we search (incrementally) through the list of hosts
			for y := range Manager[x].hosts {
				// find a host that is marked false (i.e. unused)
				if !Manager[x].addressManager[Manager[x].hosts[y]] {
					// Mark it to used
					Manager[x].addressManager[Manager[x].hosts[y]] = true
					return Manager[x].hosts[y], nil
				}
			}
			// If we have found the manager for this namespace and not returned an address then we've expired the range
			return "", fmt.Errorf("No addresses available in [%s] range [%s]", namespace, ipRange)

		}
	}
	ah, err := buildHostsFromRange(ipRange)
	if err != nil {
		return "", err
	}
	// If it doesn't exist then it will need adding
	newManager := ipManager{
		namespace:      namespace,
		addressManager: make(map[string]bool),
		hosts:          ah,
		ipRange:        ipRange,
	}
	Manager = append(Manager, newManager)

	for x := range newManager.hosts {
		if !Manager[x].addressManager[newManager.hosts[x]] {
			Manager[x].addressManager[newManager.hosts[x]] = true
			return newManager.hosts[x], nil
		}
	}
	return "", fmt.Errorf("No addresses available in [%s] range [%s]", namespace, ipRange)

}

// FindAvailableHostFromCidr - will look through the cidr and the address Manager and find a free address (if possible)
func FindAvailableHostFromCidr(namespace, cidr string) (string, error) {

	// Look through namespaces and update one if it exists
	for x := range Manager {
		if Manager[x].namespace == namespace {
			// Check that the address range is the same
			if Manager[x].cidr != cidr {
				// If not rebuild the available hosts
				ah, err := buildHostsFromCidr(cidr)
				if err != nil {
					return "", err
				}
				Manager[x].hosts = ah
			}
			// TODO - currently we search (incrementally) through the list of hosts
			for y := range Manager[x].hosts {
				// find a host that is marked false (i.e. unused)
				if !Manager[x].addressManager[Manager[x].hosts[y]] {
					// Mark it to used
					Manager[x].addressManager[Manager[x].hosts[y]] = true
					return Manager[x].hosts[y], nil
				}
			}
			// If we have found the manager for this namespace and not returned an address then we've expired the range
			return "", fmt.Errorf("No addresses available in [%s] range [%s]", namespace, cidr)

		}
	}
	ah, err := buildHostsFromCidr(cidr)
	if err != nil {
		return "", err
	}
	// If it doesn't exist then it will need adding
	newManager := ipManager{
		namespace:      namespace,
		addressManager: make(map[string]bool),
		hosts:          ah,
		cidr:           cidr,
	}
	Manager = append(Manager, newManager)

	for x := range newManager.hosts {
		if !Manager[x].addressManager[newManager.hosts[x]] {
			Manager[x].addressManager[newManager.hosts[x]] = true
			return newManager.hosts[x], nil
		}
	}
	return "", fmt.Errorf("No addresses available in [%s] range [%s]", namespace, cidr)

}

// ReleaseAddress - removes the mark on an address
func ReleaseAddress(namespace, address string) error {
	for x := range Manager {
		if Manager[x].namespace == namespace {
			Manager[x].addressManager[address] = false
			return nil
		}
	}
	return fmt.Errorf("Unable to release address [%s] in namespace [%s]", address, namespace)
}

// buildHostsFromCidr - Builds a list of addresses in the cidr
func buildHostsFromCidr(cidr string) ([]string, error) {
	var ips []string

	// Split the ipranges (comma separated)
	cidrs := strings.Split(cidr, ",")
	if len(cidrs) == 0 {
		return nil, fmt.Errorf("Unable to parse IP cidrs [%s]", cidr)
	}

	for x := range cidrs {

		ip, ipnet, err := net.ParseCIDR(cidrs[x])
		if err != nil {
			return nil, err
		}

		var cidrips []string
		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			cidrips = append(cidrips, ip.String())
		}

		// remove network address and broadcast address
		lenIPs := len(cidrips)
		switch {
		case lenIPs < 2:
			ips = append(ips, cidrips...)

		default:
			ips = append(ips, cidrips[1:len(cidrips)-1]...)
		}
	}
	return removeDuplicateAddresses(ips), nil
}

// IPStr2Int - Converts the IP address in string format to an integer
func IPStr2Int(ip string) uint {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0
	}
	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24
}

//IPInt2Str - Converts the IP address in integer format to an string
func IPInt2Str(i uint) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)
	return ip.String()
}

// buildHostsFromRange - Builds a list of addresses in the cidr
func buildHostsFromRange(ipRangeString string) ([]string, error) {
	var ips []string

	// Split the ipranges (comma separated)
	ranges := strings.Split(ipRangeString, ",")
	if len(ranges) == 0 {
		return nil, fmt.Errorf("Unable to parse IP ranges [%s]", ipRangeString)
	}

	for x := range ranges {
		ipRange := strings.Split(ranges[x], "-")
		// Make sure we have x.x.x.x-x.x.x.x
		if len(ipRange) != 2 {
			return nil, fmt.Errorf("Unable to parse IP range [%s]", ranges[x])
		}

		firstIP := IPStr2Int(ipRange[0])
		lastIP := IPStr2Int(ipRange[1])
		fmt.Printf("firstIP=%d, lastIP=%d\n", firstIP, lastIP)
		if firstIP > lastIP {
			// swap
			firstIP, lastIP = lastIP, firstIP
		}

		for ip := firstIP; ip <= lastIP; ip++ {
			ips = append(ips, IPInt2Str(ip))
		}

		klog.Infof("Rebuilding addresse cache, [%d] addresses exist", len(ips))
	}
	return removeDuplicateAddresses(ips), nil
}

func removeDuplicateAddresses(arr []string) []string {
	addresses := map[string]bool{}
	uniqueAddresses := []string{} // Keep all keys from the map into a slice.

	for i := range arr {
		addresses[arr[i]] = true
	}
	for j := range addresses {
		uniqueAddresses = append(uniqueAddresses, j)
	}
	return uniqueAddresses
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
