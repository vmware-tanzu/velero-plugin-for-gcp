/*
Copyright 2017 the Velero contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	velerotest "github.com/vmware-tanzu/velero/pkg/test"
)

func TestGetVolumeID(t *testing.T) {
	b := &VolumeSnapshotter{}

	pv := &unstructured.Unstructured{
		Object: map[string]interface{}{},
	}

	// missing spec.gcePersistentDisk -> no error
	volumeID, err := b.GetVolumeID(pv)
	require.NoError(t, err)
	assert.Equal(t, "", volumeID)

	// missing spec.gcePersistentDisk.pdName -> error
	gce := map[string]interface{}{}
	pv.Object["spec"] = map[string]interface{}{
		"gcePersistentDisk": gce,
	}
	volumeID, err = b.GetVolumeID(pv)
	assert.Error(t, err)
	assert.Equal(t, "", volumeID)

	// valid
	gce["pdName"] = "abc123"
	volumeID, err = b.GetVolumeID(pv)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", volumeID)
}

func TestGetVolumeIDForCSI(t *testing.T) {
	b := &VolumeSnapshotter{
		log: logrus.New(),
	}

	cases := []struct {
		name    string
		csiJSON string
		want    string
		wantErr bool
	}{
		{
			name: "gke csi driver",
			csiJSON: `{
				"driver": "pd.csi.storage.gke.io",
				"fsType": "ext4",
				"volumeAttributes": {
					"storage.kubernetes.io/csiProvisionerIdentity": "1637243273131-8081-pd.csi.storage.gke.io"
				},
				"volumeHandle": "projects/velero-gcp/zones/us-central1-f/disks/pvc-a970184f-6cc1-4769-85ad-61dcaf8bf51d"
			}`,
			want: "pvc-a970184f-6cc1-4769-85ad-61dcaf8bf51d",
			wantErr: false,
		},
		{
			name: "gke csi driver with invalid handle name",
			csiJSON: `{
				"driver": "pd.csi.storage.gke.io",
				"fsType": "ext4",
				"volumeHandle": "pvc-a970184f-6cc1-4769-85ad-61dcaf8bf51d"
			}`,
			want: "",
			wantErr: true,
		},
		{
			name: "unknown driver",
			csiJSON: `{
				"driver": "xxx.csi.storage.gke.io",
				"fsType": "ext4",
				"volumeHandle": "pvc-a970184f-6cc1-4769-85ad-61dcaf8bf51d"
			}`,
			want: "",
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res := &unstructured.Unstructured{
				Object: map[string]interface{}{},
			}
			csi := map[string]interface{}{}
			json.Unmarshal([]byte(tt.csiJSON), &csi)
			res.Object["spec"] = map[string]interface{}{
				"csi": csi,
			}
			volumeID, err := b.GetVolumeID(res)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, volumeID)
		})
	}

}

func TestSetVolumeID(t *testing.T) {
	b := &VolumeSnapshotter{}

	pv := &unstructured.Unstructured{
		Object: map[string]interface{}{},
	}

	// missing spec.gcePersistentDisk -> error
	updatedPV, err := b.SetVolumeID(pv, "abc123")
	require.Error(t, err)

	// happy path
	gce := map[string]interface{}{}
	pv.Object["spec"] = map[string]interface{}{
		"gcePersistentDisk": gce,
	}
	updatedPV, err = b.SetVolumeID(pv, "123abc")
	require.NoError(t, err)

	res := new(v1.PersistentVolume)
	require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(updatedPV.UnstructuredContent(), res))
	require.NotNil(t, res.Spec.GCEPersistentDisk)
	assert.Equal(t, "123abc", res.Spec.GCEPersistentDisk.PDName)
}

func TestSetVolumeIDForCSI(t *testing.T) {
	b := &VolumeSnapshotter{
		log: logrus.New(),
	}

	cases := []struct {
		name     string
		csiJSON  string
		volumeID string
		wantErr  bool
	}{
		{
			name: "set ID to CSI with GKE pd CSI driver",
			csiJSON: `{
				 "driver": "pd.csi.storage.gke.io",
				 "fsType": "ext4",
				 "volumeHandle": "projects/velero-gcp/zones/us-central1-f/disks/pvc-a970184f-6cc1-4769-85ad-61dcaf8bf51d"
			}`,
			volumeID: "restore-fd9729b5-868b-4544-9568-1c5d9121dabc",
			wantErr:  false,
		},
		{
			name: "set ID to CSI with GKE pd CSI driver, but the volumeHandle is invalid",
			csiJSON: `{
				 "driver": "pd.csi.storage.gke.io",
				 "fsType": "ext4",
				 "volumeHandle": "pvc-a970184f-6cc1-4769-85ad-61dcaf8bf51d"
			}`,
			volumeID: "restore-fd9729b5-868b-4544-9568-1c5d9121dabc",
			wantErr:  true,
		},
		{
			name: "set ID to CSI with unknown driver",
			csiJSON: `"{
				 "driver": "xxx.csi.storage.gke.io",
				 "fsType": "ext4",
				 "volumeHandle": "projects/velero-gcp/zones/us-central1-f/disks/pvc-a970184f-6cc1-4769-85ad-61dcaf8bf51d"
			}`,
			volumeID: "restore-fd9729b5-868b-4544-9568-1c5d9121dabc",
			wantErr:  true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res := &unstructured.Unstructured{
				Object: map[string]interface{}{},
			}
			csi := map[string]interface{}{}
			json.Unmarshal([]byte(tt.csiJSON), &csi)
			res.Object["spec"] = map[string]interface{}{
				"csi": csi,
			}
			originalVolHanle, _ := csi["volumeHandle"].(string)
			newRes, err := b.SetVolumeID(res, tt.volumeID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				newPV := new(v1.PersistentVolume)
				require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(newRes.UnstructuredContent(), newPV))
				ind := strings.LastIndex(newPV.Spec.CSI.VolumeHandle, "/")
				assert.Equal(t, tt.volumeID, newPV.Spec.CSI.VolumeHandle[ind+1:])
				assert.Equal(t, originalVolHanle[:ind], newPV.Spec.CSI.VolumeHandle[:ind])
			}
		})
	}
}

func TestGetSnapshotTags(t *testing.T) {
	tests := []struct {
		name            string
		veleroTags      map[string]string
		diskDescription string
		expected        string
	}{
		{
			name:            "degenerate case (no tags)",
			veleroTags:      nil,
			diskDescription: "",
			expected:        "",
		},
		{
			name: "velero tags only get applied",
			veleroTags: map[string]string{
				"velero-key1": "velero-val1",
				"velero-key2": "velero-val2",
			},
			diskDescription: "",
			expected:        `{"velero-key1":"velero-val1","velero-key2":"velero-val2"}`,
		},
		{
			name:            "disk tags only get applied",
			veleroTags:      nil,
			diskDescription: `{"gcp-key1":"gcp-val1","gcp-key2":"gcp-val2"}`,
			expected:        `{"gcp-key1":"gcp-val1","gcp-key2":"gcp-val2"}`,
		},
		{
			name:            "non-overlapping velero and disk tags both get applied",
			veleroTags:      map[string]string{"velero-key": "velero-val"},
			diskDescription: `{"gcp-key":"gcp-val"}`,
			expected:        `{"velero-key":"velero-val","gcp-key":"gcp-val"}`,
		},
		{
			name: "when tags overlap, velero tags take precedence",
			veleroTags: map[string]string{
				"velero-key":      "velero-val",
				"overlapping-key": "velero-val",
			},
			diskDescription: `{"gcp-key":"gcp-val","overlapping-key":"gcp-val"}`,
			expected:        `{"velero-key":"velero-val","gcp-key":"gcp-val","overlapping-key":"velero-val"}`,
		},
		{
			name:            "if disk description is invalid JSON, apply just velero tags",
			veleroTags:      map[string]string{"velero-key": "velero-val"},
			diskDescription: `THIS IS INVALID JSON`,
			expected:        `{"velero-key":"velero-val"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := getSnapshotTags(test.veleroTags, test.diskDescription, velerotest.NewLogger())

			if test.expected == "" {
				assert.Equal(t, test.expected, res)
				return
			}

			var actualMap map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(res), &actualMap))

			var expectedMap map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(test.expected), &expectedMap))

			assert.Equal(t, len(expectedMap), len(actualMap))
			for k, v := range expectedMap {
				assert.Equal(t, v, actualMap[k])
			}
		})
	}
}

func TestRegionHelpers(t *testing.T) {
	tests := []struct {
		name                string
		volumeAZ            string
		expectedRegion      string
		expectedIsMultiZone bool
		expectedError       error
	}{
		{
			name:                "valid multizone(2) tag",
			volumeAZ:            "us-central1-a__us-central1-b",
			expectedRegion:      "us-central1",
			expectedIsMultiZone: true,
			expectedError:       nil,
		},
		{
			name:                "valid multizone(4) tag",
			volumeAZ:            "us-central1-a__us-central1-b__us-central1-f__us-central1-e",
			expectedRegion:      "us-central1",
			expectedIsMultiZone: true,
			expectedError:       nil,
		},
		{
			name:                "valid single zone tag",
			volumeAZ:            "us-central1-a",
			expectedRegion:      "us-central1",
			expectedIsMultiZone: false,
			expectedError:       nil,
		},
		{
			name:                "invalid single zone tag",
			volumeAZ:            "us^central1^a",
			expectedRegion:      "",
			expectedIsMultiZone: false,
			expectedError:       errors.Errorf("failed to parse region from zone: %q", "us^central1^a"),
		},
		{
			name:                "invalid multizone tag",
			volumeAZ:            "us^central1^a__us^central1^b",
			expectedRegion:      "",
			expectedIsMultiZone: true,
			expectedError:       errors.Errorf("failed to parse region from zone: %q", "us^central1^a__us^central1^b"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedIsMultiZone, isMultiZone(test.volumeAZ))
			region, err := parseRegion(test.volumeAZ)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			}
			assert.Equal(t, test.expectedRegion, region)
		})
	}
}
