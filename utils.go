/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : utils.go

* Purpose :

* Creation Date : 08-30-2017

* Last Modified : Wed 30 Aug 2017 02:36:34 AM UTC

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
	return s[:6]
}
