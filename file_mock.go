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
// 			MakeFilePublicFunc: func(fileId string) (string, error) {
// 				panic("mock out the MakeFilePublic method")
// 			},
// 			UploadFileFunc: func(file []byte, folderPath string, fileName string) (string, error) {
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

	// MakeFilePublicFunc mocks the MakeFilePublic method.
	MakeFilePublicFunc func(fileId string) (string, error)

	// UploadFileFunc mocks the UploadFile method.
	UploadFileFunc func(file []byte, folderPath string, fileName string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetPageURL holds details about calls to the GetPageURL method.
		GetPageURL []struct {
		}
		// MakeFilePublic holds details about calls to the MakeFilePublic method.
		MakeFilePublic []struct {
			// FileId is the fileId argument value.
			FileId string
		}
		// UploadFile holds details about calls to the UploadFile method.
		UploadFile []struct {
			// File is the file argument value.
			File []byte
			// FolderPath is the folderPath argument value.
			FolderPath string
			// FileName is the fileName argument value.
			FileName string
		}
	}
	lockGetPageURL     sync.RWMutex
	lockMakeFilePublic sync.RWMutex
	lockUploadFile     sync.RWMutex
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

// MakeFilePublic calls MakeFilePublicFunc.
func (mock *IHubspotFileAPIMock) MakeFilePublic(fileId string) (string, error) {
	if mock.MakeFilePublicFunc == nil {
		panic("IHubspotFileAPIMock.MakeFilePublicFunc: method is nil but IHubspotFileAPI.MakeFilePublic was just called")
	}
	callInfo := struct {
		FileId string
	}{
		FileId: fileId,
	}
	mock.lockMakeFilePublic.Lock()
	mock.calls.MakeFilePublic = append(mock.calls.MakeFilePublic, callInfo)
	mock.lockMakeFilePublic.Unlock()
	return mock.MakeFilePublicFunc(fileId)
}

// MakeFilePublicCalls gets all the calls that were made to MakeFilePublic.
// Check the length with:
//     len(mockedIHubspotFileAPI.MakeFilePublicCalls())
func (mock *IHubspotFileAPIMock) MakeFilePublicCalls() []struct {
	FileId string
} {
	var calls []struct {
		FileId string
	}
	mock.lockMakeFilePublic.RLock()
	calls = mock.calls.MakeFilePublic
	mock.lockMakeFilePublic.RUnlock()
	return calls
}

// UploadFile calls UploadFileFunc.
func (mock *IHubspotFileAPIMock) UploadFile(file []byte, folderPath string, fileName string) (string, error) {
	if mock.UploadFileFunc == nil {
		panic("IHubspotFileAPIMock.UploadFileFunc: method is nil but IHubspotFileAPI.UploadFile was just called")
	}
	callInfo := struct {
		File       []byte
		FolderPath string
		FileName   string
	}{
		File:       file,
		FolderPath: folderPath,
		FileName:   fileName,
	}
	mock.lockUploadFile.Lock()
	mock.calls.UploadFile = append(mock.calls.UploadFile, callInfo)
	mock.lockUploadFile.Unlock()
	return mock.UploadFileFunc(file, folderPath, fileName)
}

// UploadFileCalls gets all the calls that were made to UploadFile.
// Check the length with:
//     len(mockedIHubspotFileAPI.UploadFileCalls())
func (mock *IHubspotFileAPIMock) UploadFileCalls() []struct {
	File       []byte
	FolderPath string
	FileName   string
} {
	var calls []struct {
		File       []byte
		FolderPath string
		FileName   string
	}
	mock.lockUploadFile.RLock()
	calls = mock.calls.UploadFile
	mock.lockUploadFile.RUnlock()
	return calls
}
