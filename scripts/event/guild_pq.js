var waitingMapID = 990000000;
var exitMapID = 990001100;

function start() {
    ctrl.setDuration("60m");

    var field = ctrl.getMap(waitingMapID);
    field.reset();
    field.clearProperties();
    field.properties()["canEnter"] = true;
    field.properties()["leader"] = ctrl.players()[0].name();
    field.properties()["entryTimestamp"] = Date.now().toString();

    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (var i = 0; i < players.length; i++) {
        players[i].warp(waitingMapID);
        players[i].showCountdown(time);
    }
}

function timeout(plr) {
    plr.warp(exitMapID);
}

function playerLeaveEvent(plr) {
    ctrl.removePlayer(plr);
    plr.warp(exitMapID);

    if (plr.isLeader() || ctrl.playerCount() < 6) {
        var players = ctrl.players();
        for (var i = 0; i < players.length; i++) {
            players[i].warp(exitMapID);
        }
        ctrl.finished();
    }
}
