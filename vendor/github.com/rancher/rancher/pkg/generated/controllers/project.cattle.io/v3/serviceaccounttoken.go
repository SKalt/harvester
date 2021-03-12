/*
Copyright 2021 Rancher Labs, Inc.

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

package v3

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v3 "github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ServiceAccountTokenHandler func(string, *v3.ServiceAccountToken) (*v3.ServiceAccountToken, error)

type ServiceAccountTokenController interface {
	generic.ControllerMeta
	ServiceAccountTokenClient

	OnChange(ctx context.Context, name string, sync ServiceAccountTokenHandler)
	OnRemove(ctx context.Context, name string, sync ServiceAccountTokenHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ServiceAccountTokenCache
}

type ServiceAccountTokenClient interface {
	Create(*v3.ServiceAccountToken) (*v3.ServiceAccountToken, error)
	Update(*v3.ServiceAccountToken) (*v3.ServiceAccountToken, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ServiceAccountToken, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ServiceAccountTokenList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ServiceAccountToken, err error)
}

type ServiceAccountTokenCache interface {
	Get(namespace, name string) (*v3.ServiceAccountToken, error)
	List(namespace string, selector labels.Selector) ([]*v3.ServiceAccountToken, error)

	AddIndexer(indexName string, indexer ServiceAccountTokenIndexer)
	GetByIndex(indexName, key string) ([]*v3.ServiceAccountToken, error)
}

type ServiceAccountTokenIndexer func(obj *v3.ServiceAccountToken) ([]string, error)

type serviceAccountTokenController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewServiceAccountTokenController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ServiceAccountTokenController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &serviceAccountTokenController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromServiceAccountTokenHandlerToHandler(sync ServiceAccountTokenHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ServiceAccountToken
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ServiceAccountToken))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *serviceAccountTokenController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ServiceAccountToken))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateServiceAccountTokenDeepCopyOnChange(client ServiceAccountTokenClient, obj *v3.ServiceAccountToken, handler func(obj *v3.ServiceAccountToken) (*v3.ServiceAccountToken, error)) (*v3.ServiceAccountToken, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *serviceAccountTokenController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *serviceAccountTokenController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *serviceAccountTokenController) OnChange(ctx context.Context, name string, sync ServiceAccountTokenHandler) {
	c.AddGenericHandler(ctx, name, FromServiceAccountTokenHandlerToHandler(sync))
}

func (c *serviceAccountTokenController) OnRemove(ctx context.Context, name string, sync ServiceAccountTokenHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromServiceAccountTokenHandlerToHandler(sync)))
}

func (c *serviceAccountTokenController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *serviceAccountTokenController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *serviceAccountTokenController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *serviceAccountTokenController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *serviceAccountTokenController) Cache() ServiceAccountTokenCache {
	return &serviceAccountTokenCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *serviceAccountTokenController) Create(obj *v3.ServiceAccountToken) (*v3.ServiceAccountToken, error) {
	result := &v3.ServiceAccountToken{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *serviceAccountTokenController) Update(obj *v3.ServiceAccountToken) (*v3.ServiceAccountToken, error) {
	result := &v3.ServiceAccountToken{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *serviceAccountTokenController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *serviceAccountTokenController) Get(namespace, name string, options metav1.GetOptions) (*v3.ServiceAccountToken, error) {
	result := &v3.ServiceAccountToken{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *serviceAccountTokenController) List(namespace string, opts metav1.ListOptions) (*v3.ServiceAccountTokenList, error) {
	result := &v3.ServiceAccountTokenList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *serviceAccountTokenController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *serviceAccountTokenController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ServiceAccountToken, error) {
	result := &v3.ServiceAccountToken{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type serviceAccountTokenCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *serviceAccountTokenCache) Get(namespace, name string) (*v3.ServiceAccountToken, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ServiceAccountToken), nil
}

func (c *serviceAccountTokenCache) List(namespace string, selector labels.Selector) (ret []*v3.ServiceAccountToken, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ServiceAccountToken))
	})

	return ret, err
}

func (c *serviceAccountTokenCache) AddIndexer(indexName string, indexer ServiceAccountTokenIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ServiceAccountToken))
		},
	}))
}

func (c *serviceAccountTokenCache) GetByIndex(indexName, key string) (result []*v3.ServiceAccountToken, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ServiceAccountToken, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ServiceAccountToken))
	}
	return result, nil
}
