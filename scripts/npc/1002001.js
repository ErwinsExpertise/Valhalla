// Quest progress variables
var Quest_Maya = 1000200

// Check if quest is at "end" state
if (plr.level() >= 15 && plr.checkQuestData(Quest_Maya, "end")) {
    npc.sendOk("Oh~ So it seems that you have delievered #bSparkling Rock#k to #p1012101#...! Thanks. Now #p1012101# can get better now. huh? Marbles It is still very beautiful~")
}
// Quest ready to turn in at m1
else if (plr.level() >= 15 && plr.checkQuestData(Quest_Maya, "m1")) {
    npc.sendBackNext("I heard that #p1012101# is sick again. Sad...", false, true)
    npc.sendBackNext("Well, I do have the #bStrange Medicine#k ... but I can't give it to you for free. If you get me the #bMarbles#k, though, then I may reconsider ... meaning I am willing to trade with you...", false, true)
    npc.sendBackNext("Do you want to know how to get that stone? I wouldn't be asking for help if I knew how ... here's the deal, how about going to the department store at #bKerning City#k, ask for the daughter of the owner, #rSophia#k, and ask her about its whereabouts? She may have a clue ...", false, true)
    if (npc.sendYesNo("Please don't tell me you have no idea how to get to #bKerning City#k. Ok, take the exit on the right from the harbor, go past #bPerion#k up northwest, and then keep going east, then you'll find #bKerning City#k. Or do you know all this??")) {
        plr.setQuestData(Quest_Maya, "m2")
    }
}
// Ready to complete stage m6 (has Marbles in inventory)
else if (plr.level() >= 15 && plr.checkQuestData(Quest_Maya, "m6") && plr.itemCount(4031004) >= 1) {
    npc.sendBackNext("I heard that #p1012101# is sick again. Sad...", false, true)
    if (npc.sendYesNo("Ohhh ... this ... it's #bMarbles#k!!! How did you get this?? That's just incredible!! How about trading that stone with me? I'll give you #bStrange Medicine#k in return!")) {
        plr.setQuestData(Quest_Maya, "m7")
        plr.removeItemsByID(4031004, 1)
        plr.giveItem(4031006, 1)
        npc.sendOk("Thank you!!! Please take this medicine instead. #r#p1012101##k from #b#m100000000##k is sick again. This will take care of the sickness a little bit...")
    }
}
// Already have medicine, give final reminder
else if (plr.level() >= 15 && plr.checkQuestData(Quest_Maya, "m7")) {
    if (plr.itemCount(4031006) >= 1) {
        npc.sendOk("Hurry and give #p1012101# the #bStrange Medicine#k that I gave you. #m100000000# is very far from here so make sure you don't lose that medicine ...!")
    } else {
        plr.giveItem(4031006, 1)
        npc.sendOk("Hurry and give #p1012101# the #bStrange Medicine#k that I gave you. #m100000000# is very far from here so make sure you don't lose that medicine ...!")
    }
}
// In-progress quests (dummy selection showing stage)
else if (plr.level() >= 15) {
    if (plr.checkQuestData(Quest_Maya, "m2")) {
        npc.sendBackNext("I heard that #p1012101# is sick again. Sad...", false, true)
        npc.sendOk("Didn't get #bMarbles#k yet? Oh well ... those two aren't the easiest things to acquire ... they look gorgeous when they shine like stars ... hurry and go see #rSophia#k from the department store at #bKerning City#k.")
    } else if ((plr.checkQuestData(Quest_Maya, "m3") || plr.checkQuestData(Quest_Maya, "m4") || plr.checkQuestData(Quest_Maya, "m5") || (plr.checkQuestData(Quest_Maya, "m6") && !plr.itemCount(4031004)))) {
        npc.sendBackNext("I heard that #p1012101# is sick again. Sad...", false, true)
        npc.sendOk("Didn't get #bMarbles#k yet? Oh well ... those two aren't the easiest things to acquire ... they look gorgeous when they shine like stars ... hurry and go see #rSophia#k from the department store at #bKerning City#k.")
    } else {
        npc.sendOk("I heard that #p1012101# is sick again. Sad...")
    }
} else {
    npc.sendOk("I heard that #p1012101# is sick again. Sad...")
}