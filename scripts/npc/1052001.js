// Dark Lord — Thief job advancements (level and job only)

// Beginner -> Thief (Level 10+)
if (plr.job() === 0) {
    if (plr.getLevel() >= 10) {
        npc.sendBackNext(
            "I'm sure you came here because you want to be a Thief, correct? I hope your heart is in this... many Beginners think they have what it takes, but run screaming the moment they see me.",
            true, true
        )
        if (!npc.sendYesNo("All right, you ready to become a Thief?")) {
            npc.sendOk("Think it over. Come back when your resolve is firm.")
        } else {
            plr.setJob(400)
            // Starter dagger
            plr.giveItem(1332063, 1)
            npc.sendBackNext("With this, you have become a Thief. Since you can use Thief skills now, open your Skill window and have a look. As you level up, you will be able to learn more skills.", true, true)
            npc.sendBackNext("A true Thief must have the stats to match! A Thief uses LUK as the main stat and DEX as the secondary stat. If you don't know how to raise stats, just use #bAuto-Assign#k.", true, true)
            npc.sendOk("Right, that's it. Take the equipment I gave you, and use it to train your skills as a Thief.")
        }
    } else {
        npc.sendOk("Train a bit more until you reach Level 10 and I can show you the way of the #rThief#k.")
    }
    // Thief -> 2nd Job (Level 30+): Assassin or Bandit
} else if (plr.job() === 400) {
    if (plr.getLevel() >= 30) {
        npc.sendBackNext("Hmmm... you seem to have gotten a whole lot stronger. Ready to take the next step?", false, true)

        // Simple branch picker: Assassin or Bandit
        var branch = npc.sendMenu(
            "Choose your 2nd job advancement.",
            "Assassin",
            "Bandit"
        )

        var jobName = (branch === 0) ? "Assassin" : "Bandit"
        var jobId   = (branch === 0) ? 410 : 420

        if (npc.sendYesNo("So you want to advance as a #b" + jobName + "#k? Once you make the decision, you can't go back. Are you sure?")) {
            plr.setJob(jobId)
            if (branch === 0) {
                npc.sendBackNext("From here on out you are an Assassin. Keep training and hone your skills in the shadows.", true, true)
            } else {
                npc.sendBackNext("From here on out you are a Bandit. Keep training and hone your skills in the shadows.", true, true)
            }
            npc.sendOk("I have also given you a little bit of #bSP#k. Open the #bSkill Menu#k to enhance your 2nd job skills. Some skills require others first, so choose wisely.")
        } else {
            npc.sendOk("Take your time. This decision is important.")
        }
    } else {
        npc.sendOk("Keep training as a Thief. Return to me at #rLevel 30#k for your next advancement.")
    }
    // Already 2nd job Thief paths
} else if (plr.job() === 410 || plr.job() === 420) {
    npc.sendOk("Walk the path you've chosen with pride. Keep training and growing stronger.")
} else {
    npc.sendOk("The progress you have made is astonishing.")
}