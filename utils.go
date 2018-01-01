/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : utils.go

* Purpose :

* Creation Date : 08-30-2017

* Last Modified : Mon 01 Jan 2018 12:05:42 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"crypto/md5"
	"encoding/hex"
)

func hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	s := hex.EncodeToString(hasher.Sum(nil))
	return s[:8]
}
