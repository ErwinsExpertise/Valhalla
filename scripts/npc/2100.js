var Quest_Heena_Sera = 1;

if (plr.checkQuestData(Quest_Heena_Sera, "2") && plr.itemCount(4031003) >= 1) {
    npc.sendOk("Haven't given #r#p2101##k the mirror yet? She should be on sitting somewhere on the west side...it's pretty close from here so it will be easy to spot her...")
} else if (plr.checkQuestData(Quest_Heena_Sera, "2") && plr.itemCount(4031003) < 1) {
    npc.sendBackNext("How am I going to hang all these up? Sigh... what? My mirror? Please don't tell me #p2101# asked you for this ...", false, true)
    plr.giveItem(4031003, 1)
    npc.sendOk("Aye...she should have come and get it herself. Seriously, she is SOOO lazy. Here's the mirror you're looking for.")
} else if (plr.checkQuestData(Quest_Heena_Sera, "1")) {
    npc.sendBackNext("How am I going to hang all these up? Sigh... what? My mirror? Please don't tell me #p2101# asked you for this ...", false, true)
    plr.setQuestData(Quest_Heena_Sera, "2")
    plr.giveItem(4031003, 1)
    npc.sendOk("Aye...she should have come and get it herself. Seriously, she is SOOO lazy. Here's the mirror you're looking for.")
} else if (plr.checkQuestData(Quest_Heena_Sera, "3")) {
    npc.sendOk("Did you give my mirror to Sarah? When will she help me to do this...")
} else {
    npc.sendOk("It is a fine day to do the laundry~ Don't you think so?")
}