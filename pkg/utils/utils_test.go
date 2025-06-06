// Copyright Contributors to the Open Cluster Management project

package utils

import (
	"os"
	"reflect"
	"testing"

	backplanev1 "github.com/stolostron/backplane-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_deduplicate(t *testing.T) {
	tests := []struct {
		name string
		have []backplanev1.ComponentConfig
		want []backplanev1.ComponentConfig
	}{
		{
			name: "unique components",
			have: []backplanev1.ComponentConfig{
				{Name: "component1", Enabled: true},
				{Name: "component2", Enabled: true},
			},
			want: []backplanev1.ComponentConfig{
				{Name: "component1", Enabled: true},
				{Name: "component2", Enabled: true},
			},
		},
		{
			name: "duplicate components",
			have: []backplanev1.ComponentConfig{
				{Name: "component1", Enabled: false},
				{Name: "component2", Enabled: true},
				{Name: "component1", Enabled: true},
			},
			want: []backplanev1.ComponentConfig{
				{Name: "component1", Enabled: true},
				{Name: "component2", Enabled: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deduplicate(tt.have); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deduplicate() = %v, want %v", got, tt.want)
			}
		})
	}
	m := &backplanev1.MultiClusterEngine{}
	yes := SetDefaultComponents(m)
	if !yes {
		t.Error("Setting default did not work")
	}

	yes = DeduplicateComponents(m)
	if yes {
		t.Error("Unexpected duplicates")
	}

	os.Setenv("NO_PROXY", "test")
	yes = ProxyEnvVarsAreSet()
	if !yes {
		t.Error("Unexpected proxy failure")
	}
	os.Unsetenv("NO_PROXY")
	yes = ProxyEnvVarsAreSet()
	if yes {
		t.Error("Unexpected proxy success")
	}

	var sample backplanev1.AvailabilityType
	sample = backplanev1.HAHigh

	yes = AvailabilityConfigIsValid(sample)
	if !yes {
		t.Error("Unexpected availabilitty config failure")
	}

	sample = "test"
	yes = AvailabilityConfigIsValid(sample)
	if yes {
		t.Error("Unexpected availabilitty config successs")
	}

	stringList := []string{"test1", "test2"}
	stringRemoveList := []string{"test2"}

	yes = Contains(stringList, "test1")
	if !yes {
		t.Error("Contains did not work")
	}
	attemptedRemove := Remove(stringList, "test1")
	if len(attemptedRemove) != len(stringRemoveList) {
		t.Error("Removes did not work")
	}
}

func TestGetHubSize(t *testing.T) {
	tests := []struct {
		name string
		mce  *backplanev1.MultiClusterEngine
		want backplanev1.HubSize
	}{
		{
			name: "get default",
			mce:  &backplanev1.MultiClusterEngine{},
			want: backplanev1.Small,
		},
		{
			name: "set hubsize Small",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						AnnotationHubSize: "Small",
					},
				},
			},
			want: backplanev1.Small,
		},
		{
			name: "set hubsize Medium",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						AnnotationHubSize: "Medium",
					},
				},
			},
			want: backplanev1.Medium,
		},
		{
			name: "set hubsize Large",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						AnnotationHubSize: "Large",
					},
				},
			},
			want: backplanev1.Large,
		},
		{
			name: "set hubsize XLarge",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						AnnotationHubSize: "XLarge",
					},
				},
			},
			want: backplanev1.XLarge,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHubSize(tt.mce); got != tt.want {
				t.Errorf("GetHubSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHubType(t *testing.T) {
	tests := []struct {
		name string
		env  string
		mce  *backplanev1.MultiClusterEngine
		want string
	}{
		{
			name: "mce",
			env:  "multicluster-engine",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "mce",
				},
			},
			want: "mce",
		},
		{
			name: "acm",
			env:  "multicluster-engine",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "mce",
					Labels: map[string]string{"multiclusterhubs.operator.open-cluster-management.io/managed-by": "true"},
				},
			},
			want: "acm",
		},
		{
			name: "stolostron-engine",
			env:  "stolostron-engine",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "mce",
				},
			},
			want: "stolostron-engine",
		},
		{
			name: "stolostron",
			env:  "stolostron-engine",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "mce",
					Labels: map[string]string{"multiclusterhubs.operator.open-cluster-management.io/managed-by": "true"},
				},
			},
			want: "stolostron",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OPERATOR_PACKAGE", tt.env)
			if got := GetHubType(tt.mce); got != tt.want {
				t.Errorf("GetHubType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetTestImages(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "should return correct number of test images",
			want: 79,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			images := GetTestImages()
			got := len(images)

			if got != tt.want {
				t.Errorf("GetTestImages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DefaultReplicaCount(t *testing.T) {
	tests := []struct {
		name string
		mce  *backplanev1.MultiClusterEngine
		want int
	}{
		{
			name: "should get default replica count for HABasic",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "engine",
				},
				Spec: backplanev1.MultiClusterEngineSpec{
					AvailabilityConfig: backplanev1.HABasic,
				},
			},
			want: 1,
		},
		{
			name: "should get default replica count for HAHigh",
			mce: &backplanev1.MultiClusterEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "engine",
				},
				Spec: backplanev1.MultiClusterEngineSpec{
					AvailabilityConfig: backplanev1.HAHigh,
				},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultReplicaCount(tt.mce); got != tt.want {
				t.Errorf("DefaultReplicaCount(tt.mce) = %v, want %v", got, tt.want)
			}
		})
	}
}
