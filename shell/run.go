package shell

import (
	"fmt"
	"github.com/BakerHub/trivial/check"
	"os/exec"
	"strings"
)

func Run(name string, args ...string) {
	cmd := exec.Command(name, args...) // #nosec
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("%s %s", name, strings.Join(args, " "))
		fmt.Printf("%s\n", out)
		check.Check(err)
	}
}
