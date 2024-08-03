package pkg

const s3MountPath string = "/s3mnt"
const s3fsPasswdFile string = "/etc/passwd-s3fs"
const cronLogFile = "/var/log/pg-bkup.log"
const tmpPath = "/tmp/backup"
const backupCronFile = "/usr/local/bin/backup_cron.sh"
const algorithm = "aes256"
const gpgExtension = "gpg"

var (
	storage            = "local"
	file               = ""
	dbPassword         = ""
	dbUserName         = ""
	dbName             = ""
	dbHost             = ""
	dbPort             = "5432"
	executionMode      = "default"
	storagePath        = "/backup"
	disableCompression = false
	encryption         = false
)

// dbHVars Required environment variables for database
var dbHVars = []string{
	"DB_HOST",
	"DB_PASSWORD",
	"DB_USERNAME",
	"DB_NAME",
}

// sshVars Required environment variables for SSH remote server storage
var sshVars = []string{
	"SSH_USER",
	"SSH_REMOTE_PATH",
	"SSH_HOST_NAME",
	"SSH_PORT",
}
