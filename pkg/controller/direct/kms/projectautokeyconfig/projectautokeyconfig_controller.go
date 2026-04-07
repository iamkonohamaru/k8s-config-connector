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

package projectautokeyconfig

import (
	"context"
	"fmt"

	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/kms/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/config"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/directbase"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/registry"

	gcp "cloud.google.com/go/kms/apiv1"

	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

const (
	ctrlName = "kms-projectautokeyconfig-controller"
)

func init() {
	registry.RegisterModel(krm.KMSProjectAutokeyConfigGVK, NewModel)
}

func NewModel(ctx context.Context, config *config.ControllerConfig) (directbase.Model, error) {
	return &model{config: *config}, nil
}

var _ directbase.Model = &model{}

type model struct {
	config config.ControllerConfig
}

func (m *model) client(ctx context.Context) (*gcp.AutokeyAdminClient, error) {
	var opts []option.ClientOption
	opts, err := m.config.RESTClientOptions()
	if err != nil {
		return nil, err
	}
	gcpClient, err := gcp.NewAutokeyAdminRESTClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("building AutokeyConfig client: %w", err)
	}
	return gcpClient, err
}

func (m *model) AdapterForObject(ctx context.Context, op *directbase.AdapterForObjectOperation) (directbase.Adapter, error) {
	u := op.GetUnstructured()
	reader := op.Reader
	obj := &krm.KMSProjectAutokeyConfig{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &obj); err != nil {
		return nil, fmt.Errorf("error converting to %T: %w", obj, err)
	}

	id, err := krm.NewProjectAutokeyConfigIdentity(ctx, reader, obj)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve project for autokeyConfig name: %s, err: %w", obj.GetName(), err)
	}

	gcpClient, err := m.client(ctx)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		id:        id,
		gcpClient: gcpClient,
		desired:   obj,
	}, nil
}

func (m *model) AdapterForURL(ctx context.Context, url string) (directbase.Adapter, error) {
	// TODO: Support URLs
	return nil, nil
}

type Adapter struct {
	id        *krm.KMSProjectAutokeyConfigIdentity
	gcpClient *gcp.AutokeyAdminClient
	desired   *krm.KMSProjectAutokeyConfig
	actual    *kmspb.AutokeyConfig
}

var _ directbase.Adapter = &Adapter{}

func (a *Adapter) Find(ctx context.Context) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("getting KMSProjectAutokeyConfig", "name", a.id)

	req := &kmspb.GetAutokeyConfigRequest{Name: a.id.String()}
	autokeyconfigpb, err := a.gcpClient.GetAutokeyConfig(ctx, req)
	if err != nil {
		return false, fmt.Errorf("getting KMSProjectAutokeyConfig %q: %w", a.id, err)
	}

	a.actual = autokeyconfigpb

	mapCtx := &direct.MapContext{}
	observedState := KMSProjectAutokeyConfigObservedState_FromProto(mapCtx, autokeyconfigpb)
	if mapCtx.Err() != nil {
		return false, mapCtx.Err()
	}
	a.desired.Status.ObservedState = observedState

	return true, nil
}

func (a *Adapter) Create(ctx context.Context, createOp *directbase.CreateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("Create operation not supported for ProjectAutokeyConfig resource.")
	return fmt.Errorf("Create operation not supported for ProjectAutokeyConfig resource")
}

func (a *Adapter) Update(ctx context.Context, updateOp *directbase.UpdateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("updating ProjectAutokeyConfig", "name", a.id)
	mapCtx := &direct.MapContext{}

	resource := KMSProjectAutokeyConfig_FromFields(mapCtx, a.id, a.desired.Spec.KeyProjectResolutionMode)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}

	var updateMask []string
	if a.actual.GetKeyProjectResolutionMode() != resource.GetKeyProjectResolutionMode() {
		updateMask = append(updateMask, "key_project_resolution_mode")
	}

	if len(updateMask) == 0 {
		return nil
	}

	req := &kmspb.UpdateAutokeyConfigRequest{
		AutokeyConfig: resource,
		UpdateMask:    &fieldmaskpb.FieldMask{Paths: updateMask},
	}

	op, err := a.gcpClient.UpdateAutokeyConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("updating ProjectAutokeyConfig %q: %w", a.id, err)
	}

	observedState := KMSProjectAutokeyConfigObservedState_FromProto(mapCtx, op)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}
	a.desired.Status.ObservedState = observedState

	return nil
}

func (a *Adapter) Export(ctx context.Context) (*unstructured.Unstructured, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (a *Adapter) Delete(ctx context.Context, deleteOp *directbase.DeleteOperation) (bool, error) {
	return false, fmt.Errorf("not implemented yet")
}
