package main

import (
	"fmt"
	"github.com/abrance/something/pkg/gen_host"
	"log"
)

func main() {
	caseGenHost()
}

func caseGenHost() {
	nh1 := gen_host.NodeHost{
		IpLs: []string{"192.168.1.2", "192.168.1.1"},
		TyId: gen_host.HostnameTypeM,
	}
	nh2 := gen_host.NodeHost{
		IpLs: []string{"172.16.0.1"},
		TyId: gen_host.HostnameSeByNodeId,
	}
	nh3 := gen_host.NodeHost{
		IpLs: []string{"172.16.0.10"},
		TyId: gen_host.HostnameSeByNodeId,
	}
	scHost := gen_host.NodeHost{
		IpLs: []string{"172.24.0.1"},
		TyId: gen_host.HostnameSC,
	}
	nh2.SetSeNodeId("$nodeid2")
	nh3.SetSeNodeId("$nodeid3")

	entries := gen_host.GenerateEntries([]gen_host.NodeHost{nh1, nh2, nh3, scHost})
	fmt.Println(entries)
	err := gen_host.GenerateHostsFile(entries, "./hosts")
	if err != nil {
		log.Fatal(err)
	}
}
