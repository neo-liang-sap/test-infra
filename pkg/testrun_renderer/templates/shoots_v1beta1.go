// Copyright 2019 Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
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

package templates

import (
	"fmt"
	"strconv"

	"github.com/gardener/test-infra/pkg/apis/testmachinery/v1beta1"
	"github.com/gardener/test-infra/pkg/common"
	"github.com/gardener/test-infra/pkg/util"
)

func stepCreateShootV1beta1(cloudprovider common.CloudProvider, name string, dependencies []string, cfg *CreateShootConfig) ([]*v1beta1.DAGStep, string, error) {
	stepConfig := defaultShootConfig(cfg)
	var generatorStep *v1beta1.DAGStep
	switch cloudprovider {
	case common.CloudProviderAWS:
		if name == "" {
			name = "create-shoot-aws"
		}
		generatorStep, stepConfig = v1beta1AWSShootConfig(name, dependencies, stepConfig)
	case common.CloudProviderGCP:
		if name == "" {
			name = "create-shoot-gcp"
		}
		generatorStep, stepConfig = v1beta1GCPShootConfig(name, dependencies, stepConfig)
	case common.CloudProviderAzure:
		if name == "" {
			name = "create-shoot-azure"
		}
		generatorStep, stepConfig = v1beta1AzureShootConfig(name, dependencies, stepConfig)
	default:
		return []*v1beta1.DAGStep{}, "", fmt.Errorf("unsupported cloudprovider %s", cloudprovider)
	}

	return []*v1beta1.DAGStep{
		generatorStep,
		{
			Name: name,
			Definition: v1beta1.StepDefinition{
				Name:   "create-shoot",
				Config: stepConfig,
			},
			UseGlobalArtifacts: false,
			DependsOn:          []string{generatorStep.Name},
			ArtifactsFrom:      "",
			Annotations:        nil,
		},
	}, name, nil
}

var defaultProviderConfig = []v1beta1.ConfigElement{
	{
		Type:  v1beta1.ConfigTypeEnv,
		Name:  ConfigControlplaneProviderPathName,
		Value: ConfigControlplaneProviderPath,
	},
	{
		Type:  v1beta1.ConfigTypeEnv,
		Name:  ConfigInfrastructureProviderPathName,
		Value: ConfigInfrastructureProviderPath,
	},
}

func defaultShootConfig(cfg *CreateShootConfig) []v1beta1.ConfigElement {
	var defaultShootConfig = []v1beta1.ConfigElement{
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigShootName,
			Value: cfg.ShootName,
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigProjectNamespaceName,
			Value: cfg.Namespace,
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigK8sVersionName,
			Value: cfg.K8sVersion,
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigSeedName,
			Value: ConfigSeedValue,
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigShootAnnotations,
			Value: util.MarshalMap(cfg.ShootAnnotations),
		},
	}

	if cfg.AllowPrivilegedContainers != nil {
		ce := v1beta1.ConfigElement{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigAllowPrivilegedContainers,
			Value: strconv.FormatBool(*cfg.AllowPrivilegedContainers),
		}
		defaultShootConfig = append(defaultShootConfig, ce)
	}

	return defaultShootConfig
}

func v1beta1GCPShootConfig(name string, dependencies []string, cfg []v1beta1.ConfigElement) (*v1beta1.DAGStep, []v1beta1.ConfigElement) {

	step := &v1beta1.DAGStep{
		Name: fmt.Sprintf("%s-gen", name),
		Definition: v1beta1.StepDefinition{
			Name: "gen-provider-gcp",
			Config: append(defaultProviderConfig,
				v1beta1.ConfigElement{
					Type:  v1beta1.ConfigTypeEnv,
					Name:  ConfigZoneName,
					Value: "europe-west1-b",
				},
			),
		},
		UseGlobalArtifacts: false,
		DependsOn:          dependencies,
		ArtifactsFrom:      "",
		Annotations:        nil,
	}

	return step, append(cfg, []v1beta1.ConfigElement{
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigCloudproviderName,
			Value: "gcp",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigProviderTypeName,
			Value: "gcp",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigCloudprofileName,
			Value: "gcp",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigSecretBindingName,
			Value: "core-gcp-gcp",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigRegionName,
			Value: "europe-west1",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigZoneName,
			Value: "europe-west1-b",
		},
	}...)
}

func v1beta1AWSShootConfig(name string, dependencies []string, cfg []v1beta1.ConfigElement) (*v1beta1.DAGStep, []v1beta1.ConfigElement) {

	step := &v1beta1.DAGStep{
		Name: fmt.Sprintf("%s-gen", name),
		Definition: v1beta1.StepDefinition{
			Name: "gen-provider-aws",
			Config: append(defaultProviderConfig,
				v1beta1.ConfigElement{
					Type:  v1beta1.ConfigTypeEnv,
					Name:  ConfigZoneName,
					Value: "eu-west-1b",
				},
			),
		},
		UseGlobalArtifacts: false,
		DependsOn:          dependencies,
		ArtifactsFrom:      "",
		Annotations:        nil,
	}

	return step, append(cfg, []v1beta1.ConfigElement{
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigCloudproviderName,
			Value: "aws",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigProviderTypeName,
			Value: "aws",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigCloudprofileName,
			Value: "aws",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigSecretBindingName,
			Value: "core-aws-aws",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigRegionName,
			Value: "eu-west-1",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigZoneName,
			Value: "eu-west-1b",
		},
	}...)
}

func v1beta1AzureShootConfig(name string, dependencies []string, cfg []v1beta1.ConfigElement) (*v1beta1.DAGStep, []v1beta1.ConfigElement) {

	step := &v1beta1.DAGStep{
		Name: fmt.Sprintf("%s-gen", name),
		Definition: v1beta1.StepDefinition{
			Name:   "gen-provider-azure",
			Config: defaultProviderConfig,
		},
		UseGlobalArtifacts: false,
		DependsOn:          dependencies,
		ArtifactsFrom:      "",
		Annotations:        nil,
	}

	return step, append(cfg, []v1beta1.ConfigElement{
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigCloudproviderName,
			Value: "azure",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigProviderTypeName,
			Value: "azure",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigCloudprofileName,
			Value: "azure",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigSecretBindingName,
			Value: "core-azure-azure",
		},
		{
			Type:  v1beta1.ConfigTypeEnv,
			Name:  ConfigRegionName,
			Value: "westeurope",
		},
	}...)
}
