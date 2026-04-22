if (plr.job() === 0) {
    if (plr.getLevel() < 10 || plr.getStr() < 35) {
        npc.sendOk("Train a bit more and I can show you the way of the #rWarrior#k.");
    } else {
        npc.sendNext("So you decided to become a #rWarrior#k?");
        npc.sendNext("It is an important and final choice. You will not be able to turn back.");
        if (!npc.sendYesNo("Do you want to become a #rWarrior#k?")) {
            npc.sendOk("Make up your mind and visit me again.");
        } else {
            plr.setJob(100);
            plr.gainItem(1402001, 1);
            npc.sendOk("So be it! Now go, and go with pride.");
        }
    }
} else if (plr.job() === 100) {
    if (plr.getLevel() < 30) {
        npc.sendOk("You have chosen wisely.");
    } else if (plr.questStarted(100003) || plr.questCompleted(100005)) {
        if (plr.questCompleted(100005)) {
            var branch = npc.sendMenu("What do you want to become?#b", "Fighter", "Page", "Spearman");
            var jobName = branch === 0 ? "Fighter" : branch === 1 ? "Page" : "Spearman";
            var jobId = branch === 0 ? 110 : branch === 1 ? 120 : 130;
            if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
                plr.setJob(jobId);
                npc.sendOk("So be it! Now go, and go with pride.");
            }
        } else {
            plr.completeQuest(100005);
            if (plr.questCompleted(100005)) {
                npc.sendNext("I see you have done well. I will allow you to take the next step on your long road.");
                var selection = npc.sendMenu("What do you want to become?#b", "Fighter", "Page", "Spearman");
                var selectedName = selection === 0 ? "Fighter" : selection === 1 ? "Page" : "Spearman";
                var selectedJob = selection === 0 ? 110 : selection === 1 ? 120 : 130;
                if (npc.sendYesNo("Do you want to become a #r" + selectedName + "#k?")) {
                    plr.setJob(selectedJob);
                    npc.sendOk("So be it! Now go, and go with pride.");
                }
            } else {
                npc.sendOk("Go and see the #rJob Instructor#k.");
            }
        }
    } else {
        npc.sendNext("The progress you have made is astonishing.");
        if (npc.sendAcceptDecline("But first I must test your skills. Are you ready?")) {
            if (plr.haveItem(4031008, 1)) {
                npc.sendOk("Please report this bug using @bug\r\nstatus = 13");
            } else {
                plr.startQuest(100003);
                npc.sendOk("Go see the #bJob Instructor#k near Perion. He will show you the way.");
            }
        }
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
