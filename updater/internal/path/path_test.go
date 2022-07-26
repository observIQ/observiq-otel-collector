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

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLatestDirFromTempDir(t *testing.T) {
	require.Equal(t, filepath.Join("tmp", "latest"), LatestDirFromTempDir("tmp"))
}

func TestBackupDirFromTempDir(t *testing.T) {
	require.Equal(t, filepath.Join("tmp", "rollback"), BackupDirFromTempDir("tmp"))
}

func TestServiceFileDir(t *testing.T) {
	installDir := filepath.Join("tmp", "rollback")
	require.Equal(t, filepath.Join(installDir, "install"), ServiceFileDir(installDir))
}

func TestBackupServiceFile(t *testing.T) {
	serviceFileDir := filepath.Join("tmp", "rollback", "install")
	require.Equal(t, filepath.Join(serviceFileDir, "backup.service"), BackupServiceFile(serviceFileDir))
}

func TestLogFile(t *testing.T) {
	installDir := filepath.Join("install")
	require.Equal(t, filepath.Join(installDir, "log", "updater.log"), LogFile(installDir))
}
