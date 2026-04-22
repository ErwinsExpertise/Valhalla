if (plr.job() === 0) {
    if (plr.getLevel() < 10 || plr.getDex() < 25) {
        npc.sendOk("Train a bit more and I can show you the way of the #rThief#k.");
    } else {
        npc.sendNext("So you decided to become a #rThief#k?");
        npc.sendNext("It is an important and final choice. You will not be able to turn back.");
        if (!npc.sendYesNo("Do you want to become a #rThief#k?")) {
            npc.sendOk("Make up your mind and visit me again.");
        } else {
            plr.setJob(400);
            plr.gainItem(1332063, 1);
            npc.sendOk("So be it! Now go, and go with pride.");
        }
    }
} else if (plr.job() === 400) {
    if (plr.getLevel() < 30) {
        npc.sendOk("You have chosen wisely.");
    } else if (plr.questStarted(100200) || plr.questCompleted(100202)) {
        if (plr.questCompleted(100202)) {
            var branch = npc.sendMenu("What do you want to become?#b", "Assassin", "Bandit");
            var jobName = branch === 0 ? "Assassin" : "Bandit";
            var jobId = branch === 0 ? 410 : 420;
            if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
                plr.setJob(jobId);
                npc.sendOk("So be it! Now go, and go with pride.");
            }
        } else if (plr.questStarted(100202)) {
            npc.sendOk("Go and find me the #rNecklace of Wisdom#k which is hidden on the Holy Ground at the Snowfield.");
        } else {
            plr.completeQuest(100200);
            if (plr.questCompleted(100200)) {
                if (npc.sendAcceptDecline("Is your mind ready to undertake the final test?")) {
                    plr.startQuest(100202);
                    npc.sendOk("Go and find me the #rNecklace of Wisdom#k which is hidden on the Holy Ground at the Snowfield.");
                }
            } else {
                npc.sendOk("Well, well. Now go and see #bthe Dark Lord#k. He will show you the way.");
            }
        }
    } else if (plr.getRemainingSP() <= (plr.getLevel() - 30) * 3) {
        npc.sendNext("#rBy Odin's beard!#k You are a strong one.");
        if (npc.sendAcceptDecline("But I can make you even stronger. Although you will have to prove not only your strength but your knowledge. Are you ready for the challenge?")) {
            plr.startQuest(100200);
            npc.sendOk("Well, well. Now go and see #bthe Dark Lord#k. He will show you the way.");
        }
    } else {
        npc.sendOk("Your time has yet to come...");
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
