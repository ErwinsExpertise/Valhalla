var mapId = plr.mapID();
var ticket = 4031711 + Math.floor(mapId / 300000000);

if (mapId === 103000100 || mapId === 600010001) {
    var destination = mapId === 103000100 ? "New Leaf City of Masteria" : "Kerning City of Victoria Island";
    if (npc.sendYesNo("Travel to " + destination + " will cost you #b5000 mesos#k. Are you sure you want to buy a #b#t" + ticket + "##k?")) {
        if (plr.getMesos() >= 5000) {
            plr.gainMesos(-5000);
            plr.gainItem(ticket, 1);
            npc.sendNext("You now have the travel ticket.");
        } else {
            npc.sendNext("You do not have enough mesos!");
        }
    }
} else if (mapId === 600010002 || mapId === 600010004) {
    if (npc.sendYesNo("You want to get out before the train leaves? There will be no refund.")) {
        npc.sendNext("All right, see you the next time.");
        plr.warp(mapId === 600010002 ? 600010001 : 103000100);
    }
}
