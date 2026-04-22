if (plr.job() === 0) {
    if (plr.getLevel() < 10 || plr.getDex() < 25) {
        npc.sendOk("Train a bit more and I can show you the way of the #rBowman#k.");
    } else {
        npc.sendNext("So you decided to become a #rBowman#k?");
        npc.sendNext("It is an important and final choice. You will not be able to turn back.");
        if (!npc.sendYesNo("Do you want to become a #rBowman#k?")) {
            npc.sendOk("Make up your mind and visit me again.");
        } else {
            plr.setJob(300);
            plr.gainItem(1452002, 1);
            plr.gainItem(2060000, 1000);
            npc.sendOk("So be it! Now go, and go with pride.");
        }
    }
} else if (plr.job() === 300) {
    if (plr.getLevel() < 30) {
        npc.sendOk("You have chosen wisely.");
    } else if (plr.questStarted(100000) || plr.questCompleted(100002)) {
        if (plr.questCompleted(100002)) {
            var branch = npc.sendMenu("What do you want to become?#b", "Hunter", "Crossbowman");
            var jobName = branch === 0 ? "Hunter" : "Crossbowman";
            var jobId = branch === 0 ? 310 : 320;
            if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
                plr.setJob(jobId);
                npc.sendOk("So be it! Now go, and go with pride.");
            }
        } else {
            plr.completeQuest(100002);
            if (plr.questCompleted(100002)) {
                npc.sendNext("I see you have done well. I will allow you to take the next step on your long road.");
                var selection = npc.sendMenu("What do you want to become?#b", "Hunter", "Crossbowman");
                var selectedName = selection === 0 ? "Hunter" : "Crossbowman";
                var selectedJob = selection === 0 ? 310 : 320;
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
            if (plr.haveItem(4031010, 1)) {
                npc.sendOk("Please report this bug using @bug\r\nstatus = 13");
            } else {
                plr.startQuest(100000);
                npc.sendOk("Go see the #bJob Instructor#k near Henesys. He will show you the way.");
            }
        }
    }
} else {
    npc.sendOk("You have chosen wisely.");
}
