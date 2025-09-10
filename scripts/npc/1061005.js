// Config
var Quest_DietMed = 1000700;
var Quest_AgingMed = 1000701;
var scrolls = [2040504, 2040505];
var selectedScroll = scrolls[Math.floor(Math.random() * 2)]

// Decide what string to show for each quest status
var menu = "#bSabitrama#k\r\n\r\nLots of medicinal herbs in this forest. Nothing makes me happier than finding new herbs here!\r\n\r\n";

// DietMed line
if (plr.getQuestStatus(Quest_DietMed) === 0 && plr.getLevel() >= 25) {
    menu += "#r#eQUEST AVAILABLE#k#n#l\r\n#L0##bSabitrama and the Diet Medicine#k#l";
} else if ((plr.checkQuestData(Quest_DietMed, "1_01") || (plr.checkQuestData(Quest_DietMed, "1_11") && plr.itemCount(4031020) === 0)) && plr.getLevel() >= 25) {
    menu += "#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bSabitrama and the Diet Medicine (In Progress)#k#l";
} else if (plr.checkQuestData(Quest_DietMed, "1_11") && plr.getLevel() >= 25 && plr.itemCount(4031020) >= 1) {
    menu += "#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bSabitrama and the Diet Medicine (Ready to complete.)#k#l";
} else if (plr.checkQuestData(Quest_DietMed, "1_00") && plr.getLevel() >= 50) {
    // DietMed done, now see Anti-AgingMed line
    if (plr.getQuestStatus(Quest_AgingMed) === 0) {
        menu += "#r#eQUEST AVAILABLE#k#n#l\r\n#L1##bSabitrama's Anti-Aging Medicine#k#l";
    } else if ((plr.checkQuestData(Quest_AgingMed, "2_01") || (plr.checkQuestData(Quest_AgingMed, "2_11") && plr.itemCount(4031032) === 0)) && plr.getLevel() >= 50) {
        menu += "#r#eQUEST IN PROGRESS#k#n#l\r\n#L1##bSabitrama's Anti-Aging Medicine (In Progress)#k#l";
    } else if (plr.checkQuestData(Quest_AgingMed, "2_11") && plr.getLevel() >= 50 && plr.itemCount(4031032) >= 1) {
        menu += "#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L1##bSabitrama's Anti-Aging Medicine (Ready to complete.)#k#l";
    } else {
        npc.sendOk("It's you!! Thanks to the herbs you got me, the medicine is well on its way. It should be done pretty soon. Thanks again for your help.");
    }
} else {
    npc.sendOk("Lots of medicinal herbs in this forest. Nothing makes me happier than finding new herbs here!");
}

npc.sendSelection(menu);
var sel = npc.selection();

if (sel === 0) { // Diet Medicine
    if (plr.getQuestStatus(Quest_DietMed) === 0) {
        npc.sendBackNext("Wait, hold on one second. I am an herb-collector traveling around the world finding herbs. I'm looking for useful medicinal herbs around this area. It's been hard finding those these days ... so, hey, have found a place where the herbs run aplenty?", false, true)
        if (npc.sendYesNo("Actually I've found a place where you can find good medicinal herbs. It's at a forest not too far from here ... lots of obstacles around the area but I can tell in the end there will be goods available for us to use ... so what do you think? Do you want to go there in place of me?")) {
            plr.startQuest(Quest_DietMed);
            plr.setQuestData(Quest_DietMed, "1_01")
            npc.sendBackNext("Thank you. The place where you can find the mysterious herb is actually the place you've been to before, which is #m101000000#. I heard someone's accepting an entrance fee at the entrance ... you have the mesos to go in, right? Sorry but I've spent all my money traveling so I'm afraid I can't help you on that ...", true, true)
            npc.sendBackNext("Yes! I'll explain to you about the herb you'll be getting for me. Remember there are similar herbs around so make sure you know this. The herb you'll need to get is #b#t4031020##k, and the flower looks like this #i4031020#. Look carefully and please get the same one as I described for you.", true, false)
        } else {
            npc.sendOk("Really. You look like you can just breeze through there ... please come back here when you have time. I'll be waiting for you.")
        }
    }

    if (plr.checkQuestData(Quest_DietMed, "1_01") || (plr.checkQuestData(Quest_DietMed, "1_11") && plr.itemCount(4031020) === 0)) {
        npc.sendOk("You haven't gotten the herb yet. The herb you need to get is #b#t4031020##k. The roots look like this #i4031020#. Remember it and get it from #p1032003# in #m101000000#.")
    }

    if (plr.checkQuestData(Quest_DietMed, "1_11") && plr.itemCount(4031020) >= 1) {
        npc.sendBackNext("Ohhh ... this is IT! With #b#t4031020##k, I can finally make the diet medicine!! Hahaha, if you ever feel like you have gained weight, feel free to find me, because by then, I may have a special medicine in place for just that!", true, true)
        if (plr.canHold(selectedScroll)) {
            plr.completeQuest(Quest_DietMed)
            plr.setQuestData(Quest_DietMed, "1_00")
            plr.removeItemsByID(4031020, 1)
            plr.giveEXP(1000)
            plr.giveFame(1)
            plr.giveItem(selectedScroll, 1)
            npc.sendOk("Oh, I almost forgot. Since you helped me out, I should thank you for your hard work. Here, take this scroll. My brother made this for me a while back, and it adds to the guarding abilities of the armor. Hopefully you'll use it well. And from here on out, #p1032003# will let you in free. Thanks for your help...\r\n\r\n#e#rREWARD:#k\r\n+1000 EXP\r\n+1 Fame\r\n+#i" + selectedScroll + "# #t" + selectedScroll + "#")
        } else {
            npc.sendOk("You don't have enough space in your inventory. Please make space and talk to me again.")
        }
    }
}

