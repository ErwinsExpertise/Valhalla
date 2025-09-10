if (plr.getQuestStatus(31000) < 2) {
    npc.sendOk("I don't think you're ready to go Chryse. I can't move you if you've never been to Chryse by visiting me.");
} else {
    if (npc.sendYesNo("You want to go to Chryse?")) {
        npc.sendNext("Okay, I am going to send you to Chryse. Get ready?");
        plr.warp(200100001);
    }
}