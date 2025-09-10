var Quest_Lucas_Letter = 4;

if (plr.checkQuestData(Quest_Lucas_Letter, "1")) {
    // Letter ready to complete
    npc.sendSelection("A letter from #r#p2103##k should be here, any minute. Is there something wrong? Hmm... \r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bLetter For Lucas (Ready to complete.)#k#l");
    npc.selection();

    if (plr.haveItem(4031000, 1, false, true)) {
        if (npc.sendYesNo("This is definitely the letter from #p2103#! Ohhh, thank you. I was beginning to get worried because the letter didn't get here. Ok! here's the #breply letter#k. Please get this to her. Just head back to #p2103# and you'll be fine.")) {
            plr.setQuestData(Quest_Lucas_Letter, "2");
            plr.gainExp(10);
            plr.gainItem(4031000, -1);
            plr.gainItem(4031001, 1);
            npc.sendOk("So you have gotten the reply from #p2103#. Thanks! and I want to give you something to show my appreciation.\r\n\r\n#e#rREWARD:#k\r\n+10 EXP\r\n+Lucas' Reply");
        } else {
            npc.sendOk("Are you really busy? Ahhhh, this is not good. If I don't get this to Maria... anyway if you have some free time please come back and talk to me.");
        }
    } else {
        npc.sendOk("Do you have the letter for me?... No? You must have lost it somewhere along the way here... I don't blame you, there are a lot of monsters around here. Go back and talk to #p2103# and she will write you another letter! Please hurry!");
    }
} else if (plr.checkQuestData(Quest_Lucas_Letter, "2")) {
    // Reply quest in progress
    npc.sendSelection("A letter from #r#p2103##k should be here, any minute. Is there something wrong? Hmm... \r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bLucas' Reply (In Progress)#k#l");
    npc.selection();

    if (plr.haveItem(4031001, 1, false, true)) {
        npc.sendOk("You didn't meet up with #r#p2103##k yet? Please get her my reply letter... if you lose the letter by any chance, come find me again... I can always write a new one.");
    } else {
        npc.sendOk("You lost my reply letter!! Should have been more careful. Oh well, there are lots of monsters around this area, so it's understandable. Anyway, here's the reply letter. Please be careful this time around.");
        plr.gainItem(4031001, 1);
    }
} else if (plr.checkQuestData(Quest_Lucas_Letter, "3")) {
    npc.sendOk("You safely delivered my reply to #r#p2103##k? Thanks~ so #r#p2103##k gave you a present? Haha... She is one of the best for making a hat....");
} else {
    npc.sendOk("A letter from #r#p2103##k should be here, any minute. Is there something wrong? Hmm... ");
}