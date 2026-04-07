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
	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/kms/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct"

	pb "cloud.google.com/go/kms/apiv1/kmspb"
)

func KMSProjectAutokeyConfigSpec_FromProto(mapCtx *direct.MapContext, in *pb.AutokeyConfig) *krm.KMSProjectAutokeyConfigSpec {
	if in == nil {
		return nil
	}
	out := &krm.KMSProjectAutokeyConfigSpec{}
	if in.GetKeyProjectResolutionMode() != pb.AutokeyConfig_KEY_PROJECT_RESOLUTION_MODE_UNSPECIFIED {
		out.KeyProjectResolutionMode = direct.Enum_FromProto(mapCtx, in.GetKeyProjectResolutionMode())
	}
	return out
}

func KMSProjectAutokeyConfig_FromFields(mapCtx *direct.MapContext, id *krm.KMSProjectAutokeyConfigIdentity, keyProjectResolutionMode *string) *pb.AutokeyConfig {
	out := &pb.AutokeyConfig{}
	out.Name = id.String()
	if keyProjectResolutionMode != nil {
		out.KeyProjectResolutionMode = direct.Enum_ToProto[pb.AutokeyConfig_KeyProjectResolutionMode](mapCtx, keyProjectResolutionMode)
	}
	return out
}

func KMSProjectAutokeyConfigObservedState_FromProto(mapCtx *direct.MapContext, in *pb.AutokeyConfig) *krm.KMSProjectAutokeyConfigObservedState {
	if in == nil {
		return nil
	}
	out := &krm.KMSProjectAutokeyConfigObservedState{}
	out.State = direct.Enum_FromProto(mapCtx, in.GetState())
	return out
}
