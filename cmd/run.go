package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"syscall"
)

func Run(ctx *cli.Context) error {
	command := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	command.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func Child(ctx *cli.Context) error {
	check(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	check(os.MkdirAll("rootfs/oldrootfs", 0700))
	check(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	check(os.Chdir("/"))

	command := exec.Command(os.Args[2], os.Args[3:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}