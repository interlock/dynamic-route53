package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dynamic-route53",
		Short: "Dynamic DNS updated for AWS Route53",
		Long: `dynamic-route53 is a cli application meant to be run periodically or with its own internal
		regular timer. It updates a single AWS Route53 A/AAAA record.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func initConfig() {
	viper.AutomaticEnv()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("network", "n", "tcp4", "tcp4 or tcp6, sets A or AAAA as appropriate")
	viper.BindPFlag("network", rootCmd.PersistentFlags().Lookup("network"))
	rootCmd.PersistentFlags().StringP("lookup", "l", "ifconfig.co", "which lookup agent to lookup")
	viper.BindPFlag("lookup", rootCmd.PersistentFlags().Lookup("lookup"))
	rootCmd.PersistentFlags().StringP("domain", "d", "", "domain to update")
	viper.BindPFlag("domain", rootCmd.PersistentFlags().Lookup("domain"))
	rootCmd.PersistentFlags().StringP("hosted-zone-id", "z", "", "hosted zone ID")
	viper.BindPFlag("hosted-zone-id", rootCmd.PersistentFlags().Lookup("hosted-zone-id"))
	rootCmd.PersistentFlags().Int64("ttl", 3600, "TTL to set on the Route53 Record")
	viper.BindPFlag("ttl", rootCmd.PersistentFlags().Lookup("ttl"))
	rootCmd.PersistentFlags().String("comment", "", "comment to RecordSet")
	viper.BindPFlag("comment", rootCmd.PersistentFlags().Lookup("comment"))
	rootCmd.PersistentFlags().Int64P("frequency", "f", 0, "frequency to refresh. 0 value is once")
	viper.BindPFlag("frequency", rootCmd.PersistentFlags().Lookup("frequency"))
	rootCmd.PersistentFlags().Bool("immediate", true, "immediately start refresh on start")
	viper.BindPFlag("immediate", rootCmd.PersistentFlags().Lookup("immediate"))
	rootCmd.PersistentFlags().Bool("profile", false, "Enable profiling")
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	rootCmd.PersistentFlags().Int32("profile-port", 6660, "Port for profiling")
	viper.BindPFlag("profile-port", rootCmd.PersistentFlags().Lookup("profile-port"))
	rootCmd.PersistentFlags().Bool("dry", false, "dry run")
	viper.BindPFlag("dry", rootCmd.PersistentFlags().Lookup("dry"))
}

var validNetworkFlags = []string{"tcp4", "tcp6"}

// var networkFlag = flag.String("network", "tcp4", "tcp4 or tcp6, sets A or AAAA as appropriate")
var validLookupFlags = []string{"ifconfig.co"}

// var lookupFlag = flag.String("lookup", "ifconfig.co", "ifconfig.co for lookup")
// var domainFlag = flag.String("domain", "", "domain to update on Route53")
// var hostedZoneIDFlag = flag.String("hosted-zone-id", "", "hosted zone id on Route53")
// var ttlFlag = flag.Int64("ttl", 3600, "dns record ttl")
// var commentFlag = flag.String("comment", "", "comment to set on the Route53 Record")
// var refreshTimeFlag = flag.Int64("refreshTime", 3600, "time in seconds to refresh the Route53 Record")
// var immediateRefreshFlag = flag.Bool("immediate-refresh", false, "immediately refresh")

func containsValue(value string, options []string) bool {
	for _, i := range options {
		if strings.Compare(value, i) == 0 {
			return true
		}
	}
	return false
}

func validateFlags() []error {
	var errors []error

	if containsValue(viper.GetString("network"), validNetworkFlags) == false {
		errors = append(errors, fmt.Errorf("Expected network flag to be %s", strings.Join(validNetworkFlags, ",")))
	}
	if containsValue(viper.GetString("lookup"), validLookupFlags) == false {
		errors = append(errors, fmt.Errorf("Expected lookup flag %s", strings.Join(validLookupFlags, ",")))
	}
	if len(viper.GetString("domain")) == 0 {
		errors = append(errors, fmt.Errorf("domain is required"))
	}
	if len(viper.GetString("hosted-zone-id")) == 0 {
		errors = append(errors, fmt.Errorf("hosted-zone-id is required"))
	}

	return errors
}
