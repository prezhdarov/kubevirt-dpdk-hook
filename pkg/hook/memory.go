package hook

import (
	"fmt"
	"regexp"
	"strconv"

	virtv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/log"
	"libvirt.org/go/libvirtxml"
)

const (
	sharedMemoryBackingAccessMode = "shared"
	memfdMemoryBackingSourceType  = "memfd"
)

func configMemory(vmiMemory *virtv1.Memory, domainMemoryBacking *libvirtxml.DomainMemoryBacking) error {

	domainMemoryBacking.MemoryAccess = &libvirtxml.DomainMemoryAccess{Mode: sharedMemoryBackingAccessMode}
	domainMemoryBacking.MemorySource = &libvirtxml.DomainMemorySource{Type: memfdMemoryBackingSourceType}

	if vmiMemory != nil &&
		vmiMemory.Hugepages != nil &&
		vmiMemory.Hugepages.PageSize != "" {

		ugePage, err := hugepageFromVMI(vmiMemory.Hugepages.PageSize)
		if err != nil {
			return err
		}

		if domainMemoryBacking.MemoryHugePages == nil {
			domainMemoryBacking.MemoryHugePages = &libvirtxml.DomainMemoryHugepages{}
		}

		if len(domainMemoryBacking.MemoryHugePages.Hugepages) < 1 {
			log.Log.Infof("No hugepages configruration find, adding one with page size of %s", vmiMemory.Hugepages.PageSize)
			domainMemoryBacking.MemoryHugePages.Hugepages = append(domainMemoryBacking.MemoryHugePages.Hugepages, ugePage)
		} else {
			log.Log.Infof("Existing hugepages configruration found, updating to page size of %s", vmiMemory.Hugepages.PageSize)
			domainMemoryBacking.MemoryHugePages.Hugepages[0] = ugePage
		}

	}

	return nil
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
