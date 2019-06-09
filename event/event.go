package event

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

var (
	timeZone = "Asia/Tokyo"
	loc      = time.FixedZone(timeZone, 9*60*60)

	contextTimeout = 15 * time.Second
)

func oneOnOne() *schema.Resource {
	return &schema.Resource{
		Create: eventCreate("1on1"),
		Read:   eventRead,
		Update: eventUpdate("1on1"),
		Delete: eventDelete,
		Schema: eventSchema,
	}
}

func nomikai() *schema.Resource {
	return &schema.Resource{
		Create: eventCreate("飲み会"),
		Read:   eventRead,
		Update: eventUpdate("飲み会"),
		Delete: eventDelete,
		Schema: eventSchema,
	}
}

func tsuribori() *schema.Resource {
	return &schema.Resource{
		Create: eventCreate("釣り堀"),
		Read:   eventRead,
		Update: eventUpdate("釣り堀"),
		Delete: eventDelete,
		Schema: eventSchema,
	}
}