if (sel === 1) { // Anti-Aging Medicine
    if (plr.getQuestStatus(Quest_AgingMed) === 0) {
        if (npc.sendYesNo("Ohhh, you're the traveler that helped me out a lot the other day! I made the diet medicine with the herbs you got me and made some money ... and this time, I'd like to make a different kind of a medicine. What do you think? Do you want to help me out one more time?")) {
            plr.startQuest(Quest_AgingMed)
            plr.setQuestData(Quest_AgingMed, "2_01")
            npc.sendBackNext("Thank you. The place where you can find the mysterious herb is actually the place you've been to before, #m101000000#. I heard someone's accepting an entrance fee at the entrance ... you have the mesos to go in, right? This time you'll be going in much deeper than before so be prepared ...", true, true)
            npc.sendBackNext("Yes! I'll explain to you about the herb you'll be getting for me. Remember there are similar herbs around so make sure you know this. The herb you'll need to get is #b#t4031032##k, and the root looks like this #i4031032#. Look carefully and please get the same one as I described for you.", true, false)
        } else {
            npc.sendOk("Really. You look like you can just breeze through there ... please come back here when you have time. I'll be waiting for you.")
        }
    }

    if (plr.checkQuestData(Quest_AgingMed, "2_01") || (plr.checkQuestData(Quest_AgingMed, "2_11") && plr.itemCount(4031032) === 0)) {
        npc.sendOk("You haven't gotten the herb yet. The herb you need to get is #b#t4031032##k. The roots look like this #i4031032#. Remember it and get it from #p1032003# in #m101000000#.")
    }

    if (plr.checkQuestData(Quest_AgingMed, "2_11") && plr.itemCount(4031032) >= 1) {
        npc.sendBackNext("Ohhh ... this is IT! With #b#t4031032##k, I can finally make the anti-aging medicine!!! Hahaha, if you ever become old and weak, find me. By then I may have a special medicine for just that!", true, true)
        if (plr.canHold(4021009)) {
            plr.completeQuest(Quest_AgingMed)
            plr.setQuestData(Quest_AgingMed, "2_00")
            plr.removeItemsByID(4031032, 1)
            plr.giveEXP(3000)
            plr.giveFame(2)
            plr.giveItem(4021009, 1)
            npc.sendOk("Oh, I almost forgot. Since you helped me out, I should thank you for your hard work ... #b#t4021009##k is something I found at the very bottom of a valley a lont time ago in the middle of a journey. It'll probably help you down the road. I also boosted up your fame level and from here on out, #p1032003# may let you in for free. Well, so long...\r\n\r\n#e#rREWARD:#k\r\n+3000 EXP\r\n+2 Fame\r\n+#i4021009# #t4021009#")
        } else {
            npc.sendOk("You don't have enough space in your etc. inventory. Please make space and talk to me again.")
        }
    }
}