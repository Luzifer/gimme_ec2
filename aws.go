package main

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cenkalti/backoff"
)

func ensureEC2(securityGroupID string) (string, string, error) {
	dio, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(cfg.InstanceName)},
			},
		},
	})

	if err == nil && len(dio.Reservations) == 1 && len(dio.Reservations[0].Instances) == 1 {
		instance := dio.Reservations[0].Instances[0]
		return *instance.InstanceId, *instance.PublicDnsName, nil
	}

	rio, err := ec2Client.RunInstances(&ec2.RunInstancesInput{
		ImageId: aws.String(cfg.ImageID),
		InstanceInitiatedShutdownBehavior: aws.String(ec2.ShutdownBehaviorTerminate),
		InstanceType:                      aws.String(cfg.InstanceType),
		KeyName:                           aws.String(cfg.KeyName),
		MaxCount:                          aws.Int64(1),
		MinCount:                          aws.Int64(1),
		SecurityGroupIds:                  []*string{aws.String(securityGroupID)},
	})

	if err != nil {
		return "", "", err
	}

	if len(rio.Instances) != 1 {
		return "", "", errors.New("Started wrong instance count")
	}

	_, err = ec2Client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{rio.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(cfg.InstanceName),
			},
		},
	})

	if err != nil {
		return "", "", err
	}

	instanceID := *rio.Instances[0].InstanceId
	dnsName := ""

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.MaxElapsedTime = 30 * time.Second

	err = backoff.Retry(func() error {
		dio, boErr := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(instanceID)},
		})

		if boErr != nil {
			return boErr
		}

		switch {
		case len(dio.Reservations) != 1:
			backoffConfig.MaxElapsedTime = 0
			return errors.New("Unable to fetch information about instance")
		case len(dio.Reservations[0].Instances) != 1:
			backoffConfig.MaxElapsedTime = 0
			return errors.New("Reservation contained wrong number of instances")
		case *dio.Reservations[0].Instances[0].PublicDnsName == "":
			return errors.New("No public DNS name was found")
		}

		dnsName = *dio.Reservations[0].Instances[0].PublicDnsName
		return nil
	}, backoffConfig)

	if err != nil {
		return "", "", err
	}

	return instanceID, dnsName, nil
}

func ensureSecurityGroup() (string, error) {
	dsgo, err := ec2Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{aws.String(cfg.SecurityGroupName)},
	})

	if err == nil && len(dsgo.SecurityGroups) > 0 {
		return *dsgo.SecurityGroups[0].GroupId, nil
	}

	csgo, err := ec2Client.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		Description: aws.String("Automatically created security group for gimme_ec2 CLI utility"),
		GroupName:   aws.String(cfg.SecurityGroupName),
	})

	if err != nil {
		return "", err
	}

	_, err = ec2Client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:     aws.String("0.0.0.0/0"),
		FromPort:   aws.Int64(cfg.SSHPort),
		ToPort:     aws.Int64(cfg.SSHPort),
		IpProtocol: aws.String("tcp"),
		GroupName:  aws.String(cfg.SecurityGroupName),
	})

	if err != nil {
		return "", err
	}

	return *csgo.GroupId, nil
}

func shutdownEC2(instanceID string) error {
	_, err := ec2Client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})

	return err
}
