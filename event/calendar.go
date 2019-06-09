package event

import (
	"context"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	calendar "google.golang.org/api/calendar/v3"
)

var (
	eventSchema = map[string]*schema.Schema{
		"start": &schema.Schema{Type: schema.TypeString, Required: true, ValidateFunc: datetime()},
		"end":   &schema.Schema{Type: schema.TypeString, Required: true, ValidateFunc: datetime()},

		"location":                    &schema.Schema{Optional: true, Type: schema.TypeString},
		"description":                 &schema.Schema{Optional: true, Type: schema.TypeString},
		"transparency":                &schema.Schema{Optional: true, Type: schema.TypeString, Default: "transparent", ValidateFunc: in([]string{"opaque", "transparent"})},
		"visibility":                  &schema.Schema{Optional: true, Type: schema.TypeString, Default: "public", ValidateFunc: in([]string{"public", "private"})},
		"guests_can_invite_others":    &schema.Schema{Optional: true, Type: schema.TypeBool, Default: true},
		"guests_can_modify":           &schema.Schema{Optional: true, Type: schema.TypeBool, Default: true},
		"guests_can_see_other_guests": &schema.Schema{Optional: true, Type: schema.TypeBool, Default: true},
		"send_notifications":          &schema.Schema{Optional: true, Type: schema.TypeBool, Default: true},
		"attendee": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"email":    &schema.Schema{Type: schema.TypeString, Required: true},
					"optional": &schema.Schema{Type: schema.TypeBool, Optional: true, Default: false},
				},
			},
		},

		"event_id":     &schema.Schema{Type: schema.TypeString, Computed: true},
		"hangout_link": &schema.Schema{Type: schema.TypeString, Computed: true},
		"html_link":    &schema.Schema{Type: schema.TypeString, Computed: true},
	}
)

func eventCreate(name string) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		event := buildCalendarEvent(name, d, meta)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		eventAPI, err := config.calendar.Events.Insert("primary", event).SendNotifications(d.Get("send_notifications").(bool)).MaxAttendees(25).Context(ctx).Do()
		if err != nil {
			return err
		}

		d.SetId(eventAPI.Id)

		return eventRead(d, meta)
	}
}

func eventRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	event, err := config.calendar.Events.Get("primary", d.Id()).Context(ctx).Do()
	if err != nil {
		return err
	}

	attendees := []map[string]interface{}{}
	if len(event.Attendees) > 0 {
		for _, v := range event.Attendees {
			attendees = append(attendees, map[string]interface{}{"email": v.Email, "optional": v.Optional})
		}
	}

	d.Set("summary", event.Summary)
	d.Set("location", event.Location)
	d.Set("description", event.Description)
	d.Set("start", event.Start)
	d.Set("end", event.End)
	if event.GuestsCanInviteOthers != nil {
		d.Set("guests_can_invite_others", *event.GuestsCanInviteOthers)
	}
	d.Set("guests_can_modify", event.GuestsCanModify)
	if event.GuestsCanSeeOtherGuests != nil {
		d.Set("guests_can_see_other_guests", *event.GuestsCanSeeOtherGuests)
	}
	d.Set("transparency", event.Transparency)
	d.Set("visibility", event.Visibility)
	d.Set("attendee", attendees)
	d.Set("event_id", event.Id)
	d.Set("hangout_link", event.HangoutLink)
	d.Set("html_link", event.HtmlLink)

	return nil
}

func eventUpdate(name string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		event := buildCalendarEvent(name, d, meta)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		eventAPI, err := config.calendar.Events.Update("primary", d.Id(), event).SendNotifications(d.Get("send_notifications").(bool)).MaxAttendees(25).Context(ctx).Do()
		if err != nil {
			return err
		}
		d.SetId(eventAPI.Id)
		return eventRead(d, meta)
	}
}

func eventDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	id := d.Id()
	sendNotifications := d.Get("send_notifications").(bool)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := config.calendar.Events.Delete("primary", id).SendNotifications(sendNotifications).Context(ctx).Do()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func ptr(b bool) *bool { return &b }

func buildCalendarEvent(summary string, d *schema.ResourceData, meta interface{}) *calendar.Event {
	// time format is already validated, so err can be ignored
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", d.Get("start").(string), loc)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", d.Get("end").(string), loc)

	attendeesRaw := d.Get("attendee").(*schema.Set)
	// vkgtaro must be invited
	// attendees := []*calendar.EventAttendee{&calendar.EventAttendee{Email: os.Getenv("VKG_EMAIL"), Optional: false}}
	attendees := []*calendar.EventAttendee{}
	if attendeesRaw.Len() > 0 {
		for _, v := range attendeesRaw.List() {
			m := v.(map[string]interface{})
			attendees = append(attendees, &calendar.EventAttendee{Email: m["email"].(string), Optional: m["optional"].(bool)})
		}
	}
	return &calendar.Event{
		Summary:                 summary,
		Start:                   &calendar.EventDateTime{DateTime: start.Format(time.RFC3339)},
		End:                     &calendar.EventDateTime{DateTime: end.Format(time.RFC3339)},
		Location:                d.Get("location").(string),
		Description:             d.Get("description").(string),
		GuestsCanModify:         d.Get("guests_can_modify").(bool),
		Transparency:            d.Get("transparency").(string),
		Visibility:              d.Get("visibility").(string),
		GuestsCanInviteOthers:   ptr(d.Get("guests_can_invite_others").(bool)),
		GuestsCanSeeOtherGuests: ptr(d.Get("guests_can_see_other_guests").(bool)),
		Source:                  &calendar.EventSource{Title: "terraform-provider-vkg", Url: "https://github.com/vkg/terraform-provider-vkg"},
		Attendees:               attendees,
	}
}
