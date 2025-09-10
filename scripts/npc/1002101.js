if (plr.getQuestStatus(22579) < 2) {
    npc.sendOk("I'm just a retired crewman. I'm focused on training powerful Explorers now.");
} else if (npc.sendYesNo("Do you want to go to the island in John's Map right now?")) {
    plr.warp(200090080);
    // server will handle map-time-limit internally, no additional API
} else {
    npc.sendOk("Ah, you still have businesses left in Lith Harbor.");
}