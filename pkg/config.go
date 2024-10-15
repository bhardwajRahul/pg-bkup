// Package pkg /
/*****
@author    Jonas Kaninda
@license   MIT License <https://opensource.org/licenses/MIT>
@Copyright © 2024 Jonas Kaninda
**/
package pkg

import (
	"fmt"
	"github.com/jkaninda/pg-bkup/utils"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Path     string `yaml:"path"`
}
type Config struct {
	Databases      []Database `yaml:"databases"`
	CronExpression string     `yaml:"cronExpression"`
}

type dbConfig struct {
	dbHost     string
	dbPort     string
	dbName     string
	dbUserName string
	dbPassword string
}
type targetDbConfig struct {
	targetDbHost     string
	targetDbPort     string
	targetDbUserName string
	targetDbPassword string
	targetDbName     string
}
type TgConfig struct {
	Token  string
	ChatId string
}
type BackupConfig struct {
	backupFileName     string
	backupRetention    int
	disableCompression bool
	prune              bool
	remotePath         string
	encryption         bool
	usingKey           bool
	passphrase         string
	publicKey          string
	storage            string
	cronExpression     string
}
type FTPConfig struct {
	host       string
	user       string
	password   string
	port       string
	remotePath string
}

// SSHConfig holds the SSH connection details
type SSHConfig struct {
	user         string
	password     string
	hostName     string
	port         string
	identifyFile string
}
type AWSConfig struct {
	endpoint       string
	bucket         string
	accessKey      string
	secretKey      string
	region         string
	disableSsl     bool
	forcePathStyle bool
}

func initDbConfig(cmd *cobra.Command) *dbConfig {
	//Set env
	utils.GetEnv(cmd, "dbname", "DB_NAME")
	dConf := dbConfig{}
	dConf.dbHost = os.Getenv("DB_HOST")
	dConf.dbPort = os.Getenv("DB_PORT")
	dConf.dbName = os.Getenv("DB_NAME")
	dConf.dbUserName = os.Getenv("DB_USERNAME")
	dConf.dbPassword = os.Getenv("DB_PASSWORD")

	err := utils.CheckEnvVars(dbHVars)
	if err != nil {
		utils.Error("Please make sure all required environment variables for database are set")
		utils.Fatal("Error checking environment variables: %s", err)
	}
	return &dConf
}

func getDatabase(database Database) *dbConfig {
	return &dbConfig{
		dbHost:     database.Host,
		dbPort:     database.Port,
		dbName:     database.Name,
		dbUserName: database.User,
		dbPassword: database.Password,
	}
}

// loadSSHConfig loads the SSH configuration from environment variables
func loadSSHConfig() (*SSHConfig, error) {
	utils.GetEnvVariable("SSH_HOST", "SSH_HOST_NAME")
	sshVars := []string{"SSH_USER", "SSH_HOST", "SSH_PORT", "REMOTE_PATH"}
	err := utils.CheckEnvVars(sshVars)
	if err != nil {
		return nil, fmt.Errorf("error missing environment variables: %w", err)
	}

	return &SSHConfig{
		user:         os.Getenv("SSH_USER"),
		password:     os.Getenv("SSH_PASSWORD"),
		hostName:     os.Getenv("SSH_HOST"),
		port:         os.Getenv("SSH_PORT"),
		identifyFile: os.Getenv("SSH_IDENTIFY_FILE"),
	}, nil
}
func initFtpConfig() *FTPConfig {
	//Initialize data configs
	fConfig := FTPConfig{}
	fConfig.host = utils.GetEnvVariable("FTP_HOST", "FTP_HOST_NAME")
	fConfig.user = os.Getenv("FTP_USER")
	fConfig.password = os.Getenv("FTP_PASSWORD")
	fConfig.port = os.Getenv("FTP_PORT")
	fConfig.remotePath = os.Getenv("REMOTE_PATH")
	err := utils.CheckEnvVars(ftpVars)
	if err != nil {
		utils.Error("Please make sure all required environment variables for FTP are set")
		utils.Fatal("Error missing environment variables: %s", err)
	}
	return &fConfig
}
func initAWSConfig() *AWSConfig {
	//Initialize AWS configs
	aConfig := AWSConfig{}
	aConfig.endpoint = utils.GetEnvVariable("AWS_S3_ENDPOINT", "S3_ENDPOINT")
	aConfig.accessKey = utils.GetEnvVariable("AWS_ACCESS_KEY", "ACCESS_KEY")
	aConfig.secretKey = utils.GetEnvVariable("AWS_SECRET_KEY", "SECRET_KEY")
	aConfig.bucket = utils.GetEnvVariable("AWS_S3_BUCKET_NAME", "BUCKET_NAME")
	aConfig.region = os.Getenv("AWS_REGION")
	disableSsl, err := strconv.ParseBool(os.Getenv("AWS_DISABLE_SSL"))
	if err != nil {
		disableSsl = false
	}
	forcePathStyle, err := strconv.ParseBool(os.Getenv("AWS_FORCE_PATH_STYLE"))
	if err != nil {
		forcePathStyle = false
	}
	aConfig.disableSsl = disableSsl
	aConfig.forcePathStyle = forcePathStyle
	err = utils.CheckEnvVars(awsVars)
	if err != nil {
		utils.Error("Please make sure all required environment variables for AWS S3 are set")
		utils.Fatal("Error checking environment variables: %s", err)
	}
	return &aConfig
}
func initBackupConfig(cmd *cobra.Command) *BackupConfig {
	utils.SetEnv("STORAGE_PATH", storagePath)
	utils.GetEnv(cmd, "cron-expression", "BACKUP_CRON_EXPRESSION")
	utils.GetEnv(cmd, "period", "BACKUP_CRON_EXPRESSION")
	utils.GetEnv(cmd, "path", "REMOTE_PATH")
	//Get flag value and set env
	remotePath := utils.GetEnvVariable("REMOTE_PATH", "SSH_REMOTE_PATH")
	storage = utils.GetEnv(cmd, "storage", "STORAGE")
	backupRetention, _ := cmd.Flags().GetInt("keep-last")
	prune, _ := cmd.Flags().GetBool("prune")
	disableCompression, _ = cmd.Flags().GetBool("disable-compression")
	_, _ = cmd.Flags().GetString("mode")
	passphrase := os.Getenv("GPG_PASSPHRASE")
	_ = utils.GetEnv(cmd, "path", "AWS_S3_PATH")
	cronExpression := os.Getenv("BACKUP_CRON_EXPRESSION")

	publicKeyFile, err := checkPubKeyFile(os.Getenv("GPG_PUBLIC_KEY"))
	if err == nil {
		encryption = true
		usingKey = true
	} else if passphrase != "" {
		encryption = true
		usingKey = false
	}
	//Initialize backup configs
	config := BackupConfig{}
	config.backupRetention = backupRetention
	config.disableCompression = disableCompression
	config.prune = prune
	config.storage = storage
	config.encryption = encryption
	config.remotePath = remotePath
	config.passphrase = passphrase
	config.publicKey = publicKeyFile
	config.usingKey = usingKey
	config.cronExpression = cronExpression
	return &config
}

