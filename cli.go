package main

import (
	"errors"
	"time"
)

type inFlags struct {
	ImageID           string `flag:"image" vardefault:"image-id" description:"Image to launch the EC2 instance from"`
	InstanceName      string `flag:"instance-name" default:"gimme-ec2-instance" description:"Name of the instance for later resume"`
	InstanceType      string `flag:"instance-type" default:"m3.large" description:"Type of the instance to start"`
	KeyName           string `flag:"key-name,k" default:"" description:"SSH key name to access the EC2 instance (must already exist)"`
	MaxSSHWait        string `flag:"ssh-wait" default:"2m" description:"How long to wait for SSH connection to become available"`
	NoShutdown        bool   `flag:"no-shutdown" default:"false" description:"Leave instance running for resuming connection later"`
	Region            string `flag:"region" vardefault:"region" description:"Region to start the EC2 in"`
	SSHPort           int64  `flag:"ssh-port" default:"22" description:"SSH port to use (default is 22)"`
	SSHUser           string `flag:"user,u" default:"ubuntu" description:"User to use for SSH connection"`
	SecurityGroupName string `flag:"security-group" default:"gimme-ec2-security" description:"Name of the EC2 security group to start the instance with"`

	VersionAndExit bool `flag:"version" default:"false" description:"Print version and exit"`
	parsedSSHWait  time.Duration
}

func (i *inFlags) validate() error {
	var err error

	if i.parsedSSHWait, err = time.ParseDuration(cfg.MaxSSHWait); err != nil {
		return errors.New("ssh-wait contained wrong value")
	}

	switch {
	case i.ImageID == "":
		return errors.New("image parameter is required")
	case i.KeyName == "":
		return errors.New("key-name parameter is required")
	case i.Region == "":
		return errors.New("region parameter is required")
	default:
		return nil
	}
}
