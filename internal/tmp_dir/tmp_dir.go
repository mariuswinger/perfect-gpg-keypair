package tmpdir

import (
	"os"
	"path/filepath"
	"time"

	logger "github.com/sirupsen/logrus"

	userinfo "perfect-gpg-keypair/internal/state/user_info"
)

type TmpDir struct {
	parent                   string
	name                     string
	parametersFileName       string
	statusFileName           string
	exportedKeysDirName      string
	RevocationCertFileName   string
	PublicMasterKeyFileName  string
	PrivateMasterKeyFileName string
	SigningSubkeyFileName    string
}

func NewTmpDir(debug bool) TmpDir {
	name := "gpg_key_generator_debug"
	if !debug {
		name = time.Now().Local().Format("06-02-01-15-04")
	}
	return TmpDir{
		parent:                   os.TempDir(),
		name:                     name,
		parametersFileName:       "parameters",
		statusFileName:           "status",
		exportedKeysDirName:      "keys",
		RevocationCertFileName:   ".revocation-certification.asc",
		PublicMasterKeyFileName:  ".public-master.gpg",
		PrivateMasterKeyFileName: ".private-master.gpg",
		SigningSubkeyFileName:    ".signing-subkey.gpg",
	}
}

func (tmpDir TmpDir) Create() error {
	logger.Debugf("Creating temporary directory at '%s'", tmpDir.Path())
	for _, dir := range []string{tmpDir.Path(), tmpDir.ExportedKeysDirPath()} {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (tmpDir TmpDir) Path() string {
	return filepath.Join(tmpDir.parent, tmpDir.name)
}

func (tmpDir TmpDir) ParametersFilePath() string {
	return filepath.Join(tmpDir.Path(), tmpDir.parametersFileName)
}

func (tmpDir TmpDir) StatusFilePath() string {
	return filepath.Join(tmpDir.Path(), tmpDir.statusFileName)
}

func (tmpDir TmpDir) RevocationCertFilePath() string {
	return filepath.Join(tmpDir.ExportedKeysDirPath(), tmpDir.RevocationCertFileName)
}

func (tmpDir TmpDir) PublicMasterKeyFilePath() string {
	return filepath.Join(tmpDir.ExportedKeysDirPath(), tmpDir.PublicMasterKeyFileName)
}

func (tmpDir TmpDir) PrivateMasterKeyFilePath() string {
	return filepath.Join(tmpDir.ExportedKeysDirPath(), tmpDir.PrivateMasterKeyFileName)
}

func (tmpDir TmpDir) SigningSubkeyFilePath() string {
	return filepath.Join(tmpDir.ExportedKeysDirPath(), tmpDir.SigningSubkeyFileName)
}

func (tmpDir TmpDir) ExportedKeysDirPath() string {
	return filepath.Join(tmpDir.Path(), tmpDir.exportedKeysDirName)
}

func (tmpDir TmpDir) CreateParametersFile(userInfo userinfo.UserInfo) error {
	return ParametersFile{Path: tmpDir.ParametersFilePath()}.Create(userInfo)
}

func (tmpDir TmpDir) ReadStatusFileKeyId() (string, error) {
	return StatusFile{Path: tmpDir.StatusFilePath()}.ReadCreatedKeyId()
}
