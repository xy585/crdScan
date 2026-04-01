package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	adminToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6InhMeHVtZnZBZGgtbzVpWVRQc3V6NF9OZFk3U21KVV92N0FKUGNjSTZRQncifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzc1MDE5NjM2LCJpYXQiOjE3NzUwMTYwMzYsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiNTJjNThmMGQtYmYwMi00OTc3LWIxOWItNmQ5NGI2OTRmMjk0Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImFkbWluIiwidWlkIjoiNmI4NmFhZWQtMzRlYi00MzY2LThmZjAtYThjY2M1MWY3ZDFkIn19LCJuYmYiOjE3NzUwMTYwMzYsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmFkbWluIn0.pI8RAnvSMWVGPRB6Q574NnULe2liLX6wU3uaHhxNtVCFnFJxDUsRXG9aAg2VgkZDSBrHB9nuKvqs0si4jzuTzeAVSiHxfJrLOEkVRwzmBG-H63gR9M2Vz1JnvFsntlR2A1LZ42UpqWwl3FPfwf-J4EZa_i24mgMe6NDzRv5DJht3kztEBB8qHlG8SMkV5cDRn0yRrEPR4qeyV23xbrRwhStjEweE_mkEB9522gqAOUni7EMTPbpH26hmCYk_WsEIrCo78yLFT6lrvG-QlAaPxUB5N7zGMmjqt-vInl-TCAXEJd8a2HyHn7mUbo7ToqenvU7Q5xfzh0-Tyf_TYKC5vw"
	apiServer  = "34.134.249.25:6443" // apiserver address
	caBundle   = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZJRENDQXdpZ0F3SUJBZ0lVSEZCeXJsQzJnamNhNFhPQkE0VE02VHVQWVlVd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0Z6RVZNQk1HQTFVRUF3d01lVzkxY2kxM1pXSm9iMjlyTUI0WERUSTJNRFF3TVRBek5UQXhOMW9YRFRJMwpNRFF3TVRBek5UQXhOMW93RnpFVk1CTUdBMVVFQXd3TWVXOTFjaTEzWldKb2IyOXJNSUlDSWpBTkJna3Foa2lHCjl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUF2ZkNyZmgwcUhsdCtJSklwNFVzYkZWcldPbGVwOWN1ZnBveDcKQ3g3QklSd2Q0MHJqWXVLak1BcGVuTXRSRFBCZitTcWVWTURNays3bXZFRE5iNlFFcVhkVzhDRzQyanB5ak9lZApTdWNWcE5ORHJ4aThIZWN2RWtxWGgwQ3lTOEZkS2ZBc3BsQmllbHM0MHhaL3N4andpM21OQW5YditrcVdCRVR3ClFWNHdGRWpnMk81TjB6WjQ2MUJTc1lBZkpnL21tOEI0MExlOWRVc2VQMmtxQ2N2Tng1b0JRZ3pqa2pVbUtXMDIKU205WDhCNU9JbktnM1o5MkJIR3BKaEpnK1g5L0xUaG5IVm90U2tlQVNoOUptcXJ4MVRJR0VEODVLTFVYak4rRwpTendySS9pVWN3VEVYdmxkbDFhczZZN3BDVzVDTDlac0svK09WWkFIekZGcHY5eUFoRG96OVEwbVdlRW1iaXB2CmxWVnF4RzJUTWpzMzNLNWVnRGFubXBzOTUvWUl6OWxYZUF5WFpEZURzYlorVnZwejd6RVJLVEVoZ3c2OHRKOWwKb0JQbUxtR0ZWZUs2YjZSNkNuTzlzYlJhZmVvQ2dGZS9NMEtnSFJ4eHFEWUhHb1plMUdTOGd0V1hCdWluSXgzaQpkcDJjU21acHptTnc1aVluMm1ZTWFSODlvcEN6SS9yTDQrQ09zeUIyelNpaUpDYVZYQUFZdU9RMUpCNUptWjlpCk1sVzlZd3V2a0RlVmJuUDJGQUtKcnpEV0RDQllnZmJWTEFqeHNtZEorOFpPakNpVnlBNndBWmFMNytNTjNZRGYKRVp4Vm9DNHFtNzEvVHVVRGxUT0VYaVkzS2pOL1BQNXNWaDBNRFAyY2U4MXJaZXBoeC9GQktoMGFjZjFUZEowdQpONThtS2ljQ0F3RUFBYU5rTUdJd0hRWURWUjBPQkJZRUZQZ2czeGhlOXQ1bFdzbnQrZkRVb0ExQlFsOEhNQjhHCkExVWRJd1FZTUJhQUZQZ2czeGhlOXQ1bFdzbnQrZkRVb0ExQlFsOEhNQThHQTFVZEV3RUIvd1FGTUFNQkFmOHcKRHdZRFZSMFJCQWd3Qm9jRUlrSm5vakFOQmdrcWhraUc5dzBCQVFzRkFBT0NBZ0VBYXd2dm10bGgxUEtXN3N3SgpYOVFocXNmUVZyMGQ1UGMza2hkNnp3cWxja1pOK2w4VDVJei9Rem1xS3lvS3QybmJiVWNDdEhmcWMxSVh3Y2NmCk4vYlV0Ky90OHJUU3E5bFhlTzZ4VDVqK1M3NG5laFRRV3V2QStSRjV1dWlJZTlOY0VhM3JCVHgxaUxpRytnMm8KVVkyUW15LzRacUdCcVZSdkU2d3gyeGFkMjUwK2ZZRVBwUWszUnZodmEzRkNidk9YUUZvSVRaZ0dYV0lIT2hlVQpZbzNTWmNienFKOGdDLy9nNytiWGd3TG54bVBnM0ZDNk5Bb21qSXZ5TFpyeXc0UWFkWTk3d1c0VldKYUNZYjFpCjU1TmN3SzRERENjeWZHS0NwaHUwWlNIL3NOc2RMUU1ibWtCc3Y4WlBvTnhoYjlya3VOYVU1RGRWM0o2cHhiZGYKRGk0ak9PM0JyNVVNQU1rb2hDTW5RWWlVZHFnanZVNEtkS2lQY0tobk5LdURKN3FtOUtCdXA3Uk1rS2ZJVG1WTApiRWR3bmIxQ1B4TlF5aTR6QVBkS1kwNWp0Y1hBOGYyZkZQYUJ3dW03MzYyQjNjODA2OGE4Zmp3OE1nL1U5cWJRCktEckZpVE9XNlVTaUxOWW42TFFiamJNSnZIUFhzSXlxd3lFbWJGbTlBZkF4WHlqeFpPWjdUNHNtQU1mU28rSTcKcU1sbGttNWZvQmVJekkvdXRwNlRzT3BzWkE3QzQ0c3d2L0JTc1B5ZXdvNXBadmtvRkdJaEVSVDJXUVE4SXEvQgo0MFVZbmRvaTdNQ3VlQlBZdEZFTXlyVzJ2RDh4TnNTZXcrcFdtMWsvV1NIUHF0Y2RvNUxKNGFTaWU2YzdtTyszCjJyeTlEL3ZucjhCcjVPS2NIeFlkdlRka0crYz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="

	ports          = []string{"22", "80", "20", "8081"}
	webhookAddress = "https://" + "34.66.103.162:443"
	crdName        = "testcr.crd.com"
	crName         = "testcr"
)

