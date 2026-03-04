package main

// Room defines the configuration for a room's entertainment system.
type Room struct {
	Name           string
	TVActionHost   string
	RokuActionHost string
	ReceiverHost   string // empty if room has no receiver
	DefaultVolume  int    // receiver default volume on input switch
	InputMap       map[string]InputConfig
}

// FamilyRoom is the configuration for the family room entertainment system.
var FamilyRoom = Room{
	Name:           "Family Room",
	TVActionHost:   "http://192.168.72.20:8080/tv/actions",
	RokuActionHost: "http://192.168.72.222:8080/systems/family-room/actions",
	ReceiverHost:   "http://192.168.72.222:8081/receiver/",
	DefaultVolume:  -30,
	InputMap:       frInputMap(),
}

// MasterBedroom is the configuration for the master bedroom entertainment system.
var MasterBedroom = Room{
	Name:           "Master Bedroom",
	TVActionHost:   "http://192.168.72.25:8080/tv/actions",
	RokuActionHost: "http://192.168.72.222:8080/systems/master-bedroom/actions",
	InputMap:       mbInputMap(),
}

func frInputMap() map[string]InputConfig {
	m := make(map[string]InputConfig)

	// Direct TV inputs (no Roku app)
	addAliases(m, InputConfig{ReceiverInput: "AV1", TVInput: "InputTV"}, "TV", "T", "V")
	addAliases(m, InputConfig{ReceiverInput: "AV1", TVInput: "HDMI2"}, "RETROPI", "RETROPIE", "RETROPOT", "RETROBY", "RETRO", "PIE", "PI")
	addAliases(m, InputConfig{ReceiverInput: "HDMI4", TVInput: "HDMI1"}, "PSTHREE", "PS3", "THREE", "3", "P.S.3")
	addAliases(m, InputConfig{ReceiverInput: "HDMI2", TVInput: "HDMI1"}, "PSFOUR", "PS4", "FOUR", "4", "P.S.4")
	addAliases(m, InputConfig{ReceiverInput: "AV1", TVInput: "HDMI3"}, "PSFIVE", "PS5", "FIVE", "5", "P.S.5")
	addAliases(m, InputConfig{ReceiverInput: "HDMI3", TVInput: "HDMI1"}, "WIIU", "WIYOU", "WILLYOU", "WEYOU", "WEEYOU", "WE", "WEE", "WE'LL")
	addAliases(m, InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1"}, "FIRETV", "FIRE", "ROKU")
	addAliases(m, InputConfig{ReceiverInput: "HDMI5", TVInput: "HDMI1"}, "SWITCH")
	addAliases(m, InputConfig{ReceiverInput: "V-AUX", TVInput: "HDMI1"}, "XBOX")

	// Roku app inputs (all through receiver HDMI1 → TV HDMI1)
	rokuHDMI1 := func(app string) InputConfig {
		return InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: app}
	}
	m["DCUNIVERSE"] = rokuHDMI1("DC Universe")
	m["DAILYBURN"] = rokuHDMI1("Daily Burn")
	addAliases(m, rokuHDMI1("Netflix"), "NETFLIX", "NET", "FLIX")
	addAliases(m, rokuHDMI1("Plex"), "PLEX", "PLAQUES")
	addAliases(m, rokuHDMI1("Prime Video"), "PRIME", "AMAZON")
	m["HBO"] = rokuHDMI1("HBO GO")
	m["CRUNCHYROLL"] = rokuHDMI1("Crunchyroll")
	m["HGTV"] = rokuHDMI1("Watch HGTV")
	m["STARS"] = rokuHDMI1("STARZ")
	m["PBS"] = rokuHDMI1("PBS Video")
	m["SHOWTIME"] = rokuHDMI1("Showtime Anytime")
	m["YOUTUBE"] = rokuHDMI1("YouTube")
	addAliases(m, rokuHDMI1("NatGeoTV"), "NATGEOTV", "NATGEO")
	m["SMITHSONIAN"] = rokuHDMI1("Smithsonian Channel")

	return m
}

func mbInputMap() map[string]InputConfig {
	m := make(map[string]InputConfig)

	// Direct TV inputs (no receiver, no Roku)
	addAliases(m, InputConfig{TVInput: "InputTV"}, "TV", "T", "V")
	addAliases(m, InputConfig{TVInput: "InputAV1"}, "PS2", "TWO", "2", "PSTWO", "PS")
	addAliases(m, InputConfig{TVInput: "InputComponent1"}, "WII", "WI", "WILL", "WE", "WEEK", "WIFI", "WEE", "WE'LL")
	addAliases(m, InputConfig{TVInput: "HDMI2"}, "SWITCH")

	// Roku app inputs (TV HDMI1)
	rokuHDMI1 := func(app string) InputConfig {
		return InputConfig{TVInput: "HDMI1", RokuApp: app}
	}
	m["DAILYBURN"] = rokuHDMI1("Daily Burn")
	addAliases(m, rokuHDMI1("Netflix"), "NETFLIX", "NET", "FLIX")
	addAliases(m, rokuHDMI1("Plex"), "PLEX", "PLAQUES")
	addAliases(m, rokuHDMI1("Prime Video"), "PRIME", "AMAZON")
	m["HBO"] = rokuHDMI1("HBO GO")
	m["CRUNCHYROLL"] = rokuHDMI1("Crunchyroll")
	m["HGTV"] = rokuHDMI1("Watch HGTV")
	m["STARS"] = rokuHDMI1("STARZ")
	m["PBS"] = rokuHDMI1("PBS Video")
	m["SHOWTIME"] = rokuHDMI1("Showtime Anytime")
	m["YOUTUBE"] = rokuHDMI1("YouTube")
	addAliases(m, rokuHDMI1("NatGeoTV"), "NATGEOTV", "NATGEO")
	m["SMITHSONIAN"] = rokuHDMI1("Smithsonian Channel")

	return m
}
