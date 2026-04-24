package channel

import (
	"log"
	"time"

	"github.com/Hucaru/Valhalla/constant"
)

const (
	heliosElevatorBoardingDuration = time.Minute
	heliosElevatorRideDuration     = time.Minute
)

func scheduleHeliosElevator(server *Server) {
	wait := func(duration time.Duration) {
		timer := time.NewTimer(duration)
		<-timer.C
		timer.Stop()
	}

	departures := map[int32]transportDestination{
		constant.MapHeliosTowerLudiWaitingRoom: {mapID: constant.MapHeliosTowerLudiElevator},
		constant.MapHeliosTowerKFTWaitingRoom:  {mapID: constant.MapHeliosTowerKFTElevator},
	}

	arrivals := map[int32]transportDestination{
		constant.MapHeliosTowerLudiElevator: {mapID: constant.MapHeliosTower99thFloor, portalName: "in00"},
		constant.MapHeliosTowerKFTElevator:  {mapID: constant.MapHeliosTower2ndFloor, portalName: "in00"},
	}

	for {
		log.Println("Helios elevator departing in", heliosElevatorBoardingDuration)
		wait(heliosElevatorBoardingDuration)

		server.dispatch <- func() {
			moveTransportPlayers(server, departures)
		}

		log.Println("Helios elevator arriving in", heliosElevatorRideDuration)
		wait(heliosElevatorRideDuration)

		server.dispatch <- func() {
			moveTransportPlayers(server, arrivals)
		}
	}
}
