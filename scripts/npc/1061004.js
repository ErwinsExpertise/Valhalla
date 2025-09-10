var Quest_Sauna_Robe = 1000600;

if (plr.checkQuestData(Quest_Sauna_Robe, "1") && plr.getLevel() >= 30) {
    npc.sendSelection("Now what exactly is in this book that makes my dad take care of it so much? I want to know what's inside, but I don't think I'll understand it one bit...\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bA Clue to the Secret Book (Ready to complete.)#k#l");
    var sel = npc.selection();
    if (sel === 0) {
        npc.sendBackNext("Hey who are you? You know my dad well? Ah... you want this red book, huh?? I see...my daddy likes this book more than he likes me! This book ... NO, I can't give you this book!! (stomach growling...) Ahhhh!", false, true);
        npc.sendBackNext("No... no... I'm NOT hungry ...! Dang, all, alright. I'll give you back the book. BUT! Not for free! I am really starving, and I need some food, so...if you get me something to eat, the book is yours. I promise!", false, true);
        npc.sendBackNext("I want... #b50 #t4000029#s#k and #p1010100#'s #bUnagi Special#k, along with a #b#t4031015##k. #p1010100# is my awesome friend that lives in #m100000000#. Ask her for the Unagi Special and she'll make it for you.", false, true);
        if (npc.sendYesNo("Oh yeah! The fairies from #m101000000# probably have #b#t4031015##k. I usually ate it at #m101000000#. If you ever get hungry on the way back and eat a couple of those... my dad's book is going to #o3230100#, so you better take care of that food!!!")) {
            plr.setQuestData(Quest_Sauna_Robe, "2");
        } else {
            npc.sendOk("Well... I guess you don't want the book then...");
        }
    }
} else if ((plr.checkQuestData(Quest_Sauna_Robe, "2")) || (plr.checkQuestData(Quest_Sauna_Robe, "3") && !plr.haveItem(4031014, 1) || !plr.haveItem(4031015, 1) || !plr.haveItem(4000029, 50, false, true)) && plr.getLevel() >= 30) {
    npc.sendSelection("Now what exactly is in this book that makes my dad take care of it so much? I want to know what's inside, but I don't think I'll understand it one bit...\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bHungry Ronnie (In Progress)#k#l");
    var sel = npc.selection();
    if (sel === 0) {
        npc.sendOk("You didn't get all my food yet?? Bring #b50 #t4000029#s#k, #p1010100# from #m100000000#'s #bUnagi Special#k, and #b#t4031015##k from the fairies of #m101000000#, and I'll give you back my dad's book. I promise~!!");
    }
} else if (plr.checkQuestData(Quest_Sauna_Robe, "3") && plr.haveItem(4031014, 1) && plr.haveItem(4031015, 1) && plr.haveItem(4000029, 50, false, true) && plr.getLevel() >= 30) {
    npc.sendSelection("Now what exactly is in this book that makes my dad take care of it so much? I want to know what's inside, but I don't think I'll understand it one bit...\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bHungry Ronnie (Ready to complete.)#k#l");
    var sel = npc.selection();
    if (sel === 0) {
        plr.setQuestData(Quest_Sauna_Robe, "4");
        plr.giveItem(4031014, -1);
        plr.giveItem(4031015, -1);
        plr.giveItem(4000029, -50);
        plr.giveItem(4031016, 1);
        plr.giveEXP(300);
        npc.sendOk("Wow...! You DID bring all that food!!!! Sweeeeeeet!!! Thank you so, so, soo much!! Oh yeah, a promise is a promise...here's my dad's book. I have no idea what this book is about, but...why is he so gaga over this anyway??");
    }
} else if (plr.checkQuestData(Quest_Sauna_Robe, "e") && plr.getLevel() >= 30) {
    npc.sendOk("You got him the book, right? I still don't think he even cares about me. I should just go back to the fairy town and stay there.");
} else {
    npc.sendOk("Now what exactly is in this book that makes my dad take care of it so much? I want to know what's inside, but I don't think I'll understand it one bit...");
}