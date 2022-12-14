package openapi

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/kaniuse/kaniuse/data-scraper/pkg/syncmap"
	"sort"
)

// GroupVersionKindAPIAvailabilityJSONContainer is the container for the final output. Because we would restore as json which
// only support string as the key.
type GroupVersionKindAPIAvailabilityJSONContainer map[GroupVersionStr]map[KindStr][]VersionAPILifeCycle

type GroupVersionKindAvailability syncmap.TypedSyncMap[GroupVersionKind, *OrderedKubernetesMinorReleaseAndAPILifeCycleTuple]

func AsJSONContainer(g GroupVersionKindAvailability) GroupVersionKindAPIAvailabilityJSONContainer {
	result := make(map[GroupVersionStr]map[KindStr][]VersionAPILifeCycle)
	typedMap := syncmap.TypedSyncMap[GroupVersionKind, *OrderedKubernetesMinorReleaseAndAPILifeCycleTuple](g)
	typedMap.Range(func(key GroupVersionKind, value *OrderedKubernetesMinorReleaseAndAPILifeCycleTuple) bool {
		groupVersionStr := key.GroupVersionString()
		kindStr := key.KindString()
		if _, ok := result[groupVersionStr]; !ok {
			result[groupVersionStr] = make(map[KindStr][]VersionAPILifeCycle)
		}
		result[groupVersionStr][kindStr] = value.AsArray()
		return true
	})

	return result
}

type APILifeCycle string

const APILifecycleUnknown APILifeCycle = "unknown"
const APILifecycleStable APILifeCycle = "stable"
const APILifecycleDeprecated APILifeCycle = "deprecated"
const APILifecycleRemoved APILifeCycle = "removed"

type GroupVersionStr string

type KindStr string

// GroupVersionKind represents a certain Kubernetes API Kind.
// For more detail: https://book.kubebuilder.io/cronjob-tutorial/gvks.html
type GroupVersionKind struct {
	// Group represents the API group for a Kubernetes API, e.g. "apps" for the Deployment kind.
	Group string `json:"group"`
	// Version represents the API version for a Kubernetes API, e.g. "v1" for the Deployment kind.
	Version string `json:"version"`
	// Kind represents the API kind for a Kubernetes API, e.g. "Deployment" for the Deployment kind.
	Kind string `json:"kind"`
}

func (g GroupVersionKind) GroupVersionString() GroupVersionStr {
	return GroupVersionStr(fmt.Sprintf("%s/%s", g.Group, g.Version))
}

func (g GroupVersionKind) KindString() KindStr {
	return KindStr(g.Kind)
}
func (g GroupVersionKind) String() string {
	return fmt.Sprintf("%s/%s %s", g.Group, g.Version, g.Kind)
}

// VersionAPILifeCycle represents the availability of a certain API on a certain Kubernetes release.
type VersionAPILifeCycle struct {
	// KubernetesVersion is the version of Kubernetes, without the v prefix, without the patch version, e.g. "1.24", "1.25"
	KubernetesMinorRelease string `json:"kubernetesMinorRelease"`
	// APILifecycle represents the availability.
	APILifecycle APILifeCycle `json:"APILifecycle"`
}

type OrderedKubernetesMinorReleaseAndAPILifeCycleTuple struct {
	base []VersionAPILifeCycle
}

func (o *OrderedKubernetesMinorReleaseAndAPILifeCycleTuple) Append(tuple VersionAPILifeCycle) {
	o.base = append(o.base, tuple)
}

func (o *OrderedKubernetesMinorReleaseAndAPILifeCycleTuple) AsArray() []VersionAPILifeCycle {
	sort.Slice(o.base, func(i, j int) bool {
		return semver.MustParse(o.base[i].KubernetesMinorRelease).LessThan(semver.MustParse(o.base[j].KubernetesMinorRelease))
	})
	return o.base
}
