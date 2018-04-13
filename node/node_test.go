package node

import (
	"testing"
)

func Test_NewNode(t *testing.T) {
	t.Log("Start to init…")

	n, err := NewNode("..\\static\\file.yaml")
	if err != nil {
		t.Errorf("Failed to init,err=%s", err)
	}

	t.Logf("%+v", n)
	t.Log("End Init!!!")
}