func main() {
	if err := loadCredentialsFromEnv(); err != nil {
		fmt.Println(err)
		return
	}

	for _, port := range ports {
		// for port := 10000; port <= 15000; port++ {
		webhookServer := fmt.Sprintf(webhookAddress+"/%s", port)
		// webhookServer := fmt.Sprintf(webhookadd+"/%s", strconv.Itoa(port))
		err := patchCRD(crdName, webhookServer, caBundle)
		if err != nil {
			fmt.Println(err)
		}
		// if err == nil {
		// 	fmt.Println(port)
		// }
		err = getReq()
		if err != nil {
			if !strings.Contains(err.Error(), "connection refused") {
				fmt.Println("-------------------")
				fmt.Println(port)
				fmt.Println(err)
				fmt.Println("-------------------")
			}
			//fmt.Println(err)
		}
	}

}

func loadCredentialsFromEnv() error {
	adminToken = strings.TrimSpace(adminToken)
	if adminToken == "" {
		adminToken = strings.TrimSpace(os.Getenv("ADMIN_TOKEN"))
	}

	caBundle = strings.TrimSpace(caBundle)
	if caBundle == "" {
		caBundle = strings.TrimSpace(os.Getenv("CA_BUNDLE"))
	}

	if adminToken == "" {
		return fmt.Errorf("adminToken is empty: set code value or environment variable ADMIN_TOKEN")
	}
	if caBundle == "" {
		return fmt.Errorf("caBundle is empty: set code value or environment variable CA_BUNDLE")
	}

	return nil
}

