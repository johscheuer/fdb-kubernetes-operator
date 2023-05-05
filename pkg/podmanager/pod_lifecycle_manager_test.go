/*
 * pod_lifecycle_manager_test.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2020-2021 Apple Inc. and the FoundationDB project authors
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

package podmanager

import (
	fdbv1beta2 "github.com/FoundationDB/fdb-kubernetes-operator/api/v1beta2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("pod_lifecycle_manager", func() {
	var manager StandardPodLifecycleManager

	DescribeTable("getting the deletion mode of the cluster",
		func(cluster *fdbv1beta2.FoundationDBCluster, expected fdbv1beta2.PodUpdateMode) {
			Expect(manager.GetDeletionMode(cluster)).To(Equal(expected))
		},
		Entry("Without a deletion mode defined",
			&fdbv1beta2.FoundationDBCluster{},
			fdbv1beta2.PodUpdateModeZone,
		),
		Entry("With deletion mode Zone",
			&fdbv1beta2.FoundationDBCluster{
				Spec: fdbv1beta2.FoundationDBClusterSpec{
					AutomationOptions: fdbv1beta2.FoundationDBClusterAutomationOptions{
						DeletionMode: fdbv1beta2.PodUpdateModeZone,
					},
				},
			},
			fdbv1beta2.PodUpdateModeZone,
		),
		Entry("With deletion mode All",
			&fdbv1beta2.FoundationDBCluster{
				Spec: fdbv1beta2.FoundationDBClusterSpec{
					AutomationOptions: fdbv1beta2.FoundationDBClusterAutomationOptions{
						DeletionMode: fdbv1beta2.PodUpdateModeAll,
					},
				},
			},
			fdbv1beta2.PodUpdateModeAll,
		),
		Entry("With deletion mode Process Group",
			&fdbv1beta2.FoundationDBCluster{
				Spec: fdbv1beta2.FoundationDBClusterSpec{
					AutomationOptions: fdbv1beta2.FoundationDBClusterAutomationOptions{
						DeletionMode: fdbv1beta2.PodUpdateModeProcessGroup,
					},
				},
			},
			fdbv1beta2.PodUpdateModeProcessGroup,
		),
	)

	Describe("GetProcessGroupIDFromProcessID", func() {
		It("can parse a process ID", func() {
			Expect(GetProcessGroupIDFromProcessID("storage-1-1")).To(Equal("storage-1"))
		})
		It("can parse a process ID with a prefix", func() {
			Expect(GetProcessGroupIDFromProcessID("dc1-storage-1-1")).To(Equal("dc1-storage-1"))
		})

		It("can handle a process group ID with no process number", func() {
			Expect(GetProcessGroupIDFromProcessID("storage-2")).To(Equal("storage-2"))
		})
	})

	Describe("ParseProcessGroupID", func() {
		Context("with a storage ID", func() {
			It("can parse the ID", func() {
				prefix, id, err := ParseProcessGroupID("storage-12")
				Expect(err).NotTo(HaveOccurred())
				Expect(prefix).To(Equal(fdbv1beta2.ProcessClassStorage))
				Expect(id).To(Equal(12))
			})
		})

		Context("with a cluster controller ID", func() {
			It("can parse the ID", func() {
				prefix, id, err := ParseProcessGroupID("cluster_controller-3")
				Expect(err).NotTo(HaveOccurred())
				Expect(prefix).To(Equal(fdbv1beta2.ProcessClassClusterController))
				Expect(id).To(Equal(3))
			})
		})

		Context("with a custom prefix", func() {
			It("parses the prefix", func() {
				prefix, id, err := ParseProcessGroupID("dc1-storage-12")
				Expect(err).NotTo(HaveOccurred())
				Expect(prefix).To(Equal(fdbv1beta2.ProcessClass("dc1-storage")))
				Expect(id).To(Equal(12))
			})
		})

		Context("with no prefix", func() {
			It("gives a parsing error", func() {
				_, _, err := ParseProcessGroupID("6")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("could not parse process group ID 6"))
			})
		})

		Context("with no numbers", func() {
			It("gives a parsing error", func() {
				_, _, err := ParseProcessGroupID("storage")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("could not parse process group ID storage"))
			})
		})

		Context("with a text suffix", func() {
			It("gives a parsing error", func() {
				_, _, err := ParseProcessGroupID("storage-bad")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("could not parse process group ID storage-bad"))
			})
		})
	})
})
