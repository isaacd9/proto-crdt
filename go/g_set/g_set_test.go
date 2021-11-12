package g_set

import (
	"log"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/status"
)

func TestGSetAddContains(t *testing.T) {
	a := New()

	Insert(a, &status.Status{Code: 1})
	Insert(a, &status.Status{Code: 1})

	els, _ := Elements(a)
	log.Printf("%+v", els)
}
