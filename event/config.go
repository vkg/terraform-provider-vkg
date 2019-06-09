package event

import (
	"context"
	"fmt"
	"runtime"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

// Config is the structure used to instantiate the Google Calendar provider.
type Config struct {
	calendar *calendar.Service
}

// loadAndValidate loads the application default credentials from the
// environment and creates a client for communicating with Google APIs.
func (c *Config) loadAndValidate() error {
	client, err := google.DefaultClient(context.Background(), calendar.CalendarScope)
	if err != nil {
		return err
	}

	client.Transport = logging.NewTransport("Google", client.Transport)

	calendarSvc, err := calendar.New(client)
	if err != nil {
		return err
	}

	calendarSvc.UserAgent = fmt.Sprintf("(%s %s) Terraform/%s", runtime.GOOS, runtime.GOARCH, terraform.VersionString())
	c.calendar = calendarSvc

	return nil
}
