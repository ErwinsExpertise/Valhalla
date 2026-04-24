package channel

import (
	"log"
	"time"

	"github.com/Hucaru/Valhalla/constant"
)

const muLungTransportRideDuration = time.Minute

func scheduleMuLungTransport(server *Server) {
	wait := func(duration time.Duration) {
		timer := time.NewTimer(duration)
		<-timer.C
		timer.Stop()
	}

	arrivals := map[int32]transportDestination{
		constant.MapTransportToMuLung: {mapID: constant.MapMuLungArrival},
		constant.MapTransportToOrbis:  {mapID: constant.MapOrbisArrival},
	}

	for {
		log.Println("Mu Lung transport arriving in", muLungTransportRideDuration)
		wait(muLungTransportRideDuration)

		server.dispatch <- func() {
			moveTransportPlayers(server, arrivals)
		}
	}
}
