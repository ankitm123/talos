// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package runtime

import (
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/resource/typed"

	"github.com/talos-systems/talos/pkg/machinery/proto"
)

// PlatformMetadataType is type of Metadata resource.
const PlatformMetadataType = resource.Type("PlatformMetadatas.talos.dev")

// PlatformMetadataID is the ID for Metadata resource platform.
const PlatformMetadataID resource.ID = "platformmetadata"

// PlatformMetadata resource holds.
type PlatformMetadata = typed.Resource[PlatformMetadataSpec, PlatformMetadataRD]

// PlatformMetadataSpec describes platform metadata properties.
//
//gotagsrewrite:gen
type PlatformMetadataSpec struct {
	Platform     string `yaml:"platform,omitempty" protobuf:"1"`
	Hostname     string `yaml:"hostname,omitempty" protobuf:"2"`
	Region       string `yaml:"region,omitempty" protobuf:"3"`
	Zone         string `yaml:"zone,omitempty" protobuf:"4"`
	InstanceType string `yaml:"instanceType,omitempty" protobuf:"5"`
	InstanceID   string `yaml:"instanceId,omitempty" protobuf:"6"`
	ProviderID   string `yaml:"providerId,omitempty" protobuf:"7"`
}

// NewPlatformMetadataSpec initializes a MetadataSpec resource.
func NewPlatformMetadataSpec(namespace resource.Namespace, id resource.ID) *PlatformMetadata {
	return typed.NewResource[PlatformMetadataSpec, PlatformMetadataRD](
		resource.NewMetadata(namespace, PlatformMetadataType, PlatformMetadataID, resource.VersionUndefined),
		PlatformMetadataSpec{},
	)
}

// PlatformMetadataRD provides auxiliary methods for PlatformMetadata.
type PlatformMetadataRD struct{}

// ResourceDefinition implements typed.ResourceDefinition interface.
func (PlatformMetadataRD) ResourceDefinition(resource.Metadata, PlatformMetadataSpec) meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             PlatformMetadataType,
		DefaultNamespace: NamespaceName,
		PrintColumns:     []meta.PrintColumn{},
	}
}

func init() {
	proto.RegisterDefaultTypes()

	err := protobuf.RegisterDynamic[PlatformMetadataSpec](PlatformMetadataType, &PlatformMetadata{})
	if err != nil {
		panic(err)
	}
}
