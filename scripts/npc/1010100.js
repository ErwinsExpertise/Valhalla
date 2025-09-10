// Constants
var Quest_Sauna_Robe = 1000600;
var Quest_Unagi = 1090901;

// Decide branch based on quest progress (stateless flow)

// CASE 1: Ready to turn in "Hungry Ronnie" only
if (plr.checkQuestData(Quest_Sauna_Robe, "2") && plr.checkQuestData(Quest_Unagi, "")) {
    npc.sendBackNext("This town is made by the group of bowmen. If you want to become a bowman, please meet with #r#p1012100##k... She will help you. What? You don't know #r#p1012100##k? She saved our town long time ago from the monsters. Of course, it is safe now. She is the hero of our town.\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bHungry Ronnie (Ready to complete.)#k#l", false, true)
    npc.sendBackNext("So you came here through a favor by #p1061004#? Hahaha ... hopefully #p1061004# is not in any trouble this time around. Anyway, he wants the #bUnagi special#k, huh? It's pretty easy, so why don't you just sit around and wait for a little bit as I make this dish...", false, true)
    npc.sendBackNext("Oh shoot. I'm lacking #b#t4000013##k and #b#t4000017##k, the most important ingredients for Unagi. Do these really go in the dish? oh of course~~just please make this a secret from #p1061004#, ok?", false, true)
    if (npc.sendYesNo("Anyway I don't have enough ingredients for Unagi. Sorry but can you gather up the ingredients for me? #b50 #t4000013#s and 5 #t4000017#s#k and then the #bUnagi Special#k will be made.")) {
        plr.startQuest(Quest_Unagi)
        plr.setQuestData(Quest_Unagi, "1")
    }
}
// CASE 2: Still gathering ingredients
else if (plr.checkQuestData(Quest_Sauna_Robe, "2") && plr.checkQuestData(Quest_Unagi, "1") && (!plr.haveItem(4000013, 50, false, true) || !plr.haveItem(4000017, 5, false, true))) {
    npc.sendBackNext("This town is made by the group of bowmen. If you want to become a bowman, please meet with #r#p1012100##k... She will help you. What? You don't know #r#p1012100##k? She saved our town long time ago from the monsters. Of course, it is safe now. She is the hero of our town.\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bSecret to Unagi Special (In Progress)#k#l", false, true)
    npc.sendOk("Please get me #b50 #b#t4000013#s and #b5 #t4000017#s#k. Then I'll get you #p1061004#'s favorite, the Unagi Special.")
}
// CASE 3: Ready to turn in "Secret to Unagi Special"
else if (plr.checkQuestData(Quest_Sauna_Robe, "2") && plr.checkQuestData(Quest_Unagi, "1") && plr.haveItem(4000013, 50, false, true) && plr.haveItem(4000017, 5, false, true)) {
    npc.sendBackNext("This town is made by the group of bowmen. If you want to become a bowman, please meet with #r#p1012100##k... She will help you. What? You don't know #r#p1012100##k? She saved our town long time ago from the monsters. Of course, it is safe now. She is the hero of our town.\r\n\r\n#r#eQUEST THAT CAN BE COMPLETED#k#n#l\r\n#L0##bSecret to Unagi Special (Ready to complete.)#k#l", false, true)
    npc.sendBackNext("You got all the ingredients!! I knew you'd be the one to do it... alright, now just wait oooone second. I, Rina, proudly present the #bUnagi special#k!", false, true)
    
    if (plr.haveItem(4031015, 1, false, true)) {
        plr.setQuestData(Quest_Sauna_Robe, "3")
    }
    plr.setQuestData(Quest_Unagi, "2")
    plr.takeItem(4000013, 50)
    plr.takeItem(4000017, 5)
    plr.giveEXP(500)
    plr.giveItem(4031014, 1)
    npc.sendOk("Ok, here it is, the #bUnagi Special#k! You should take this to \r\n#p1061004# before it gets cold. It's #p1061004#'s favorite.")
}
// CASE 4: Player already has Unagi Special but hasn't turned it in
else if ((plr.checkQuestData(Quest_Sauna_Robe, "2") || plr.checkQuestData(Quest_Sauna_Robe, "3")) && plr.checkQuestData(Quest_Unagi, "2") && !plr.haveItem(4031014, 1, false, true)) {
    plr.giveItem(4031014, 1)
    npc.sendOk("Ok, here it is, the #bUnagi Special#k! You should take this to \r\n#p1061004# before it gets cold. It's #p1061004#'s favorite.")
}
else if ((plr.checkQuestData(Quest_Sauna_Robe, "2") || plr.checkQuestData(Quest_Sauna_Robe, "3")) && plr.checkQuestData(Quest_Unagi, "2") && plr.haveItem(4031014, 1, false, true)) {
    npc.sendOk("This town is made by the group of bowmen. If you want to become a bowman, please meet with #r#p1012100##k... She will help you. What? You don't know #r#p1012100##k? She saved our town long time ago from the monsters. Of course, it is safe now. She is the hero of our town.\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bHungry Ronnie (In Progress)#k#l")
    npc.sendOk("You should take the #bUnagi Special#k to #p1061004# before it gets cold. It's #p1061004#'s favorite.")
}
// CASE 5: Fully complete
else if (plr.checkQuestData(Quest_Sauna_Robe, "e")) {
    npc.sendOk("Oh... You are the one who gave the stuff back? So what's up? Is #p1012102# doing fine? If you get to #m100000000# someday, please say hello to #p1012102# for me.")
}
// CASE 6: Default greeting
else {
    npc.sendOk("This town is made by the group of bowmen. If you want to become a bowman, please meet with #r#p1012100##k... She will help you. What? You don't know #r#p1012100##k? She saved our town long time ago from the monsters. Of course, it is safe now. She is the hero of our town.")
}