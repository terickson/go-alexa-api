package main

import "testing"

func TestFRInputAliases(t *testing.T) {
	tests := []struct {
		alias string
		want  InputConfig
	}{
		// TV
		{"TV", InputConfig{ReceiverInput: "AV1", TVInput: "InputTV"}},
		{"T", InputConfig{ReceiverInput: "AV1", TVInput: "InputTV"}},
		{"V", InputConfig{ReceiverInput: "AV1", TVInput: "InputTV"}},
		// RetroPi
		{"RETROPI", InputConfig{ReceiverInput: "AV1", TVInput: "HDMI2"}},
		{"RETROPIE", InputConfig{ReceiverInput: "AV1", TVInput: "HDMI2"}},
		{"RETRO", InputConfig{ReceiverInput: "AV1", TVInput: "HDMI2"}},
		{"PI", InputConfig{ReceiverInput: "AV1", TVInput: "HDMI2"}},
		// PS3
		{"PS3", InputConfig{ReceiverInput: "HDMI4", TVInput: "HDMI1"}},
		{"PSTHREE", InputConfig{ReceiverInput: "HDMI4", TVInput: "HDMI1"}},
		{"3", InputConfig{ReceiverInput: "HDMI4", TVInput: "HDMI1"}},
		// PS4
		{"PS4", InputConfig{ReceiverInput: "HDMI2", TVInput: "HDMI1"}},
		{"PSFOUR", InputConfig{ReceiverInput: "HDMI2", TVInput: "HDMI1"}},
		{"4", InputConfig{ReceiverInput: "HDMI2", TVInput: "HDMI1"}},
		// PS5
		{"PS5", InputConfig{ReceiverInput: "AV1", TVInput: "HDMI3"}},
		{"PSFIVE", InputConfig{ReceiverInput: "AV1", TVInput: "HDMI3"}},
		{"5", InputConfig{ReceiverInput: "AV1", TVInput: "HDMI3"}},
		// WiiU
		{"WIIU", InputConfig{ReceiverInput: "HDMI3", TVInput: "HDMI1"}},
		{"WE", InputConfig{ReceiverInput: "HDMI3", TVInput: "HDMI1"}},
		// FireTV/Roku
		{"FIRETV", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1"}},
		{"FIRE", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1"}},
		{"ROKU", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1"}},
		// Switch
		{"SWITCH", InputConfig{ReceiverInput: "HDMI5", TVInput: "HDMI1"}},
		// Xbox
		{"XBOX", InputConfig{ReceiverInput: "V-AUX", TVInput: "HDMI1"}},
		// Roku apps
		{"NETFLIX", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Netflix"}},
		{"NET", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Netflix"}},
		{"FLIX", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Netflix"}},
		{"PLEX", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Plex"}},
		{"PLAQUES", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Plex"}},
		{"PRIME", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Prime Video"}},
		{"AMAZON", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Prime Video"}},
		{"HBO", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "HBO GO"}},
		{"YOUTUBE", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "YouTube"}},
		{"HGTV", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Watch HGTV"}},
		{"STARS", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "STARZ"}},
		{"PBS", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "PBS Video"}},
		{"SMITHSONIAN", InputConfig{ReceiverInput: "HDMI1", TVInput: "HDMI1", RokuApp: "Smithsonian Channel"}},
	}

	for _, tt := range tests {
		cfg, ok := FamilyRoom.InputMap[tt.alias]
		if !ok {
			t.Errorf("alias %q not found in FamilyRoom InputMap", tt.alias)
			continue
		}
		if cfg != tt.want {
			t.Errorf("alias %q: got %+v, want %+v", tt.alias, cfg, tt.want)
		}
	}
}

func TestMBInputAliases(t *testing.T) {
	tests := []struct {
		alias string
		want  InputConfig
	}{
		// TV
		{"TV", InputConfig{TVInput: "InputTV"}},
		{"T", InputConfig{TVInput: "InputTV"}},
		{"V", InputConfig{TVInput: "InputTV"}},
		// PS2
		{"PS2", InputConfig{TVInput: "InputAV1"}},
		{"TWO", InputConfig{TVInput: "InputAV1"}},
		{"2", InputConfig{TVInput: "InputAV1"}},
		{"PS", InputConfig{TVInput: "InputAV1"}},
		// Wii
		{"WII", InputConfig{TVInput: "InputComponent1"}},
		{"WI", InputConfig{TVInput: "InputComponent1"}},
		{"WEE", InputConfig{TVInput: "InputComponent1"}},
		// Switch
		{"SWITCH", InputConfig{TVInput: "HDMI2"}},
		// Roku apps
		{"NETFLIX", InputConfig{TVInput: "HDMI1", RokuApp: "Netflix"}},
		{"NET", InputConfig{TVInput: "HDMI1", RokuApp: "Netflix"}},
		{"FLIX", InputConfig{TVInput: "HDMI1", RokuApp: "Netflix"}},
		{"PLEX", InputConfig{TVInput: "HDMI1", RokuApp: "Plex"}},
		{"PRIME", InputConfig{TVInput: "HDMI1", RokuApp: "Prime Video"}},
		{"HBO", InputConfig{TVInput: "HDMI1", RokuApp: "HBO GO"}},
		{"YOUTUBE", InputConfig{TVInput: "HDMI1", RokuApp: "YouTube"}},
	}

	for _, tt := range tests {
		cfg, ok := MasterBedroom.InputMap[tt.alias]
		if !ok {
			t.Errorf("alias %q not found in MasterBedroom InputMap", tt.alias)
			continue
		}
		if cfg != tt.want {
			t.Errorf("alias %q: got %+v, want %+v", tt.alias, cfg, tt.want)
		}
	}
}

func TestUnknownInputFallback(t *testing.T) {
	_, ok := FamilyRoom.InputMap["NONEXISTENT"]
	if ok {
		t.Error("expected NONEXISTENT to not be in FamilyRoom InputMap")
	}

	_, ok = MasterBedroom.InputMap["NONEXISTENT"]
	if ok {
		t.Error("expected NONEXISTENT to not be in MasterBedroom InputMap")
	}
}
