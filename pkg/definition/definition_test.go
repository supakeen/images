package definition_test

import (
	"testing"

	"github.com/osbuild/images/pkg/definition"
	"github.com/stretchr/testify/assert"
)

func TestPackageSpecifierValidate(t *testing.T) {
	t.Run("empty package specifier", func(t *testing.T) {
		ps := definition.NewPackageSpecifier("")
		assert.Error(t, ps.Validate())
	})
}

func TestDistributionDefinitionValidate(t *testing.T) {
	t.Run("empty distribution definition", func(t *testing.T) {
	})
}

func TestRegistryGet(t *testing.T) {
	t.Run("empty registry", func(t *testing.T) {
		rt := definition.RegistryName{"foo", "foo", "foo", "foo", "foo", "foo"}
		r := definition.NewRegistry()
		_, err := r.Get(rt)
		assert.Error(t, err)
	})
}

func TestIsValidRegistryName(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		assert.Error(t, definition.IsValidRegistryName(""))
	})

	t.Run("name with too few parts", func(t *testing.T) {
		assert.Error(t, definition.IsValidRegistryName("foo"))
		assert.Error(t, definition.IsValidRegistryName("foo"+definition.REGISTRY_NAME_SEPARATOR+"bar"))
	})

	t.Run("happy", func(t *testing.T) {
		assert.NoError(t, definition.IsValidRegistryName("fedora-40--foo--bar--baz--five--six"))
	})
}
func TestIsValidRegistryNamePart(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		assert.Error(t, definition.IsValidRegistryNamePart(""))
	})

	t.Run("name with registry separator", func(t *testing.T) {
		assert.Error(t, definition.IsValidRegistryNamePart("foo"+definition.REGISTRY_NAME_SEPARATOR+"bar"))
	})

	t.Run("happy", func(t *testing.T) {
		assert.NoError(t, definition.IsValidRegistryNamePart("fedora-40"))
	})
}

func TestDefinition(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		d := definition.NewDefinition()
		d.Distribution = definition.NewDistributionDefinition("foo")
		d.Architecture = definition.NewArchitectureDefinition("foo")
		d.Variant = definition.NewVariantDefinition("foo")
		d.Platform = definition.NewPlatformDefinition("foo")
		d.Workload = definition.NewWorkloadDefinition("foo")
		d.Format = definition.NewFormatDefinition("foo")

		assert.NoError(t, d.Validate())
	})
}
