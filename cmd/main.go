package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	vmiJSON   = flag.String("vmi", "", "VMI to change in JSON format")
	domainXML = flag.String("domain", "", "Domain spec in XML format")
)

func main() {

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Parse()

	if *vmiJSON == "" || *domainXML == "" {
		log.Printf("--vmi and --domain cannot be undefined")
		os.Exit(1)
	}

	log.Printf(*vmiJSON)

	log.Printf(*domainXML)

	fmt.Printf(*domainXML)

}
