const Quest_Nina_Sen = 3;

if (plr.checkQuestData(Quest_Nina_Sen, "")) {
    npc.sendSelection(
        "What is #p2001# doing, now?...\r\n\r\n" +
        "#r#eQUEST AVAILABLE#k#n#l\r\n#L0##bWhat Sen wants to eat#k#l"
    );
    
    if (npc.selection() === 0) {
        if (npc.sendYesNo("Oh, a traveler!! Nice, right on time... I have a favor to ask, will you do it for me? Go a little more to the right and you'll find a #bhouse with the orange roof#k.")) {
            plr.startQuest(Quest_Nina_Sen);
            plr.setQuestData(Quest_Nina_Sen, "1");
            npc.sendOk("That's my house. I have a little brother #r#p2001##k that's at home, so can you please ask him what he wants for dinner? Stand in front of the door, press the #bup arrow#k and then you'll be able to enter the house.");
        } else {
            npc.sendOk("Oh, you must be busy. Wouldn't it be fun to get to know some others, though?");
        }
    }
}
else if (plr.checkQuestData(Quest_Nina_Sen, "1")) {
    npc.sendSelection(
        "What is #p2001# doing, now?...\r\n\r\n" +
        "#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bWhat Sen wants to eat (In Progress)#k#l"
    );
    
    if (npc.selection() === 0) {
        npc.sendOk("Haven't met #p2001# yet? Press the up arrow in front of the door!");
    }
}
else if (plr.checkQuestData(Quest_Nina_Sen, "2")) {
    npc.sendSelection(
        "What is #p2001# doing, now?...\r\n\r\n" +
        "#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bReturning to Nina (Ready to complete.)#k#l"
    );
    
    if (npc.selection() === 0) {
        npc.sendOk("He wants the mushroom soup? I guess that's our dinner right there then. Thanks for doing me a favor.\r\n\r\n#e#rREWARD:#k\r\n+20 EXP");
        plr.setQuestData(Quest_Nina_Sen, "3");
        plr.completeQuest(Quest_Nina_Sen);
        plr.giveEXP(20);
    }
}
else if (plr.checkQuestData(Quest_Nina_Sen, "3")) {
    npc.sendOk("I will make a mushroom soup for #p2001#. Is there a fresh mushroom?");
}
else {
    npc.sendOk("What is #p2001# doing, now?");
}