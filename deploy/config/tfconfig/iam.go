// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfconfig

import (
	"encoding/json"
	"fmt"
)

// ProjectIAMMembers represents multiple Terraform project IAM members.
// It is used to wrap and merge multiple IAM members into a single IAM member when being marshalled to JSON.
type ProjectIAMMembers struct {
	Members   []*ProjectIAMMember
	DependsOn []string
	project   string
}

// ProjectIAMMember represents a Terraform project IAM member.
type ProjectIAMMember struct {
	Role   string `json:"role"`
	Member string `json:"member"`

	// The following fields should not be set by users.

	// ForEach is used to let a single iam member expand to reference multiple iam members
	// through the use of terraform's for_each iterator.
	ForEach   map[string]*ProjectIAMMember `json:"for_each,omitempty"`
	Project   string                       `json:"project,omitempty"`
	DependsOn []string                     `json:"depends_on,omitempty"`
}

// Init initializes the resource.
func (ms *ProjectIAMMembers) Init(projectID string) error {
	ms.project = projectID
	return nil
}

// ID returns the resource unique identifier.
// It is hardcoded to return "project" as there is at most one of this resource in a deployment.
func (ms *ProjectIAMMembers) ID() string {
	return "project"
}

// ResourceType returns the resource terraform provider type.
func (ms *ProjectIAMMembers) ResourceType() string {
	return "google_project_iam_member"
}

// MarshalJSON marshals the list of members into a single member.
// The single member will set a for_each block to expand to multiple iam members in the terraform call.
func (ms *ProjectIAMMembers) MarshalJSON() ([]byte, error) {
	forEach := make(map[string]*ProjectIAMMember)
	for _, m := range ms.Members {
		key := fmt.Sprintf("%s %s", m.Role, m.Member)
		forEach[key] = m
	}

	return json.Marshal(&ProjectIAMMember{
		ForEach:   forEach,
		Project:   ms.project,
		Role:      "${each.value.role}",
		Member:    "${each.value.member}",
		DependsOn: ms.DependsOn,
	})
}

// UnmarshalJSON unmarshals the bytes to a list of members.
func (ms *ProjectIAMMembers) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &ms.Members)
}

// ServiceAccount represents a Terraform service account.
type ServiceAccount struct {
	AccountID   string `json:"account_id"`
	Project     string `json:"project"`
	DisplayName string `json:"display_name"`
}

// Init initializes the resource.
func (a *ServiceAccount) Init(projectID string) error {
	if a.Project != "" {
		return fmt.Errorf("project must not be set: %v", a.Project)
	}
	a.Project = projectID
	return nil
}

// ID returns the resource unique identifier.
func (a *ServiceAccount) ID() string {
	return a.AccountID
}

// ResourceType returns the resource terraform provider type.
func (a *ServiceAccount) ResourceType() string {
	return "google_service_account"
}
