package route_rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
)

func TestCheckerWithPodsMatching(t *testing.T) {
	assert := assert.New(t)
	conf := config.NewConfig()
	config.Set(conf)

	// Setup mocks
	podList := []v1.Pod{
		fakePodsForLabels(map[string]string{"app": "reviews", "version": "v1", "stage": "production"}),
		fakePodsForLabels(map[string]string{"app": "reviews", "version": "v2", "stage": "production"}),
	}

	validations, valid := VersionPresenceChecker{"bookinfo", podList, fakeCorrectVersions()}.Check()

	// Well configured object
	assert.Empty(validations)
	assert.True(valid)
}

func fakePodsForLabels(labels labels.Set) v1.Pod {
	return v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   "reviews-12345-hello",
			Labels: labels,
		},
	}
}

func fakeCorrectVersions() kubernetes.IstioObject {
	validRouteRule := (&kubernetes.RouteRule{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "reviews",
		},
		Spec: map[string]interface{}{
			"destination": map[string]interface{}{
				"name": "reviews",
			},
			"route": []map[string]interface{}{
				map[string]interface{}{
					"weight": uint64(55),
					"labels": map[string]interface{}{
						"version": "v1",
						"stage":   "production",
					},
				},
				map[string]interface{}{
					"weight": uint64(45),
					"labels": map[string]interface{}{
						"version": "v2",
						"stage":   "production",
					},
				},
			},
		},
	}).DeepCopyIstioObject()

	return validRouteRule
}

func TestNoMatchingPods(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	podList := []v1.Pod{
		fakePodsForLabels(map[string]string{"version": "v1"}),
		fakePodsForLabels(map[string]string{"version": "v2"}),
	}
	validations, valid := VersionPresenceChecker{"bookinfo", podList, fakeNoPodsVersion()}.Check()

	// There are no pods no deployments
	assert.False(valid)
	assert.NotEmpty(validations)
	assert.Len(validations, 2)
	assert.Equal(validations[0].Message, "No pods found for this selector")
	assert.Equal(validations[0].Severity, "warning")
	assert.Equal(validations[0].Path, "spec/route[0]/labels")
	assert.Equal(validations[1].Message, "No pods found for this selector")
	assert.Equal(validations[1].Severity, "warning")
	assert.Equal(validations[1].Path, "spec/route[1]/labels")
}

func fakeNoPodsVersion() kubernetes.IstioObject {
	validRouteRule := (&kubernetes.RouteRule{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "reviews",
		},
		Spec: map[string]interface{}{
			"destination": map[string]interface{}{
				"name": "reviews",
			},
			"route": []map[string]interface{}{
				map[string]interface{}{
					"weight": uint64(55),
					"labels": map[string]interface{}{
						"version": "not-v1",
						"stage":   "production",
					},
				},
				map[string]interface{}{
					"weight": uint64(45),
					"labels": map[string]interface{}{
						"version": "not-v2",
						"stage":   "production",
					},
				},
			},
		},
	}).DeepCopyIstioObject()

	return validRouteRule
}

func TestRouteRuleWithoutDestination(t *testing.T) {
	assert := assert.New(t)
	conf := config.NewConfig()
	config.Set(conf)

	// Setup mocks
	podList := []v1.Pod{
		fakePodsForLabels(map[string]string{"app": "reviews", "version": "v1", "stage": "production"}),
		fakePodsForLabels(map[string]string{"app": "reviews", "version": "v2", "stage": "production"}),
	}

	validations, valid := VersionPresenceChecker{"bookinfo", podList, fakeNilDestination()}.Check()

	// There are no pods no deployments
	assert.False(valid)
	assert.NotEmpty(validations)
	assert.Len(validations, 2)
	assert.Equal(validations[0].Message, "No pods found for this selector")
	assert.Equal(validations[0].Severity, "warning")
	assert.Equal(validations[0].Path, "spec/route[0]/labels")
	assert.Equal(validations[1].Message, "No pods found for this selector")
	assert.Equal(validations[1].Severity, "warning")
	assert.Equal(validations[1].Path, "spec/route[1]/labels")
}

func fakeNilDestination() kubernetes.IstioObject {
	validRouteRule := (&kubernetes.RouteRule{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "reviews",
		},
		Spec: map[string]interface{}{
			"route": []map[string]interface{}{
				map[string]interface{}{
					"weight": uint64(55),
					"labels": map[string]interface{}{
						"version": "v1",
						"stage":   "production",
					},
				},
				map[string]interface{}{
					"weight": uint64(45),
					"labels": map[string]interface{}{
						"version": "v2",
						"stage":   "production",
					},
				},
			},
		},
	}).DeepCopyIstioObject()

	return validRouteRule
}

func TestRouteRuleWithBadLabels(t *testing.T) {
	assert := assert.New(t)
	conf := config.NewConfig()
	config.Set(conf)

	// Setup mocks
	podList := []v1.Pod{
		fakePodsForLabels(map[string]string{"app": "reviews", "version": "v1", "stage": "production"}),
		fakePodsForLabels(map[string]string{"app": "reviews", "version": "v2", "stage": "production"}),
	}

	validations, valid := VersionPresenceChecker{"bookinfo", podList, fakeBadLabels()}.Check()

	// There are no pods no deployments
	assert.False(valid)
	assert.NotEmpty(validations)
	assert.Len(validations, 2)
	assert.Equal(validations[0].Message, "No pods found for this selector")
	assert.Equal(validations[0].Severity, "warning")
	assert.Equal(validations[0].Path, "spec/route[0]/labels")
	assert.Equal(validations[1].Message, "No pods found for this selector")
	assert.Equal(validations[1].Severity, "warning")
	assert.Equal(validations[1].Path, "spec/route[1]/labels")
}

func fakeBadLabels() kubernetes.IstioObject {
	validRouteRule := (&kubernetes.RouteRule{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "reviews",
		},
		Spec: map[string]interface{}{
			"destination": map[string]interface{}{
				"name": "reviews",
			},
			"route": []map[string]interface{}{
				map[string]interface{}{
					"weight": uint64(55),
					"labels": "label1",
				},
				map[string]interface{}{
					"weight": uint64(45),
					"labels": "label2",
				},
			},
		},
	}).DeepCopyIstioObject()

	return validRouteRule
}

func TestRouteRuleWithoutSpec(t *testing.T) {
	assert := assert.New(t)
	conf := config.NewConfig()
	config.Set(conf)

	// Setup mocks
	podList := []v1.Pod{
		fakePodsForLabels(map[string]string{"app": "reviews", "version": "v1", "stage": "production"}),
		fakePodsForLabels(map[string]string{"app": "reviews", "version": "v2", "stage": "production"}),
	}

	validations, valid := VersionPresenceChecker{"bookinfo", podList, fakeBadSpec()}.Check()

	assert.True(valid)
	assert.Empty(validations)
}

func fakeBadSpec() kubernetes.IstioObject {
	validRouteRule := (&kubernetes.RouteRule{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "reviews",
		},
	}).DeepCopyIstioObject()

	return validRouteRule
}
