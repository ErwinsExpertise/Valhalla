var startMapID = 109080000;
var exitWinnerMapID = 109050000;
var exitMapID = 109050001;
var teams = {};
var teamScores = { red: 0, blue: 0 };
var teamOrder = ["red", "blue"];
var scoreLimit = 30;
var ended = false;

function start() {
    ctrl.setDuration("5m");

    var field = ctrl.getMap(startMapID);
    field.reset();
    field.clearProperties();

    var players = ctrl.players();
    var time = ctrl.remainingTime();

    for (let i = 0; i < players.length; i++) {
        var team = teamOrder[i % teamOrder.length];
        teams[players[i].name()] = team;
        players[i].warp(startMapID);
        players[i].showCountdown(time);
    }

    ctrl.schedule("endRound", "5m");
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

function onReactorHit(plr, reactorName) {
    if (ended || reactorName.indexOf("coconut") !== 0) {
        return;
    }

    var team = teams[plr.name()];
    if (!team) {
        return;
    }

    teamScores[team] += 1;

    var players = ctrl.players();
    for (let i = 0; i < players.length; i++) {
        players[i].sendMessage(plr.name() + " scored for " + team + " team.");
    }

    if (teamScores[team] >= scoreLimit) {
        finishEvent(team);
    }
}

function endRound() {
    var winner = teamScores.red === teamScores.blue
        ? teamOrder[Math.floor(Math.random() * teamOrder.length)]
        : (teamScores.red > teamScores.blue ? "red" : "blue");
    finishEvent(winner);
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
