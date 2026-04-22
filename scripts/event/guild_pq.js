var waitingMapID = 990000000;
var exitMapID = 990001100;
var GQItems = [1032033, 4001024, 4001025, 4001026, 4001027, 4001028, 4001031, 4001032, 4001033, 4001034, 4001035, 4001037];

function clearGQItems(plr) {
    for (var i = 0; i < GQItems.length; i++) {
        plr.removeAll(GQItems[i]);
    }
}

function start() {
    ctrl.setDuration("60m");

    var field = ctrl.getMap(waitingMapID);
    field.reset();
    field.clearProperties();
    field.properties()["canEnter"] = true;
    field.properties()["leader"] = ctrl.players()[0].name();
    field.properties()["entryTimestamp"] = Date.now().toString();
    ctrl.schedule("begin", "1m")

    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (var i = 0; i < players.length; i++) {
        players[i].warp(waitingMapID);
        players[i].showCountdown(time);
    }
}

function begin() {
    var field = ctrl.getMap(waitingMapID);
    field.properties()["canEnter"] = false;
    ctrl.schedule("earringcheck", "15s")

    var players = ctrl.players();
    for (var i = 0; i < players.length; i++) {
        players[i].sendMessage("[Guild Quest] The quest has begun!");
    }
}

function earringcheck() {
    var players = ctrl.players();
    for (var i = 0; i < players.length; i++) {
        var plr = players[i];
        if (plr.mapID() > 990000200 && plr.itemCount(1032033) < 1) {
            plr.setHP(1);
            plr.sendMessage("[Guild Quest] You were struck down for entering without the proper earrings.");
        }
    }
    ctrl.schedule("earringcheck", "15s")
}

function timeout(plr) {
    clearGQItems(plr);
    plr.hideCountdown();
    plr.warp(exitMapID);
}

function playerLeaveEvent(plr) {
    ctrl.removePlayer(plr);
    clearGQItems(plr);
    plr.hideCountdown();
    plr.warp(exitMapID);

    if (plr.isLeader() || ctrl.playerCount() < 6) {
        var players = ctrl.players();
        for (var i = 0; i < players.length; i++) {
            clearGQItems(players[i]);
            players[i].hideCountdown();
            players[i].warp(exitMapID);
        }
        ctrl.finished();
    }
}
