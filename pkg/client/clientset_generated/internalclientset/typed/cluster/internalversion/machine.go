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
package internalversion

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	cluster "sigs.k8s.io/cluster-api/pkg/apis/cluster"
	scheme "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/internalclientset/scheme"
)

// MachinesGetter has a method to return a MachineInterface.
// A group's client should implement this interface.
type MachinesGetter interface {
	Machines(namespace string) MachineInterface
}

// MachineInterface has methods to work with Machine resources.
type MachineInterface interface {
	Create(*cluster.Machine) (*cluster.Machine, error)
	Update(*cluster.Machine) (*cluster.Machine, error)
	UpdateStatus(*cluster.Machine) (*cluster.Machine, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*cluster.Machine, error)
	List(opts v1.ListOptions) (*cluster.MachineList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *cluster.Machine, err error)
	MachineExpansion
}

// machines implements MachineInterface
type machines struct {
	client rest.Interface
	ns     string
}

// newMachines returns a Machines
func newMachines(c *ClusterClient, namespace string) *machines {
	return &machines{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the machine, and returns the corresponding machine object, and an error if there is any.
func (c *machines) Get(name string, options v1.GetOptions) (result *cluster.Machine, err error) {
	result = &cluster.Machine{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("machines").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Machines that match those selectors.
func (c *machines) List(opts v1.ListOptions) (result *cluster.MachineList, err error) {
	result = &cluster.MachineList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("machines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested machines.
func (c *machines) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("machines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a machine and creates it.  Returns the server's representation of the machine, and an error, if there is any.
func (c *machines) Create(machine *cluster.Machine) (result *cluster.Machine, err error) {
	result = &cluster.Machine{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("machines").
		Body(machine).
		Do().
		Into(result)
	return
}

// Update takes the representation of a machine and updates it. Returns the server's representation of the machine, and an error, if there is any.
func (c *machines) Update(machine *cluster.Machine) (result *cluster.Machine, err error) {
	result = &cluster.Machine{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("machines").
		Name(machine.Name).
		Body(machine).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *machines) UpdateStatus(machine *cluster.Machine) (result *cluster.Machine, err error) {
	result = &cluster.Machine{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("machines").
		Name(machine.Name).
		SubResource("status").
		Body(machine).
		Do().
		Into(result)
	return
}

// Delete takes name of the machine and deletes it. Returns an error if one occurs.
func (c *machines) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("machines").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *machines) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("machines").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched machine.
func (c *machines) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *cluster.Machine, err error) {
	result = &cluster.Machine{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("machines").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
