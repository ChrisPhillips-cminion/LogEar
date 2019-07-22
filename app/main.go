/*
Copyright 2017 The Kubernetes Authors.
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
	"bytes"
	"flag"
	"fmt"
	"io"
	"k8s.io/klog"
	"net/http"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
)

var version string

func main() {

	var podname string
	var namespace string
	var userSet string
	var passwordSet string

	flag.StringVar(&podname, "podname", "podname", "PodName")
	flag.StringVar(&namespace, "namespace", "default", "NameSpace")
	flag.StringVar(&userSet, "username", "default", "unset")
	flag.StringVar(&passwordSet, "password", "default", "unset")
	flag.Parse()

	fmt.Printf("Version: %v", version)
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	if err != nil {
		klog.Fatal(err)
	}
	podLogOpts := corev1.PodLogOptions{}
	req := clientset.CoreV1().Pods("default").GetLogs("kubernetes-bootcamp-6bf84cb898-hsg6w", &podLogOpts)
	podLogs, err := req.Stream()
	if err != nil {
		klog.Fatal(err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		klog.Fatal(err)
	}
	str := buf.String()


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if (userSet != "unset" && passwordSet  != "unset") {
			user, password, ok := r.BasicAuth()
			// fmt.Printf("%v",user)
			// fmt.Printf("%v",password)
			// fmt.Printf("%v",ok)
			if !ok || user != userSet || password != passwordSet {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+namespace+"-"+podname+`"`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorised.\n"))

			} else {
				fmt.Fprintf(w, str)
			}
		} else {
			fmt.Fprintf(w, str)
		}
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Healthy")
	})
	fmt.Printf(str)
	http.ListenAndServe(":80", nil)
}
