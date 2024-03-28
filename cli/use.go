// Copyright 2022 Tristan Isham. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/tristanisham/zvm/cli/meta"
)

func (z *ZVM) Use(ver string) error {
	err := z.getVersion(ver)
	if errors.Is(err, os.ErrNotExist) {

		// We might want to add an option in settings to toggle these in case some
		// version map uses "default" as a version name.
		switch ver {
		case "default":
			if err = z.Settings.ResetVersionMap(); err != nil {
				return err
			}
		case "mach":
			if err := z.Settings.SetVersionMapUrl("https://machengine.org/zig/index.json"); err != nil {
				log.Info("Run `-vmu default` to reset your version map.")
			}
		default:
			fmt.Printf("It looks like %s isn't installed. Would you like to install it? [y/n]\n", ver)
			if getConfirmation() {
				err = z.Install(ver)
			} else {
				return fmt.Errorf("version %s is not installed", ver)
			}
		}

	} else if errors.Is(err, ErrUnsupportedVersion) {
		log.Debug("User entered VMU alias")

	}

	if err != nil {
		return err
	}

	return z.setBin(ver)
}

func (z *ZVM) setBin(ver string) error {
	version_path := filepath.Join(z.baseDir, ver)
	if err := os.Remove(filepath.Join(z.baseDir, "bin")); err != nil {
		log.Warn(err)
	}

	if err := meta.Symlink(filepath.Join(z.baseDir, ver), filepath.Join(z.baseDir, "bin")); err != nil {
		return err
	}

	if err := meta.Symlink(version_path, filepath.Join(z.baseDir, "bin")); err != nil {
		return err
	}

	return nil
}

func getConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	answer := strings.TrimSpace(strings.ToLower(text))
	return answer == "y" || answer == "ye" || answer == "yes"

}
