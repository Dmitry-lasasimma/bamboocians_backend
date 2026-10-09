package database

import (
	"bamboocians/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Demo accounts are public on purpose so anyone can explore the platform.
// Every demo user shares DemoPassword.
const DemoPassword = "demo1234"

const (
	DemoOrganizerEmail = "organizer@demo.bamboocians.com"
	DemoTalentEmail    = "talent@demo.bamboocians.com"
	DemoVenueEmail     = "venue@demo.bamboocians.com"
)

// SeedDemo inserts demo users, profiles, events, bookings, a contract and
// messages. It does nothing if the demo organizer already exists.
func SeedDemo() {
	var count int64
	DB.Model(&models.User{}).Where("email = ?", DemoOrganizerEmail).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(DemoPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("demo seed skipped: %v", err)
		return
	}
	pw := string(hash)

	newUser := func(name, email, role string) models.User {
		u := models.User{Name: name, Email: email, Password: pw, Role: role}
		DB.Create(&u)
		return u
	}

	// ── Main demo accounts (one per role) ────────────────────────────────────
	org := newUser("Maya Chen", DemoOrganizerEmail, models.RoleOrganizer)
	DB.Create(&models.OrganizerProfile{
		UserID: org.ID, CompanyName: "Bamboo Events Co.", Location: "Bangkok, Thailand",
		Bio: "Boutique agency planning weddings, launches and festivals.", Phone: "+66 2 000 0000",
	})

	talent := newUser("Leo Rivera", DemoTalentEmail, models.RoleTalent)
	DB.Create(&models.TalentProfile{
		UserID: talent.ID, StageName: "DJ Leo", Genre: "House / Deep House", HourlyRate: 180,
		Location: "Bangkok, Thailand", Skills: "DJ set, weddings, corporate events, own equipment",
		Bio: "10 years behind the decks at rooftop bars, weddings and brand launches.",
	})

	venue := newUser("Siam Riverside Hall", DemoVenueEmail, models.RoleVenue)
	DB.Create(&models.VenueProfile{
		UserID: venue.ID, VenueName: "Siam Riverside Hall", Capacity: 350, HourlyRate: 450,
		Location: "Riverside, Bangkok", Amenities: "Parking, Catering, AV system, Wi-Fi, Terrace",
		Description: "Riverside ballroom with a terrace overlooking the Chao Phraya.",
	})

	// ── Extra marketplace listings ───────────────────────────────────────────
	extraTalents := []models.TalentProfile{
		{StageName: "Aria Strings Quartet", Genre: "Classical", HourlyRate: 320, Location: "Chiang Mai, Thailand",
			Skills: "string quartet, ceremonies, cocktail hour", Bio: "Elegant live strings for ceremonies and dinners."},
		{StageName: "Nok the Comedian", Genre: "Stand-up / MC", HourlyRate: 150, Location: "Bangkok, Thailand",
			Skills: "MC, host, corporate events", Bio: "Bilingual host who keeps any crowd laughing."},
	}
	for i, p := range extraTalents {
		u := newUser(p.StageName, []string{"aria@demo.bamboocians.com", "nok@demo.bamboocians.com"}[i], models.RoleTalent)
		p.UserID = u.ID
		DB.Create(&p)
	}

	extraVenues := []models.VenueProfile{
		{VenueName: "The Bamboo Garden", Capacity: 120, HourlyRate: 220, Location: "Sukhumvit, Bangkok",
			Amenities: "Garden, Bar, Wi-Fi", Description: "Leafy garden venue for intimate parties."},
		{VenueName: "Skyline Rooftop 54", Capacity: 200, HourlyRate: 600, Location: "Silom, Bangkok",
			Amenities: "Rooftop, Bar, Sound system, Lounge", Description: "54th-floor rooftop with city views."},
	}
	for i, p := range extraVenues {
		u := newUser(p.VenueName, []string{"garden@demo.bamboocians.com", "skyline@demo.bamboocians.com"}[i], models.RoleVenue)
		p.UserID = u.ID
		DB.Create(&p)
	}

	// ── Events for the demo organizer ────────────────────────────────────────
	now := time.Now()
	wedding := models.Event{
		OrganizerID: org.ID, Title: "Anna & Tom Wedding", EventType: models.EventTypeWedding,
		Description: "Riverside ceremony and reception for 150 guests.", Location: "Siam Riverside Hall",
		Date: now.AddDate(0, 1, 0), Budget: 25000, Capacity: 150, Status: models.EventStatusPublished,
	}
	launch := models.Event{
		OrganizerID: org.ID, Title: "TechNova Product Launch", EventType: models.EventTypeCorporate,
		Description: "Evening launch party with DJ and cocktails.", Location: "Bangkok",
		Date: now.AddDate(0, 2, 0), Budget: 12000, Capacity: 200, Status: models.EventStatusDraft,
	}
	DB.Create(&wedding)
	DB.Create(&launch)

	for _, g := range []models.Guest{
		{Name: "Anna Lee", Email: "anna@example.com", Status: models.GuestStatusConfirmed},
		{Name: "Tom Becker", Email: "tom@example.com", Status: models.GuestStatusConfirmed},
		{Name: "Priya Shah", Email: "priya@example.com", Status: models.GuestStatusInvited},
		{Name: "Ken Watanabe", Email: "ken@example.com", Status: models.GuestStatusDeclined},
	} {
		g.EventID = wedding.ID
		DB.Create(&g)
	}

	// ── Bookings ─────────────────────────────────────────────────────────────
	// Venue booking already confirmed, with a contract waiting for the venue to sign
	venueBooking := models.Booking{
		EventID: wedding.ID, BookedByID: org.ID, BookedToID: venue.ID, BookingType: models.BookingTypeVenue,
		Price: 3600, Date: wedding.Date, Status: models.BookingStatusConfirmed, Notes: "Ballroom + terrace, 6pm to 2am",
	}
	DB.Create(&venueBooking)
	DB.Create(&models.Contract{
		BookingID: venueBooking.ID, OrganizerSigned: true, Status: models.ContractStatusPending,
		Content: "Siam Riverside Hall provides the ballroom and terrace for \"Anna & Tom Wedding\" from 6pm to 2am for $3,600, including AV and parking.",
	})

	// Talent booking for the wedding, confirmed (contract not created yet: try it from Bookings)
	DB.Create(&models.Booking{
		EventID: wedding.ID, BookedByID: org.ID, BookedToID: talent.ID, BookingType: models.BookingTypeTalent,
		Price: 900, Date: wedding.Date, Status: models.BookingStatusConfirmed, Notes: "5-hour reception set",
	})

	// Talent booking for the launch, still pending: the talent can accept or decline
	DB.Create(&models.Booking{
		EventID: launch.ID, BookedByID: org.ID, BookedToID: talent.ID, BookingType: models.BookingTypeTalent,
		Price: 720, Date: launch.Date, Status: models.BookingStatusPending, Notes: "4-hour set, brand-friendly playlist",
	})

	// ── Messages ─────────────────────────────────────────────────────────────
	for _, m := range []models.Message{
		{SenderID: org.ID, ReceiverID: talent.ID, Content: "Hi Leo! Are you free for the TechNova launch?", Read: true, CreatedAt: now.Add(-3 * time.Hour)},
		{SenderID: talent.ID, ReceiverID: org.ID, Content: "Hi Maya, yes I'm free. Sending over my setup needs.", CreatedAt: now.Add(-2 * time.Hour)},
		{SenderID: org.ID, ReceiverID: venue.ID, Content: "Contract for the wedding is ready for your signature.", CreatedAt: now.Add(-1 * time.Hour)},
	} {
		DB.Create(&m)
	}

	log.Println("demo data seeded")
}
