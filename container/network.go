package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"syscall"
)

const (
	NsRunDir  = "/var/run/netns"
	SelfNetNs = "/proc/self/ns/net"
)

func main() {
	netNsPath := path.Join(NsRunDir, "myns")
	os.Mkdir(NsRunDir, 0755)

	if err := syscall.Mount(NsRunDir, NsRunDir, "none", syscall.MS_BIND, ""); err != nil {
		log.Fatalf("Could not create Network namespace: %s", err)
	}

	fd, err := syscall.Open(netNsPath, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_EXCL, 0)
	if err != nil {
		log.Fatalf("Could not create Network namespace: %s", err)
	}
	syscall.Close(fd)

	if err := syscall.Unshare(syscall.CLONE_NEWNET); err != nil {
		log.Fatalf("Could not clone new Network namespace: %s", err)
	}

	if err := syscall.Mount(SelfNetNs, netNsPath, "none", syscall.MS_BIND, ""); err != nil {
		log.Fatalf("Could not Mount Network namespace: %s", err)
	}

	if err := syscall.Unmount(netNsPath, syscall.MNT_DETACH); err != nil {
		log.Fatalf("Could not Unmount new Network namespace: %s", err)
	}

	if err := syscall.Unlink(netNsPath); err != nil {
		log.Fatalf("Could not Unlink new Network namespace: %s", err)
	}

	ifcs, _ := net.Interfaces()
	for _, ifc := range ifcs {
		fmt.Printf("%#v\n", ifc)
	}
}
