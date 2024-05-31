package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/kaniuse/kaniuse/data-scraper/pkg/openapi"
	"github.com/kaniuse/kaniuse/data-scraper/pkg/syncmap"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const ScrapMinVersion = "1.5"
const ScrapMaxVersion = "1.30"

type apiLifecycleCmdOptions struct {
	Write      string
	MinVersion string
	MaxVersion string
}

func NewApiLifecycleCommand() (*cobra.Command, error) {
	options := apiLifecycleCmdOptions{}
	apiLifecycle := cobra.Command{
		Use:   "api-lifecycle",
		Short: "scrape Kubernetes API data for api-lifecycle",
		RunE: func(cmd *cobra.Command, args []string) error {
			return apiLifecycleRun(options)
		},
	}
	apiLifecycle.Flags().StringVarP(&options.Write, "write", "w", "", "write to file")
	apiLifecycle.Flags().StringVarP(&options.MinVersion, "min-version", "m", ScrapMinVersion, "min kubernetes version to scrape")
	apiLifecycle.Flags().StringVarP(&options.MaxVersion, "max-version", "M", ScrapMaxVersion, "max kubernetes version to scrape")
	return &apiLifecycle, nil
}

func apiLifecycleRun(options apiLifecycleCmdOptions) error {
	minVersion, err := semver.NewVersion(options.MinVersion)
	if err != nil {
		return errors.Wrap(err, "parse scrape min version")
	}
	maxVersion, err := semver.NewVersion(options.MaxVersion)
	if err != nil {
		return errors.Wrap(err, "parse scrape max version")
	}
	multiVersion, err := fetchAPIs(context.Background(), *minVersion, *maxVersion)
	if err != nil {
		return err
	}
	container := openapi.AsJSONContainer(multiVersion)
	bytes, _ := json.Marshal(container)
	if options.Write != "" {
		err := os.WriteFile(options.Write, bytes, 0644)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(string(bytes))
	}
	return nil
}

func fetchAPIs(ctx context.Context, minVersion semver.Version, maxVersion semver.Version) (openapi.GroupVersionKindAvailability, error) {

	var versions []string
	for version := minVersion; version.LessThan(&maxVersion) || version.Equal(&maxVersion); version = version.IncMinor() {
		versions = append(versions, fmt.Sprintf("%d.%d.%d", version.Major(), version.Minor(), version.Patch()))
	}

	var wg sync.WaitGroup
	wg.Add(len(versions))
	var multiVersion openapi.GroupVersionKindAvailability = syncmap.NewSyncMapWrapper[openapi.GroupVersionKind, *openapi.OrderedKubernetesMinorReleaseAndAPILifeCycleTuple]()
	for _, version := range versions {
		version := version
		go func() {
			defer wg.Done()
			fetcher, err := openapi.NewRepoSwaggerFetcher(version)
			if err != nil {
				panic(err)
			}
			apis, err := fetcher.ListAPI(ctx)
			if err != nil {
				panic(err)
			}
			minorVersion, err := fetcher.KubernetesMinorVersion()
			if err != nil {
				panic(err)
			}
			for gvk, lifecycle := range apis {
				lifecycles, ok := multiVersion.Load(gvk)
				if ok {
					lifecycles.Append(openapi.VersionAPILifeCycle{
						KubernetesMinorRelease: minorVersion,
						APILifecycle:           lifecycle,
					})
				} else {
					lifecycles = &openapi.OrderedKubernetesMinorReleaseAndAPILifeCycleTuple{}
					lifecycles.Append(openapi.VersionAPILifeCycle{
						KubernetesMinorRelease: minorVersion,
						APILifecycle:           lifecycle,
					})
					multiVersion.Store(gvk, lifecycles)
				}
			}
		}()
	}
	wg.Wait()
	return multiVersion, nil
}
