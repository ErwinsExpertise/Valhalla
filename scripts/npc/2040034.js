// Alishar - LudiPQ Entry NPC

var maps = [922010100, 922010200, 922010300, 922010400, 922010500, 922010600, 922010700, 922010800, 922010900];

if (plr.isPartyLeader()) {
    var plrs = plr.partyMembersOnMap();

    var badLevel = false;

    for (let i = 0; i < plrs.length; i++) {
        if (plrs[i].level() < 35 || plrs[i].level() > 50) {
            badLevel = true;
            break;
        }
    }

    if (plr.partyMembersOnMapCount() < 3) {
        npc.sendOk("You need to be a party of at least 3 on the same map");
    } else if (badLevel) {
        npc.sendOk("Someone in your party is not the correct level (35-50)");
    } else {
        for (let instance = 0; instance < 1; instance++) {
            var count = 0;

            for(let i = 0; i < maps.length; i++) {
                var m = map.getMap(maps[i], instance);
                count += m.playerCount();
            }

            if (count == 0) {
                plr.startPartyQuest("ludibrium_pq", instance);
            } else {
                npc.sendOk("A party is already doing the quest, please come back another time");
            }
        }
    }
} else {
    npc.sendOk("You need to be party leader to start a party quest");
}