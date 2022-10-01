package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Masterminds/semver"
)

func main() {
	var versions []string

	for version := semver.MustParse("1.5"); version.Minor() <= 25; {
		versions = append(versions, fmt.Sprintf("%d.%d.%d", version.Major(), version.Minor(), version.Patch()))
		tmp := version.IncMinor()
		version = &tmp
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(len(versions))
	var multiVersion groupVersionKindAvailability = newSyncMapWrapper[GroupVersionKind, *OrderedKubernetesMinorReleaseAndAPILifeCycleTuple]()
	for _, version := range versions {
		version := version
		go func() {
			defer wg.Done()
			fetcher, err := NewRepoSwaggerFetcher(version)
			if err != nil {
				panic(err)
			}
			apis, err := fetcher.ListAPI(ctx)
			if err != nil {
				panic(err)
			}
			minVersion, err := fetcher.KubernetesMinorVersion()
			if err != nil {
				panic(err)
			}
			for gvk, lifecycle := range apis {
				lifecycles, ok := multiVersion.Load(gvk)
				if ok {
					lifecycles.Append(KubernetesMinorReleaseAndAPILifeCycleTuple{
						KubernetesMinorRelease: minVersion,
						APILifecycle:           lifecycle,
					})
				} else {
					lifecycles = &OrderedKubernetesMinorReleaseAndAPILifeCycleTuple{}
					lifecycles.Append(KubernetesMinorReleaseAndAPILifeCycleTuple{
						KubernetesMinorRelease: minVersion,
						APILifecycle:           lifecycle,
					})
					multiVersion.Store(gvk, lifecycles)
				}
			}
		}()
	}
	wg.Wait()
	container := AsJSONContainer(multiVersion)
	bytes, _ := json.Marshal(container)
	fmt.Print(string(bytes))
}
