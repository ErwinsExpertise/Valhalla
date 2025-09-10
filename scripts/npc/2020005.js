// Alcaster's rare items shop (stateless conversion)

var item = [2050003, 2050004, 4006000, 4006001];
var cost = [300, 400, 5000, 5000];

npc.sendSimple("What is it? \r\n#L0##bI want to buy something really rare.");

if (npc.selection() === 0) {
    if (plr.getLevel() < 30) {
        npc.sendNext("I am Alcaster the Magician. I have been studying all kinds of magic for over 300 years.");
    } else if (plr.getQuestStatus(3035) !== 2) {
        npc.sendNext("If you decide to help me out, then in return, I'll make the item available for sale.");
    } else {
        let menu = "Thanks to you, #bThe Book of Ancient#k is safely sealed. As a result, I used up about half of the power I have accumulated over the last 800 years...but can now die in peace. Would you happen to be looking for rare items by any chance? As a sign of appreciation for your hard work, I'll sell some items in my possession to you and ONLY you. Pick out the one you want! #b";
        for (let i = 0; i < item.length; i++) {
            menu += "\r\n#L" + i + "##t" + item[i] + "#(Price : " + cost[i] + " mesos)#l";
        }
        npc.sendSimple(menu);

        const select = npc.selection();
        const text1 = [
            "So the item you need is #bHoly Water#k, right? That's The item that cures the state of being sealed and cursed. It's not an easy item to get, but for you, I'll sell it for cheap. It'll cost you #b300 mesos#k per. How many would you like to buy?",
            "So the item you need is #bAll Cure Potion#k, right? That's The item that cures all. It's not an easy item to get, but for you, I'll sell it for cheap. It'll cost you #b400 mesos#k per. How many would you like to buy?",
            "So the item you need is #bThe Magic Rock#k, right? That's The item that possesses magical power and is used for high-quality skills. It's not an easy item to get, but for you, I'II sell it for cheap. It'll cost you #b5000 mesos#k per. How many would you like to buy?",
            "So the item you need is #bThe Summoning Rock#k, right? That's The item that possesses summoning power and is used for high-quality skills. It's not an easy item to get, but for you, I'll sell it for cheap. It'll cost you #b5000 mesos#k per. How many would you like to buy?"
        ];

        const num = npc.askNumber(text1[select], 1, 1, 100);
        if (num > 0) {
            const totalCost = cost[select] * num;
            if (npc.sendYesNo("Are you sure you want to buy #r" + num + " #t" + item[select] + "#(s)#k? It'll cost you " + cost[select] + " mesos per #t" + item[select] + "#, which will cost you #r" + totalCost + "#k mesos total.")) {
                if (plr.mesos() < totalCost) {
                    npc.sendNext("Are you sure you have enough mesos? Please check if you have at least #r" + totalCost + "#k mesos.");
                } else {
                    plr.giveMesos(-totalCost);
                    plr.giveItem(item[select], num);
                    npc.sendNext("Thank you. If you need anything else, come see me anytime. I may have lost a lot of power, but I can still make magical items!");
                }
            } else {
                npc.sendNext("I see. Well, please understand that I carry many different items here. I'm only selling these items to you, so I won't be ripping you off in any way shape or form.");
            }
        }
    }
}