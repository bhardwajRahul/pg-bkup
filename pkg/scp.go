package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/jkaninda/pg-bkup/utils"
	"golang.org/x/crypto/ssh"
	"os"
)

func CopyToRemote(fileName, remotePath string) error {
	sshUser := os.Getenv("SSH_USER")
	sshPassword := os.Getenv("SSH_PASSWORD")
	sshHostName := os.Getenv("SSH_HOST_NAME")
	sshPort := os.Getenv("SSH_PORT")
	sshIdentifyFile := os.Getenv("SSH_IDENTIFY_FILE")

	clientConfig, _ := auth.PasswordKey(sshUser, sshPassword, ssh.InsecureIgnoreHostKey())
	if sshIdentifyFile != "" && utils.FileExists(sshIdentifyFile) {
		clientConfig, _ = auth.PrivateKey(sshUser, sshIdentifyFile, ssh.InsecureIgnoreHostKey())

	} else {
		if sshPassword == "" {
			return errors.New("SSH_PASSWORD environment variable is required if SSH_IDENTIFY_FILE is empty\n")
		}
		clientConfig, _ = auth.PasswordKey(sshUser, sshPassword, ssh.InsecureIgnoreHostKey())

	}
	// Create a new SCP client
	client := scp.NewClient(fmt.Sprintf("%s:%s", sshHostName, sshPort), &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		return errors.New("Couldn't establish a connection to the remote server\n")
	}

	// Open a file
	file, _ := os.Open(fileName)

	// Close client connection after the file has been copied
	defer client.Close()
	// Close the file after it has been copied
	defer file.Close()
	// the context can be adjusted to provide time-outs or inherit from other contexts if this is embedded in a larger application.
	err = client.CopyFromFile(context.Background(), *file, remotePath, "0655")
	if err != nil {
		fmt.Println("Error while copying file ")
		return err
	}
	return nil
}

func CopyFromRemote(fileName, remotePath string) error {
	sshUser := os.Getenv("SSH_USER")
	sshPassword := os.Getenv("SSH_PASSWORD")
	sshHostName := os.Getenv("SSH_HOST_NAME")
	sshPort := os.Getenv("SSH_PORT")
	sshIdentifyFile := os.Getenv("SSH_IDENTIFY_FILE")

	clientConfig, _ := auth.PasswordKey(sshUser, sshPassword, ssh.InsecureIgnoreHostKey())
	if sshIdentifyFile != "" && utils.FileExists(sshIdentifyFile) {
		clientConfig, _ = auth.PrivateKey(sshUser, sshIdentifyFile, ssh.InsecureIgnoreHostKey())

	} else {
		if sshPassword == "" {
			return errors.New("SSH_PASSWORD environment variable is required if SSH_IDENTIFY_FILE is empty\n")
		}
		clientConfig, _ = auth.PasswordKey(sshUser, sshPassword, ssh.InsecureIgnoreHostKey())

	}

	// Create a new SCP client
	client := scp.NewClient(fmt.Sprintf("%s:%s", sshHostName, sshPort), &clientConfig)
	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		return errors.New("Couldn't establish a connection to the remote server\n")
	}
	// Close client connection after the file has been copied
	defer client.Close()
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Couldn't open the output file")
	}
	defer file.Close()

	// the context can be adjusted to provide time-outs or inherit from other contexts if this is embedded in a larger application.
	err = client.CopyFromRemote(context.Background(), file, remotePath)

	if err != nil {
		fmt.Println("Error while copying file ", err)
		return err
	}
	return nil

}
