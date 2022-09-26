package main

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

func main() {

	versions := []string{"1.16", "1.17", "1.18", "1.19", "1.20", "1.21", "1.22", "1.23", "1.24", "1.25"}
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
			minVersion, err := fetcher.KubernetesMinorVersion()
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
	file, _ := os.OpenFile("gvk_api_lifecycle.json", os.O_CREATE|os.O_WRONLY, 0644)
	file.Truncate(0)
	file.Write(bytes)
}
