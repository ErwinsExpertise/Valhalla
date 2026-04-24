package channel

import (
	"log"
	"time"

	"github.com/Hucaru/Valhalla/constant"
)

const (
	subwayWaitingDuration = 5 * time.Minute
	subwayRideDuration    = time.Minute
)

func scheduleSubway(server *Server) {
	wait := func(duration time.Duration) {
		timer := time.NewTimer(duration)
		<-timer.C
		timer.Stop()
	}

	departures := map[int32]transportDestination{
		constant.MapNLCToKerningWaitingRoom: {mapID: constant.MapNLCToKerningTrain},
		constant.MapKerningToNLCWaitingRoom: {mapID: constant.MapKerningToNLCTrain},
	}

	arrivals := map[int32]transportDestination{
		constant.MapNLCToKerningTrain: {mapID: constant.MapKerningSubwayStation},
		constant.MapKerningToNLCTrain: {mapID: constant.MapNLCSubwayStation},
	}

	for {
		log.Println("Subway departing in", subwayWaitingDuration)
		wait(subwayWaitingDuration)

		server.dispatch <- func() {
			moveTransportPlayers(server, departures)
		}

		log.Println("Subway arriving in", subwayRideDuration)
		wait(subwayRideDuration)

		server.dispatch <- func() {
			moveTransportPlayers(server, arrivals)
		}
	}
}
