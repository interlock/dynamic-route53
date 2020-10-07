package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	flag.Parse()
	errs := validateFlags()
	if len(errs) != 0 {
		for _, e := range errs {
			log.Print(e)
		}
		os.Exit(1)
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		log.Fatal(err)
	}
	ip, err := lookup(*networkFlag)
	if err != nil {
		log.Fatal(err)
	}
	svc := route53.New(sess)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(*domainFlag),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
						TTL:  aws.Int64(*ttlFlag),
						Type: aws.String("A"),
					},
				},
			},
			Comment: aws.String(*commentFlag),
		},
		HostedZoneId: aws.String(*hostedZoneIDFlag),
	}
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
		return
	}

	fmt.Println(result)

}
