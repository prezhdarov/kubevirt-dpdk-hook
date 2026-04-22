package hook

import (
	"fmt"
	"regexp"
	"strconv"

	"libvirt.org/go/libvirtxml"
)

func havePCIControllers(controllers []libvirtxml.DomainController) bool {
	for _, device := range controllers {

		if device.Type == "pci" {
			return true
		}
	}
	return false
}

func hugepageFromVMI(pagesize string) (libvirtxml.DomainMemoryHugepage, error) {

	var pagesizeRegex = regexp.MustCompile(`^(\d+)([A-Za-z]+)$`)

	pagesizeMatch := pagesizeRegex.FindStringSubmatch(pagesize)
	if len(pagesizeMatch) != 3 {
		return libvirtxml.DomainMemoryHugepage{}, fmt.Errorf("invalid pagesize: %s", pagesize)
	}

	size, err := strconv.ParseUint(pagesizeMatch[1], 10, 64)
	if err != nil {
		return libvirtxml.DomainMemoryHugepage{}, err
	}

	return libvirtxml.DomainMemoryHugepage{
		Size: uint(size),
		Unit: pagesizeMatch[2] + "B",
	}, nil
}
