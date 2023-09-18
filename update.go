package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/viper"
)

func doUpdate() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		log.Fatal(err)
	}
	network := viper.GetString(FLAG_NETWORK)
	ip, err := lookup(network)
	if err != nil {
		log.Fatal("lookup", err)
	}
	log.Printf("discovered IP: %s", ip)
	recordSetType := "A"
	if strings.Compare(network, "tcp6") == 0 {
		recordSetType = "AAAA"
	}
	svc := route53.New(sess)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(viper.GetString(FLAG_DOMAIN)),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
						TTL:  aws.Int64(viper.GetInt64(FLAG_TTL)),
						Type: aws.String(recordSetType),
					},
				},
			},
			Comment: aws.String(viper.GetString(FLAG_COMMENT)),
		},
		HostedZoneId: aws.String(viper.GetString(FLAG_HOSTED_ZONE_ID)),
	}
	if viper.GetBool(FLAG_DRY) {
		log.Printf("%v", input)
	} else {
		result, err := svc.ChangeResourceRecordSets(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case route53.ErrCodeNoSuchHostedZone:
					fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
				case route53.ErrCodeNoSuchHealthCheck:
					fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
				case route53.ErrCodeInvalidChangeBatch:
					fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
				case route53.ErrCodeInvalidInput:
					fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
				case route53.ErrCodePriorRequestNotComplete:
					fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			os.Exit(2)
		}

		log.Printf("finished Request: %s", *result.ChangeInfo.Id)
	}
	return nil
}