type RestoreConfig struct {
	s3Path     string
	remotePath string
	storage    string
	file       string
	bucket     string
	usingKey   bool
	passphrase string
	privateKey string
}

func initRestoreConfig(cmd *cobra.Command) *RestoreConfig {
	utils.SetEnv("STORAGE_PATH", storagePath)
	utils.GetEnv(cmd, "path", "REMOTE_PATH")

	//Get flag value and set env
	s3Path := utils.GetEnv(cmd, "path", "AWS_S3_PATH")
	remotePath := utils.GetEnvVariable("REMOTE_PATH", "SSH_REMOTE_PATH")
	storage = utils.GetEnv(cmd, "storage", "STORAGE")
	file = utils.GetEnv(cmd, "file", "FILE_NAME")
	bucket := utils.GetEnvVariable("AWS_S3_BUCKET_NAME", "BUCKET_NAME")
	passphrase := os.Getenv("GPG_PASSPHRASE")
	privateKeyFile, err := checkPrKeyFile(os.Getenv("GPG_PRIVATE_KEY"))
	if err == nil {
		usingKey = true
	} else if passphrase != "" {
		usingKey = false
	}

	//Initialize restore configs
	rConfig := RestoreConfig{}
	rConfig.s3Path = s3Path
	rConfig.remotePath = remotePath
	rConfig.storage = storage
	rConfig.bucket = bucket
	rConfig.file = file
	rConfig.storage = storage
	rConfig.passphrase = passphrase
	rConfig.usingKey = usingKey
	rConfig.privateKey = privateKeyFile
	return &rConfig
}
func initTargetDbConfig() *targetDbConfig {
	tdbConfig := targetDbConfig{}
	tdbConfig.targetDbHost = os.Getenv("TARGET_DB_HOST")
	tdbConfig.targetDbPort = utils.EnvWithDefault("TARGET_DB_PORT", "5432")
	tdbConfig.targetDbName = os.Getenv("TARGET_DB_NAME")
	tdbConfig.targetDbUserName = os.Getenv("TARGET_DB_USERNAME")
	tdbConfig.targetDbPassword = os.Getenv("TARGET_DB_PASSWORD")

	err := utils.CheckEnvVars(tdbRVars)
	if err != nil {
		utils.Error("Please make sure all required environment variables for the target database are set")
		utils.Fatal("Error checking target database environment variables: %s", err)
	}
	return &tdbConfig
}
func loadConfigFile() (string, error) {
	backupConfigFile, err := checkConfigFile(os.Getenv("BACKUP_CONFIG_FILE"))
	if err == nil {
		return backupConfigFile, nil
	}
	return "", fmt.Errorf("backup config file not found")
}
