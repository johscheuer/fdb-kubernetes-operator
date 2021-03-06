/*
 * update_sidecar_versions.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2019 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controllers

import (
	ctx "context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"

	fdbtypes "github.com/FoundationDB/fdb-kubernetes-operator/api/v1beta1"
)

// UpdateSidecarVersions provides a reconciliation step for upgrading the
// sidecar.
type UpdateSidecarVersions struct {
}

// Reconcile runs the reconciler's work.
func (u UpdateSidecarVersions) Reconcile(r *FoundationDBClusterReconciler, context ctx.Context, cluster *fdbtypes.FoundationDBCluster) (bool, error) {
	instances, err := r.PodLifecycleManager.GetInstances(r, cluster, context, getPodListOptions(cluster, "", "")...)
	if err != nil {
		return false, err
	}
	upgraded := false
	for _, instance := range instances {
		if instance.Pod == nil {
			return false, MissingPodError(instance, cluster)
		}

		image, err := getSidecarImage(cluster, instance)
		if err != nil {
			return false, err
		}

		for containerIndex, container := range instance.Pod.Spec.Containers {
			if container.Name == "foundationdb-kubernetes-sidecar" && container.Image != image {
				log.Info("Upgrading sidecar", "namespace", cluster.Namespace, "cluster", cluster.Name, "pod", instance.Pod.Name, "oldImage", container.Image, "newImage", image)
				err = r.PodLifecycleManager.UpdateImageVersion(r, context, cluster, instance, containerIndex, image)
				if err != nil {
					return false, err
				}
				upgraded = true
			}
		}
	}

	if upgraded {
		r.Recorder.Event(cluster, corev1.EventTypeNormal, "SidecarUpgraded", fmt.Sprintf("New version: %s", cluster.Spec.Version))
	}

	return true, nil
}

func getSidecarImage(cluster *fdbtypes.FoundationDBCluster, instance FdbInstance) (string, error) {
	settings := cluster.GetProcessSettings(instance.GetProcessClass())

	image := ""
	if settings.PodTemplate != nil {
		for _, container := range settings.PodTemplate.Spec.Containers {
			if container.Name == "foundationdb-kubernetes-sidecar" && container.Image != "" {
				image = container.Image
			}
		}
	}

	return getImage(cluster.Spec.SidecarContainer.ImageName, image, "foundationdb/foundationdb-kubernetes-sidecar", cluster.GetFullSidecarVersion(false), settings.GetAllowTagOverride())
}

// RequeueAfter returns the delay before we should run the reconciliation
// again.
func (u UpdateSidecarVersions) RequeueAfter() time.Duration {
	return 0
}
