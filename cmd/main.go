package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	vmiJSON   = flag.String("vmi", "", "VMI to change in JSON format")
	domainXML = flag.String("domain", "", "Domain spec in XML format")
)

func localAddresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		}

		for _, a := range addrs {
			fmt.Printf("%v - %v\n", i.Name, a)
		}
	}
}

func main() {

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Parse()

	//if *vmiJSON == "" || *domainXML == "" {
	//	log.Printf("--vmi and --domain cannot be undefined")
	//	os.Exit(1)
	//}

	log.Println(*vmiJSON)

	log.Println(*domainXML)

	fmt.Println(*domainXML)

	for _, env := range os.Environ() {
		fmt.Println(env)
	}

	localAddresses()
	os.Exit(1)
}
