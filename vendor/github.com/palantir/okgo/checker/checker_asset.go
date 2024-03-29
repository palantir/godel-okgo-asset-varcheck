// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checker

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/palantir/godel/v2/framework/pluginapi"
	"github.com/palantir/okgo/okgo"
	"github.com/pkg/errors"
)

type assetChecker struct {
	assetPath       string
	cfgYML          string
	checkerType     okgo.CheckerType
	checkerPriority okgo.CheckerPriority
	checkerMultiCPU okgo.CheckerMultiCPU
}

func (c *assetChecker) Type() (okgo.CheckerType, error) {
	return c.checkerType, nil
}

func (c *assetChecker) Priority() (okgo.CheckerPriority, error) {
	return c.checkerPriority, nil
}

func (c *assetChecker) MultiCPU() (okgo.CheckerMultiCPU, error) {
	return c.checkerMultiCPU, nil
}

func (c *assetChecker) VerifyConfig() error {
	verifyConfigCmd := exec.Command(c.assetPath, verifyConfigCmdName,
		"--"+commonCmdConfigYMLFlagName, c.cfgYML,
	)
	if _, err := runCommand(verifyConfigCmd); err != nil {
		return err
	}
	return nil
}

func (c *assetChecker) Check(pkgs []string, projectDir string, stdout io.Writer) {
	checkCmd := exec.Command(c.assetPath, append([]string{
		checkCmdName,
		"--" + commonCmdConfigYMLFlagName, c.cfgYML,
		"--" + pluginapi.ProjectDirFlagName, projectDir,
	}, pkgs...)...)
	checkCmd.Stdout = stdout
	checkCmd.Stderr = stdout

	if err := checkCmd.Run(); err != nil {
		// if running check failed, write failure as its own issue
		okgo.WriteErrorAsIssue(err, stdout)
	}
}

func (c *assetChecker) RunCheckCmd(args []string, stdout io.Writer) {
	execArgs := []string{
		runCheckCmdCmdName,
	}
	if len(args) > 0 {
		// if arguments were provided, prepend with a standalone "--" so that the arguments are not interpreted by the
		// CLI. This allows flags to to be properly passed through.
		execArgs = append(execArgs, "--")
	}
	execArgs = append(execArgs, args...)

	runCheckCmdCmd := exec.Command(c.assetPath, execArgs...)
	runCheckCmdCmd.Stdout = stdout
	runCheckCmdCmd.Stderr = stdout

	if err := runCheckCmdCmd.Run(); err != nil {
		// if running check failed, write failure
		if _, ok := err.(*exec.ExitError); !ok {
			_, _ = fmt.Fprintf(stdout, "command %v failed with error %v\n", runCheckCmdCmd.Args, err)
		}
	}
}

func runCommand(cmd *exec.Cmd) ([]byte, error) {
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return outputBytes, errors.New(strings.TrimSpace(strings.TrimPrefix(string(outputBytes), "Error: ")))
	}
	return outputBytes, nil
}
