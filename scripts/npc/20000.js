var Quest_John1 = 5500;
var Quest_John2 = 5501;
var Quest_John3 = 5502;

var text = "Is there anyone who can help me? Well... I would like to go by myself, but as you can see, I have lots of stuff to do...";

if (plr.level() >= 15) {
    if (plr.checkQuestData(Quest_John1, "1")) {
        if (plr.itemCount(4031025) >= 10) {
            text += "\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bJohn's Pink Flower Basket (Ready to complete.)#k#l";
        } else {
            text += "\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bJohn's Pink Flower Basket (In Progress)#k#l";
        }
    } else if (plr.checkQuestData(Quest_John1, "")) {
        text += "\r\n\r\n#r#eQUEST AVAILABLE#k#n#l\r\n#L0##bJohn's Pink Flower Basket#k#l";
    } else if (plr.checkQuestData(Quest_John1, "2") && plr.level() <= 29) {
        npc.sendOk("You are the one, who brought the flower to me. Again, thanks a lot.. Feel free to stay in this town.");
    }
}
if (plr.level() >= 30) {
    if (plr.checkQuestData(Quest_John2, "1")) {
        if (plr.itemCount(4031026) >= 20) {
            text += "\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L1##bJohn's Present (Ready to complete.)#k#l";
        } else {
            text += "\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L1##bJohn's Present (In Progress)#k#l";
        }
    } else if (plr.checkQuestData(Quest_John2, "")) {
        text += "\r\n\r\n#r#eQUEST AVAILABLE#k#n#l\r\n#L1##bJohn's Present#k#l";
    } else if (plr.checkQuestData(Quest_John1, "2") && plr.checkQuestData(Quest_John2, "2") && plr.level() <= 59) {
        npc.sendOk("You are the one, who brought the flower to me. Again, thanks a lot.. Feel free to stay in this town.");
    }
}
if (plr.level() >= 60) {
    if (plr.checkQuestData(Quest_John3, "1")) {
        if (plr.itemCount(4031028) >= 30) {
            text += "\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L2##bJohn's Last Present (Ready to complete.)#k#l";
        } else {
            text += "\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L2##bJohn's Last Present (In Progress)#k#l";
        }
    } else if (plr.checkQuestData(Quest_John3, "")) {
        text += "\r\n\r\n#r#eQUEST AVAILABLE#k#n#l\r\n#L2##bJohn's Last Present#k#l";
    } else if (plr.checkQuestData(Quest_John1, "2") && plr.checkQuestData(Quest_John2, "2") && plr.checkQuestData(Quest_John3, "2")) {
        npc.sendOk("You are the one, who brought the flower to me. Again, thanks a lot.. Feel free to stay in this town.");
    }
}

npc.sendSelection(text);
var sel = npc.selection();

if (sel == 0 && plr.checkQuestData(Quest_John1, "2") && plr.level() <= 29) {
    npc.sendOk("You are the one, who brought the flower to me. Again, thanks a lot.. Feel free to stay in this town.");
}
if (sel == 1 && plr.checkQuestData(Quest_John1, "2") && plr.checkQuestData(Quest_John2, "2") && plr.level() <= 59) {
    npc.sendOk("You are the one, who brought the flower to me. Again, thanks a lot.. Feel free to stay in this town.");
}
if (sel == 2 && plr.checkQuestData(Quest_John1, "2") && plr.checkQuestData(Quest_John2, "2") && plr.checkQuestData(Quest_John3, "2")) {
    npc.sendOk("You are the one, who brought the flower to me. Again, thanks a lot.. Feel free to stay in this town.");
}

if (sel == 0 && plr.checkQuestData(Quest_John1, "")) {
    if (npc.sendYesNo("How's traveling these days? I actually have a favor to ask you ... this time, my wedding anniversary is coming up and I need flowers. Can you get them for me?")) {
        plr.startQuest(Quest_John1);
        plr.setQuestData(Quest_John1, "1");
        npc.sendNext("Thank you! This time I'd like to give my wife #b#t4031025##k ... It has a very pleasant scent, and I heard it's found deep in the forest ... I heard the place where it exists doesn't let everyone in; only a select few, I think. Something about #p1061006# at #m105040300# and something something ...");
    } else {
        npc.sendOk("I understand ... my wedding anniversary is coming up and I am screwed! Please come back if you have some spare time.");
    }
} else if (sel == 0 && plr.checkQuestData(Quest_John1, "1")) {
    if (plr.itemCount(4031025) >= 10) {
        npc.sendNext("Ohhh ... you got me #b10 #b#t4031025#s#k~! This is awesome ... I can't believe you went deep into the forest and got these flowers ... there's a story about this flower where it supposedly doesn't die for 100 years. With this, I can make my wife happy.");
    } else {
        npc.sendOk("You haven't gotten #b#t4031025##k yet. There's #p1061006# at #m105040300# and I heard that with that you can go to the place where #t4031025##k is. Please go into the forest and collect #t4031025##k for me. I need 10 to make my wedding anniversary");
    }
}

