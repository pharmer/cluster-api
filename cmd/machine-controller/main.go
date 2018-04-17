/*
Copyright 2018 The Kubernetes Authors.

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
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/util/logs"

	//"context"
	//"sigs.k8s.io/cluster-api/cloud"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/cluster-api/pkg/controller/config"
	"sigs.k8s.io/cluster-api/pkg/controller/machine"
	"sigs.k8s.io/cluster-api/pkg/controller/sharedinformers"
	"fmt"
	"sigs.k8s.io/cluster-api/cloud"
	"context"
	_ "sigs.k8s.io/cluster-api/cloud/providers"
)

var (
	kubeadmToken = pflag.String("token", "", "Kubeadm token to use to join new machines")
	Provider     = pflag.String("provider", "gce", "Cloud provider name")
)

func init() {
	config.ControllerConfig.AddFlags(pflag.CommandLine)
}

func main() {
	pflag.Parse()

	logs.InitLogs()
	defer logs.FlushLogs()

	config, err := controller.GetConfig(config.ControllerConfig.Kubeconfig)
	if err != nil {
		glog.Fatalf("Could not create Config for talking to the apiserver: %v", err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create client for talking to the apiserver: %v", err)
	}
	fmt.Println(client)

	cm, err := cloud.GetCloudManager(*Provider, context.Background())
	if err != nil {
		glog.Fatalf("Could not found provider: %v", err)
	}
	err = cm.PreparedActuator(*kubeadmToken, client.ClusterV1alpha1().Machines(corev1.NamespaceDefault))
	if err != nil {
		glog.Fatalf("Could not create machine actuator: %v", err)
	}

	shutdown := make(chan struct{})
	si := sharedinformers.NewSharedInformers(config, shutdown)
	// If this doesn't compile, the code generator probably
	// overwrote the customized NewMachineController function.
	c := machine.NewMachineController(config, si, cm)
	c.Run(shutdown)
	select {}
}