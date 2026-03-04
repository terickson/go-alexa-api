package main

import (
	"strings"
	"testing"
)

func TestFamilyRoomConfig(t *testing.T) {
	if FamilyRoom.Name != "Family Room" {
		t.Errorf("expected name 'Family Room', got %q", FamilyRoom.Name)
	}
	if !strings.Contains(FamilyRoom.TVActionHost, "192.168.72.20") {
		t.Errorf("unexpected TV host: %s", FamilyRoom.TVActionHost)
	}
	if !strings.Contains(FamilyRoom.RokuActionHost, "family-room") {
		t.Errorf("unexpected Roku host: %s", FamilyRoom.RokuActionHost)
	}
	if FamilyRoom.ReceiverHost == "" {
		t.Error("FamilyRoom should have a receiver")
	}
	if !strings.Contains(FamilyRoom.ReceiverHost, "8081") {
		t.Errorf("unexpected receiver host: %s", FamilyRoom.ReceiverHost)
	}
	if FamilyRoom.DefaultVolume != -30 {
		t.Errorf("expected default volume -30, got %d", FamilyRoom.DefaultVolume)
	}
	if len(FamilyRoom.InputMap) == 0 {
		t.Error("FamilyRoom InputMap should not be empty")
	}
}

func TestMasterBedroomConfig(t *testing.T) {
	if MasterBedroom.Name != "Master Bedroom" {
		t.Errorf("expected name 'Master Bedroom', got %q", MasterBedroom.Name)
	}
	if !strings.Contains(MasterBedroom.TVActionHost, "192.168.72.25") {
		t.Errorf("unexpected TV host: %s", MasterBedroom.TVActionHost)
	}
	if !strings.Contains(MasterBedroom.RokuActionHost, "master-bedroom") {
		t.Errorf("unexpected Roku host: %s", MasterBedroom.RokuActionHost)
	}
	if MasterBedroom.ReceiverHost != "" {
		t.Errorf("MasterBedroom should not have a receiver, got %q", MasterBedroom.ReceiverHost)
	}
	if len(MasterBedroom.InputMap) == 0 {
		t.Error("MasterBedroom InputMap should not be empty")
	}
}

func TestNoStaleIPAddresses(t *testing.T) {
	staleIP := "192.168.72.91"
	hosts := []string{
		FamilyRoom.TVActionHost,
		FamilyRoom.RokuActionHost,
		FamilyRoom.ReceiverHost,
		MasterBedroom.TVActionHost,
		MasterBedroom.RokuActionHost,
	}
	for _, host := range hosts {
		if strings.Contains(host, staleIP) {
			t.Errorf("found stale IP %s in host %q", staleIP, host)
		}
	}
}
