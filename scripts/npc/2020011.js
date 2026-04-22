if (!(plr.job() === 410 || plr.job() === 420 || plr.job() === 400)) {
	npc.sendOk("The shadows do not open for just anyone. Return when you truly understand the thief's road.");
} else if (plr.job() === 410 || plr.job() === 420) {
	npc.sendOk("A thief survives by staying sharp. Keep honing your instincts.");
} else if (plr.questCompleted(100102)) {
	npc.sendNext("You've brought back proof that both your nerve and your judgment held firm. That's what I wanted to see.");
	var branch = npc.sendMenu("What do you want to become?#b", "Hermit", "Chief Bandit");
	var jobName = branch === 0 ? "Hermit" : "Chief Bandit";
	var jobId = branch === 0 ? 411 : 421;
	if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
		plr.setJob(jobId);
		plr.giveAP(5);
		npc.sendOk("Good. From here on, you are a #b" + jobName + "#k.");
	}
} else if (plr.questStarted(100102)) {
	npc.sendOk("The Holy Stone still waits for you. Bring me the #rNecklace of Wisdom#k and prove your judgment is as sharp as your blade.");
} else if (plr.questCompleted(100100)) {
	npc.sendNext("You've already proven your hands are fast enough. Now prove your mind is just as sharp.");
	if (npc.sendAcceptDecline("Is your mind ready to undertake the final test?")) {
		plr.startQuest(100102);
		npc.sendOk("Go to the Holy Ground at the Snowfield and bring me the #rNecklace of Wisdom#k.");
	}
} else if (plr.questStarted(100100)) {
	npc.sendOk("The Dark Lord is waiting. Finish the trial he set for you, then return here.");
} else if (plr.job() === 400 && plr.getLevel() >= 70 && plr.getRemainingSP() <= (plr.getLevel() - 70) * 3) {
	npc.sendNext("You've made it this far because your instincts are good. The next step demands more than instinct.");
	if (npc.sendAcceptDecline("Will you prove that you have the patience, cunning, and judgment to rise even higher?")) {
		plr.startQuest(100100);
		npc.sendOk("Go and see #bthe Dark Lord#k. He'll decide whether you're ready for the next trial.");
	}
} else {
	npc.sendOk("Not yet. Keep training until your hands and your head are both ready.");
}
