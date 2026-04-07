var startMapID = 109060000;
var exitWinnerMapID = 109050000;
var exitMapID = 109050001;
var teams = {};
var progress = { maple: 0, story: 0 };
var teamOrder = ["maple", "story"];
var targetProgress = 100;
var ended = false;
var started = false;

function start() {
    ctrl.setDuration("10m");

    var field = ctrl.getMap(startMapID);
    if (field.getMapID() !== 0) {
        field.reset();
        field.clearProperties();
        field.portalEnabled(false, "start00");
    }

    var players = ctrl.players();
    var time = ctrl.remainingTime();

    for (let i = 0; i < players.length; i++) {
        var team = teamOrder[i % teamOrder.length];
        teams[players[i].name()] = team;
        players[i].warp(startMapID);
        players[i].showCountdown(time);
    }

    ctrl.schedule("beginRound", "10m");
}

function beginRound() {
    if (ended || started) {
        return;
    }
    started = true;
    ctrl.setDuration("60m");
    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (const player of players) {
        player.showCountdown(time);
    }
    var field = ctrl.getMap(startMapID);
    if (field.getMapID() !== 0) {
        field.portalEnabled(true, "start00");
    }
    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (let i = 0; i < players.length; i++) {
        players[i].showCountdown(time);
    }
    ctrl.schedule("tick", "5s");
    ctrl.schedule("endRound", "60m");
}

function beforePortal(plr, src, dst) {
    return true;
}

function afterPortal(plr, dst) {
    plr.showCountdown(ctrl.remainingTime());
}

function timeout(plr) {
    if (started) {
        plr.warp(exitMapID);
    }
}

function playerLeaveEvent(plr) {
    ctrl.removePlayer(plr);
    plr.warp(exitMapID);

    if (ctrl.playerCount() <= 0) {
        ctrl.finished();
    }
}

function onReactorHit(plr, reactorName) {
    if (!started || ended || reactorName.indexOf("snow") !== 0) {
        return;
    }

    var team = teams[plr.name()];
    if (!team) {
        return;
    }

    progress[team] += 2;
    checkFinish();
}

function tick() {
    if (ended || !started) {
        return;
    }

    var map = ctrl.getMap(startMapID);
    var players = ctrl.players();

    for (let i = 0; i < players.length; i++) {
        var team = teams[players[i].name()];
        if (!team) {
            continue;
        }
        var areaId = team === "maple" ? 0 : 1;
        if (map.isPlayerInArea(players[i], areaId)) {
            progress[team] += 1;
        }
    }

    checkFinish();

    if (!ended) {
        ctrl.schedule("tick", "5s");
    }
}

function endRound() {
    var winner = progress.maple === progress.story
        ? teamOrder[Math.floor(Math.random() * teamOrder.length)]
        : (progress.maple > progress.story ? "maple" : "story");
    finishEvent(winner);
}

function checkFinish() {
    if (progress.maple >= targetProgress) {
        finishEvent("maple");
    } else if (progress.story >= targetProgress) {
        finishEvent("story");
    }
}

function finishEvent(winner) {
    if (ended) {
        return;
    }
    ended = true;

    var players = ctrl.players();
    for (let i = 0; i < players.length; i++) {
        var team = teams[players[i].name()];
        if (team === winner) {
            players[i].warp(exitWinnerMapID);
        } else {
            players[i].warp(exitMapID);
        }
    }

    ctrl.finished();
}
