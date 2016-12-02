[![Build Status](https://travis-ci.org/jboelter/ddr53.png?branch=master)](https://travis-ci.org/jboelter/ddr53)

ddr53 provides a tool to update an Amazon Route53 resource record in the spirit of dyndns. I use this on a Ubiquiti EdgeLite router to maintain an entry for my external IP that runs in `/etc/dhcp/dhclient-exit-hooks.d/` as a hook.

It uses the Amazon Go SDK (https://github.com/aws/aws-sdk-go) and relies on the credentials sourced by the SDK.

Dependencies

    go get -u github.com/aws/aws-sdk-go/aws
    or
    gvt restore

QuickStart:

    AWS_ACCESS_KEY_ID=xxx AWS_SECRET_ACCESS_KEY=zzz ddr53
    --zoneid ZA9XF3OWSDQP1
    --type A
    --value 192.168.1.1
    --fqdn foo.example.com
    --ttl 300

Output:

    2016/06/19 18:34:52 DynDnsRoute53 Tool v1.0
    2016/06/19 18:34:52 Updating foo.example.com in ZA9XF3OWSDQP1 to 192.168.1.1 with type=A and ttl=300
    2016/06/19 18:34:54 {
    ChangeInfo: {
      Id: "/change/C1RAZW3Y13BPG4",
      Status: "PENDING",
      SubmittedAt: 2016-06-20 01:34:53.975 +0000 UTC
      }
    }

Create a limited AWS IAM user and assign the policy below (changing the Hosted Zone ID) to limit the account.

    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "route53:ChangeResourceRecordSets"
          ],
          "Resource": [
            "arn:aws:route53:::hostedzone/ZA9XF3OWSDQP1"
          ]
        }
      ]
    }

Example /etc/dhcp/dhclient-exit-hooks.d/ script
````
user@host:/etc/dhcp/dhclient-exit-hooks.d
> cat dyndns-route53.sh
#!/bin/bash

IP="$(/sbin/ifconfig eth0 | grep 'inet addr' | cut -d: -f2 | awk '{print $1}')"

echo $IP

AWS_ACCESS_KEY_ID=key AWS_SECRET_ACCESS_KEY=secret /config/scripts/ddr53 --zoneid ZA9XF3OWSDQP1 --type A --value ${IP} --fqdn foo.example.com >> /var/log/ddr53.log 2>&1

logger ${0##*/}: 'dyndns-route53 attempted to update foo.example.com to ' $IP
````
