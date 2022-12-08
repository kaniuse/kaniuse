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
	"sync"
)

func NewFieldsCmd() (*cobra.Command, error) {
	options := fieldsCmdOptions{}
	fieldsCmd := cobra.Command{
		Use:   "fields",
		Short: "scrape Kubernetes API data for fields",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fieldsRun(options)
		}}

	fieldsCmd.Flags().StringVarP(&options.Write, "write", "w", "", "write to file")
	fieldsCmd.Flags().StringVarP(&options.MinVersion, "min-version", "m", ScrapMinVersion, "min kubernetes version to scrape")
	fieldsCmd.Flags().StringVarP(&options.MaxVersion, "max-version", "M", ScrapMaxVersion, "max kubernetes version to scrape")
	return &fieldsCmd, nil
}

type fieldsCmdOptions struct {
	Write      string
	MinVersion string
	MaxVersion string
}

func fieldsRun(options fieldsCmdOptions) error {

	minVersion, err := semver.NewVersion(options.MinVersion)
	if err != nil {
		return errors.Wrap(err, "parse scrape min version")
	}
	maxVersion, err := semver.NewVersion(options.MaxVersion)
	if err != nil {
		return errors.Wrap(err, "parse scrape max version")
	}
	type VersionTypeAPILifeCycle struct {
		// KubernetesVersion is the version of Kubernetes, without the v prefix, without the patch version, e.g. "1.24", "1.25"
		KubernetesMinorRelease string
		// APILifecycle represents the availability.
		APILifecycle openapi.APILifeCycle
		FieldType    string
	}

	ctx := context.Background()
	type FieldSummary struct {
		// the field path, e.g. "io.k8s.api.core.v1.Namespace.spec"
		FieldPath  string
		LifeCycles []VersionTypeAPILifeCycle
	}
	type FieldGVKSummary struct {
		GVK            openapi.GroupVersionKind
		KindLifeCycles []openapi.VersionAPILifeCycle
		// key is the field path
		Summary map[string]FieldSummary
	}
	rwlock := sync.RWMutex{}
	// key is GVK string
	result := make(map[string]FieldGVKSummary)

	var versions []string
	for version := *minVersion; version.LessThan(maxVersion) || version.Equal(maxVersion); version = version.IncMinor() {
		versions = append(versions, fmt.Sprintf("%d.%d.%d", version.Major(), version.Minor(), version.Patch()))
	}

	var wg sync.WaitGroup
	wg.Add(len(versions))
	for _, version := range versions {
		version := version
		go func() {
			defer wg.Done()
			fetcher, err := openapi.NewRepoSwaggerFetcher(version)
			if err != nil {
				panic(err)
			}
			gvks, err := fetcher.ListAPI(ctx)
			if err != nil {
				panic(err)
			}
			fields, err := fetcher.ListFields(ctx)
			if err != nil {
				panic(err)
			}
			for gvk, kindLifeCycle := range gvks {
				fieldLifeCycle := fields[gvk]
				func() {
					rwlock.Lock()
					defer rwlock.Unlock()
					if _, ok := result[gvk.String()]; !ok {
						result[gvk.String()] = FieldGVKSummary{
							GVK:            gvk,
							KindLifeCycles: []openapi.VersionAPILifeCycle{},
							Summary:        make(map[string]FieldSummary),
						}
					}
					gvkSummary := result[gvk.String()]
					gvkSummary.KindLifeCycles = append(result[gvk.String()].KindLifeCycles, openapi.VersionAPILifeCycle{
						KubernetesMinorRelease: version,
						APILifecycle:           kindLifeCycle,
					})
					result[gvk.String()] = gvkSummary
					for _, flattenField := range fieldLifeCycle {
						if _, ok := result[gvk.String()].Summary[flattenField.FieldPath]; !ok {
							gvkSummary.Summary[flattenField.FieldPath] = FieldSummary{
								FieldPath:  flattenField.FieldPath,
								LifeCycles: []VersionTypeAPILifeCycle{},
							}
						}

						fieldSummary := gvkSummary.Summary[flattenField.FieldPath]
						fieldSummary.LifeCycles =
							append(
								result[gvk.String()].Summary[flattenField.FieldPath].LifeCycles,
								VersionTypeAPILifeCycle{
									FieldType:              flattenField.FieldType,
									KubernetesMinorRelease: version,
									APILifecycle:           flattenField.Lifecycle,
								})
						gvkSummary.Summary[flattenField.FieldPath] = fieldSummary
					}
				}()
			}
		}()
	}
	wg.Wait()
	for _, gvkSummary := range result {
		sort.Slice(gvkSummary.KindLifeCycles, func(i, j int) bool {
			return semVerLessThan(gvkSummary.KindLifeCycles[i].KubernetesMinorRelease, gvkSummary.KindLifeCycles[j].KubernetesMinorRelease)
		})
		for _, fieldSummary := range gvkSummary.Summary {
			sort.Slice(fieldSummary.LifeCycles, func(i, j int) bool {
				return semVerLessThan(fieldSummary.LifeCycles[i].KubernetesMinorRelease, fieldSummary.LifeCycles[j].KubernetesMinorRelease)
			})
		}
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "marshal result")
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
