// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common"
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KMSProjectAutokeyConfigIdentity struct {
	parent *KMSProjectAutokeyConfigParent
}

func (i *KMSProjectAutokeyConfigIdentity) String() string {
	return i.parent.String() + "/autokeyConfig"
}

func (r *KMSProjectAutokeyConfigIdentity) Parent() *KMSProjectAutokeyConfigParent {
	return r.parent
}

type KMSProjectAutokeyConfigParent struct {
	ProjectID string
}

func (p *KMSProjectAutokeyConfigParent) String() string {
	return "projects/" + p.ProjectID
}

func NewProjectAutokeyConfigIdentity(ctx context.Context, reader client.Reader, obj *KMSProjectAutokeyConfig) (*KMSProjectAutokeyConfigIdentity, error) {
	// Get Parent
	projectRef, err := refsv1beta1.ResolveProject(ctx, reader, obj.GetNamespace(), obj.Spec.ProjectRef)

	if err != nil {
		return nil, err
	}
	projectID := projectRef.ProjectID
	externalRef := common.ValueOf(obj.Status.ExternalRef)
	if externalRef != "" {
		actualIdentity, err := ParseKMSProjectAutokeyConfigExternal(externalRef)
		if err != nil {
			return nil, err
		}
		if actualIdentity.parent.ProjectID != projectID {
			return nil, fmt.Errorf("spec.projectRef changed, expect %s, got %s", actualIdentity.parent.ProjectID, projectID)
		}
	}

	return &KMSProjectAutokeyConfigIdentity{
		parent: &KMSProjectAutokeyConfigParent{ProjectID: projectID},
	}, nil
}

func ParseKMSProjectAutokeyConfigExternal(external string) (parent *KMSProjectAutokeyConfigIdentity, err error) {
	external = strings.TrimPrefix(external, "/")
	tokens := strings.Split(external, "/")
	if len(tokens) != 3 || tokens[0] != "projects" || tokens[2] != "autokeyConfig" {
		return nil, fmt.Errorf("format of KMSProjectAutokeyConfig external=%q was not known (use projects/<projectID>/autokeyConfig)", external)
	}
	return &KMSProjectAutokeyConfigIdentity{parent: &KMSProjectAutokeyConfigParent{
		ProjectID: tokens[1],
	}}, nil
}
