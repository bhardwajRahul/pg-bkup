/*
 *  MIT License
 *
 * Copyright (c) 2024 Jonas Kaninda
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package pkg

import (
	"fmt"
	"github.com/jkaninda/go-storage/pkg/s3"
	goutils "github.com/jkaninda/go-utils"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/pg-bkup/utils"

	"os"
	"path/filepath"
	"time"
)

func s3Backup(db *dbConfig, config *BackupConfig) {

	logger.Info("Backup database to s3 storage")
	// Backup database
	err := BackupDatabase(db, config.backupFileName, disableCompression, config.all, config.allInOne)
	if err != nil {
		recoverMode(err, "Error backing up database")
		return
	}
	finalFileName := config.backupFileName
	if config.encryption {
		encryptBackup(config)
		finalFileName = fmt.Sprintf("%s.%s", config.backupFileName, "gpg")
	}
	logger.Info("Uploading backup archive to remote storage S3 ... ")
	awsConfig := initAWSConfig()
	if config.remotePath == "" {
		config.remotePath = awsConfig.remotePath
	}
	logger.Info(fmt.Sprintf("Backup name is %s", finalFileName))
	s3Storage, err := s3.NewStorage(s3.Config{
		Endpoint:       awsConfig.endpoint,
		Bucket:         awsConfig.bucket,
		AccessKey:      awsConfig.accessKey,
		SecretKey:      awsConfig.secretKey,
		Region:         awsConfig.region,
		DisableSsl:     awsConfig.disableSsl,
		ForcePathStyle: awsConfig.forcePathStyle,
		RemotePath:     config.remotePath,
		LocalPath:      tmpPath,
	})
	if err != nil {
		logger.Fatal("Error creating s3 storage", "error", err)
	}
	err = s3Storage.Copy(finalFileName)
	if err != nil {
		logger.Fatal("Error uploading backup file", "error", err)
	}
	// Get backup info
	fileInfo, err := os.Stat(filepath.Join(tmpPath, finalFileName))
	if err != nil {
		logger.Error("Error getting backup info", "error", err)
	}
	backupSize = fileInfo.Size()

	// Delete backup file from tmp folder
	err = utils.DeleteFile(filepath.Join(tmpPath, config.backupFileName))
	if err != nil {
		fmt.Println("Error deleting file: ", err)

	}
	// Delete old backup
	if config.prune {
		err := s3Storage.Prune(config.backupRetention)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Error deleting old backup from %s storage: %s ", config.storage, err))
		}
	}

	duration := goutils.FormatDuration(time.Since(startTime), 2)
	logger.Info("Backup file uploaded to  S3 storage", "file", finalFileName, "destination", storagePath)
	logger.Info("Backup completed", "file", finalFileName, "size", goutils.ConvertBytes(uint64(backupSize)), "duration", duration)
	// Send notification
	utils.NotifySuccess(&utils.NotificationData{
		File:           finalFileName,
		BackupSize:     goutils.ConvertBytes(uint64(backupSize)),
		Database:       db.dbName,
		Storage:        string(config.storage),
		BackupLocation: filepath.Join(config.remotePath, finalFileName),
		Duration:       duration,
	})
	// Delete temp
	deleteTemp()
	logger.Info(fmt.Sprintf("The backup of the %s database has been completed in %s", db.dbName, duration))

}
func s3Restore(db *dbConfig, conf *RestoreConfig) {
	logger.Info("Restore database from s3")
	awsConfig := initAWSConfig()
	if conf.remotePath == "" {
		conf.remotePath = awsConfig.remotePath
	}
	s3Storage, err := s3.NewStorage(s3.Config{
		Endpoint:       awsConfig.endpoint,
		Bucket:         awsConfig.bucket,
		AccessKey:      awsConfig.accessKey,
		SecretKey:      awsConfig.secretKey,
		Region:         awsConfig.region,
		DisableSsl:     awsConfig.disableSsl,
		ForcePathStyle: awsConfig.forcePathStyle,
		RemotePath:     conf.remotePath,
		LocalPath:      tmpPath,
	})
	if err != nil {
		logger.Fatal("Error creating s3 storage", "error", err)
	}
	err = s3Storage.CopyFrom(conf.file)
	if err != nil {
		logger.Fatal("Error download file from S3 storage", "error", err)
	}
	RestoreDatabase(db, conf)
}
