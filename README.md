[![Download on GoBuilder](http://badge.luzifer.io/v1/badge?title=Download%20on&text=GoBuilder)](https://gobuilder.me/github.com/Luzifer/gimme_ec2)
[![License: Apache v2.0](https://badge.luzifer.io/v1/badge?color=5d79b5&title=license&text=Apache+v2.0)](http://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/gimme_ec2)](https://goreportcard.com/report/github.com/Luzifer/gimme_ec2)

# Luzifer / gimme\_ec2

`gimme_ec2` is a small utility to start up an EC2-Classic instance, SSH into it and shutting it down again after the SSH connection is closed gracefully. The main purpose for me to write it was I sometimes need a machine to do things with a good internet connection or to test things on a plain Linux machine (maybe in the USA or other countries).

Basically this is a "SSH me into a throw-away-instance"-utility

## Features

- By default start newest Ubuntu 15.10 AMI (AMIs are fetched on startup)
- Bring-your-own-AMI: You can start every AMI supporting SSH access
- The utility takes care about creating a security group and the machine, then waits for SSH to become available
- Resume your previous machine based on name matching (use `--no-shutdown` flag)
- By default the utility takes care about removing the instance after you close your connection

## Usage

```bash
# Usage of ./gimme_ec2:
      --image="ami-3455d547": Image to launch the EC2 instance from
      --instance-name="gimme-ec2-instance": Name of the instance for later resume
      --instance-type="m3.large": Type of the instance to start
  -k, --key-name="": SSH key name to access the EC2 instance (must already exist)
      --no-shutdown[=false]: Leave instance running for resuming connection later
      --region="eu-west-1": Region to start the EC2 in
      --security-group="gimme-ec2-security": Name of the EC2 security group to start the instance with
      --ssh-port=22: SSH port to use (default is 22)
      --ssh-wait="2m": How long to wait for SSH connection to become available
  -u, --user="ubuntu": User to use for SSH connection
      --version[=false]: Print version and exit
```

You need to set typical AWS environment variables, by example using [awsenv](https://github.com/Luzifer/awsenv):

```bash
# awsenv run private -- gimme_ec2 -k mykey
2016/04/17 15:37:52 Started instance i-acbc1c20 with hostname ec2-54-228-163-53.eu-west-1.compute.amazonaws.com, trying to open SSH connection now...
Welcome to Ubuntu 15.10 (GNU/Linux 4.2.0-35-generic x86_64)

 * Documentation:  https://help.ubuntu.com/

  Get cloud support with Ubuntu Advantage Cloud Guest:
    http://www.ubuntu.com/business/services/cloud

0 packages can be updated.
0 updates are security updates.


Last login: Sun Apr 17 13:35:56 2016 from 31.18.130.46
To run a command as administrator (user "root"), use "sudo <command>".
See "man sudo_root" for details.

ubuntu@ip-10-12-175-37:~$
```
