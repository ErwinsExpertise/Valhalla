if (!(plr.job() === 110 || plr.job() === 120 || plr.job() === 130 || plr.job() === 100)) {
    npc.sendOk("Only those who have already walked the warrior's path belong here.");
} else if (plr.job() === 110 || plr.job() === 120 || plr.job() === 130) {
    npc.sendOk("A warrior's training never ends. Keep sharpening your body and your will.");
} else if (plr.questCompleted(100005)) {
    npc.sendNext("You've returned with proof that your strength is real and your judgment is steady. That's enough for me.");
    var branch = npc.sendMenu("Choose the path that best suits the kind of warrior you intend to become.#b", "Crusader", "White Knight", "Dragon Knight");
    var jobName = branch === 0 ? "Crusader" : branch === 1 ? "White Knight" : "Dragon Knight";
    var jobId = branch === 0 ? 111 : branch === 1 ? 121 : 131;
    if (npc.sendYesNo("Do you wish to become a #r" + jobName + "#k?")) {
        plr.setJob(jobId);
        plr.giveAP(5);
        npc.sendOk("Good. Then carry that title with pride. You are now a #b" + jobName + "#k.");
    }
} else if (plr.questStarted(100005)) {
    npc.sendOk("The Holy Stone still waits for you. Bring me the #rNecklace of Wisdom#k and prove your resolve is as firm as your body.");
} else if (plr.questCompleted(100004)) {
    npc.sendNext("Your strength has already been proven. Now I need to know whether your judgment can bear the same weight.");
    if (npc.sendAcceptDecline("Are you ready for the final test?")) {
        plr.startQuest(100005);
        npc.sendOk("Go to the Holy Ground at the Snowfield and bring me the #rNecklace of Wisdom#k.");
    }
} else if (plr.questStarted(100004)) {
    npc.sendOk("Dances with Balrog is waiting for you. Finish the trial he set before you return here.");
} else if (plr.job() === 100 && plr.getLevel() >= 70 && plr.getRemainingSP() <= (plr.getLevel() - 70) * 3) {
    npc.sendNext("You've built a strong foundation. The next step will demand more than raw strength from you.");
    if (npc.sendAcceptDecline("Will you prove that your power, discipline, and resolve are worthy of a higher rank?")) {
        plr.startQuest(100004);
        npc.sendOk("Go and see #bDances with Balrog#k. He will measure whether your strength is truly battle-forged.");
    }
} else {
    npc.sendOk("You are not ready yet. Train until both your body and your technique can support the next step.");
}
