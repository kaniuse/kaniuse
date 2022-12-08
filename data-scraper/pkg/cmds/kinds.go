package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/kaniuse/kaniuse/data-scraper/pkg/openapi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"sort"
)

type kindsCmdOptions struct {
	Write      string
	MinVersion string
	MaxVersion string
}

func NewKindsCmd() (*cobra.Command, error) {
	options := kindsCmdOptions{}
	kindsCmd := cobra.Command{
		Use:   "kinds",
		Short: "scrape Kubernetes API data for kinds",
		RunE: func(cmd *cobra.Command, args []string) error {
			return kindsRun(options)
		},
	}
	kindsCmd.Flags().StringVarP(&options.Write, "write", "w", "", "write to file")
	kindsCmd.Flags().StringVarP(&options.MinVersion, "min-version", "m", ScrapMinVersion, "min kubernetes version to scrape")
	kindsCmd.Flags().StringVarP(&options.MaxVersion, "max-version", "M", ScrapMaxVersion, "max kubernetes version to scrape")
	return &kindsCmd, nil
}

func kindsRun(options kindsCmdOptions) error {
	minVersion, err := semver.NewVersion(options.MinVersion)
	if err != nil {
		return errors.Wrap(err, "parse scrape min version")
	}
	maxVersion, err := semver.NewVersion(options.MaxVersion)
	if err != nil {
		return errors.Wrap(err, "parse scrape max version")
	}
	apis, err := fetchAPIs(context.Background(), *minVersion, *maxVersion)
	if err != nil {
		return err
	}
	type Kind string
	type GVK string
	type GVKAndMinorVersionAndAPILifecycles struct {
		DisplayGVK string
		GVK        openapi.GroupVersionKind
		Lifecycles []openapi.KubernetesMinorReleaseAndAPILifeCycleTuple
	}
	result := make(map[string][]GVKAndMinorVersionAndAPILifecycles)
	apis.Range(
		func(key openapi.GroupVersionKind, value *openapi.OrderedKubernetesMinorReleaseAndAPILifeCycleTuple) bool {

			entry := GVKAndMinorVersionAndAPILifecycles{
				DisplayGVK: key.String(),
				GVK:        key,
				Lifecycles: value.AsArray(),
			}
			result[key.Kind] = append(result[key.Kind], entry)
			return true
		})

	// sort result with the supported version
	for _, lifecycles := range result {
		sort.Slice(lifecycles, func(i, j int) bool {
			itemI := lifecycles[i]
			itemJ := lifecycles[j]
			// order with group, version conventions
			if itemI.GVK.Group != itemJ.GVK.Group {
				ci := containsInLegacyGroups(itemI.GVK.Group)
				cj := containsInLegacyGroups(itemJ.GVK.Group)
				if ci && !cj {
					return true
				}
				if !ci && cj {
					return false
				}
			}
			apiVersionI, err := ParseAPIVersion(itemI.GVK.Version)
			if err != nil {
				panic(err)
			}
			apiVersionJ, err := ParseAPIVersion(itemJ.GVK.Version)
			if err != nil {
				panic(err)
			}
			return apiVersionI.LessThan(*apiVersionJ)
		})
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
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

func filterStables(origin []openapi.KubernetesMinorReleaseAndAPILifeCycleTuple) []openapi.KubernetesMinorReleaseAndAPILifeCycleTuple {
	result := make([]openapi.KubernetesMinorReleaseAndAPILifeCycleTuple, 0)
	for _, lifecycle := range origin {
		if lifecycle.APILifecycle == openapi.APILifecycleStable {
			result = append(result, lifecycle)
		}
	}
	return result
}

func filterDeprecated(origin []openapi.KubernetesMinorReleaseAndAPILifeCycleTuple) []openapi.KubernetesMinorReleaseAndAPILifeCycleTuple {
	result := make([]openapi.KubernetesMinorReleaseAndAPILifeCycleTuple, 0)
	for _, lifecycle := range origin {
		if lifecycle.APILifecycle == openapi.APILifecycleDeprecated {
			result = append(result, lifecycle)
		}
	}
	return result
}

var legacyGroups = []string{"extensions", "core"}

func containsInLegacyGroups(group string) bool {
	for _, legacyGroup := range legacyGroups {
		if group == legacyGroup {
			return true
		}
	}
	return false
}
