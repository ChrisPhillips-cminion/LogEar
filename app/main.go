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
	"strings"
	"flag"
	"fmt"
	"io"
	"log"
	"k8s.io/klog"
	"net/http"

	corev1 "k8s.io/api/core/v1"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	log.Printf("Version: %v", version)
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("%v", err)
		klog.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("%v", err)
		klog.Fatal(err)
	}

	// deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)




	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ok:=true

		var user string
		var password string
		log.Printf("username '%v' password '%v'", userSet, passwordSet)
		if (userSet != "unset" && passwordSet  != "unset") {
			log.Printf("Requesting basic auth")
			user, password, ok = r.BasicAuth()
		}
			log.Printf("%v",user)
			log.Printf("%v",password)
			log.Printf("%v",ok)
			if !ok ||  ( ok && (user != userSet || password != passwordSet)) {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+namespace+"\t"+podname+`"`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorised.\n"))

			} else {

				fmt.Fprintf(w, reverseLineOrder(getAndPrintLog(namespace,podname,clientset)))
			}

	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Healthy")
	})
	// fmt.Printf(str)
	http.ListenAndServe(":80", nil)
}

func reverseLineOrder(str string) string {
	array := strings.Split(str,"\n")
	array2 := make([]string,len(array))
	count:=len(array);
	for _,v := range array {
		count--
		array2[count]=v
	}
	return strings.Join(array2,"\n")
}

func getAndPrintLog(namespace string, podname string, clientset *kubernetes.Clientset) string {
	req := clientset.CoreV1().Pods(namespace).GetLogs(podname, &corev1.PodLogOptions{})
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
	return buf.String()
}
