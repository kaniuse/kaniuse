package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

// CertainVersionAPILifecycleFetcher fetch all available API resources for a certain Kubernetes minor version.
type CertainVersionAPILifecycleFetcher interface {
	KubernetesMinorVersion() (string, error)
	ListAPI(ctx context.Context) (map[GroupVersionKind]APILifecycle, error)
}

var _ error = (*ErrVersionStartWithV)(nil)

type ErrVersionStartWithV struct {
	version string
}

func (e ErrVersionStartWithV) Error() string {
	return fmt.Sprintf("version %s start with v", e.version)
}

var _ CertainVersionAPILifecycleFetcher = (*RepoSwaggerFetcher)(nil)

// RepoSwaggerFetcher would download the swagger.json from Kubernetes GitHub repo, then parse it.
type RepoSwaggerFetcher struct {
	// the target version of Kubernetes, including the patch version, without the "v" prefix, e.g. "1.24.0", "1.25.1"
	version *semver.Version
}

func NewRepoSwaggerFetcher(version string) (*RepoSwaggerFetcher, error) {
	if strings.HasPrefix(version, "v") {
		return nil, errors.WithStack(ErrVersionStartWithV{version: version})
	}

	parsed, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	return &RepoSwaggerFetcher{version: parsed}, nil
}

func (r *RepoSwaggerFetcher) KubernetesMinorVersion() (string, error) {
	return fmt.Sprintf("%d.%d", r.version.Major(), r.version.Minor()), nil
}

func (r *RepoSwaggerFetcher) ListAPI(ctx context.Context) (map[GroupVersionKind]APILifecycle, error) {
	downloader := swaggerJsonDownloader{gitTag: fmt.Sprintf("v%s", r.version.String())}
	content, err := downloader.FetchSwaggerJSONContent(ctx)
	if err != nil {
		return nil, err
	}

	document, err := loads.Analyzed(content, "")
	if err != nil {
		return nil, err
	}

	var result = make(map[GroupVersionKind]APILifecycle)

	for _, item := range document.Analyzer.AllDefinitions() {
		if strings.HasPrefix(item.Name, "io.k8s.api.") {
			gvk, ok := parseGVKFromDefinitionExtension(item.Schema.VendorExtensible.Extensions)
			if !ok {
				continue
			}
			// FIXME: use a logger
			fmt.Printf("%s %s %s\n", r.version.String(),
				gvk.GroupVersionString(), gvk.KindString())
			result[*gvk] = fetchAPILifeCycleFromSchema(item.Schema)
		}
	}

	return result, nil
}

func parseGVKFromDefinitionExtension(extension map[string]interface{}) (*GroupVersionKind, bool) {
	if extension == nil {
		return nil, false
	}
	if gvkAnnotation, ok := extension["x-kubernetes-group-version-kind"]; ok {
		if gvkAnnotationSlice, ok := gvkAnnotation.([]interface{}); ok {
			if len(gvkAnnotationSlice) > 0 {
				gvkAnnotation := gvkAnnotationSlice[0].(map[string]interface{})
				group := gvkAnnotation["group"].(string)
				if group == "" {
					group = "core"
				}
				return &GroupVersionKind{
					Group:   group,
					Version: gvkAnnotation["version"].(string),
					Kind:    gvkAnnotation["kind"].(string),
				}, true
			}
			if len(gvkAnnotationSlice) > 1 {
				panic("more than one x-kubernetes-group-version-kind")
			}
		}
	}
	return nil, false
}

func fetchAPILifeCycleFromSchema(schema *spec.Schema) APILifecycle {
	if schema == nil {
		return APILifecycleUnknown
	}
	if strings.Contains(strings.ToLower(schema.Description), "deprecated") {
		return APILifecycleDeprecated
	}
	return APILifecycleStable
}

type swaggerJsonDownloader struct {
	gitTag string
}

const swaggerJsonURLTemplate = "https://raw.githubusercontent.com/kubernetes/kubernetes/%s/api/openapi-spec/swagger.json"

func (s swaggerJsonDownloader) FetchSwaggerJSONContent(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf(swaggerJsonURLTemplate, s.gitTag)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	return io.ReadAll(response.Body)
}
