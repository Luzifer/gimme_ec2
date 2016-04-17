package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cenkalti/backoff"

	"github.com/Luzifer/rconfig"
)

const (
	ubuntuReleaseTableURL = "https://cloud-images.ubuntu.com/locator/ec2/releasesTable"
	defaultRegion         = "eu-west-1"
	defaultUbuntuVersion  = "wily"
)

var (
	cfg     = &inFlags{}
	version = "dev"

	ec2Client *ec2.EC2
)

func init() {
	imageID, err := findRecentUbuntuAMI()
	if err != nil {
		log.Printf("Unable to get default parameters: %s", err)
		os.Exit(1)
	}
	rconfig.SetVariableDefaults(map[string]string{
		"image-id": imageID,
		"region":   defaultRegion,
	})
	rconfig.Parse(cfg)

	if err := cfg.validate(); err != nil {
		log.Printf("Input validation failed: %s", err)
		os.Exit(1)
	}

	if cfg.VersionAndExit {
		fmt.Printf("gimme_ec2 %s", version)
		os.Exit(1)
	}

	ec2Client = ec2.New(aws.NewConfig().WithRegion(cfg.Region).WithCredentials(credentials.NewEnvCredentials()))
}

func main() {
	securityGroupID, err := ensureSecurityGroup()
	if err != nil {
		log.Printf("Unable to find / create security group: %s", err)
		os.Exit(1)
	}

	instanceID, dnsName, err := ensureEC2(securityGroupID)
	if err != nil {
		log.Printf("An error occurred while starting the instance: %s\nPlease look into your AWS account!", err)
		os.Exit(1)
	}

	log.Printf("Started instance %s with hostname %s, trying to open SSH connection now...", instanceID, dnsName)

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.MaxElapsedTime = cfg.parsedSSHWait

	err = backoff.Retry(func() error {
		cmd := exec.Command("/usr/bin/ssh",
			"-o", "PasswordAuthentication=no",
			"-o", "StrictHostKeyChecking=no",
			"-p", strconv.FormatInt(cfg.SSHPort, 10),
			cfg.SSHUser+"@"+dnsName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}, backoffConfig)

	if err != nil {
		log.Printf("No SSH connection possible or ended not gracefully, leaving instance active")
		os.Exit(1)
	}

	if cfg.NoShutdown {
		log.Printf("You decided to leave the instance running, keep in mind it generates cost.")
		log.Printf("Start again without --no-shutdown to terminate the instance afterwards")
		os.Exit(0)
	} else {
		log.Printf("SSH connection exitted gracefully, now shutting down instance")
		if err := shutdownEC2(instanceID); err != nil {
			log.Printf("Unable to terminate instance: %s", err)
			os.Exit(1)
		}
	}
}
