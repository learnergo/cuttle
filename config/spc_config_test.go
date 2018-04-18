package config

import (
	"testing"
)

func Test_NewSpeConfig(t *testing.T) {
	t.Log("Start to init…")

	speConfig, err := NewSpeConfig("../static/cuttle.yaml")
	if err != nil {
		t.Errorf("Failed to init,err=%s", err)
	}

	t.Logf("%+v", speConfig)
	t.Log("End Init!!!")

	err = speConfig.Marshal("../static/cuttle1.yaml")
	if err != nil {
		t.Errorf("Failed to Marshal,err=%s", err)
	}
}
