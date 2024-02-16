// `definitions` are external definitions used by `images` to define image types and everything necessary
// to put those together. `definitions` are shipped separately from `images` and are a declarative format.

// Definitions are maintained by distributions and consumed by `images` and are meant to be human readable.
// In the Image Builder ecosystem they exist between blueprints which are end-user customizations to a defined
// image and manifests which are a low-level assembly-like that is used by `osbuild` to do the heavy lifting.

// The definitions are consumed by images, any blueprint customizations are applied, and they are then transformed
// into a manifest which can be given to `osbuild` to produce output artifacts.

// The definitions are merged together based on a 6-tuple of values which are applied in order, each of the
// next can override or extend values from the previous.
//
// - distribution
// - architecture
// - variant
// - platform
// - workload
// - format
//
// Each layer of the definition contains:
//
// - packages, the set of packages to be included and excluded
//
// Aside from the definition layers there is a separate build definition which defines, for a given 6-tuple, how to
// put it together.

package definition

import (
	"errors"
	"fmt"
	"strings"

	"github.com/osbuild/images/pkg/rpmmd"
)

const REGISTRY_NAME_SIZE = 6
const REGISTRY_NAME_SEPARATOR = "--"

type PackageSpecifier string

func NewPackageSpecifier(n string) *PackageSpecifier {
	ps := PackageSpecifier(n)
	return &ps
}

func (ps PackageSpecifier) Validate() error {
	if len(ps) == 0 {
		return fmt.Errorf("PackageSpecifier '%s' did not validate due to length", ps)
	}

	return nil
}

type PackageDefinition struct {
	Include []PackageSpecifier `toml:"include"`
	Exclude []PackageSpecifier `toml:"exclude"`
}

func (pd *PackageDefinition) Validate() error {
	errs := []error{}

	for _, inc := range pd.Include {
		errs = append(errs, inc.Validate())
	}

	for _, exc := range pd.Exclude {
		errs = append(errs, exc.Validate())
	}

	return errors.Join(errs...)
}

type RepositorySpecifier struct{}

func (rs *RepositorySpecifier) Validate() error {
	return nil
}

type RepositoryDefinition struct {
	Repositories []RepositorySpecifier
}

func (rd *RepositoryDefinition) Validate() error {
	errs := []error{}

	for _, rep := range rd.Repositories {
		errs = append(errs, rep.Validate())
	}

	return errors.Join(errs...)
}

// A distribution is put together from a PackageDefinition (its contents) and a RepositoryDefinition (where to
// get those packages from).
type DistributionDefinition struct {
	Name       string
	Package    PackageDefinition    `toml:"packages"`
	Repository RepositoryDefinition `toml:"repositories"`
}

func NewDistributionDefinition(n string) *DistributionDefinition {
	return &DistributionDefinition{Name: n}
}

func (dd *DistributionDefinition) Validate() error {
	return errors.Join(
		IsValidRegistryNamePart(dd.Name),
		dd.Package.Validate(),
		dd.Repository.Validate(),
	)
}

type ArchitectureDefinition struct {
	Name    string
	Package PackageDefinition
}

func NewArchitectureDefinition(n string) *ArchitectureDefinition {
	return &ArchitectureDefinition{Name: n}
}

func (ad *ArchitectureDefinition) Validate() error {
	return errors.Join(
		IsValidRegistryNamePart(ad.Name),
		ad.Package.Validate(),
	)
}

type VariantDefinition struct {
	Name    string
	Package PackageDefinition
}

func NewVariantDefinition(n string) *VariantDefinition {
	return &VariantDefinition{Name: n}
}

func (vd *VariantDefinition) Validate() error {
	return errors.Join(
		IsValidRegistryNamePart(vd.Name),
		vd.Package.Validate(),
	)
}

type PlatformDefinition struct {
	Name    string
	Package PackageDefinition
}

func NewPlatformDefinition(n string) *PlatformDefinition {
	return &PlatformDefinition{Name: n}
}

