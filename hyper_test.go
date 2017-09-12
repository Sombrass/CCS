package shadowbins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConverBin(t *testing.T) {
	confStr := `
[hyper]
port = 5348

[hyper.binaries]
code = "/usr/bin/code"
subl = "/usr/bin/subl"
`
	conf, _ := decodeConfig(confStr)
	launcher := new(Launcher)
	launcher.conf = &conf

	assert := assert.New(t)

	bin, err := launcher.convertBin("code")
	assert.Equal(bin, "/usr/bin/code")

	bin, err = launcher.convertBin("subl")
	assert.Equal(bin, "/usr/bin/subl")

	bin, err = launcher.convertBin("notepad++")
	assert.NotNil(err)
}

func TestConverArgs(t *testing.T) {
	confStr := `
[shadows.wsl]
    auto_alias = ['subl', 'code']

    [[shadows.wsl.dir_map]]
    hyper_path = "C:\\"
    shadow_path = "/mnt/c/"

    [[shadows.wsl.dir_map]]
    hyper_path = "D:\\"
    shadow_path = "/mnt/d/"

[shadows.vagrant]
    auto_alias = ['subl', 'code']

    [[shadows.vagrant.dir_map]]
    hyper_path = "/mnt/vagrant/"
    shadow_path = "~/vagrant/"
`
	conf, _ := decodeConfig(confStr)
	launcher := new(Launcher)
	launcher.conf = &conf

	assert := assert.New(t)

	args, _ := launcher.convertArgs("vagrant", []string{"ls", "-la"})
	assert.Equal(args, []string{"ls", "-la"})

	args, _ = launcher.convertArgs("vagrant", []string{"~/vagrant/src/foo", "--debug"})
	assert.Equal(args, []string{"/mnt/vagrant/src/foo", "--debug"})

	args, _ = launcher.convertArgs("vagrant", []string{"--src", "~/vagrant/bar.go", "--dst", "~/vagrant/foo.go"})
	assert.Equal(args, []string{"--src", "/mnt/vagrant/bar.go", "--dst", "/mnt/vagrant/foo.go"})

	args, _ = launcher.convertArgs("wsl", []string{"/mnt/c/bar/test.go"})
	assert.Equal(args, []string{"C:\\bar/test.go"})

	args, _ = launcher.convertArgs("wsl", []string{"/mnt/d/bar"})
	assert.Equal(args, []string{"D:\\bar"})

}
