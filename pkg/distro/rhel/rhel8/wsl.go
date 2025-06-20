package rhel8

import (
	"github.com/osbuild/images/internal/common"
	"github.com/osbuild/images/pkg/distro"
	"github.com/osbuild/images/pkg/distro/rhel"
)

func mkWslImgType() *rhel.ImageType {
	it := rhel.NewImageType(
		"wsl",
		"image.wsl",
		"application/x-tar",
		packageSetLoader,
		rhel.TarImage,
		[]string{"build"},
		[]string{"os", "archive"},
		[]string{"archive"},
	)

	it.Compression = "xz"
	it.DefaultImageConfig = &distro.ImageConfig{
		Locale:    common.ToPtr("en_US.UTF-8"),
		NoSElinux: common.ToPtr(true),
		WSLConfig: &distro.WSLConfig{
			BootSystemd: true,
		},
	}

	return it
}
