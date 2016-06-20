// Package ddr53 provides a tool to update an Amazon Route53
// resource record in the spirit of dyndns.
//
// It uses the Amazon Go SDK (https://github.com/aws/aws-sdk-go) and
// relies on the credentials sourced by the SDK.
//
// See the readme at github.com/jboelter/ddr53 for a quick start guide
package main

import (
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var (
	zoneID string
	fqdn   string
	rType  string
	rValue string
	ttl    int64
)

func init() {
	flag.StringVar(&zoneID, "zoneid", "", "the Route53 hosted zone ID to update (e.g. ZA9XF3OWSDQP1)")
	flag.StringVar(&fqdn, "fqdn", "", "the FQDN to update (e.g. foo.example.com)")
	flag.StringVar(&rType, "type", "A", "the record type")
	flag.StringVar(&rValue, "value", "", "the record value (e.g. 192.168.1.1)")
	flag.Int64Var(&ttl, "ttl", 300, "the ttl value")

	flag.Parse()
}

func main() {
	log.Println("DynDnsRoute53 Tool v1.0")
	if len(zoneID) == 0 || len(fqdn) == 0 || len(rType) == 0 || len(rValue) == 0 || ttl < 0 {
		flag.Usage()
		os.Exit(2)
	}

	log.Printf("Updating %v in %v to %v with type=%v and ttl=%v\n", fqdn, zoneID, rValue, rType, ttl)

	awsSession := session.New()
	r53 := route53.New(awsSession)

	change, err := r53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(route53.ChangeActionUpsert),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(fqdn),
						Type: aws.String(rType),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(rValue),
							},
						},
						TTL: aws.Int64(ttl),
					},
				},
			},
		},
		HostedZoneId: aws.String(zoneID),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(change)
}
