npc.sendBackNext("Welcome to the Thieves' Hideout. Only those who are invited will ever find it. Try not to get lost on the way out.", false, true)

if (plr.getQuestStatus(2351) != 1 || plr.job() == 400) {
    npc.sendOk("A secret conversation? Thieves may trade in secrets, but such things are reserved for their enemies.")

}

npc.sendBackNext("I'm sure you came here because you want to be a Thief, correct? I hope your heart is in this...many Beginners think they have what it takes, but run screaming the moment they see me. They must really be afraid of the challenges Thieves face...", true, true)

if (!npc.sendYesNo("All right, you ready to become a Thief?")) {

}

if (plr.getLevel() < 10) {
    npc.sendOk("Train a bit more until you reach the base requirements and I can show you the way of the #rThief#k.")

}

// Check Equip tab has at least 3 free slots (generic free-slot check not directly exposed; skip or stub)
if (plr.getQuestStatus(7635) < 1) {
    plr.startQuest(7635)
    plr.setJob(400)
    // resetStats call not available in new API; must leave as-is
    // expandInventory function unavailable; skipped
    plr.giveItem(1332063, 1)
}

npc.sendBackNext("With this, you have become a Thief. Since you can use Thief skills now, open your Skill window and have a look. As you level up, you will be able to learn more skills.", true, true)
npc.sendBackNext("But skills aren't enough, right? A true Thief must have the stats to match! A Thief uses LUK as the main stat and DEX as the secondary stat. If you don't know how to raise stats, just use #bAuto-Assign#k.", true, true)
npc.sendBackNext("Oh, I gave you a little gift, too. I expanded a few slots in your Equip and ETC item tabs. Bigger Inventory, better life, I always say.", true, true)
npc.sendBackNext("Now a word of warning. Everyone loses some of their earned EXP when they fall in battle. Be careful. You don't want to lose anything you worked to get, eh?", true, true)
npc.sendBackNext("Right, that's it. Take the equipment I gave you, and use it to train your skills as a Thief.", true, true)