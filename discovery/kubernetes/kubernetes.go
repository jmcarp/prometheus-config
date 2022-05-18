// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/common/version"

	"github.com/prometheus/prometheus/discovery"
)

const (
	// metaLabelPrefix is the meta prefix used for all meta labels.
	// in this discovery.
	metaLabelPrefix  = model.MetaLabelPrefix + "kubernetes_"
	namespaceLabel   = metaLabelPrefix + "namespace"
	metricsNamespace = "prometheus_sd_kubernetes"
	presentValue     = model.LabelValue("true")
)

var (
	// Http header
	userAgent = fmt.Sprintf("Prometheus/%s", version.Version)
	// Custom events metric
	eventCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "events_total",
			Help:      "The number of Kubernetes events handled.",
		},
		[]string{"role", "event"},
	)
	// DefaultSDConfig is the default Kubernetes SD configuration
	DefaultSDConfig = SDConfig{
		HTTPClientConfig: config.DefaultHTTPClientConfig,
	}
)

func init() {
	discovery.RegisterConfig(&SDConfig{})
	prometheus.MustRegister(eventCount)
}

// Role is role of the service in Kubernetes.
type Role string

// The valid options for Role.
const (
	RoleNode          Role = "node"
	RolePod           Role = "pod"
	RoleService       Role = "service"
	RoleEndpoint      Role = "endpoints"
	RoleEndpointSlice Role = "endpointslice"
	RoleIngress       Role = "ingress"
)

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *Role) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal((*string)(c)); err != nil {
		return err
	}
	switch *c {
	case RoleNode, RolePod, RoleService, RoleEndpoint, RoleEndpointSlice, RoleIngress:
		return nil
	default:
		return errors.Errorf("unknown Kubernetes SD role %q", *c)
	}
}

// SDConfig is the configuration for Kubernetes service discovery.
type SDConfig struct {
	APIServer          config.URL              `yaml:"api_server,omitempty"`
	Role               Role                    `yaml:"role"`
	KubeConfig         string                  `yaml:"kubeconfig_file"`
	HTTPClientConfig   config.HTTPClientConfig `yaml:",inline"`
	NamespaceDiscovery NamespaceDiscovery      `yaml:"namespaces,omitempty"`
	Selectors          []SelectorConfig        `yaml:"selectors,omitempty"`
	AttachMetadata     AttachMetadataConfig    `yaml:"attach_metadata,omitempty"`
}

// Name returns the name of the Config.
func (*SDConfig) Name() string { return "kubernetes" }

// NewDiscoverer returns a Discoverer for the Config.
func (c *SDConfig) NewDiscoverer(opts discovery.DiscovererOptions) (discovery.Discoverer, error) {
	// We drop the kubernetes discoverer and its constructor to avoid dependency conflicts,
	// so return nil here.
	return nil, nil
}

// SetDirectory joins any relative file paths with dir.
func (c *SDConfig) SetDirectory(dir string) {
	c.HTTPClientConfig.SetDirectory(dir)
	c.KubeConfig = config.JoinDir(dir, c.KubeConfig)
}

type roleSelector struct {
	node          resourceSelector
	pod           resourceSelector
	service       resourceSelector
	endpoints     resourceSelector
	endpointslice resourceSelector
	ingress       resourceSelector
}

type SelectorConfig struct {
	Role  Role   `yaml:"role,omitempty"`
	Label string `yaml:"label,omitempty"`
	Field string `yaml:"field,omitempty"`
}

type resourceSelector struct {
	label string
	field string
}

// AttachMetadataConfig is the configuration for attaching additional metadata
// coming from nodes on which the targets are scheduled.
type AttachMetadataConfig struct {
	Node bool `yaml:"node"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *SDConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultSDConfig
	type plain SDConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if c.Role == "" {
		return errors.Errorf("role missing (one of: pod, service, endpoints, endpointslice, node, ingress)")
	}
	err = c.HTTPClientConfig.Validate()
	if err != nil {
		return err
	}
	if c.APIServer.URL != nil && c.KubeConfig != "" {
		// Api-server and kubeconfig_file are mutually exclusive
		return errors.Errorf("cannot use 'kubeconfig_file' and 'api_server' simultaneously")
	}
	if c.KubeConfig != "" && !reflect.DeepEqual(c.HTTPClientConfig, config.DefaultHTTPClientConfig) {
		// Kubeconfig_file and custom http config are mutually exclusive
		return errors.Errorf("cannot use a custom HTTP client configuration together with 'kubeconfig_file'")
	}
	if c.APIServer.URL == nil && !reflect.DeepEqual(c.HTTPClientConfig, config.DefaultHTTPClientConfig) {
		return errors.Errorf("to use custom HTTP client configuration please provide the 'api_server' URL explicitly")
	}
	if c.APIServer.URL != nil && c.NamespaceDiscovery.IncludeOwnNamespace {
		return errors.Errorf("cannot use 'api_server' and 'namespaces.own_namespace' simultaneously")
	}
	if c.KubeConfig != "" && c.NamespaceDiscovery.IncludeOwnNamespace {
		return errors.Errorf("cannot use 'kubeconfig_file' and 'namespaces.own_namespace' simultaneously")
	}

	foundSelectorRoles := make(map[Role]struct{})
	allowedSelectors := map[Role][]string{
		RolePod:           {string(RolePod)},
		RoleService:       {string(RoleService)},
		RoleEndpointSlice: {string(RolePod), string(RoleService), string(RoleEndpointSlice)},
		RoleEndpoint:      {string(RolePod), string(RoleService), string(RoleEndpoint)},
		RoleNode:          {string(RoleNode)},
		RoleIngress:       {string(RoleIngress)},
	}

	for _, selector := range c.Selectors {
		if _, ok := foundSelectorRoles[selector.Role]; ok {
			return errors.Errorf("duplicated selector role: %s", selector.Role)
		}
		foundSelectorRoles[selector.Role] = struct{}{}

		if _, ok := allowedSelectors[c.Role]; !ok {
			return errors.Errorf("invalid role: %q, expecting one of: pod, service, endpoints, endpointslice, node or ingress", c.Role)
		}
		var allowed bool
		for _, role := range allowedSelectors[c.Role] {
			if role == string(selector.Role) {
				allowed = true
				break
			}
		}

		if !allowed {
			return errors.Errorf("%s role supports only %s selectors", c.Role, strings.Join(allowedSelectors[c.Role], ", "))
		}
	}
	return nil
}

// NamespaceDiscovery is the configuration for discovering
// Kubernetes namespaces.
type NamespaceDiscovery struct {
	IncludeOwnNamespace bool     `yaml:"own_namespace"`
	Names               []string `yaml:"names"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *NamespaceDiscovery) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = NamespaceDiscovery{}
	type plain NamespaceDiscovery
	return unmarshal((*plain)(c))
}
