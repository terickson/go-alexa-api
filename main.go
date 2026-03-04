package main

import (
	"os"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

var applications = map[string]interface{}{
	"/echo/mbr": alexa.EchoApplication{
		AppID:    os.Getenv("MBR_APP_ID"),
		OnIntent: handleIntent(MasterBedroom),
		OnLaunch: handleIntent(MasterBedroom),
	},
	"/echo/fr": alexa.EchoApplication{
		AppID:    os.Getenv("FR_APP_ID"),
		OnIntent: handleIntent(FamilyRoom),
		OnLaunch: handleIntent(FamilyRoom),
	},
}

func main() {
	alexa.SetVerifyAWSCerts(false)
	alexa.Run(applications, "8000")
}
