if (!gmEventActive) {
    npc.sendOk("There is no GM event running right now.");
} else if (plr.channel() != gmEventChannel) {
    plr.sendNotice("The GM event is on channel " + gmEventChannel + ". Please switch channels first.");
    npc.sendOk("The GM event is on channel " + gmEventChannel + ". Please switch channels first.");
} else if (npc.sendYesNo("A GM event is running on channel " + gmEventChannel + ". Enter now?")) {
    plr.warp(gmEventMap);
}
