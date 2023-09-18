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

const (
	FLAG_NETWORK        = "network"
	FLAG_LOOKUP         = "lookup"
	FLAG_DOMAIN         = "domain"
	FLAG_HOSTED_ZONE_ID = "hosted_zone_id"
	FLAG_TTL            = "ttl"
	FLAG_COMMENT        = "comment"
	FLAG_FREQUENCY      = "frequency"
	FLAG_IMMEDIATE      = "immediate"
	FLAG_PROFILE        = "profile"
	FLAG_PROFILE_PORT   = "profile_port"
	FLAG_DRY            = "dry"
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP(FLAG_NETWORK, "n", "tcp4", "tcp4 or tcp6, sets A or AAAA as appropriate")
	viper.BindPFlag(FLAG_NETWORK, rootCmd.PersistentFlags().Lookup(FLAG_NETWORK))
	rootCmd.PersistentFlags().StringP(FLAG_LOOKUP, "l", "ifconfig.co", "which lookup agent to lookup")
	viper.BindPFlag(FLAG_LOOKUP, rootCmd.PersistentFlags().Lookup(FLAG_LOOKUP))
	rootCmd.PersistentFlags().StringP(FLAG_DOMAIN, "d", "", "domain to update")
	viper.BindPFlag(FLAG_DOMAIN, rootCmd.PersistentFlags().Lookup(FLAG_DOMAIN))
	rootCmd.PersistentFlags().StringP(FLAG_HOSTED_ZONE_ID, "z", "", "hosted zone ID")
	viper.BindPFlag(FLAG_HOSTED_ZONE_ID, rootCmd.PersistentFlags().Lookup(FLAG_HOSTED_ZONE_ID))
	rootCmd.PersistentFlags().Int64(FLAG_TTL, 3600, "TTL to set on the Route53 Record")
	viper.BindPFlag(FLAG_TTL, rootCmd.PersistentFlags().Lookup(FLAG_TTL))
	rootCmd.PersistentFlags().String(FLAG_COMMENT, "", "comment to RecordSet")
	viper.BindPFlag(FLAG_COMMENT, rootCmd.PersistentFlags().Lookup(FLAG_COMMENT))
	rootCmd.PersistentFlags().Int64P(FLAG_FREQUENCY, "f", 0, "frequency to refresh. 0 value is once")
	viper.BindPFlag(FLAG_FREQUENCY, rootCmd.PersistentFlags().Lookup(FLAG_FREQUENCY))
	rootCmd.PersistentFlags().Bool(FLAG_IMMEDIATE, true, "immediately start refresh on start")
	viper.BindPFlag(FLAG_IMMEDIATE, rootCmd.PersistentFlags().Lookup(FLAG_IMMEDIATE))
	rootCmd.PersistentFlags().Bool(FLAG_PROFILE, false, "Enable profiling")
	viper.BindPFlag(FLAG_PROFILE, rootCmd.PersistentFlags().Lookup(FLAG_PROFILE))
	rootCmd.PersistentFlags().Int32(FLAG_PROFILE_PORT, 6660, "Port for profiling")
	viper.BindPFlag(FLAG_PROFILE_PORT, rootCmd.PersistentFlags().Lookup(FLAG_PROFILE_PORT))
	rootCmd.PersistentFlags().Bool(FLAG_DRY, false, "dry run")
	viper.BindPFlag(FLAG_DRY, rootCmd.PersistentFlags().Lookup(FLAG_DRY))
}

var validNetworkFlags = []string{"tcp4", "tcp6"}

// var networkFlag = flag.String("network", "tcp4", "tcp4 or tcp6, sets A or AAAA as appropriate")
var validLookupFlags = []string{"ifconfig.co"}

// var lookupFlag = flag.String("lookup", "ifconfig.co", "ifconfig.co for lookup")
// var domainFlag = flag.String("domain", "", "domain to update on Route53")
// var hostedZoneIDFlag = flag.String("hosted-zone-id", "", "hosted zone id on Route53")
// var ttlFlag = flag.Int64(FLAG_TTL, 3600, "dns record ttl")
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

	if containsValue(viper.GetString(FLAG_NETWORK), validNetworkFlags) == false {
		errors = append(errors, fmt.Errorf("Expected network flag to be %s", strings.Join(validNetworkFlags, ",")))
	}
	if containsValue(viper.GetString(FLAG_NETWORK), validLookupFlags) == false {
		errors = append(errors, fmt.Errorf("Expected lookup flag %s", strings.Join(validLookupFlags, ",")))
	}
	if len(viper.GetString(FLAG_DOMAIN)) == 0 {
		errors = append(errors, fmt.Errorf("domain is required"))
	}
	if len(viper.GetString(FLAG_HOSTED_ZONE_ID)) == 0 {
		errors = append(errors, fmt.Errorf("hosted_zone_id is required"))
	}

	return errors
}
