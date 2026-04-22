var mapId = plr.mapID();

if (mapId === 920010000) {
    if (plr.isLeader()) {
        npc.sendOk("Thank you for rescuing me. I will now teleport you to the Door of the Goddess's Tower.");
        plr.warpEventMembers(920010100);
    } else {
        npc.sendOk("Please ask your party leader to speak with me.");
    }
} else if (mapId === 920010100) {
    var props = map.properties();
    var scars = ["scar1", "scar2", "scar3", "scar4", "scar5", "scar6"];
    var fixed = true;
    for (var i = 0; i < scars.length; i++) {
        if (!props[scars[i]]) {
            fixed = false;
            break;
        }
    }

    if (!plr.isLeader()) {
        npc.sendOk(fixed ? "Please return to the Center Tower and continue the quest." : "Please ask your party leader to speak with me.");
    } else if (fixed) {
        npc.sendOk("The Goddess's statue has been restored. I'll send the party onward.");
        plr.warpEventMembers(920010800);
    } else {
        npc.sendOk("Fix the Goddess statue before speaking with me again.");
    }
} else if (mapId === 920011200) {
    plr.warp(200080101);
} else {
    npc.sendOk("Keep going. The Tower of the Goddess still holds more trials ahead.");
}
