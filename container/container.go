package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

func makedev(major int, minor int) int {
	return (minor & 0xff) | (major&0xfff)<<8
}

func createDirInContainer(combinedDir string) {
	// reset umask befor do anything
	oldUmask := syscall.Umask(0)

	// create and mount all dirs
	if err := os.Mkdir(path.Join(combinedDir, "/sys"),
		syscall.S_IRUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	if err := os.Mkdir(path.Join(combinedDir, "/proc"),
		syscall.S_IRUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	if err := os.Mkdir(path.Join(combinedDir, "/tmp"),
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IWOTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	devDir := path.Join(combinedDir, "/dev")
	os.RemoveAll(devDir)
	if err := os.Mkdir(devDir,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}

	if err := syscall.Mknod(path.Join(devDir, "/full"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|
			syscall.S_IROTH|syscall.S_IWOTH|
			syscall.S_IFCHR, makedev(1, 7)); err != nil {
		log.Printf("mknod /dev/full: %s", err)
	}

	if err := syscall.Mknod(path.Join(devDir, "/ptmx"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|
			syscall.S_IROTH|syscall.S_IWOTH|
			syscall.S_IFCHR, makedev(5, 2)); err != nil {
		log.Printf("mknod /dev/ptmx: %s", err)
	}

	if err := syscall.Mknod(path.Join(devDir, "/random"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IROTH|
			syscall.S_IFCHR, makedev(1, 8)); err != nil {
		log.Printf("mknod /dev/random: %s", err)
	}

	if err := syscall.Mknod(path.Join(devDir, "/urandom"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IROTH|
			syscall.S_IFCHR, makedev(1, 9)); err != nil {
		log.Printf("mknod /dev/urandom: %s", err)
	}

	if err := syscall.Mknod(path.Join(devDir, "/zero"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|
			syscall.S_IROTH|syscall.S_IWOTH|
			syscall.S_IFCHR, makedev(1, 5)); err != nil {
		log.Printf("mknod /dev/zero: %s", err)
	}

	if err := syscall.Mknod(path.Join(devDir, "/tty"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|
			syscall.S_IROTH|syscall.S_IWOTH|
			syscall.S_IFCHR, makedev(5, 0)); err != nil {
		log.Printf("mknod /dev/tty: %s", err)
	}

	// go back with rights
	syscall.Umask(oldUmask)
}

func mountEverythingAndRun(combinedDir string, argv0 string, argv []string) {
	log.Printf("I am going to run: %+v\n\n", strings.Join(argv, " "))
	if err := syscall.Unshare(syscall.CLONE_FILES | syscall.CLONE_FS | syscall.CLONE_NEWPID | syscall.CLONE_SYSVSEM); err != nil {
		log.Fatalf("Could not clone new fs and proc: %s", err)
	}
	if err := syscall.Chroot(combinedDir); err != nil {
		log.Fatalf("Could not change root: %s", err)
	}

	if err := syscall.Mount("proc", "/proc", "proc",
		syscall.MS_NODEV|syscall.MS_NOEXEC|syscall.MS_NOSUID, ""); err != nil {
		log.Fatalf("mount proc: %s", err)
	}
	defer syscall.Unmount("/proc", syscall.MNT_DETACH)

	var procInfo syscall.SysProcAttr
	procInfo.Chroot = "/" // combinedDir
	var env syscall.ProcAttr
	env.Env = []string{"PATH=/usr/sbin:/usr/bin:/sbin:/bin"}
	// TODO: hackish way, but ok for now
	env.Files = []uintptr{0, 1, 2}
	env.Sys = &procInfo

	pid, err := syscall.ForkExec(argv0, argv, &env)
	if err != nil {
		log.Fatalf("Issues with run: %s", err)
	}

	syscall.Wait4(pid, nil, 0, nil)
	log.Printf("Wait 10 seconds before revert everything.")
	time.Sleep(10 * time.Second)
}

func main() {
	var commandList []string
	commandList = os.Args[1:]
	if len(commandList) == 0 {
		commandList = []string{"/bin/sh"}
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cloudify error: %s\n", err.Error())
	}

	// create dirs for overlayfs
	baseDir := path.Join(dir, "base")
	if err := os.Mkdir(baseDir,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}
	log.Printf("As Operation System filesystem will be used: %s\n", baseDir)

	dataDir := path.Join(dir, "data")
	if err := os.Mkdir(dataDir,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}
	log.Printf("Data changes will be stored in: %s\n", dataDir)

	workDir := path.Join(dir, "work")
	if err := os.Mkdir(workDir,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}
	// try to delete, on error
	defer os.RemoveAll(workDir)

	combinedDir := path.Join(dir, "overlay")
	if err := os.Mkdir(combinedDir,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		log.Printf("Not critical: %s\n", err.Error())
	}
	// try to delete, on error
	defer os.RemoveAll(combinedDir)

	// https://www.kernel.org/doc/Documentation/filesystems/overlayfs.txt
	mountOptions := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", baseDir, dataDir, workDir)
	// mount overlayfs
	if err := syscall.Mount("overlay", combinedDir, "overlay", 0, mountOptions); err != nil {
		log.Printf("Not critical, already merged: %s", err)
	}
	// try to delete, on error
	defer syscall.Unmount(combinedDir, syscall.MNT_DETACH)

	createDirInContainer(combinedDir)

	// real work
	mountEverythingAndRun(combinedDir, commandList[0], commandList)
}
