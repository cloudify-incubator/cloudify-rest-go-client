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

	if err := syscall.Unshare(syscall.CLONE_FILES | syscall.CLONE_FS | syscall.CLONE_NEWPID | syscall.CLONE_SYSVSEM); err != nil {
		log.Fatalf("Could not clone new fs and proc: %s", err)
	}
	if err := syscall.Chroot(dir); err != nil {
		log.Fatalf("Could not change root: %s", err)
	}

	// create and mount all dirs
	if err := syscall.Mount("sysfs", "/sys", "sysfs", syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID, ""); err != nil {
		log.Fatalf("mount sys: %s", err)
	}
	defer syscall.Unmount("/sys", syscall.MNT_DETACH)

	if err := syscall.Mount("proc", "/proc", "proc", syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID, ""); err != nil {
		log.Fatalf("mount proc: %s", err)
	}
	defer syscall.Unmount("/proc", syscall.MNT_DETACH)

	if err := syscall.Mount("tmpfs", "/tmp", "tmpfs", 0, "size=65536k"); err != nil {
		log.Fatalf("mount tmp: %s", err)
	}
	defer syscall.Unmount("/tmp", syscall.MNT_DETACH)

	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", 0, "size=65536k"); err != nil {
		log.Fatalf("mount dev: %s", err)
	}
	defer syscall.Unmount("/dev", syscall.MNT_DETACH)
}