func (pd *PlatformDefinition) Validate() error {
	return errors.Join(
		IsValidRegistryNamePart(pd.Name),
		pd.Package.Validate(),
	)
}

type WorkloadDefinition struct {
	Name    string
	Package PackageDefinition
}

func NewWorkloadDefinition(n string) *WorkloadDefinition {
	return &WorkloadDefinition{Name: n}
}

func (wd *WorkloadDefinition) Validate() error {
	return errors.Join(
		IsValidRegistryNamePart(wd.Name),
		wd.Package.Validate(),
	)
}

type FormatDefinition struct {
	Name    string
	Package PackageDefinition
}

func NewFormatDefinition(n string) *FormatDefinition {
	return &FormatDefinition{Name: n}
}

func (fd *FormatDefinition) Validate() error {
	return errors.Join(
		IsValidRegistryNamePart(fd.Name),
		fd.Package.Validate(),
	)
}

type Definition struct {
	Distribution *DistributionDefinition
	Architecture *ArchitectureDefinition
	Variant      *VariantDefinition
	Platform     *PlatformDefinition
	Workload     *WorkloadDefinition
	Format       *FormatDefinition
}

func NewDefinition() *Definition {
	return &Definition{}
}

func (d *Definition) Validate() error {
	return errors.Join(
		d.Distribution.Validate(),
		d.Architecture.Validate(),
		d.Variant.Validate(),
		d.Platform.Validate(),
		d.Workload.Validate(),
		d.Format.Validate(),
	)
}

// Take a definition and turn it into its RegistryName.
func (d *Definition) GetRegistryName() RegistryName {
	return RegistryName{
		d.Distribution.Name,
		d.Architecture.Name,
		d.Variant.Name,
		d.Platform.Name,
		d.Workload.Name,
		d.Format.Name,
	}
}

func (d *Definition) String() string {
	s := ""
	r := d.GetRegistryName()

	for i, p := range r {
		s = s + p

		if i != len(r)-1 {
			s += REGISTRY_NAME_SEPARATOR
		}
	}

	return s
}

// Flattens a definition into a resolved object called an Image which is used internally by `images`.
func (d *Definition) Resolve() (*Image, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}

	return nil, nil
}

type Image struct {
	// The definition that this Image is resolved from.
	Definition Definition

	Packages     rpmmd.PackageSet
	Repositories rpmmd.RepoConfig
}

// A 6-tuple describing an image to get from the registry.
type RegistryName [REGISTRY_NAME_SIZE]string

// Use this to read in a directory structure of definitions and get available definitions as Image's.
type Registry struct {
	Path string
	Map  map[RegistryName]*Image
}

func NewRegistry() Registry {
	return Registry{}
}

func (r *Registry) Add(it RegistryName, i *Image) error {
	if _, err := r.Get(it); err == nil {
		return fmt.Errorf("registry already contains %v", it)
	}

	r.Map[it] = i

	return nil
}

// Get an image by its RegistryName
func (r *Registry) Get(it RegistryName) (*Image, error) {
	v, ok := r.Map[it]

	if !ok {
		return nil, fmt.Errorf("could not find %v in registry", it)
	}

	return v, nil
}

// Validate a name, it needs to contain all parts and all parts need to be valid
func IsValidRegistryName(n string) error {
	if len(n) == 0 {
		return fmt.Errorf("empty name")
	}

	ps := strings.Split(n, REGISTRY_NAME_SEPARATOR)

	if len(ps) != REGISTRY_NAME_SIZE {
		return fmt.Errorf("name not of correct length")
	}

	errs := []error{}

	for _, p := range ps {
		errs = append(errs, IsValidRegistryNamePart(p))
	}

	return errors.Join(errs...)
}

// Validate a name part, it is not allowed to be empty nor contain the registry tuple separator.
func IsValidRegistryNamePart(n string) error {
	if len(n) == 0 {
		return fmt.Errorf("empty name")
	}

	if strings.Contains(n, REGISTRY_NAME_SEPARATOR) {
		return fmt.Errorf("name contains the registry tuple separator")
	}

	return nil
}
