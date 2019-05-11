package main

import (
	"github.com/Odania-IT/terraless/support"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/gofrs/flock"
	"os"
	"os/user"
	"path/filepath"
)


func writeSessionProfile(credentials sts.Credentials, targetProfile string, region string) {
	awsProfileWriter := AwsProfileWriter{
		awsConfigFile:   getAwsConfigFile(),
		credentialsFile: getCredentialsFile(),
		lock:            flock.New(filepath.Join(os.TempDir(), "terraless-provider-aws.lock")),
	}

	awsProfileWriter.lockAndWriteAwsCredentials(credentials, targetProfile, region)
}

func (pw AwsProfileWriter) lockAndWriteAwsCredentials(credentials sts.Credentials, targetProfile string, region string) {
	defer func() {
		cerr := pw.lock.Unlock()
		if cerr != nil {
			fatal("[Provider: AWS] Failed to unlock Error %s\n", cerr)
		}
	}()

	locked, err := pw.lock.TryLock()

	if err != nil {
		fatal("Failed aquiring lock for updating AWS credentials! %s\n", err)
	}

	if locked {
		pw.writeAwsCredentials(credentials, targetProfile)
		pw.writeAwsConfig(region, targetProfile)

		return
	}

	fatal("AWS credentials lock already locked!")
}

func (pw AwsProfileWriter) writeAwsCredentials(credentials sts.Credentials, targetProfile string) {
	logger.Debug("Loading credentials file %s\n", pw.credentialsFile)
	support.WriteToFileIfNotExists(pw.credentialsFile, "[default]")
	cfg, err := ini.Load(pw.credentialsFile)

	if err != nil {
		fatal("Error loading aws credentials! %s\n", err)
	}

	section := cfg.Section(targetProfile)

	if section == nil {
		section, err = cfg.NewSection(targetProfile)

		if err != nil {
			fatal("Failed creating section in aws credentials file! Error: %s\n", err)
		}
	} else {
		section.DeleteKey("aws_access_key_id")
		section.DeleteKey("aws_secret_access_key")
		section.DeleteKey("aws_session_token")
	}

	writeKeyToSection(section, "aws_access_key_id", *credentials.AccessKeyId)
	writeKeyToSection(section, "aws_secret_access_key", *credentials.SecretAccessKey)
	writeKeyToSection(section, "aws_session_token", *credentials.SessionToken)

	err = cfg.SaveTo(pw.credentialsFile)
	if err != nil {
		fatal("Failed writing config file %s! Error: %s\n", pw.credentialsFile, err)
	}

	logger.Debug("Wrote session token for profile %s\n", targetProfile)
	logger.Debug("Token is valid until: %v\n", credentials.Expiration)
}

func (pw AwsProfileWriter) writeAwsConfig(region string, targetProfile string) {
	logger.Debug("Loading config file %s\n", pw.awsConfigFile)
	support.WriteToFileIfNotExists(pw.awsConfigFile, "[default]")
	cfg, err := ini.Load(pw.awsConfigFile)

	if err != nil {
		fatal("Error loading aws config! %s\n", err)
	}

	section := cfg.Section(targetProfile)

	if section == nil {
		section, err = cfg.NewSection(targetProfile)

		if err != nil {
			fatal("Failed creating section in aws config file! Error: %s\n", err)
		}
	} else {
		section.DeleteKey("region")
	}

	writeKeyToSection(section, "region", region)

	err = cfg.SaveTo(pw.awsConfigFile)
	if err != nil {
		fatal("Failed writing config file %s! Error: %s\n", pw.credentialsFile, err)
	}

	logger.Debug("Wrote aws config section for profile %s\n", targetProfile)
}

func writeKeyToSection(section *ini.Section, key string, val string) {
	_, err := section.NewKey(key, val)

	if err != nil {
		fatal("Failed writting key %s to aws profile section\n", err)
	}
}

func getCredentialsFile() string {
	credentialsPath := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")

	if credentialsPath != "" {
		return credentialsPath
	}

	usr, err := user.Current()
	if err != nil {
		fatal("Error fetching home dir: %s", err)
	}

	return filepath.Join(usr.HomeDir, ".aws", "credentials")
}

func getAwsConfigFile() string {
	usr, err := user.Current()
	if err != nil {
		fatal("Error fetching home dir: %s", err)
	}

	return filepath.Join(usr.HomeDir, ".aws", "config")
}
