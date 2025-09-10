var Quest_Maya = 1000200;

// Check quest state and offer
if (plr.level() >= 15 && plr.checkQuestStatus(Quest_Maya, 0)) {
    npc.sendNext("Cough... Cough... Ah... Headache... Can somebody help me?...\r\n\r\n#r#eQUEST AVAILABLE#k#n#l\r\n#L0##bMaya of Henesys#k#l")

    npc.sendNext("Cough ... cough ... ah ... oh, hello stranger. Sorry, but may I ask you for a favor? I've been suffering from sickness for a while, and the doctors can't do anything about it. Lately, it has gotten so bad I can't even take care of myself.")

    if (npc.sendYesNo("Sorry to ask, but is there any way you can get me the #b#t4031006##k? I am not sure exactly how to get that medicine, but #r#p1002001##k from #b#m104000000##k may know a thing or two about it. Please help me out.")) {
        plr.startQuest(Quest_Maya)
        plr.setQuestData(Quest_Maya, "m1")
        npc.sendOk("#p1002001# from #b#m104000000##k can definitely help you find some Weird Medicine. Please...talk to him for me...")
    }

} else if (plr.level() >= 15 && plr.checkQuestData(Quest_Maya, "m7") && plr.itemCount(4031006) >= 1) {
    npc.sendNext("Cough... Cough... Ah... Headache... Can somebody help me?...\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bDelivering the Weird Medicine (Ready to complete.)#k#l")

    npc.sendNext("Darn ... my whole body's aching ... what, oh my ... isn't that #b#t4031006##k?? How did you get it?? wow, you must be amazing.")

    if (npc.sendBackNext("Um ... is it okay if I get that medicine? I'll give you something that I don't really need ... I urge you ... please ... I need that medicine ...", true, true)) {
        plr.setQuestData(Quest_Maya, "end")
        plr.completeQuest(Quest_Maya)
        plr.giveItem(4031006, -1)
        plr.giveItem(1002026, 1)
        plr.giveEXP(2000)
        plr.giveMesos(5000)
        npc.sendOk("Thank you so much ... this may cure my longtime sickness afterall ... here's something I don't really need ... hopefully it'll help you through your journey ... here are some mesos also as a sign of my appreciation ...\r\n\r\n#e#rREWARD:#k\r\n+5,000 mesos\r\n+2,000 exp\r\n+#i1002026# #t1002026#")
    }

} else if (plr.level() >= 15 && (plr.checkQuestData(Quest_Maya, "m1") || plr.checkQuestData(Quest_Maya, "m2") || plr.checkQuestData(Quest_Maya, "m3") || plr.checkQuestData(Quest_Maya, "m4") || plr.checkQuestData(Quest_Maya, "m5") || plr.checkQuestData(Quest_Maya, "m6"))) {
    npc.sendOk("You haven't met up with #r#p1002001##k, yet? #p1002001# from #b#m104000000##k can definitely help you find some Weird Medicine. Please...talk\r\nto him for me...")

} else if (plr.level() >= 15 && plr.checkQuestData(Quest_Maya, "end")) {
    npc.sendOk("Thanks for the last time. Now I feel much better. Thanks for everything.")

} else {
    npc.sendOk("Cough... Cough... Ah... Headache... Can somebody help me?...")
}