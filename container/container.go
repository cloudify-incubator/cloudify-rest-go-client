package main

import (
	"os"
	"path"
	"log"
	"syscall"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cloudify error: %s\n", err.Error())
	}

	// create and mount all dirs
	os.Mkdir(path.Join(dir, "sys"),
				syscall.S_IRUSR | syscall.S_IXUSR |
				syscall.S_IRGRP | syscall.S_IXGRP |
				syscall.S_IROTH | syscall.S_IXOTH)

	os.Mkdir(path.Join(dir, "proc"),
				syscall.S_IRUSR | syscall.S_IXUSR |
				syscall.S_IRGRP | syscall.S_IXGRP |
				syscall.S_IROTH | syscall.S_IXOTH)

	os.Mkdir(path.Join(dir, "tmp"),
				syscall.S_IRUSR | syscall.S_IWUSR | syscall.S_IXUSR |
				syscall.S_IRGRP | syscall.S_IWGRP | syscall.S_IXGRP |
				syscall.S_IROTH | syscall.S_IWOTH | syscall.S_IXOTH)

	os.Mkdir(path.Join(dir, "dev"),
				syscall.S_IRUSR | syscall.S_IWUSR | syscall.S_IXUSR |
				syscall.S_IRGRP | syscall.S_IXGRP |
				syscall.S_IROTH | syscall.S_IXOTH)

	if err := syscall.Unshare(syscall.CLONE_FS | syscall.CLONE_FILES | syscall.CLONE_NEWPID); err != nil {
		log.Fatalf("Could not clone new fs and proc: %s", err)
	}
	if err := syscall.Chroot(dir); err != nil {
		log.Fatalf("Could not change root: %s", err)
	}
}
