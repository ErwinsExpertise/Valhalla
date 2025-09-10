// Sen on Maple Island — Nina & Sen quest (ID 3)

var Quest_Nina_Sen = 3;

if (plr.checkQuestData(Quest_Nina_Sen, "1")) {
    // Ready to complete
    npc.sendSelection(
        "There is nothing to eat in here~ oh...\r\n\r\n"
        + "#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n"
        + "#L0##bWhat Sen wants to eat (Ready to complete.)#k#l"
    )
    var sel = npc.selection()
    if (sel === 0) {
        npc.sendBackNext("Ahh, soooo hungry. Where's my sister?!! I was gonna ask her to make me a mushroom soup. Soooo hungry!!", false, true)
        if (npc.sendYesNo("Please tell my sister I really really want #bmushroom soup#k for dinner!\r\n\r\n#e#rREWARD:#k\r\n+7 EXP")) {
            plr.setQuestData(Quest_Nina_Sen, "2")
            plr.giveEXP(7)
        } else {
            npc.sendOk("Awww... You aren't going to tell her?")
        }
    }
} else if (plr.checkQuestData(Quest_Nina_Sen, "2")) {
    // In progress returning to Nina
    npc.sendSelection(
        "There is nothing to eat in here~ oh...\r\n\r\n"
        + "#r#eQUEST IN PROGRESS#k#n#l\r\n"
        + "#L0##bReturning to Nina (In Progress)#k#l"
    )
    var sel = npc.selection()
    if (sel === 0) {
        npc.sendOk("What did my sister say? oh... I am so hungry...")
    }
} else if (plr.checkQuestData(Quest_Nina_Sen, "3")) {
    npc.sendOk("What did my sister say? oh... I am so hungry...")
} else {
    npc.sendOk("There is nothing to eat in here~ oh...")
}