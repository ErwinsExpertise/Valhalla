var Quest_Maya = 1000200;

if (plr.getLevel() >= 15 && plr.checkQuestData(Quest_Maya, "m2")) {
    npc.sendSelection(
        "Man, I want to travel around and stuff. I don't want to be stuck here working. This stinks!! I'm stuck here making potions everyday thanks to my mom opening up a convenient store. This is definitely NOT fun.\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bFinding Sophia (Ready to complete.)#k#l"
    );
    if (npc.selection() === 0) {
        if (npc.sendYesNo("Hmm... So you want to have #b#t4031004##k? That stone is very rare... Alright, get me the items, I tell you... Ready?")) {
            plr.setQuestData(Quest_Maya, "m3");
            npc.sendOk("50#b#e#n #t4000004#s#k, 50#b#e#n #t4000005#s#k, 20#b#e#n #t4000006#s#k, and 1#b#e#n \r\n#t4031005##k. Everything else should be easy to obtain. As for the #t4031005#...  you should ask #r#p1022002##k. He is somewhere around the town.");
        }
    }
} else if (plr.getLevel() >= 15 && (plr.checkQuestData(Quest_Maya, "m3") || plr.checkQuestData(Quest_Maya, "m4"))) {
    npc.sendSelection(
        "Man, I want to travel around and stuff. I don't want to be stuck here working. This stinks!! I'm stuck here making potions everyday thanks to my mom opening up a convenient store. This is definitely NOT fun.\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bMaking a Sparkling Rock (In Progress)#k#l"
    );
    if (npc.selection() === 0) {
        npc.sendOk("If you have 50#b#e#n #t4000004#s#k, 50#b#e#n #t4000005#s#k, 20#b#e#n #t4000006#s#k, and 1#b#e#n #t4031005##k, you can make 1 #b#e#n#t4031004##k... You should ask #r#p1022002##k about #t4031005#... I think he is somewhere around #m102000000#");
    }
} else if (plr.getLevel() >= 15 && plr.checkQuestData(Quest_Maya, "m5")) {
    npc.sendSelection(
        "Man, I want to travel around and stuff. I don't want to be stuck here working. This stinks!! I'm stuck here making potions everyday thanks to my mom opening up a convenient store. This is definitely NOT fun.\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bGetting Arcon's Blood (Ready to complete.)#k#l"
    );
    if (npc.selection() === 0) {
        if (npc.sendYesNo("How did you get all these items?? You must be really good! Especially #t4031005#. Wow... Anyway, good job! Now we can make a #b#t4031004##k")) {
            plr.setQuestData(Quest_Maya, "m6");
            plr.removeItemsByID(4000004, 50);
            plr.removeItemsByID(4000005, 50);
            plr.removeItemsByID(4000006, 20);
            plr.removeItemsByID(4031005, 1);
            plr.giveItem(4031004, 1);
            npc.sendOk("Here, take this, the #b#t4031004##k. By the way, what do you plan on doing with that stone? It is a special item, indeed, but ... unless you're collecting stones, this may be of no use...");
        }
    }
} else if (plr.getLevel() >= 15 && plr.checkQuestData(Quest_Maya, "m6")) {
    if (!plr.itemCount(4031004)) {
        plr.giveItem(4031004, 1);
        npc.sendOk("Here, take this, the #b#t4031004##k. By the way, what do you plan on doing with that stone? It is a special item, indeed, but ... unless you're collecting stones, this may be of no use...");
    } else {
        npc.sendOk("What did you do with #b#t4031004##k? That's a coveted rock and all, but it requires so many items to make that you're the first one to actually gather them all up! Anyway, hope this is put to good use.");
    }
} else if (plr.getLevel() >= 15 && (plr.checkQuestData(Quest_Maya, "m7") || plr.checkQuestData(Quest_Maya, "end"))) {
    npc.sendOk("What did you do with #b#t4031004##k? That's a coveted rock and all, but it requires so many items to make that you're the first one to actually gather them all up! Anyway, hope this is put to good use.");
} else {
    npc.sendOk("Man, I want to travel around and stuff. I don't want to be stuck here working. This stinks!! I'm stuck here making potions everyday thanks to my mom opening up a convenient store. This is definitely NOT fun.");
}