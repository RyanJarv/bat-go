// Code generated by MockGen. DO NOT EDIT.
// Source: cbr.go

// Package mock_cbr is a generated GoMock package.
package mock_cbr

import (
	x "."
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreateIssuer mocks base method
func (m *MockClient) CreateIssuer(ctx context.Context, issuer string, maxTokens int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIssuer", ctx, issuer, maxTokens)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIssuer indicates an expected call of CreateIssuer
func (mr *MockClientMockRecorder) CreateIssuer(ctx, issuer, maxTokens interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIssuer", reflect.TypeOf((*MockClient)(nil).CreateIssuer), ctx, issuer, maxTokens)
}

// GetIssuer mocks base method
func (m *MockClient) GetIssuer(ctx context.Context, issuer string) (*x.IssuerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIssuer", ctx, issuer)
	ret0, _ := ret[0].(*x.IssuerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIssuer indicates an expected call of GetIssuer
func (mr *MockClientMockRecorder) GetIssuer(ctx, issuer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIssuer", reflect.TypeOf((*MockClient)(nil).GetIssuer), ctx, issuer)
}

// SignCredentials mocks base method
func (m *MockClient) SignCredentials(ctx context.Context, issuer string, creds []string) (*x.CredentialsIssueResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignCredentials", ctx, issuer, creds)
	ret0, _ := ret[0].(*x.CredentialsIssueResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignCredentials indicates an expected call of SignCredentials
func (mr *MockClientMockRecorder) SignCredentials(ctx, issuer, creds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignCredentials", reflect.TypeOf((*MockClient)(nil).SignCredentials), ctx, issuer, creds)
}

// RedeemCredential mocks base method
func (m *MockClient) RedeemCredential(ctx context.Context, issuer, preimage, signature, payload string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RedeemCredential", ctx, issuer, preimage, signature, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// RedeemCredential indicates an expected call of RedeemCredential
func (mr *MockClientMockRecorder) RedeemCredential(ctx, issuer, preimage, signature, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedeemCredential", reflect.TypeOf((*MockClient)(nil).RedeemCredential), ctx, issuer, preimage, signature, payload)
}
