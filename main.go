package main

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

var Config struct {
	KeepRunning bool `short:"k" long:"keep-running" env:"KEEP_RUNNING" description:"Should we do this once or do it every $SLEEP_INTERVAL seconds"`
	Every       int  `short:"n" long:"sleep" env:"SLEEP_INTERVAL" default:"60" description:"How long between loops should we sleep before rerunning"`

	Tags                map[string]string `long:"tag" env:"TAGS" env-delim:"," required:"1" description:"Tag to apply to the spot fleet instances in the form 'k:v' (repeat --tag to set multiple)"`
	SpotFleetRequestIds []string          `long:"id" env:"SPOT_FLEET_REQUEST_IDS" env-delim:"," required:"1" description:"SpotFleetRequestId used to find the instances that need tagging (repeat --id to target multiple)"`
}

func main() {
	_, err := flags.Parse(&Config)
	if err != nil {
		return
	}
	log.Printf("Tagging SpotFleetRequests: %q with %q", Config.SpotFleetRequestIds, Config.Tags)

	log.Println("Setting up AWS Client")
	ec2Client := ec2.New(session.New(), &aws.Config{Region: aws.String("us-west-2")})

	// Isn't a range because I want it to run the first time without sleeping
	sleepInterval := time.Duration(Config.Every) * time.Second
	for t := time.Tick(sleepInterval); ; <-t {
		metrics.iterations.Inc()

		err := tagSpotFleetRequestIds(ec2Client, Config.SpotFleetRequestIds, Config.Tags)
		if err != nil {
			metrics.errors.Inc()
			log.Println(errors.Wrap(err, "failed to tag spot instances"))
		}

		log.Printf("Sleeping for %q", sleepInterval)
		if !Config.KeepRunning {
			break
		}
	}
	log.Println("Quitting")
}
