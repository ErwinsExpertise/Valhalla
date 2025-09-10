// Brown Skullcap, Green Skullcap, Red Headband, Green Headband, Yellow Headband, Blue Headband
var hats = [1002008, 1002053, 1002014, 1002070, 1002068, 1002071];
var rnd = Math.floor(Math.random() * 6);
var selectedHat = hats[rnd];
var Quest_Lucas_Letter = 4;

if (plr.checkQuestData(Quest_Lucas_Letter, "")) {
    // Offer quest
    npc.sendBackNext("Ahh...I'm getting worried...because I need to get this letter to #p12000# fast. It's an urgent matter so I need to let him know of this ASAP. Too bad I have things to do here for a while so I won't be leaving this spot anytime soon...", false, true);
    if (npc.sendYesNo("I'm sorry but can you get this #bletter#k to #r#p12000##k from #b#m1010000##k? I have a lot of things to do here so I have to stay here for now. It'll only take a minute...")) {
        plr.startQuest(Quest_Lucas_Letter);
        plr.setQuestData(Quest_Lucas_Letter, "1");
        plr.giveItem(4031000, 1);
        npc.sendBackNext("You're gonna do it? Thank goodness. Now I can breathe a sigh of relief. Here's my letter, and please get this to the town chief that's at #b#m1010000##k.", true, true);
        npc.sendOk("Head northeast and soon you'll find #b#m1010000##k. #p12000# is the town chief of #m1010000#. He should be in front of the department store probably taking a walk. Please get #p12000#'s reply letter fast!");
    } else {
        npc.sendOk("Are you really busy? Ahhhh, this is not good. If I don't get this to the town chief... anyway if you have some free time please come back and talk to me.");
    }
} else if (plr.checkQuestData(Quest_Lucas_Letter, "1")) {
    if (plr.haveItem(4031000, 1)) {
        npc.sendOk("Haven't met #r#p12000##k from #m1010000# yet? Please send him my letter. It's urgent. I need to get a reply from #p12000# quickly...");
    } else {
        npc.sendOk("You lost my letter!! You should have been more careful. Here's the letter again. Please make sure you don't lose the letter, since there are a lot monsters around this area.");
        plr.giveItem(4031000, 1);
    }
} else if (plr.checkQuestData(Quest_Lucas_Letter, "2")) {
    if (plr.haveItem(4031001, 1)) {
        npc.sendBackNext("Here is a hat. It has a Lv. limitation, but I think you are strong enough to wear this. I hope this can help you. Thanks!!!\r\n\r\n#e#rREWARD:#k\r\n+10 EXP\r\n+1 Random Level 5 Hat", false, true);
        plr.giveEXP(10);
        plr.giveItem(selectedHat, 1);
        plr.removeItemsByID(4031001, 1);
        plr.setQuestData(Quest_Lucas_Letter, "3");
        plr.completeQuest(Quest_Lucas_Letter);
    } else {
        npc.sendOk("The letter from #r#p2103##k should be here by now. What happened...? Someone please let me know what's going on here...");
    }
} else if (plr.checkQuestData(Quest_Lucas_Letter, "3")) {
    npc.sendOk("Now I can deliver important news to the town. Thanks a lot. Are you using the hat, I gave you? How is that? Pretty good, huh?");
} else {
    npc.sendOk("I should bring this to #m1010000#... What should I do... hmm...");
}