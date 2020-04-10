package main

import (
	"fmt"
	"os"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	eventing "knative.dev/eventing/pkg/client/clientset/versioned"
	"knative.dev/eventing/pkg/utils"
	"knative.dev/pkg/kmeta"
)

func main() {
	out("Starting upgrade-check v0.13.x to v0.14.x.")

	k, _ := os.LookupEnv("KUBECONFIG")
	namespace, _ := os.LookupEnv("NAMESPACE")

	cfg, err := clientcmd.BuildConfigFromFlags("", k)
	if err != nil {
		panic(err.Error())
	}

	c := eventing.NewForConfigOrDie(cfg)
	r := checkTriggers(c, namespace)

	out("Found %d subscriptions that need upgrade.", len(r))
	for _, re := range r {
		printSubscriptionNeedRecreate(re.oldName, re.newName, re.namespace, re.found)
	}
}

type reportEntry struct {
	namespace string
	oldName   string
	newName   string
	found     bool
}

type report []reportEntry

func checkTriggers(c eventing.Interface, namespace string) report {
	triggers, err := c.EventingV1alpha1().Triggers(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	report := make([]reportEntry, 0)

	for _, trigger := range triggers.Items {

		deprecatedName := utils.GenerateFixedName(&trigger, fmt.Sprintf("%s-%s", trigger.Spec.Broker, trigger.Name))
		newName := kmeta.ChildName(fmt.Sprintf("%s-%s-", trigger.Spec.Broker, trigger.Name), string(trigger.GetUID()))

		if deprecatedName != newName {
			_, err := c.MessagingV1alpha1().Subscriptions(trigger.Namespace).Get(deprecatedName, metav1.GetOptions{})
			if err != nil && !apierrors.IsNotFound(err) {
				out("Error retrieving subscription %s/%s: %s", trigger.Namespace, deprecatedName, err.Error())
			}
			report = append(report, reportEntry{
				namespace: trigger.Namespace,
				oldName:   deprecatedName,
				newName:   newName,
				found:     err == nil,
			})
		}
	}
	return report
}

func printSubscriptionNeedRecreate(old, new, namespace string, found bool) {
	out(`Subscription needs update:
	namespace: %s
	old name: %s
	new name: %s
	found: %t`, namespace, old, new, found)
}

func out(message string, v ...interface{}) {
	fmt.Printf(message+"\n", v...)
}
