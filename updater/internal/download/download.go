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
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	archiver "github.com/mholt/archiver/v3"
)

const extractFolder = "latest"

// Downloads the file into the outPath, truncating the file if it already exists
func downloadFile(downloadURL string, outPath string) error {
	//#nosec G107 HTTP request must be dynamic based on input
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("could not GET url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("got non-200 status code (%d)", resp.StatusCode)
	}

	outPathClean := filepath.Clean(outPath)
	f, err := os.OpenFile(outPathClean, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Default().Printf("Failed to close file: %s", err.Error())
		}
	}()

	if _, err = io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("failed to copy request body to file: %w", err)
	}

	return nil
}

// getOutputFilePath gets the output path relative to the base dir for the archive from the given URL.
func getOutputFilePath(basePath, downloadURL string) (string, error) {
	url, err := url.Parse(downloadURL)
	if err != nil {
		return "", fmt.Errorf("cannot parse url: %w", err)
	}

	if url.Path == "" {
		return "", errors.New("input url must have path")
	}

	return filepath.Join(basePath, filepath.Base(url.Path)), nil
}

func verifyContentHash(contentPath, hexExpectedContentHash string) error {
	expectedContentHash, err := hex.DecodeString(hexExpectedContentHash)
	if err != nil {
		return fmt.Errorf("failed to decode content hash: %w", err)
	}

	// Hash file at contentPath using sha256
	fileHash := sha256.New()
	contentPathClean := filepath.Clean(contentPath)

	f, err := os.Open(contentPathClean)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Default().Printf("Failed to close file: %s", err.Error())
		}
	}()

	if _, err = io.Copy(fileHash, f); err != nil {
		return fmt.Errorf("failed to calculate file hash: %w", err)
	}

	actualContentHash := fileHash.Sum(nil)
	if subtle.ConstantTimeCompare(expectedContentHash, actualContentHash) == 0 {
		return errors.New("content hashes were not equal")
	}

	return nil
}

// FetchAndExtractArchive fetches the archive at the specified URL, placing it into dir.
// It then checks to see if it matches the "expectedHash", a hex-encoded string representing the expected sha256 sum of the file.
// If it matches, the archive is extracted into the $dir/latest directory.
// If the archive cannot be extracted, downloaded, or verified, then an error is returned.
func FetchAndExtractArchive(url, dir, expectedHash string) error {
	archiveFilePath, err := getOutputFilePath(dir, url)
	if err != nil {
		return fmt.Errorf("failed to determine archive download path: %w", err)
	}

	if err := downloadFile(url, archiveFilePath); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	extractPath := filepath.Join(dir, extractFolder)

	if err := verifyContentHash(archiveFilePath, expectedHash); err != nil {
		return fmt.Errorf("content hash could not be verified: %w", err)
	}

	// Clean the "latest" dir before extraction
	if err := os.RemoveAll(extractPath); err != nil {
		return fmt.Errorf("error cleaning archive extraction target path: %w", err)
	}

	if err := archiver.Unarchive(archiveFilePath, extractPath); err != nil {
		return fmt.Errorf("failed to extract file: %w", err)
	}

	return nil
}
