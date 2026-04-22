if (!(plr.job() === 210 || plr.job() === 220 || plr.job() === 230 || plr.job() === 200)) {
	npc.sendOk("Arcane study is not for the impatient. Return when you truly walk the path of magic.");
} else if (plr.job() === 210 || plr.job() === 220 || plr.job() === 230) {
	npc.sendOk("A mage's journey never truly ends. Keep refining your craft.");
} else if (plr.questCompleted(100102)) {
	npc.sendNext("You have returned with proof that both your power and your judgment have matured. That is what I wished to see.");
	var branch = npc.sendMenu("What do you want to become?#b", "Fire/Poison Mage", "Ice/Lightning Mage", "Priest");
	var jobName = branch === 0 ? "Fire/Poison Mage" : branch === 1 ? "Ice/Lightning Mage" : "Priest";
	var jobId = branch === 0 ? 211 : branch === 1 ? 221 : 231;
	if (npc.sendYesNo("Do you want to become a #r" + jobName + "#k?")) {
		plr.setJob(jobId);
		plr.giveAP(5);
		npc.sendOk("Then take the next step. From this point on, you are a #b" + jobName + "#k.");
	}
} else if (plr.questStarted(100102)) {
	npc.sendOk("The Holy Stone awaits you still. Bring me the #rNecklace of Wisdom#k and I will know your mind is ready.");
} else if (plr.questCompleted(100100)) {
	npc.sendNext("Your raw power is no longer in doubt. Now you must prove you can guide that power with discipline.");
	if (npc.sendAcceptDecline("Is your mind ready to undertake the final test?")) {
		plr.startQuest(100102);
		npc.sendOk("Go to the Holy Ground at the Snowfield and bring me the #rNecklace of Wisdom#k.");
	}
} else if (plr.questStarted(100100)) {
	npc.sendOk("Your next step begins with Grendel. See him, and complete the trial he sets before you.");
} else if (plr.job() === 200 && plr.getLevel() >= 70 && plr.getRemainingSP() <= (plr.getLevel() - 70) * 3) {
	npc.sendNext("You have advanced far enough to seek deeper magic.");
	if (npc.sendAcceptDecline("Power alone is not enough for a mage. Will you prove both your ability and your wisdom?")) {
		plr.startQuest(100100);
		npc.sendOk("Go and see #bGrendel the Really Old#k. He will judge whether your spellcraft is truly ready.");
	}
} else {
	npc.sendOk("You are not ready yet. Continue your studies and return when your foundation is complete.");
}
