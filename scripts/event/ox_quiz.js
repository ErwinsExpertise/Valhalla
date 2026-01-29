var startMapID = 109020001;
var exitWinnerMapID = 109050000;
var exitMapID = 109050001;
var questions = [
    { q: "The sun rises in the east.", a: 0 },
    { q: "MapleStory is a racing game.", a: 1 },
    { q: "Snow is cold.", a: 0 },
    { q: "2 + 2 equals 5.", a: 1 },
    { q: "Water is wet.", a: 0 },
    { q: "Birds can fly.", a: 0 },
    { q: "A triangle has four sides.", a: 1 },
    { q: "Fire is hot.", a: 0 },
    { q: "The earth is flat.", a: 1 },
    { q: "Fish live in water.", a: 0 }
];
var currentAnswer = 0;

function start() {
    ctrl.setDuration("15m");

    var field = ctrl.getMap(startMapID);
    field.reset();
    field.clearProperties();

    var players = ctrl.players();
    var time = ctrl.remainingTime();

    for (let i = 0; i < players.length; i++) {
        players[i].warp(startMapID);
        players[i].showCountdown(time);
    }

    ctrl.schedule("askQuestion", "5s");
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

function askQuestion() {
    if (questions.length === 0) {
        finishQuiz();
        return;
    }

    var next = questions.shift();
    currentAnswer = next.a;

    var players = ctrl.players();
    for (let i = 0; i < players.length; i++) {
        players[i].sendMessage(next.q + " (O = True, X = False)");
    }

    ctrl.schedule("checkAnswer", "30s");
}

function checkAnswer() {
    var map = ctrl.getMap(startMapID);
    var players = ctrl.players();

    for (let i = players.length - 1; i >= 0; i--) {
        var plr = players[i];
        if (!map.isPlayerInArea(plr, currentAnswer)) {
            plr.warp(exitMapID);
            ctrl.removePlayer(plr);
        }
    }

    if (ctrl.playerCount() <= 1 || questions.length === 0) {
        finishQuiz();
        return;
    }

    ctrl.schedule("askQuestion", "10s");
}

function finishQuiz() {
    var players = ctrl.players();
    for (let i = 0; i < players.length; i++) {
        players[i].warp(exitWinnerMapID);
    }
    ctrl.finished();
}
