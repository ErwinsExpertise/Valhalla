var level1Maps = [109030001, 109030101, 109030201, 109030301, 109030401];
var level2Maps = [109030002, 109030102, 109030202, 109030302, 109030402];
var level3Maps = [109030003, 109030103, 109030203, 109030303, 109030403];
var exitWinnerMapID = 109050000;
var exitMapID = 109050001;
var ended = false;

function start() {
    ctrl.setDuration("6m");

    var maps = level1Maps.concat(level2Maps, level3Maps);
    for (let i = 0; i < maps.length; i++) {
        var field = ctrl.getMap(maps[i]);
        field.reset();
        field.clearProperties();
    }

    var players = ctrl.players();
    var time = ctrl.remainingTime();

    for (let i = 0; i < players.length; i++) {
        var mapID = level1Maps[i % level1Maps.length];
        players[i].warp(mapID);
        players[i].showCountdown(time);
    }
}

function beforePortal(plr, src, dst) {
    return true;
}

function afterPortal(plr, dst) {
    plr.showCountdown(ctrl.remainingTime());
}

function timeout(plr) {
    plr.warp(exitMapID);
}

function playerLeaveEvent(plr) {
    ctrl.removePlayer(plr);
    plr.warp(exitMapID);

    if (ctrl.playerCount() <= 0) {
        ctrl.finished();
    }
}

function onMapChange(plr, dst) {
    var mapID = dst.getMapID();

    if (level3Maps.indexOf(mapID) !== -1) {
        finishWinner(plr);
        return;
    }

    if (level1Maps.indexOf(mapID) === -1 && level2Maps.indexOf(mapID) === -1 && level3Maps.indexOf(mapID) === -1) {
        ctrl.removePlayer(plr);
        plr.warp(exitMapID);

        if (ctrl.playerCount() <= 0) {
            ctrl.finished();
        }
    }
}

function finishWinner(winner) {
    if (ended) {
        return;
    }
    ended = true;

    var players = ctrl.players();

    for (let i = 0; i < players.length; i++) {
        if (players[i].name() === winner.name()) {
            players[i].warp(exitWinnerMapID);
        } else {
            players[i].warp(exitMapID);
        }
    }

    ctrl.finished();
}