if (sel == 1 && plr.checkQuestData(Quest_John2, "")) {
    if (npc.sendYesNo("Ohhhh, you're the one that helped me out the other day. You look much stronger now. How's traveling these days? I actually have another favor to ask you ... this time, my wife's birthday is coming up and I need more flowers. Can you get them for me?")) {
        plr.startQuest(Quest_John2);
        plr.setQuestData(Quest_John2, "1");
        npc.sendNext("Thank you! This time I'd like to give my wife #b#t4031026##k ... It has a very pleasant scent, and I heard it's found deep in the forest ... I heard the place where it exists doesn't let everyone in; only a selected few, I think. Something about #p1061006# at #m105040300# and something something ...");
    } else {
        npc.sendOk("I understand ... my wife's birthday is coming up and I am screwed! Please come back if you have some spare time.");
    }
} else if (sel == 1 && plr.checkQuestData(Quest_John2, "1")) {
    if (plr.itemCount(4031026) >= 20) {
        npc.sendNext("Ohhh ... you got me #b20 #b#t4031026#s#k~! This is awesome ... I can't believe you went deep into the forest and got these flowers ... there's a story about this flower where it supposedly doesn't die for 500 years. With this, I can make the whole house smell like flowers.");
    } else {
        npc.sendOk("You haven't gotten the #b#t4031026##k yet. There's #p1061006# at #m105040300# and I heard that with that you can go to the place where #t4031026##k is. Please go into the forest and collect #t4031026##k for me. I need 20 to make my wife's birthday present.");
    }
}

if (sel == 2 && plr.checkQuestData(Quest_John3, "")) {
    if (npc.sendYesNo("Ohhh...you're the person that did me huge favors a while ago. You look so much stronger now that I can't even recognize you anymore. By now it looks like you have gone pretty much everywhere. I have one last favor to ask you. Well, my mother passed away a few days ago of old age. I need a special kind of flower for her on her grave ... can you get them for me?")) {
        plr.startQuest(Quest_John3);
        plr.setQuestData(Quest_John3, "1");
        npc.sendNext("Thank you so much! The flowers I want on her grave are called #b#t4031028##k and it's a very rare kind. I heard it's found deep in the forest ... I heard the place where it exists doesn't let everyone in; only a select few, I think. Something about #p1061006# at #m105040300# and something something ...");
    } else {
        npc.sendOk("I see ... my mother loved looking at those flowers while she was alive ... I wished that you could get them for me and for her ... I understand ...");
    }
} else if (sel == 2 && plr.checkQuestData(Quest_John3, "1")) {
    if (plr.itemCount(4031028) >= 30) {
        npc.sendNext("Ohhh ... you got me all #b30 #t4031028#s#k! This is awesome ... I can't believe you went deep into the forest and got these flowers... there's a story about this flower where it supposedly doesn't die for 1000 years and it glows on its own. I can make a nice wreath out of this and bring it to my mother's grave...");
    } else {
        npc.sendOk("You haven't gotten #b#t4031028##k yet. There's #p1061006# at #m105040300# and I heard that with that you can go to the place where #t4031028##k is. Please go into the forest and collect #t4031028##k for me. I need 30 to make a wreath");
    }
}

// Completion logic for Pink Flower Basket
if (sel == 0 && plr.checkQuestData(Quest_John1, "1") && plr.itemCount(4031025) >= 10) {
    plr.setQuestData(Quest_John1, "2");
    plr.removeItemsByID(4031025, 10);
    plr.giveItem(4003000, 30);
    plr.giveEXP(300);
    npc.sendOk("If you have time, why not try going back into the forest? You may find an important item in there. I can't guarantee it since obviously I've never been there before, so please don't come back complaining if all you can find is trash.");
}

// Completion logic for John's Present
if (sel == 1 && plr.checkQuestData(Quest_John2, "1") && plr.itemCount(4031026) >= 20) {
    var jobId = plr.job();
    plr.setQuestData(Quest_John2, "2");
    plr.removeItemsByID(4031026, 20);
    plr.giveEXP(2000);
    if (jobId == 0) {
        plr.giveItem(1082002, 1);
    } else if (jobId >= 100 && jobId <= 132) {
        plr.giveItem(1082036, 1);
    } else if (jobId >= 200 && jobId <= 232) {
        plr.giveItem(1082056, 1);
    } else if (jobId >= 300 && jobId <= 322) {
        plr.giveItem(1082070, 1);
    } else if (jobId >= 400 && jobId <= 422) {
        plr.giveItem(1082045, 1);
    }
    npc.sendOk("If you have time, why not try going back into the forest? You may find an important item in there. I can't guarantee it since obviously I've never been there before, so please don't come back complaining if all you can find is trash.");
}

// Completion logic for John's Last Present
if (sel == 2 && plr.checkQuestData(Quest_John3, "1") && plr.itemCount(4031028) >= 30) {
    plr.setQuestData(Quest_John3, "2");
    plr.removeItemsByID(4031028, 30);
    plr.giveItem(1032014, 1);
    plr.giveEXP(4000);
    npc.sendOk("If you have time, why not try going back into the forest? You may find an important item in there. I can't guarantee it since obviously I've never been there before, so please don't come back complaining if all you can find is trash.");
}