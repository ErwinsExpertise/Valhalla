if (!(plr.job() === 310 || plr.job() === 320 || plr.job() === 300)) {
	npc.sendOk("A bowman must learn patience, precision, and restraint. Return when you truly walk that path.");
} else if (plr.job() === 310 || plr.job() === 320) {
	npc.sendOk("Your aim has improved, but there is always a farther mark to hit.");
} else if (plr.questCompleted(100102)) {
	npc.sendNext("You've returned with proof that your eyes are sharp and your judgment is steady. That is what a true bowman needs.");
	var branch = npc.sendMenu("What do you want to become?#b", "Ranger", "Sniper");
	var jobName = branch === 0 ? "Ranger" : "Sniper";
	var jobId = branch === 0 ? 311 : 321;
	if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
		plr.setJob(jobId);
		plr.giveAP(5);
		npc.sendOk("Then let your arrows speak for you. You are now a #b" + jobName + "#k.");
	}
} else if (plr.questStarted(100102)) {
	npc.sendOk("The Holy Stone still waits for you. Bring me the #rNecklace of Wisdom#k and I will know your judgment is ready.");
} else if (plr.questCompleted(100100)) {
	npc.sendNext("Your aim and your resolve have been proven. Now you must show me that your mind is just as steady.");
	if (npc.sendAcceptDecline("Is your mind ready to undertake the final test?")) {
		plr.startQuest(100102);
		npc.sendOk("Go to the Holy Ground at the Snowfield and bring me the #rNecklace of Wisdom#k.");
	}
} else if (plr.questStarted(100100)) {
	npc.sendOk("Athena Pierce is waiting for you. Complete the trial she gives you, then return here.");
} else if (plr.job() === 300 && plr.getLevel() >= 70 && plr.getRemainingSP() <= (plr.getLevel() - 70) * 3) {
	npc.sendNext("You've trained well. Your next step is not about strength alone, but control.");
	if (npc.sendAcceptDecline("Will you prove that your eyes, your hands, and your judgment are all worthy of the next rank?")) {
		plr.startQuest(100100);
		npc.sendOk("Go and see #bAthena Pierce#k. She will put your skill to the test.");
	}
} else {
	npc.sendOk("You are not ready yet. Keep training until your aim never wavers.");
}
