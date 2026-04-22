if (plr.job() === 0) {
    if (plr.getLevel() < 8 || plr.getInt() < 20) {
        npc.sendOk("Train a bit more and I can show you the way of the #rMagician#k.");
    } else {
        npc.sendNext("So you decided to become a #rMagician#k?");
        npc.sendNext("It is an important and final choice. You will not be able to turn back.");
        if (!npc.sendYesNo("Do you want to become a #rMagician#k?")) {
            npc.sendOk("Make up your mind and visit me again.");
        } else {
            plr.setJob(200);
            plr.gainItem(1372043, 1);
            npc.sendOk("So be it! Now go, and go with pride.");
        }
    }
} else if (plr.job() === 200) {
    if (plr.getLevel() < 30) {
        npc.sendOk("You have chosen wisely.");
    } else if (plr.questStarted(100100) || plr.questCompleted(100102)) {
        if (plr.questCompleted(100102)) {
            var branch = npc.sendMenu("What do you want to become?#b", "Wizard (Fire, Poison)", "Wizard (Ice, Lightning)", "Cleric");
            var jobName = branch === 0 ? "Wizard of Fire and Poison" : branch === 1 ? "Wizard of Ice and Lightning" : "Cleric";
            var jobId = branch === 0 ? 210 : branch === 1 ? 220 : 230;
            if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
                plr.setJob(jobId);
                npc.sendOk("So be it! Now go, and go with pride.");
            }
        } else if (plr.questStarted(100102)) {
            npc.sendOk("Go and find me the #rNecklace of Wisdom#k which is hidden on the Holy Ground at the Snowfield.");
        } else {
            plr.completeQuest(100100);
            if (plr.questCompleted(100100)) {
                if (npc.sendAcceptDecline("Is your mind ready to undertake the final test?")) {
                    plr.startQuest(100102);
                    npc.sendOk("Go and find me the #rNecklace of Wisdom#k which is hidden on the Holy Ground at the Snowfield.");
                }
            } else {
                npc.sendOk("Well, well. Now go and see #bGrendel the Really Old#k. He will show you the way.");
            }
        }
    } else if (plr.getRemainingSP() <= (plr.getLevel() - 30) * 3) {
        npc.sendNext("#rBy Odin's beard!#k You are a strong one.");
        if (npc.sendAcceptDecline("But I can make you even stronger. Although you will have to prove not only your strength but your knowledge. Are you ready for the challenge?")) {
            plr.startQuest(100100);
            npc.sendOk("Well, well. Now go and see #bGrendel the Really Old#k. He will show you the way.");
        }
    } else {
        npc.sendOk("Your time has yet to come...");
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
