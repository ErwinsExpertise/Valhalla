var maps = [920010000, 920010100, 920010200, 920010300, 920010400, 920010500, 920010600, 920010700, 920010800, 920010900, 920011000, 920011100, 920011300];
var entryMapID = 920010000;
var exitMapID = 920011200;
var finishMapID = 920011300;

function start() {
    ctrl.setDuration("60m");

    for (var i = 0; i < maps.length; i++) {
        var field = ctrl.getMap(maps[i]);
        field.reset();
        field.clearProperties();
    }

    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (var j = 0; j < players.length; j++) {
        players[j].warp(entryMapID);
        players[j].showCountdown(time);
    }
}

function timeout(plr) {
    plr.warp(exitMapID);
}

function playerLeaveEvent(plr) {
    ctrl.removePlayer(plr);
    plr.warp(exitMapID);

    if (plr.isPartyLeader() || ctrl.playerCount() < 1) {
        var players = ctrl.players();
        for (var i = 0; i < players.length; i++) {
            players[i].warp(exitMapID);
        }
        ctrl.finished();
    }
}
