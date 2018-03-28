package main

import (
	"log"
	"os"
	"path"
	"syscall"
	"time"
)

func makedev(major int, minor int) int {
	return (minor & 0xff) | (major&0xfff)<<8
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cloudify error: %s\n", err.Error())
	}

	// reset umask befor do anything
	oldUmask := syscall.Umask(0)

	// create dirs for overlayfs
	baseDir := path.Join(dir, "base")
	if err := os.Mkdir(baseDir,
		syscall.S_IRUSR|syscall.S_IXUSR|
		syscall.S_IRGRP|syscall.S_IXGRP|
		syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}
	if err := os.Mkdir(path.Join(dir, "data"),
		syscall.S_IRUSR|syscall.S_IXUSR|
		syscall.S_IRGRP|syscall.S_IXGRP|
		syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	combinedDir := path.Join(dir, "overlay")
	if err := os.Mkdir(combinedDir,
		syscall.S_IRUSR|syscall.S_IXUSR|
		syscall.S_IRGRP|syscall.S_IXGRP|
		syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	// mount overlayfs
	/*if err := syscall.Mount("sysfs", "/sys", "sysfs",
			syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID, ""); err != nil {
		log.Fatalf("mount overlayfs: %s", err)
	}*/
	//defer syscall.Unmount("/sys", syscall.MNT_DETACH)

	// create and mount all dirs
	if err := os.Mkdir(path.Join(combinedDir, "sys"),
		syscall.S_IRUSR|syscall.S_IXUSR|
		syscall.S_IRGRP|syscall.S_IXGRP|
		syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	if err := os.Mkdir(path.Join(combinedDir, "proc"),
		syscall.S_IRUSR|syscall.S_IXUSR|
		syscall.S_IRGRP|syscall.S_IXGRP|
		syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	if err := os.Mkdir(path.Join(combinedDir, "tmp"),
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
		syscall.S_IRGRP|syscall.S_IWGRP|syscall.S_IXGRP|
		syscall.S_IROTH|syscall.S_IWOTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	if err := os.Mkdir(path.Join(combinedDir, "dev"),
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
		syscall.S_IRGRP|syscall.S_IXGRP|
		syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	if err := syscall.Unshare(syscall.CLONE_FILES | syscall.CLONE_FS | syscall.CLONE_NEWPID | syscall.CLONE_SYSVSEM); err != nil {
		log.Fatalf("Could not clone new fs and proc: %s", err)
	}
	if err := syscall.Chroot(combinedDir); err != nil {
		log.Fatalf("Could not change root: %s", err)
	}

	// create and mount all dirs
	if err := syscall.Mount("sysfs", "/sys", "sysfs",
		syscall.MS_NODEV|syscall.MS_NOEXEC|syscall.MS_NOSUID, ""); err != nil {
		log.Fatalf("mount sys: %s", err)
	}
	defer syscall.Unmount("/sys", syscall.MNT_DETACH)

	if err := syscall.Mount("proc", "/proc", "proc",
		syscall.MS_NODEV|syscall.MS_NOEXEC|syscall.MS_NOSUID, ""); err != nil {
		log.Fatalf("mount proc: %s", err)
	}
	defer syscall.Unmount("/proc", syscall.MNT_DETACH)

	if err := syscall.Mount("tmpfs", "/tmp", "tmpfs",
		0, "size=65536k,mode=0755"); err != nil {
		log.Fatalf("mount tmp: %s", err)
	}
	defer syscall.Unmount("/tmp", syscall.MNT_DETACH)

	if err := syscall.Mount("tmpfs", "/dev", "tmpfs",
		0, "size=65536k,mode=0755"); err != nil {
		log.Fatalf("mount dev: %s", err)
	}
	defer syscall.Unmount("/dev", syscall.MNT_DETACH)

	if err := syscall.Mknod("/dev/full",
		syscall.S_IRUSR|syscall.S_IWUSR|
		syscall.S_IRGRP|syscall.S_IWGRP|
		syscall.S_IROTH|syscall.S_IWOTH|
		syscall.S_IFCHR, makedev(1, 7)); err != nil {
		log.Fatalf("mknod /dev/full: %s", err)
	}

	if err := syscall.Mknod("/dev/ptmx",
		syscall.S_IRUSR|syscall.S_IWUSR|
		syscall.S_IRGRP|syscall.S_IWGRP|
		syscall.S_IROTH|syscall.S_IWOTH|
		syscall.S_IFCHR, makedev(5, 2)); err != nil {
		log.Fatalf("mknod /dev/ptmx: %s", err)
	}

	if err := syscall.Mknod("/dev/random",
		syscall.S_IRUSR|syscall.S_IWUSR|
		syscall.S_IRGRP|syscall.S_IROTH|
		syscall.S_IFCHR, makedev(1, 8)); err != nil {
		log.Fatalf("mknod /dev/random: %s", err)
	}

	if err := syscall.Mknod("/dev/urandom",
		syscall.S_IRUSR|syscall.S_IWUSR|
		syscall.S_IRGRP|syscall.S_IROTH|
		syscall.S_IFCHR, makedev(1, 9)); err != nil {
		log.Fatalf("mknod /dev/urandom: %s", err)
	}

	if err := syscall.Mknod("/dev/zero",
		syscall.S_IRUSR|syscall.S_IWUSR|
		syscall.S_IRGRP|syscall.S_IWGRP|
		syscall.S_IROTH|syscall.S_IWOTH|
		syscall.S_IFCHR, makedev(1, 5)); err != nil {
		log.Fatalf("mknod /dev/zero: %s", err)
	}

	if err := syscall.Mknod("/dev/tty",
		syscall.S_IRUSR|syscall.S_IWUSR|
		syscall.S_IRGRP|syscall.S_IWGRP|
		syscall.S_IROTH|syscall.S_IWOTH|
		syscall.S_IFCHR, makedev(5, 0)); err != nil {
		log.Fatalf("mknod /dev/tty: %s", err)
	}
	// go back with rights
	syscall.Umask(oldUmask)
	log.Printf("Wait 30 seconds before revert everything.")
	time.Sleep(30 * time.Second)
}
