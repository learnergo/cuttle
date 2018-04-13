package config

import (
	"testing"
)

func Test_NewCryptoConfig(t *testing.T) {
	t.Log("Start to init…")

	cConfig, err := NewCryptoConfig("..\\static\\file.yaml")
	if err != nil {
		t.Errorf("Failed to init,err=%s", err)
	}

	t.Logf("%+v", cConfig)
	t.Log("End Init!!!")
}
