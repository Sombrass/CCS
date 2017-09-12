package shadowbins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeConfig(t *testing.T) {
	var sampleConfig = `
shadow_identifier = "vagrant"

[hyper]
port = 5348

[hyper.binaries]
code = "/usr/bin/code"
subl = "/usr/bin/subl"

[shadows.wsl]
    auto_alias = ['subl', 'code']

    [[shadows.wsl.dir_map]]
    hyper_path = "/mnt/c/"
    shadow_path = "C:\\"

    [[shadows.wsl.dir_map]]
    hyper_path = "/mnt/d/"
    shadow_path = "D:\\"


[shadows.vagrant]
    auto_alias = ['subl', 'code']

    [[shadows.vagrant.dir_map]]
    hyper_path = "/mnt/vagrant/"
    shadow_path = "~/vagrant/"

`
	conf, err := decodeConfig(sampleConfig)
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.Equal(conf.ShadowIdentifier, "vagrant")
	assert.NotEmpty(conf.Hyper)
	assert.NotEmpty(conf.Shadows)

	hyper := conf.Hyper
	assert.Equal(hyper.Port, 5348)
	assert.Contains(hyper.Binaries, "code")

	shadow := conf.Shadows["wsl"]
	assert.Contains(shadow.AutoAlias, "subl")
	assert.Equal(shadow.DirMap[0]["hyper_path"], "/mnt/c/")
	assert.Equal(shadow.DirMap[1]["shadow_path"], "D:\\")
}
