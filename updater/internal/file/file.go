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

package file

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// CopyFile copies the file from pathIn to pathOut.
// If the file does not exist, it is created. If the file does exist, it is truncated before writing.
func CopyFile(logger *zap.Logger, pathIn, pathOut string, overwrite bool) error {
	pathInClean := filepath.Clean(pathIn)

	// Open the input file for reading.
	inFile, err := os.Open(pathInClean)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() {
		err := inFile.Close()
		if err != nil {
			logger.Info("Failed to close input file", zap.Error(err))
		}
	}()

	pathOutClean := filepath.Clean(pathOut)
	fileMode := fs.FileMode(0600)

	flags := os.O_CREATE | os.O_WRONLY
	if overwrite {
		// If we are OK to overwrite, we will truncate the file on open
		flags |= os.O_TRUNC

		// Save old file's permissions and delete it first if it exists
		fileOutInfo, _ := os.Stat(pathOutClean)
		if fileOutInfo != nil {
			fileMode = fileOutInfo.Mode()
		}
		_ = os.Remove(pathOutClean)
	} else {
		// This flag will make OpenFile error if the file already exists
		flags |= os.O_EXCL
	}

	// Open the output file, creating it if it does not exist and truncating it.
	//#nosec G304 -- out file is cleaned; this is a general purpose copy function
	outFile, err := os.OpenFile(pathOutClean, flags, fileMode)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer func() {
		err := outFile.Close()
		if err != nil {
			logger.Info("Failed to close output file", zap.Error(err))
		}
	}()

	// Copy the input file to the output file.
	if _, err := io.Copy(outFile, inFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
