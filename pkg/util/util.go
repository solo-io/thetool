package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
)

func RunCmd(verbose, dryRun bool, binary string, args ...string) error {
	return RunCmdContext(nil, verbose, dryRun, binary, args...)
}

func RunCmdContext(ctx context.Context, verbose, dryRun bool, binary string, args ...string) error {
	if verbose {
		fmt.Println(binary, args)
	}
	if dryRun {
		return nil
	}

	var cmd *exec.Cmd
	if ctx != nil {
		cmd = exec.CommandContext(ctx, binary, args...)
	} else {
		cmd = exec.Command(binary, args...)
	}
	if verbose {
		cmdStdout, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Wrapf(err, "unable to create StdOut pipe for %s", binary)
		}
		stdoutScanner := bufio.NewScanner(cmdStdout)
		go func() {
			prefix := binary + ": "
			for stdoutScanner.Scan() {
				fmt.Println(prefix, stdoutScanner.Text())
			}
		}()

		cmdStderr, err := cmd.StderrPipe()
		if err != nil {
			return errors.Wrapf(err, "unable to create StdErr pipe for %s", binary)
		}
		stderrScanner := bufio.NewScanner(cmdStderr)
		go func() {
			prefix := binary + ": "
			for stderrScanner.Scan() {
				fmt.Println(prefix, stderrScanner.Text())
			}
		}()
	}
	err := cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "unable to start %s", binary)
	}
	err = cmd.Wait()
	if err != nil {
		return errors.Wrapf(err, "unable to run %s", binary)
	}
	return nil
}

func DockerRun(verbose, dryRun bool, containerName string, args ...string) error {
	ctx, cancel := dockerContext(containerName)
	if dryRun {
		defer cancel()
	}
	return RunCmdContext(ctx, verbose, dryRun, "docker", args...)
}

func dockerContext(containerName string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(name string) {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalCh:
			cancel()
			stopContainer(name)
		case <-ctx.Done():
			return
		}
	}(containerName)
	return ctx, cancel
}

func stopContainer(name string) {
	err := RunCmd(false, false, "docker", "stop", name)
	if err != nil {
		fmt.Println("error stopping docker container ", name)
	}
}

func Copy(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}
