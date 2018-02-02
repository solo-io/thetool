package util

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

func RunCmd(verbose, dryRun bool, binary string, args ...string) error {
	if verbose {
		fmt.Println(binary, args)
	}
	if dryRun {
		return nil
	}

	cmd := exec.Command(binary, args...)
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
