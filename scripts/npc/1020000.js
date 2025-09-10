var Quest_Blackbull1 = 1000100;
var Quest_Blackbull2 = 1000101;
var shields = [1092001, 1092000];
var scrolls = [2044002, 2043002, 2043102, 2043202, 2044102, 2044202, 2044302, 2044402, 2043702, 2043802, 2044502, 2044602, 2043302, 2044702];
var shieldRnd = Math.floor(Math.random() * 2);
var selectedShield = shields[shieldRnd];
var selectedScroll;
var dialog = "Okay, now choose the scroll of your liking ... The odds of winning are 10% each.\r\n";

if (plr.checkQuestData(Quest_Blackbull1, "")) {
    npc.sendSelection(
        "Our family grew, and I'll have to fix the house to make it bigger, but I need materials to do so...\r\n\r\n#r#eQUEST AVAILABLE#k#n#l\r\n#L0##bFixing \"Blackbull's\" House#k#l"
    );
    var sel = npc.selection();
    if (sel === 0) {
        if (npc.sendYesNo("Can you get me #b#e30#n #b#t4000003#es#k and #b#e50#n #t4000018#s#k? I'm trying to remodel my house and make it bigger ... If you can do it, I'll hook you up with a nice #bshield#k that I don't really need ... You'll get plenty if you take down the ones that look like trees.")) {
            plr.startQuest(Quest_Blackbull1);
            plr.setQuestData(Quest_Blackbull1, "w");
        }
    }
} else if (plr.checkQuestData(Quest_Blackbull1, "w") && (!plr.haveItem(4000003, 30, false) || !plr.haveItem(4000018, 50, false))) {
    npc.sendSelection(
        "Our family grew, and I'll have to fix the house to make it bigger, but I need materials to do so...\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bFixing \"Blackbull's\" House (In Progress)#k#l"
    );
    var sel = npc.selection();
    if (sel === 0) {
        npc.sendOk("Looks like you haven't gotten all the materials needed. Please get #b30 #t4000003#es#k and #b50 #t4000018#s#k.");
    }
} else if (plr.checkQuestData(Quest_Blackbull1, "w") && plr.haveItem(4000003, 30, false) && plr.haveItem(4000018, 50, false)) {
    npc.sendSelection(
        "Our family grew, and I'll have to fix the house to make it bigger, but I need materials to do so...\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bFixing \"Blackbull's\" House (Ready to complete.)#k#l"
    );
    var sel = npc.selection();
    if (sel === 0) {
        npc.sendOk("Incredible! You must be someone special to get that many. Hmm ... alright, the shield is yours. It's my favorite one. Please take good care of it.\r\n\r\n#e#rREWARD:#k\r\n+50 EXP\r\n+1 Random Warrior Shield (Lv. 15 or Lv. 25)");
        plr.setQuestData(Quest_Blackbull1, "end");
        plr.completeQuest(Quest_Blackbull1);
        plr.gainItem(4000003, -30);
        plr.gainItem(4000018, -50);
        plr.gainItem(selectedShield, 1);
        plr.giveExp(50);
    }
} else if (plr.checkQuestData(Quest_Blackbull1, "end") && plr.checkQuestData(Quest_Blackbull2, "") && plr.getLevel() >= 30) {
    npc.sendSelection(
        "Our family grew, and I'll have to fix the house to make it bigger, but I need materials to do so...\r\n\r\n#r#eQUEST AVAILABLE#k#n#l\r\n#L0##bBuilding a New House For \"Blackbull\"#k#l"
    );
    var sel = npc.selection();
    if (sel === 0) {
        npc.sendBackNext("Hey, it's you! Got pretty famous since the last time I saw you, huh? Well thanks to you, I got my house fixed just fine. But hmm ... there's a problem ... all my relatives from #m100000000# want to move to this town. I need to build a new house for them, but I don't even have the materials to build with ...", false, true);
        npc.sendBackNext("I'll need #b100 #t4000022#s#k, #b30 #t4003000#s#k, and #b30 #b#t4003001#s#k. But with only these ... well a couple of days ago, a deed to the land that I purchased disappeared ... my son had it on the way to #m105040300# when he got attacked by the monsters.", true, true);
        if (npc.sendYesNo("A group of reptiles that live in the forests called #r#o3230100##k took the deed to the land. Can you help me get that and the necessary materials to build the house? If so, then you'll be handsomely rewarded for your work ... good luck!")) {
            plr.startQuest(Quest_Blackbull2);
            plr.setQuestData(Quest_Blackbull2, "p0");
        }
    }
} else if (plr.checkQuestData(Quest_Blackbull1, "end") && plr.checkQuestData(Quest_Blackbull2, "p0") && plr.getLevel() >= 30 && (!plr.haveItem(4000022, 100, false) || !plr.haveItem(4003000, 30, false) || !plr.haveItem(4003001, 30, false) || !plr.haveItem(4001004, 1, false))) {
    npc.sendSelection(
        "Our family grew, and I'll have to fix the house to make it bigger, but I need materials to do so...\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bBuilding a New House For \"Blackbull\" (In Progress)#k#l"
    );
    var sel = npc.selection();
    if (sel === 0) {
        npc.sendOk("Looks like you haven't gotten all the materials needed. Please get #b100 #t4000022#s, 30 #b#t4003000#s, 30 #b#t4003001#s and the lost deed to the land#k. Do it fast, before they eat it ...");
    }
} else if (plr.checkQuestData(Quest_Blackbull1, "end") && plr.checkQuestData(Quest_Blackbull2, "p0") && plr.getLevel() >= 30 && plr.haveItem(4000022, 100, false) && plr.haveItem(4003000, 30, false) && plr.haveItem(4003001, 30, false) && plr.haveItem(4001004, 1, false)) {
    npc.sendSelection(
        "Our family grew, and I'll have to fix the house to make it bigger, but I need materials to do so...\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bBuilding a New House For \"Blackbull\" (Ready to complete.)#k#l"
    );
    var sel = npc.selection();
    if (sel === 0) {
        npc.sendNext("THIS is the deed to the land that my son lost! And you even brought all the necessary materials to build the house! Thank you so much ... my relatives can all move in and live in #m102000000#! As a sign of appreciation ...");
        
        var job = plr.job();
        var pickList = "";
        var base = 0;
        if (job == 0 || job == 100 || job == 110 || job == 111 || job == 112 || job == 120 || job == 121 || job == 122 || job == 130 || job == 131 || job == 132) {
            pickList = "#L0##b#i2044002##t2044002##k#l\r\n#L1##b#i2043002##t2043002##k#l\r\n#L2##b#i2043102##t2043102##k#l\r\n#L3##b#i2043202##t2043202##k#l\r\n#L4##b#i2044102##t2044102##k#l\r\n#L5##b#i2044202##t2044202##k#l\r\n#L6##b#i2044302##t2044302##k#l\r\n#L7##b#i2044402##t2044402##k#l";
            base = 0;
        } else if (job == 200 || job == 210 || job == 211 || job == 212 || job == 220 || job == 221 || job == 222 || job == 230 || job == 231 || job == 232) {
            pickList = "#L0##b#i2043702##t2043702##k#l\r\n#L1##b#i2043802##t2043802##k#l";
            base = 8;
        } else if (job == 300 || job == 310 || job == 311 || job == 312 || job == 320 || job == 321 || job == 322) {
            pickList = "#L0##b#i2044502##t2044502##k#l\r\n#L1##b#i2044602##t2044602##k#l";
            base = 10;
        } else if (job == 400 || job == 410 || job == 411 || job == 412 || job == 420 || job == 421 || job == 422) {
            pickList = "#L0##b#i2043302##t2043302##k#l\r\n#L1##b#i2044702##t2044702##k#l";
            base = 12;
        } else {
            npc.sendOk("I'm a GM and I deserve nothing because I can just give any item to myself anyway");
        }
        
        npc.sendSelection("Okay, now choose the scroll of your liking ... The odds of winning are 10% each.\r\n" + pickList);
        var choice = npc.selection();
        selectedScroll = scrolls[choice + base];
        
        plr.setQuestData(Quest_Blackbull2, "pe");
        plr.completeQuest(Quest_Blackbull2);
        plr.gainItem(4000022, -100);
        plr.gainItem(4003000, -30);
        plr.gainItem(4003001, -30);
        plr.gainItem(4001004, -1);
        plr.gainItem(selectedScroll, 1);
        plr.giveExp(1000);
        plr.giveMesos(15000);
        plr.giveFame(2);
        npc.sendOk("Hopefully the scroll I gave you helped. Here's also a little bit of money if that helps. I'll never forget the fact that you helped me. For that, I'll be spreading the good news about your good deed all over the town. What do you think?? Anyway thank you so much for helping me out. We'll probably meet again ...\r\n\r\n#e#rREWARD:#k\r\n+1000 EXP\r\n+2 Fame\r\n#i" + selectedScroll + "##t" + selectedScroll + "#");
    }
} else if (plr.checkQuestData(Quest_Blackbull1, "end") && plr.checkQuestData(Quest_Blackbull2, "pe")) {
    npc.sendOk("Hey, it's you! Thanks to you, the building of the house for my cousins are well on their way. You should come check it out when it's completed.");
} else {
    npc.sendOk("Our family grew, and I'll have to fix the house to make it bigger, but I need materials to do so...");
}