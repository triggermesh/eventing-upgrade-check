package main

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"knative.dev/eventing/pkg/apis/eventing/v1alpha1"
	fakeeventingclient "knative.dev/eventing/pkg/client/clientset/versioned/fake"
)

func TestCheckTriggers(t *testing.T) {
	table := []struct {
		TestName       string
		Namespace      string
		Objects        []runtime.Object
		ExpectedReport report
	}{
		{
			TestName:       "No triggers",
			Namespace:      metav1.NamespaceAll,
			Objects:        []runtime.Object{},
			ExpectedReport: report{},
		},
		{
			TestName:  "Subscription does not need update",
			Namespace: metav1.NamespaceAll,
			Objects: []runtime.Object{
				makeTrigger("namespace1", "kek", "49c3088f-32f9-11ea-be82-42010a800192", "default"),
			},
			ExpectedReport: report{},
		},
		{
			TestName:  "Subscription needs update",
			Namespace: metav1.NamespaceAll,
			Objects: []runtime.Object{
				makeTrigger("namespace1", "test-trigger-from-yaml", "56d31d3e-32f9-11ea-be82-42010a800192", "default"),
			},
			ExpectedReport: report{makeReportEntry("default-test-trigger-from--56d31d3e-32f9-11ea-be82-42010a800192", "default-test-trigger-from-yaml-653c61f8de5d23c94739755596ff8e6a", "namespace1", false)},
		},
	}

	for _, tc := range table {
		t.Run(tc.TestName, func(t *testing.T) {
			c := fakeeventingclient.NewSimpleClientset(tc.Objects...)
			report := checkTriggers(c, tc.Namespace)

			if !reflect.DeepEqual(tc.ExpectedReport, report) {
				t.Errorf("Expected report: %+v\nGot report: %+v", tc.ExpectedReport, report)
			}
		})
	}
}

func makeTrigger(namespace, name, UID, brokerName string) *v1alpha1.Trigger {
	return &v1alpha1.Trigger{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "eventing.knative.dev/v1alpha1",
			Kind:       "Trigger",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID:       types.UID(UID),
		},
		Spec: v1alpha1.TriggerSpec{
			Broker: brokerName,
		},
	}
}

func makeReportEntry(old, new, namespace string, found bool) reportEntry {
	return reportEntry{
		namespace: namespace,
		oldName:   old,
		newName:   new,
		found:     found,
	}
}
