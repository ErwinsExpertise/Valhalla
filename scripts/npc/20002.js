var Quest_Biggs = 100;
var daggers = [1332005, 1332007]; // Razor, Fruit Knife
var rnd = Math.floor(Math.random() * 2);
var selectedDagger = daggers[rnd];

function infoText() {
	var lines = [];
	lines.push("I can't stay in this town forever. Is there anyone who can help me?");
	return lines.join("\r\n");
}

if (plr.checkQuestData(Quest_Biggs, "end")) {
	npc.sendOk("Thanks for the last time. Now, if I get just little more money, I probably can start the business. Are you still using the weapon, which I gave you?");
} else if (plr.checkQuestData(Quest_Biggs, "f")) {
	if (plr.itemCount(4000001) >= 10 && plr.itemCount(4000000) >= 30) {
		plr.setQuestData(Quest_Biggs, "end");
		plr.completeQuest(Quest_Biggs);
		plr.removeItemsByID(4000001, 10);
		plr.removeItemsByID(4000000, 30);
		plr.giveEXP(30);
		plr.giveItem(selectedDagger, 1);
		npc.sendOk("Oh wow! You brought them all!! Sweet! Here's an item like I promised. I don't really need it anyway, so take it!\r\n\r\n#e#rREWARD:#k\r\n+30 EXP\r\n+#i" + selectedDagger + "# #t" + selectedDagger + "#");
	} else {
		npc.sendSelection(infoText() + "\r\n\r\n#r#eQUEST IN PROGRESS#k#n#l\r\n#L0##bBigg's Collection of Items (In Progress)#k#l");
		if (npc.selection() === 0) {
			npc.sendOk("You don't have #e10#n #b#t4000001#s#k and #e30#n#k #b#t4000000#s#k, right? Don't worry too much about it. I'll be staying here for a while.");
		}
	}
} else if (plr.checkQuestData(Quest_Biggs, "") || plr.checkQuestData(Quest_Biggs, "info")) {
	npc.sendSelection(infoText() + "\r\n\r\n#r#eQUEST AVAILABLE#k#n#l\r\n#L0##bBigg's Collection of Items#k#l");
	if (npc.selection() === 0) {
		if (npc.sendYesNo("I'll give you something nice if you get me #b#e10#n #t4000001#s#k and #b#e30#n #t4000000#s#k! You can get it by taking down the monsters, but ... looking at you, I'm not sure if you're up for the challenge...")) {
			plr.startQuest(Quest_Biggs);
			plr.setQuestData(Quest_Biggs, "f");
			npc.sendOk("Thanks... I'll be waiting here!");
		} else {
			npc.sendOk("Hmph.. I was right. You aren't up for the challenge!");
		}
	}
} else {
	npc.sendOk("I can't stay in this town forever. Is there anyone who can help me?");
}