func patchCRD(crdName string, webhookURL string, ca string) error {
	clientset := getAPIExtensionsClientSet()

	patchOps := []JSONPatchOperation{
		{
			Op:    "replace",
			Path:  "/spec/conversion/webhook/clientConfig/url",
			Value: webhookURL,
		},
		{
			Op:    "replace",
			Path:  "/spec/conversion/webhook/clientConfig/caBundle",
			Value: ca,
		},
	}

	patchData, err := json.Marshal(patchOps)
	if err != nil {
		return err
	}

	_, err = clientset.ApiextensionsV1().CustomResourceDefinitions().Patch(context.TODO(), crdName, types.JSONPatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

func getReq() error {
	config := getRestConfig()
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create dynamic client failed: %w", err)
	}

	clientset := getClientSet()
	gvr, namespaced, err := resolveResourceGVR(clientset.Discovery(), crName)
	if err != nil {
		return err
	}

	if namespaced {
		_, err = dynamicClient.Resource(gvr).Namespace("default").List(context.TODO(), metav1.ListOptions{})
	} else {
		_, err = dynamicClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	}
	if err != nil {
		return fmt.Errorf("get %s failed: %w", crName, err)
	}

	return nil
}

type JSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func getClientSet() *kubernetes.Clientset {
	config := getRestConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating kubernetes client: %v\n", err)
	}
	return clientset
}

func getAPIExtensionsClientSet() *apiextensionsclientset.Clientset {
	config := getRestConfig()

	clientset, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating apiextensions client: %v\n", err)
	}
	return clientset
}

func getRestConfig() *rest.Config {
	config := &rest.Config{}
	config.TLSClientConfig.Insecure = true
	config.BearerToken = adminToken
	config.Host = "https://" + apiServer
	return config
}

func resolveResourceGVR(discoveryClient discovery.DiscoveryInterface, resource string) (schema.GroupVersionResource, bool, error) {
	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return schema.GroupVersionResource{}, false, fmt.Errorf("discover resources failed: %w", err)
	}

	for _, apiResourceList := range apiResourceLists {
		gv, parseErr := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if parseErr != nil {
			continue
		}

		for _, apiResource := range apiResourceList.APIResources {
			if apiResource.Name == resource || apiResource.SingularName == resource {
				return schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: apiResource.Name}, apiResource.Namespaced, nil
			}
			for _, shortName := range apiResource.ShortNames {
				if shortName == resource {
					return schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: apiResource.Name}, apiResource.Namespaced, nil
				}
			}
		}
	}

	return schema.GroupVersionResource{}, false, fmt.Errorf("resource %q not found via discovery", resource)
}
