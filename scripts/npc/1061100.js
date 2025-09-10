npc.sendBackNext("Welcome. We're the #m105000000# Hotel. Our hotel works hard to serve you the best at all times. If you are tired and worn out from hunting, how about a relaxing stay at our hotel?", false, true);

var chat = "We offer two kinds of rooms for service. Please choose the one of your liking.";
chat += "\r\n#L0##b#m105000011#(499 Mesos per use)#l\r\n";
chat += "#L1##b#m105000012#(999 Mesos per use)#l";
npc.sendSelection(chat);

var select = npc.selection();
var confirm;
if (select === 0) {
    confirm = npc.sendYesNo("You've chosen the regular sauna. Your HP and MP will recover fast and you can even purchase some items there. Are you sure you want to go in?");
} else {
    confirm = npc.sendYesNo("You've chosen the VIP sauna. your HP and MP will recover even faster than that of the regular sauna and you can even find a special item in there. Are you sure you want to go in?");
}

if (confirm) {
    var cost = select === 0 ? 499 : 999;
    if (plr.mesos() < cost) {
        npc.sendNext("I'm sorry. It looks like you don't have mesos. It will cost you " + (select === 0 ? "at least " : "") + "" + cost + " Mesos to stay at our hotel.");
    } else {
        plr.takeMesos(cost);
        plr.warp(select === 0 ? 105000011 : 105000012);
    }
} else {
    npc.sendNext("We offer other kinds of services, too, so please think carefully and then make your decision.");
}