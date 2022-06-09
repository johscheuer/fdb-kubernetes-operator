/*
 * backup_client.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2022 Apple Inc. and the FoundationDB project authors
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

package fdbclient

import (
	"encoding/json"
	"fmt"
	fdbv1beta2 "github.com/FoundationDB/fdb-kubernetes-operator/api/v1beta2"
	"github.com/FoundationDB/fdb-kubernetes-operator/internal"
)

// StartBackup starts a new backup.
func (client *cliAdminClient) StartBackup(url string, snapshotPeriodSeconds int) error {
	_, err := client.runCommand(cliCommand{
		binary: "fdbbackup",
		args: []string{
			"start",
			"-d",
			url,
			"-s",
			fmt.Sprintf("%d", snapshotPeriodSeconds),
			"-z",
		},
	})
	return err
}

// StopBackup stops a backup.
func (client *cliAdminClient) StopBackup(_ string) error {
	_, err := client.runCommand(cliCommand{
		binary: "fdbbackup",
		args: []string{
			"discontinue",
		},
	})
	return err
}

// PauseBackups pauses the backups.
func (client *cliAdminClient) PauseBackups() error {
	_, err := client.runCommand(cliCommand{
		binary: "fdbbackup",
		args: []string{
			"pause",
		},
	})
	return err
}

// ResumeBackups resumes the backups.
func (client *cliAdminClient) ResumeBackups() error {
	_, err := client.runCommand(cliCommand{
		binary: "fdbbackup",
		args: []string{
			"resume",
		},
	})
	return err
}

// ModifyBackup updates the backup parameters.
func (client *cliAdminClient) ModifyBackup(snapshotPeriodSeconds int) error {
	_, err := client.runCommand(cliCommand{
		binary: "fdbbackup",
		args: []string{
			"modify",
			"-s",
			fmt.Sprintf("%d", snapshotPeriodSeconds),
		},
	})
	return err
}

// GetBackupStatus gets the status of the current backup.
func (client *cliAdminClient) GetBackupStatus() (*fdbv1beta2.FoundationDBLiveBackupStatus, error) {
	statusString, err := client.runCommand(cliCommand{
		binary: "fdbbackup",
		args: []string{
			"status",
			"--json",
		},
	})

	if err != nil {
		return nil, err
	}

	status := &fdbv1beta2.FoundationDBLiveBackupStatus{}
	statusString, err = internal.RemoveWarningsInJSON(statusString)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(statusString), &status)
	if err != nil {
		return nil, err
	}

	return status, nil
}
