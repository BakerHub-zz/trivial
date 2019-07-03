package check

import "log"

func Check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
