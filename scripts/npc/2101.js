var Quest_Heena_Sera = 1;

// Helper to pick state
var qData = plr.questData(Quest_Heena_Sera);
var baseTalk = "You must be a new traveler. I will give some of the instruction, which will be very useful. If you want to speak with us, you can just simply double-click us. You can move by pressing #bLeft, Right Key#k and jump by pressing #bSpace Bar#k. Come on~ Try! Also, sometimes, you had to climb up the ladder or use the rope to get to the destination, where you want. You can do that by pressing #bup arrow#k. Please keep this in mind.";

if (qData === "") {
    // QUEST AVAILABLE
    npc.sendSelection(baseTalk + "\r\n\r\n#r#eQUEST AVAILABLE#k#n#l\r\n#L0##bBorrowing Sera's Mirror#k#l");
} else if (qData === "1") {
    // QUEST IN PROGRESS
    npc.sendSelection(baseTalk + "\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bBorrowing Sera's Mirror (In Progress)#k#l");
} else if (qData === "2") {
    // READY TO COMPLETE
    npc.sendSelection(baseTalk + "\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bBringing a Mirror to Heena (Ready to complete.)#k#l");
} else if (qData === "3") {
    npc.sendOk("Now I can take a look at my face with this mirror. It looks alright.");
    end;
} else {
    npc.sendOk(baseTalk);
    end;
}

var sel = npc.selection();
if (sel !== 0) end;

// Branch by quest state
if (qData === "2" && plr.itemCount(4031003) >= 1) {
    npc.sendNext("Oh wow! You brought #p2100#'s mirror! Thank you so so much. Let's see... no skin damage whatsoever...\r\n\r\n#e#rREWARD:#k\r\n+1 EXP");
    plr.setQuestData(Quest_Heena_Sera, "3");
    plr.completeQuest(Quest_Heena_Sera);
    plr.removeItemsByID(4031003, 1);
    plr.giveEXP(1);
    npc.sendOk("If you go right, you will see the shiny spot. We call that a \"Portal\". If you press #bup-arrow#k, you will get to the next place. So long!");
} else if (qData === "2" && plr.itemCount(4031003) < 1) {
    npc.sendOk("Did you lose the mirror? Ask her for it once more.");
} else if (qData === "1") {
    npc.sendOk("Haven't met #r#p2100##k yet? She should be on a hill down on east side...it's pretty close from here so it will be easy to spot her...");
} else if (qData === "") {
    npc.sendNext("You must be the new traveler. Still foreign to this, huh? I'll be giving you important information here and there so please listen carefully and follow along. First if you want to talk to us, #bdouble-click#k us with the mouse.");
    npc.sendBackNext("#bLeft, right arrow#k will allow you to move. Press #bAlt#k to jump. Jump diagonally by combining it with the directional cursors. Try it later.", true, true);
    if (npc.sendYesNo("Man... the sun is literally burning my beautiful skin! It's a scorching day today. Can I ask you for a favor? Can you get me a #bmirror#k from #r#p2100##k, please?")) {
        plr.startQuest(Quest_Heena_Sera);
        plr.setQuestData(Quest_Heena_Sera, "1");
        npc.sendOk("Thank you... #r#p2100##k will be on the hill down on the east side hanging up the laundry. The mirror looks like this #i4031003#.");
    } else {
        npc.sendOk("Don't want to? Hmmm... come back when you change your mind.");
    }
}