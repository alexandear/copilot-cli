// Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"testing"

	climocks "github.com/aws/amazon-ecs-cli-v2/internal/pkg/cli/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const githubRepo = "https://github.com/badGoose/chaOS"
const githubToken = "hunter2"

func TestInitPipelineOpts_Ask(t *testing.T) {
	testCases := map[string]struct {
		inEnvironments      []string
		inGitHubRepo        string
		inGitHubAccessToken string
		inProjectEnvs       []string

		mockPrompt func(m *climocks.Mockprompter)

		expectedGitHubRepo        string
		expectedGitHubAccessToken string
		expectedEnvironments      []string
		expectedError             error
	}{
		"prompts for all input": {
			inEnvironments:      []string{},
			inGitHubRepo:        "",
			inGitHubAccessToken: "",
			inProjectEnvs:       []string{"test", "prod"},

			mockPrompt: func(m *climocks.Mockprompter) {
				m.EXPECT().Confirm(pipelineAddEnvPrompt, gomock.Any()).Return(true, nil).Times(3)
				m.EXPECT().SelectOne(pipelineSelectEnvPrompt, gomock.Any(), []string{"test", "prod"}).Return("test", nil).Times(1)
				m.EXPECT().SelectOne(pipelineSelectEnvPrompt, gomock.Any(), []string{"prod"}).Return("prod", nil).Times(1)

				m.EXPECT().Get(gomock.Eq(pipelineEnterGitHubRepoPrompt), gomock.Any(), gomock.Any()).Return(githubRepo, nil).Times(1)
				m.EXPECT().GetSecret(gomock.Eq("Please enter your GitHub Personal Access Token for your repository: https://github.com/badGoose/chaOS"), gomock.Any()).Return(githubToken, nil).Times(1)
			},

			expectedGitHubRepo:        githubRepo,
			expectedGitHubAccessToken: githubToken,
			expectedEnvironments:      []string{"test", "prod"},
			expectedError:             nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPrompt := climocks.NewMockprompter(ctrl)

			opts := &InitPipelineOpts{
				Environments:      tc.inEnvironments,
				GitHubRepo:        tc.inGitHubRepo,
				GitHubAccessToken: tc.inGitHubAccessToken,

				prompt: mockPrompt,

				projectEnvs: tc.inProjectEnvs,
			}

			tc.mockPrompt(mockPrompt)

			// WHEN
			err := opts.Ask()

			// THEN
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expectedGitHubRepo, opts.GitHubRepo)
				require.Equal(t, tc.expectedGitHubRepo, opts.GitHubRepo)
				require.Equal(t, tc.expectedGitHubAccessToken, opts.GitHubAccessToken)
				require.ElementsMatch(t, tc.expectedEnvironments, opts.Environments)
			}
		})
	}
}

func TestInitPipelineOpts_Validate(t *testing.T) {
	testCases := map[string]struct {
		inProjectEnvs []string
		inProjectName string

		expectedError error
	}{
		"errors if no environments initialized": {
			inProjectEnvs: []string{},
			inProjectName: "badgoose",

			expectedError: errNoEnvsInProject,
		},

		"invalid project name": {
			inProjectName: "",
			expectedError: errNoProjectInWorkspace,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			opts := &InitPipelineOpts{
				projectEnvs: tc.inProjectEnvs,

				globalOpts: globalOpts{projectName: tc.inProjectName},
			}

			// WHEN
			err := opts.Validate()

			// THEN
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestInitPipelineOpts_Execute(t *testing.T) {
	testCases := map[string]struct {
		inProjectEnvs []string
		inProjectName string

		expectedError error
	}{
		"errors if no environments initialized": {
			inProjectEnvs: []string{},
			inProjectName: "badgoose",

			expectedError: errNoEnvsInProject,
		},

		"invalid project name": {
			inProjectName: "",
			expectedError: errNoProjectInWorkspace,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			opts := &InitPipelineOpts{
				projectEnvs: tc.inProjectEnvs,

				globalOpts: globalOpts{projectName: tc.inProjectName},
			}

			// WHEN
			err := opts.Execute()

			// THEN
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
