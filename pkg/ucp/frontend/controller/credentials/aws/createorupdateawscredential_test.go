/*
Copyright 2023 The Radius Authors.

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
package aws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	armrpc_rest "github.com/project-radius/radius/pkg/armrpc/rest"
	"github.com/project-radius/radius/pkg/to"
	"github.com/project-radius/radius/pkg/ucp/api/v20220901privatepreview"
	"github.com/project-radius/radius/pkg/ucp/secret"
	"github.com/project-radius/radius/pkg/ucp/store"
	"github.com/project-radius/radius/test/testutil"

	"github.com/golang/mock/gomock"
	armrpc_controller "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	"github.com/stretchr/testify/require"
)

func Test_AWS_Credential(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStorageClient := store.NewMockStorageClient(mockCtrl)
	mockSecretClient := secret.NewMockClient(mockCtrl)

	credentialCtrl, err := NewCreateOrUpdateAWSCredential(armrpc_controller.Options{StorageClient: mockStorageClient}, mockSecretClient)
	require.NoError(t, err)

	tests := []struct {
		name       string
		filename   string
		headerfile string
		url        string
		expected   armrpc_rest.Response
		fn         func(mockStorageClient store.MockStorageClient, mockSecretClient secret.MockClient)
		err        error
	}{
		{
			name:       "test_credential_creation",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2022-09-01-privatepreview",
			expected:   getAwsResponse(),
			fn:         setupCredentialSuccessMocks,
			err:        nil,
		},
		{
			name:       "test_invalid_version_credential_resource",
			filename:   "aws-credential.json",
			headerfile: testHeaderFileWithBadAPIVersion,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=bad",
			expected:   nil,
			fn:         setupEmptyMocks,
			err:        v1.ErrUnsupportedAPIVersion,
		},
		{
			name:       "test_invalid_credential_request",
			filename:   "invalid-request-aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2022-09-01-privatepreview",
			expected:   nil,
			fn:         setupEmptyMocks,
			err: &v1.ErrModelConversion{
				PropertyName: "$.properties",
				ValidValue:   "not nil",
			},
		},
		{
			name:       "test_credential_created",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2022-09-01-privatepreview",
			expected:   getAwsResponse(),
			fn:         setupCredentialNotFoundMocks,
			err:        nil,
		},
		{
			name:       "test_credential_notFoundError",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2022-09-01-privatepreview",
			fn:         setupCredentialNotFoundErrorMocks,
			err:        errors.New("Error"),
		},
		{
			name:       "test_credential_get_failure",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2022-09-01-privatepreview",
			fn:         setupCredentialGetFailMocks,
			err:        errors.New("Failed Get"),
		},
		{
			name:       "test_credential_secret_save_failure",
			filename:   "aws-credential.json",
			headerfile: testHeaderFile,
			url:        "/planes/aws/awscloud/providers/System.AWS/credentials/default?api-version=2022-09-01-privatepreview",
			fn:         setupCredentialSecretSaveFailMocks,
			err:        errors.New("Secret Save Failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn(*mockStorageClient, *mockSecretClient)

			credentialVersionedInput := &v20220901privatepreview.AWSCredentialResource{}
			credentialInput := testutil.ReadFixture(tt.filename)
			err = json.Unmarshal(credentialInput, credentialVersionedInput)
			require.NoError(t, err)

			request, err := testutil.GetARMTestHTTPRequest(context.Background(), http.MethodPut, tt.headerfile, credentialVersionedInput)
			require.NoError(t, err)

			ctx := testutil.ARMTestContextFromRequest(request)

			response, err := credentialCtrl.Run(ctx, nil, request)
			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, response)
			}
		})
	}

}

func getAwsResponse() armrpc_rest.Response {
	return armrpc_rest.NewOKResponseWithHeaders(&v20220901privatepreview.AWSCredentialResource{
		Location: to.Ptr("West US"),
		ID:       to.Ptr("/planes/aws/awscloud/providers/System.AWS/credentials/default"),
		Name:     to.Ptr("default"),
		Type:     to.Ptr("System.AWS/credentials"),
		Tags: map[string]*string{
			"env": to.Ptr("dev"),
		},
		Properties: &v20220901privatepreview.AWSAccessKeyCredentialProperties{
			AccessKeyID: to.Ptr("00000000-0000-0000-0000-000000000000"),
			Kind:        to.Ptr("AccessKey"),
			Storage: &v20220901privatepreview.InternalCredentialStorageProperties{
				Kind:       to.Ptr(string(v20220901privatepreview.CredentialStorageKindInternal)),
				SecretName: to.Ptr("aws-awscloud-default"),
			},
		},
	}, map[string]string{"ETag": ""})
}

func setupCredentialSuccessMocks(mockStorageClient store.MockStorageClient, mockSecretClient secret.MockClient) {
	mockStorageClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
		return nil, &store.ErrNotFound{}
	})
	mockSecretClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockStorageClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
}

func setupEmptyMocks(mockStorageClient store.MockStorageClient, mockSecretClient secret.MockClient) {
}

func setupCredentialNotFoundMocks(mockStorageClient store.MockStorageClient, mockSecretClient secret.MockClient) {
	mockStorageClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id string, options ...store.GetOptions) (*store.Object, error) {
			return nil, &store.ErrNotFound{}
		}).Times(1)
	mockSecretClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockStorageClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
}

func setupCredentialNotFoundErrorMocks(mockStorageClient store.MockStorageClient, mockSecretClient secret.MockClient) {
	mockStorageClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id string, options ...store.GetOptions) (*store.Object, error) {
			return nil, errors.New("Error")
		}).Times(1)
}

func setupCredentialGetFailMocks(mockStorageClient store.MockStorageClient, mockSecretClient secret.MockClient) {
	mockStorageClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id string, options ...store.GetOptions) (*store.Object, error) {
			return nil, errors.New("Failed Get")
		}).Times(1)
}

func setupCredentialSecretSaveFailMocks(mockStorageClient store.MockStorageClient, mockSecretClient secret.MockClient) {
	mockStorageClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
		DoAndReturn(func(ctx context.Context, id string, options ...store.GetOptions) (*store.Object, error) {
			return nil, &store.ErrNotFound{}
		}).Times(1)
	mockSecretClient.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Secret Save Failure")).Times(1)
}
