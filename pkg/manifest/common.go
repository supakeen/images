package manifest

import (
	"github.com/osbuild/images/pkg/disk"
	"github.com/osbuild/images/pkg/osbuild"
)

// filesystemConfigStages generates the org.osbuild.fstab stage,
// org.osbuild.systemd.unit.create stages for .mount and .swap units
// or no stages at all. The last option is used for images that want
// to use `systemd-gpt-auto-generator` and thus need/want no overrides
// in the form of mount units or an fstab.
func filesystemConfigStages(pt *disk.PartitionTable, mountUnits, fstab bool) ([]*osbuild.Stage, error) {
	stages := []*osbuild.Stage{}

	if mountUnits {
		mountStages, err := osbuild.GenSystemdMountStages(pt)
		if err != nil {
			return nil, err
		}
		stages = append(stages, mountStages...)
	}

	if fstab {
		opts, err := osbuild.NewFSTabStageOptions(pt)
		if err != nil {
			return nil, err
		}
		stages = append(stages, osbuild.NewFSTabStage(opts))
	}

	return stages, nil
}
