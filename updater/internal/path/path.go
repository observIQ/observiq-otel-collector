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

package path

import "path/filepath"

const (
	latestDirFragment         = "latest"
	rollbackDirFragment       = "rollback"
	serviceFileDirFragment    = "install"
	serviceFileBackupFilename = "backup.service"
)

// LatestDirFromTempDir gets the path to the "latest" dir, where the new artifacts are,
// from the temporary directory
func LatestDirFromTempDir(tmpDir string) string {
	return filepath.Join(tmpDir, latestDirFragment)
}

// BackupDirFromTempDir gets the path to the "rollback" dir, where current artifacts are backed up,
// from the temporary directory
func BackupDirFromTempDir(tmpDir string) string {
	return filepath.Join(tmpDir, rollbackDirFragment)
}

// ServiceFileDir gets the directory of the service file definitions from the install dir
func ServiceFileDir(installBaseDir string) string {
	return filepath.Join(installBaseDir, serviceFileDirFragment)
}

// BackupServiceFile returns the full path to the backup service file from the service file directory path
func BackupServiceFile(serviceFileDir string) string {
	return filepath.Join(serviceFileDir, serviceFileBackupFilename)
}
