// khan
// https://github.com/jpholanda/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package cmd_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/jpholanda/khan/cmd"
)

func dropDB() error {
	cmd := exec.Cmd{
		Dir:  "../",
		Path: "/usr/bin/make",
		Args: []string{
			"drop-test",
		},
	}
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

var _ = Describe("Migrate Command", func() {
	BeforeEach(func() {
		err := dropDB()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Migrate Cmd", func() {
		It("Should run migrations up", func() {
			ConfigFile = "../config/test.yaml"
			InitConfig()
			err := RunMigrations(-1)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
