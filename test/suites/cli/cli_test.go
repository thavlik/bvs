package cli

import (
	"fmt"
	"testing"

	"github.com/thavlik/bvs/test/util"
)

func TestFoo(t *testing.T) {
	w := util.GetWallet(t)
	fmt.Println(w.Addreses)
}
