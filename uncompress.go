/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : uncompress.go

* Purpose :

* Creation Date : 09-05-2017

* Last Modified : Fri 17 Nov 2017 11:27:36 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func uncompress(dir, file, pass string) error {
	ext := filepath.Ext(file)
	name := file[:len(file)-len(ext)]

	var putpass, cmd string
	cmd = fmt.Sprintf(`cd "%s" && `, dir)

	switch ext {
	case ".zip":
		if len(pass) > 0 {
			putpass = fmt.Sprint(`-P '%s' `, pass)
		}
		cmd += fmt.Sprintf(`yes|unzip %s"%s" -d "%s"`, putpass, file, name)
	case ".rar":
		if len(pass) > 0 {
			putpass = fmt.Sprintf(`-p%s `, pass)
		}
		if strings.Contains(name, "part1") {
			cmd = fmt.Sprintf(`unrar %sx "%s"`, putpass, file)
		} else {
			cmd += fmt.Sprintf(`mkdir "%s" && yes|unrar %sx "%s" "%s"`, name, putpass, file, name)
		}
	default:
		log.Println("not support")
	}
	log.Println(cmd)

	c := exec.Command("/bin/bash", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}
