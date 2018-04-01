/*
Copyright (c) 2018 GigaSpaces Technologies Ltd. All rights reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package container

import (
	"fmt"
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
	sysDir := path.Join(combinedDir, "/sys")
	os.RemoveAll(sysDir)
	if err := os.Mkdir(sysDir,
		syscall.S_IRUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		fmt.Printf("Not critical: %s\n", err.Error())
	}

	procDir := path.Join(combinedDir, "/proc")
	os.RemoveAll(procDir)
	if err := os.Mkdir(procDir,
		syscall.S_IRUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		fmt.Printf("Not critical: %s\n", err.Error())
	}

	tmpDir := path.Join(combinedDir, "/tmp")
	if _, err := os.Stat(tmpDir); err != nil {
		if err := os.Mkdir(tmpDir,
			syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
				syscall.S_IRGRP|syscall.S_IWGRP|syscall.S_IXGRP|
				syscall.S_IROTH|syscall.S_IWOTH|syscall.S_IXOTH); err != nil {
			fmt.Printf("Not critical: %s\n", err.Error())
		}
	}

	devDir := path.Join(combinedDir, "/dev")
	os.RemoveAll(devDir)
	if err := os.Mkdir(devDir,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		fmt.Printf("Not critical: %s\n", err.Error())
	}

	if err := syscall.Mknod(path.Join(devDir, "/full"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|
			syscall.S_IROTH|syscall.S_IWOTH|
			syscall.S_IFCHR, makedev(1, 7)); err != nil {
		fmt.Printf("mknod /dev/full: %s\n", err.Error())
	}

	if err := syscall.Mknod(path.Join(devDir, "/ptmx"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|
			syscall.S_IROTH|syscall.S_IWOTH|
			syscall.S_IFCHR, makedev(5, 2)); err != nil {
		fmt.Printf("mknod /dev/ptmx: %s\n", err.Error())
	}

	if err := syscall.Mknod(path.Join(devDir, "/random"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IROTH|
			syscall.S_IFCHR, makedev(1, 8)); err != nil {
		fmt.Printf("mknod /dev/random: %s\n", err.Error())
	}

	if err := syscall.Mknod(path.Join(devDir, "/urandom"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IROTH|
			syscall.S_IFCHR, makedev(1, 9)); err != nil {
		fmt.Printf("mknod /dev/urandom: %s\n", err.Error())
	}

	if err := syscall.Mknod(path.Join(devDir, "/zero"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|
			syscall.S_IROTH|syscall.S_IWOTH|
			syscall.S_IFCHR, makedev(1, 5)); err != nil {
		fmt.Printf("mknod /dev/zero: %s\n", err.Error())
	}

	if err := syscall.Mknod(path.Join(devDir, "/tty"),
		syscall.S_IRUSR|syscall.S_IWUSR|
			syscall.S_IRGRP|syscall.S_IWGRP|
			syscall.S_IROTH|syscall.S_IWOTH|
			syscall.S_IFCHR, makedev(5, 0)); err != nil {
		fmt.Printf("mknod /dev/tty: %s\n", err.Error())
	}

	// go back with rights
	syscall.Umask(oldUmask)
}

func mountEverythingAndRun(combinedDir string, argv0 string, argv []string) int {
	fmt.Printf("I am going to run: %+v\n", strings.Join(argv, " "))

	procDir := path.Join(combinedDir, "/proc")
	if err := syscall.Mount("proc", procDir, "proc",
		syscall.MS_NODEV|syscall.MS_NOEXEC|syscall.MS_NOSUID, ""); err != nil {
		fmt.Printf("mount proc: %s\n", err.Error())
		return 1
	}
	defer syscall.Unmount(procDir, syscall.MNT_DETACH)

	var procInfo syscall.SysProcAttr
	procInfo.Chroot = combinedDir
	var env syscall.ProcAttr
	env.Env = []string{"PATH=/usr/sbin:/usr/bin:/sbin:/bin"}
	// TODO: hackish way, but ok for now
	env.Files = []uintptr{0, 1, 2}
	env.Sys = &procInfo
	env.Dir = "/"

	pid, err := syscall.ForkExec(argv0, argv, &env)
	if err != nil {
		fmt.Printf("Issues with run: %s\n", err.Error())
		return 1
	}

	syscall.Wait4(pid, nil, 0, nil)
	fmt.Printf("Wait 10 seconds before revert everything.\n")
	time.Sleep(10 * time.Second)
	return 0
}

// Run - execute command inside controller
func Run(baseDir, dataDir, tempDir string, commandList []string) int {
	fmt.Printf("As Operation System filesystem will be used: %s\n", baseDir)

	if _, err := os.Stat(dataDir); err != nil {
		if err := os.Mkdir(dataDir,
			syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
				syscall.S_IRGRP|syscall.S_IXGRP|
				syscall.S_IROTH|syscall.S_IXOTH); err != nil {
			fmt.Printf("Not critical: %s\n", err.Error())
		}
	}
	fmt.Printf("Data changes will be stored in: %s\n", dataDir)

	workDir := path.Join(tempDir, "work")
	if err := os.Mkdir(workDir,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		fmt.Printf("Not critical: %s\n", err.Error())
	}
	// try to delete, on error
	defer os.RemoveAll(workDir)

	combinedDir := path.Join(tempDir, "overlay")
	if err := os.Mkdir(combinedDir,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR|
			syscall.S_IRGRP|syscall.S_IXGRP|
			syscall.S_IROTH|syscall.S_IXOTH); err != nil {
		fmt.Printf("Not critical: %s\n", err.Error())
	}
	// try to delete, on error
	defer os.RemoveAll(combinedDir)

	// https://www.kernel.org/doc/Documentation/filesystems/overlayfs.txt
	mountOptions := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", baseDir, dataDir, workDir)
	// mount overlayfs
	if err := syscall.Mount("overlay", combinedDir, "overlay", 0, mountOptions); err != nil {
		fmt.Printf("Overlay fs: %s\n", err.Error())
		return 1
	}
	// try to delete, on error
	defer syscall.Unmount(combinedDir, syscall.MNT_DETACH)

	createDirInContainer(combinedDir)

	// real work
	return mountEverythingAndRun(combinedDir, commandList[0], commandList)
}
