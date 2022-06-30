// Copyright  observIQ, Inc.
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

package download

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadFile(t *testing.T) {
	t.Run("Downloads File Over HTTP", func(t *testing.T) {
		tmpDir := t.TempDir()

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Invalid request method: %s", r.Method)
				return
			}

			w.Write([]byte("Hello"))
		}))
		defer s.Close()

		outPath := filepath.Join(tmpDir, "out.txt")

		err := downloadFile(s.URL, outPath)
		require.NoError(t, err)

		b, err := os.ReadFile(outPath)
		require.NoError(t, err)
		assert.Equal(t, []byte("Hello"), b)
	})

	t.Run("Output file is existing directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Invalid request method: %s", r.Method)
				return
			}

			w.Write([]byte("Hello"))
		}))
		defer s.Close()

		err := downloadFile(s.URL, tmpDir)
		require.ErrorContains(t, err, "failed to open file:")
	})

	t.Run("Invalid URL", func(t *testing.T) {
		tmpDir := t.TempDir()
		outPath := filepath.Join(tmpDir, "out.txt")

		err := downloadFile("http://localhost:9999999", outPath)
		require.ErrorContains(t, err, "could not GET url")
	})

	t.Run("Server returns 404", func(t *testing.T) {
		tmpDir := t.TempDir()

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer s.Close()

		outPath := filepath.Join(tmpDir, "out.txt")

		err := downloadFile(s.URL, outPath)
		require.ErrorContains(t, err, "got non-200 status code (404)")
	})
}

func TestGetOutputFilePath(t *testing.T) {
	testCases := []struct {
		name        string
		basepath    string
		url         string
		out         string
		expectedErr string
	}{
		{
			name:     "Input url is valid zip",
			basepath: "/tmp/observiq-otel-collector-update",
			url:      "http://example.com/some-file.zip",
			out:      "/tmp/observiq-otel-collector-update/some-file.zip",
		},
		{
			name:     "Input url is valid tar",
			basepath: "/tmp/observiq-otel-collector-update",
			url:      "http://example.com/some-file.tar.gz",
			out:      "/tmp/observiq-otel-collector-update/some-file.tar.gz",
		},
		{
			name:        "Input url is invalid",
			basepath:    "/tmp/observiq-otel-collector-update",
			url:         "http://local\thost/some-file.zip",
			expectedErr: "cannot parse url",
		},
		{
			name:        "Input url has no path",
			basepath:    "/tmp/observiq-otel-collector-update",
			url:         "http://example.com",
			expectedErr: "input url must have path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := getOutputFilePath(tc.basepath, tc.url)
			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.out, out)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}

func TestVerifyContentHash(t *testing.T) {
	testCases := []struct {
		name        string
		contentPath string
		hash        string
		expectedErr string
	}{
		{
			name:        "Content hash matches",
			contentPath: filepath.Join("testdata", "test.txt"),
			hash:        "c87e2ca771bab6024c269b933389d2a92d4941c848c52f155b9b84e1f109fe35",
		},
		{
			name:        "File does not exist",
			contentPath: filepath.Join("testdata", "non-existant-file.txt"),
			hash:        "c87e2ca771bab6024c269b933389d2a92d4941c848c52f155b9b84e1f109fe35",
			expectedErr: "failed to open file",
		},
		{
			name:        "Content hash does not match",
			contentPath: filepath.Join("testdata", "test.txt"),
			hash:        "7e4ead2053637d9fcb7f3316e748becb8af163c6f851446eeef878a994ae5c4b",
			expectedErr: "content hashes were not equal",
		},
		{
			name:        "Content hash is not hex encoded",
			contentPath: filepath.Join("testdata", "test.txt"),
			hash:        "c87e2ca771bab6024c269b933389d2a92d4941c848c52f155b9b84e1f109fe3z",
			expectedErr: "failed to decode content hash:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyContentHash(tc.contentPath, tc.hash)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}

func TestDownloadAndVerifyExtraction(t *testing.T) {
	testCases := []struct {
		name         string
		archivePath  string
		expectedHash string
		expectedErr  string
	}{
		{
			name:         "Download and extracts tar.gz files",
			archivePath:  filepath.Join("testdata", "test.tar.gz"),
			expectedHash: "d3bf2375be7372b34eae9bc16296ce9e40e53f5b79b329e23056c4aaf77eb47c",
		},
		{
			name:         "Download and extracts zip files",
			archivePath:  filepath.Join("testdata", "test.zip"),
			expectedHash: "5594349d022f7f374fa3ee777ded15f4f06a47aa08eec300bd06cdb0d2688fac",
		},
		{
			name:         "Fails to extract non-archive",
			archivePath:  filepath.Join("testdata", "not-actually-tar.tar.gz"),
			expectedHash: "e7045ebfc48a850a8ac2d342c172099f8c937a4265c55cd93cb39908278952b4",
			expectedErr:  "failed to extract file",
		},
		{
			name:         "Hash does not match downloaded hash",
			archivePath:  filepath.Join("testdata", "test.tar.gz"),
			expectedHash: "e7045ebfc48a850a8ac2d342c172099f8c937a4265c55cd93cb39908278952b4",
			expectedErr:  "content hash could not be verified",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				f, err := os.Open(tc.archivePath)
				if err != nil {
					t.Errorf("Failed to open archive for sending over http: %s", err)
				}
				defer f.Close()

				_, err = io.Copy(w, f)
				if err != nil {
					t.Errorf("Failed to copy archive for sending over http: %s", err)
				}
			}))
			defer s.Close()

			err := FetchAndExtractArchive(fmt.Sprintf("%s/%s", s.URL, tc.archivePath), tmpDir, tc.expectedHash)
			if tc.expectedErr == "" {
				require.NoError(t, err)

				// Make sure test.txt exists in the output dir
				expectedBytes, err := os.ReadFile(filepath.Join("testdata", "test.txt"))
				require.NoError(t, err)

				actualBytes, err := os.ReadFile(filepath.Join(tmpDir, extractFolder, "test.txt"))
				require.NoError(t, err)

				require.Equal(t, expectedBytes, actualBytes)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}

func TestDownloadAndVerifyHTTPFailure(t *testing.T) {
	tmpDir := t.TempDir()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer s.Close()

	err := FetchAndExtractArchive(fmt.Sprintf("%s/%s", s.URL, "some-archive.tar.gz"), tmpDir, "")
	require.ErrorContains(t, err, "failed to download file:")
}

func TestDownloadAndVerifyInvalidURL(t *testing.T) {
	tmpDir := t.TempDir()
	err := FetchAndExtractArchive("http://\t/some-archive.tar.gz", tmpDir, "")
	require.ErrorContains(t, err, "failed to determine archive download path:")
}