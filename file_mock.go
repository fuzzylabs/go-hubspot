// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package go_hubspot

import (
	"sync"
)

// Ensure, that IHubspotFileAPIMock does implement IHubspotFileAPI.
// If this is not the case, regenerate this file with moq.
var _ IHubspotFileAPI = &IHubspotFileAPIMock{}

// IHubspotFileAPIMock is a mock implementation of IHubspotFileAPI.
//
// 	func TestSomethingThatUsesIHubspotFileAPI(t *testing.T) {
//
// 		// make and configure a mocked IHubspotFileAPI
// 		mockedIHubspotFileAPI := &IHubspotFileAPIMock{
// 			GetPageURLFunc: func() string {
// 				panic("mock out the GetPageURL method")
// 			},
// 			UploadFileFunc: func(file string, folderPath string, fileName string, options string) error {
// 				panic("mock out the UploadFile method")
// 			},
// 		}
//
// 		// use mockedIHubspotFileAPI in code that requires IHubspotFileAPI
// 		// and then make assertions.
//
// 	}
type IHubspotFileAPIMock struct {
	// GetPageURLFunc mocks the GetPageURL method.
	GetPageURLFunc func() string

	// UploadFileFunc mocks the UploadFile method.
	UploadFileFunc func(file string, folderPath string, fileName string, options string) error

	// calls tracks calls to the methods.
	calls struct {
		// GetPageURL holds details about calls to the GetPageURL method.
		GetPageURL []struct {
		}
		// UploadFile holds details about calls to the UploadFile method.
		UploadFile []struct {
			// File is the file argument value.
			File string
			// FolderPath is the folderPath argument value.
			FolderPath string
			// FileName is the fileName argument value.
			FileName string
			// Options is the options argument value.
			Options string
		}
	}
	lockGetPageURL sync.RWMutex
	lockUploadFile sync.RWMutex
}

// GetPageURL calls GetPageURLFunc.
func (mock *IHubspotFileAPIMock) GetPageURL() string {
	if mock.GetPageURLFunc == nil {
		panic("IHubspotFileAPIMock.GetPageURLFunc: method is nil but IHubspotFileAPI.GetPageURL was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetPageURL.Lock()
	mock.calls.GetPageURL = append(mock.calls.GetPageURL, callInfo)
	mock.lockGetPageURL.Unlock()
	return mock.GetPageURLFunc()
}

// GetPageURLCalls gets all the calls that were made to GetPageURL.
// Check the length with:
//     len(mockedIHubspotFileAPI.GetPageURLCalls())
func (mock *IHubspotFileAPIMock) GetPageURLCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetPageURL.RLock()
	calls = mock.calls.GetPageURL
	mock.lockGetPageURL.RUnlock()
	return calls
}

// UploadFile calls UploadFileFunc.
func (mock *IHubspotFileAPIMock) UploadFile(file string, folderPath string, fileName string, options string) error {
	if mock.UploadFileFunc == nil {
		panic("IHubspotFileAPIMock.UploadFileFunc: method is nil but IHubspotFileAPI.UploadFile was just called")
	}
	callInfo := struct {
		File       string
		FolderPath string
		FileName   string
		Options    string
	}{
		File:       file,
		FolderPath: folderPath,
		FileName:   fileName,
		Options:    options,
	}
	mock.lockUploadFile.Lock()
	mock.calls.UploadFile = append(mock.calls.UploadFile, callInfo)
	mock.lockUploadFile.Unlock()
	return mock.UploadFileFunc(file, folderPath, fileName, options)
}

// UploadFileCalls gets all the calls that were made to UploadFile.
// Check the length with:
//     len(mockedIHubspotFileAPI.UploadFileCalls())
func (mock *IHubspotFileAPIMock) UploadFileCalls() []struct {
	File       string
	FolderPath string
	FileName   string
	Options    string
} {
	var calls []struct {
		File       string
		FolderPath string
		FileName   string
		Options    string
	}
	mock.lockUploadFile.RLock()
	calls = mock.calls.UploadFile
	mock.lockUploadFile.RUnlock()
	return calls
}