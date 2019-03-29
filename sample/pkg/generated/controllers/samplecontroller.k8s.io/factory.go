/*
Copyright The Kubernetes Authors.

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

// Code generated by main. DO NOT EDIT.

package samplecontroller

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/generic"
	clientset "github.com/rancher/wrangler/sample/pkg/generated/clientset/versioned"
	informers "github.com/rancher/wrangler/sample/pkg/generated/informers/externalversions"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type Factory struct {
	synced            bool
	informerFactory   informers.SharedInformerFactory
	clientset         clientset.Interface
	controllerManager *generic.ControllerManager
	threadiness       map[schema.GroupVersionKind]int
}

func NewFactoryFromConfigOrDie(config *rest.Config) *Factory {
	f, err := NewFactoryFromConfig(config)
	if err != nil {
		panic(err)
	}
	return f
}

func NewFactoryFromConfig(config *rest.Config) (*Factory, error) {
	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, 2*time.Hour)
	return NewFactory(clientset, informerFactory), nil
}

func NewFactory(clientset clientset.Interface, informerFactory informers.SharedInformerFactory) *Factory {
	return &Factory{
		threadiness:       map[schema.GroupVersionKind]int{},
		controllerManager: &generic.ControllerManager{},
		clientset:         clientset,
		informerFactory:   informerFactory,
	}
}

func (c *Factory) SetThreadiness(gvk schema.GroupVersionKind, threadiness int) {
	c.threadiness[gvk] = threadiness
}

func (c *Factory) Sync(ctx context.Context) error {
	c.informerFactory.Start(ctx.Done())
	c.informerFactory.WaitForCacheSync(ctx.Done())
	return nil
}

func (c *Factory) Start(ctx context.Context, defaultThreadiness int) error {
	if err := c.Sync(ctx); err != nil {
		return err
	}

	return c.controllerManager.Start(ctx, defaultThreadiness, c.threadiness)
}

func (c *Factory) Samplecontroller() Interface {
	return New(c.controllerManager, c.informerFactory.Samplecontroller(), c.clientset)
}
