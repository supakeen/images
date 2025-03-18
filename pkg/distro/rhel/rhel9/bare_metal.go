package rhel9

import (
	"github.com/osbuild/images/internal/common"
	"github.com/osbuild/images/pkg/distro"
	"github.com/osbuild/images/pkg/distro/packagesets"
	"github.com/osbuild/images/pkg/distro/rhel"
	"github.com/osbuild/images/pkg/rpmmd"
)

func mkTarImgType() *rhel.ImageType {
	return rhel.NewImageType(
		"tar",
		"root.tar.xz",
		"application/x-tar",
		map[string]rhel.PackageSetFunc{
			rhel.OSPkgsKey: packageSetLoader,
		},
		rhel.TarImage,
		[]string{"build"},
		[]string{"os", "archive"},
		[]string{"archive"},
	)
}

func mkImageInstallerImgType() *rhel.ImageType {
	it := rhel.NewImageType(
		"image-installer",
		"installer.iso",
		"application/x-iso9660-image",
		map[string]rhel.PackageSetFunc{
			rhel.OSPkgsKey: func(t *rhel.ImageType) rpmmd.PackageSet {
				return common.Must(packagesets.Load(t, "bare-metal", nil))
			},
			rhel.InstallerPkgsKey: packageSetLoader,
		},
		rhel.ImageInstallerImage,
		[]string{"build"},
		[]string{"anaconda-tree", "rootfs-image", "efiboot-tree", "os", "bootiso-tree", "bootiso"},
		[]string{"bootiso"},
	)

	it.BootISO = true
	it.Bootable = true
	it.ISOLabelFn = distroISOLabelFunc

	it.DefaultInstallerConfig = &distro.InstallerConfig{
		AdditionalDracutModules: []string{
			"nvdimm", // non-volatile DIMM firmware (provides nfit, cuse, and nd_e820)
			"prefixdevname",
			"prefixdevname-tools",
			"ifcfg",
		},
		AdditionalDrivers: []string{
			"cuse",
			"ipmi_devintf",
			"ipmi_msghandler",
		},
	}

	it.DefaultImageConfig = &distro.ImageConfig{
		Locale: common.ToPtr("C.UTF-8"),
	}

	return it
}
