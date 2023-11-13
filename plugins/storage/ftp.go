package storage

import (
	"bytes"
	"io"
	"io/fs"
	"mime/multipart"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/ericmarcelinotju/gram/config"
)

// FtpStorage provides an abstraction on top of file manager system
type FtpStorage struct {
	path       string
	ftpManager *BasicSFTPManager
}

func NewFtpStorage(configuration *config.Storage) (*FtpStorage, error) {
	return NewFTPManager(
		configuration.Path,
		configuration.Host,
		configuration.Username,
		configuration.Password,
	)
}

func NewFTPManager(basePath string, host, username, password string) (*FtpStorage, error) {

	var auths []ssh.AuthMethod

	// Try to use $SSH_AUTH_SOCK which contains the path of the unix file socket that the sshd agent uses
	// for communication with other processes.
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}

	// Use password authentication if provided
	if password != "" {
		auths = append(auths, ssh.Password(password))
	}
	// Initialize client configuration
	config := ssh.ClientConfig{
		User: username,
		Auth: auths,
		// Ignore host key check
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// // Connect to server
	// conn, err := ssh.Dial("tcp", f.host, &config)
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to [%s]: %v", f.host, err)
	// }

	// // Create new SFTP client
	// sc, err := sftp.NewClient(conn)
	// if err != nil {
	// 	return fmt.Errorf("unable to start SFTP subsystem: %v", err)
	// }
	// client = sc
	ftpManager := NewBasicSFTPManager(host, &config)

	return &FtpStorage{path: basePath, ftpManager: ftpManager}, nil
}

func (f *FtpStorage) Upload(file multipart.File, fileName string) (err error) {
	var conn *SFTPConn
	var client *sftp.Client

	conn, err = f.ftpManager.GetConnection()
	if err != nil {
		return err
	}
	client = conn.GetClient()

	filePath := filepath.Join(f.path, fileName)
	dir := strings.Replace(filepath.Dir(filePath), "\\", "/", -1)

	// Check path exist
	_, err = client.Stat(dir)
	if err != nil {
		if err := client.MkdirAll(dir); err != nil {
			return err
		}
	}

	// Create the file
	out, err := client.OpenFile(filePath, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, file)
	return err
}

func (f *FtpStorage) Download(fileName string) ([]byte, error) {
	var err error
	var conn *SFTPConn
	var client *sftp.Client

	conn, err = f.ftpManager.GetConnection()
	if err != nil {
		return nil, err
	}
	client = conn.GetClient()

	pwd, err := client.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(pwd, f.path, fileName)
	file, err := client.OpenFile(filePath, (os.O_RDONLY))
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f *FtpStorage) Remove(fileName string) error {
	var err error
	var conn *SFTPConn
	var client *sftp.Client

	conn, err = f.ftpManager.GetConnection()
	if err != nil {
		return err
	}
	client = conn.GetClient()

	pwd, err := client.Getwd()
	if err != nil {
		return err
	}
	filePath := filepath.Join(pwd, f.path, fileName)
	err = client.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func (f *FtpStorage) IsExist(fileName string) bool {
	var err error
	var conn *SFTPConn
	var client *sftp.Client

	conn, err = f.ftpManager.GetConnection()
	if err != nil {
		return false
	}
	client = conn.GetClient()

	filePath := f.path + fileName

	// Check file exist
	stat, err := client.Stat(filePath)

	return err == nil && !stat.IsDir()
}

func (f *FtpStorage) Path() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, f.path)
}

func (f *FtpStorage) List(folder string) ([]fs.FileInfo, error) {
	var err error
	var conn *SFTPConn
	var client *sftp.Client
	var fileInfos []fs.FileInfo = make([]fs.FileInfo, 0)

	conn, err = f.ftpManager.GetConnection()
	if err != nil {
		return nil, err
	}
	client = conn.GetClient()

	pwd, err := client.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(pwd, f.path, folder)
	w := client.Walk(filePath)
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		fileInfos = append(fileInfos, w.Stat())
	}

	return fileInfos, nil
}
