// Toy Factory Assistant – stateless flow
if (plr.getQuestStatus(3239) !== 1) {
    npc.sendOk("Lately the mechanical parts have been missing at the Toy Factory, and that really concerns me. I want to ask for help, but you don't seem strong enough to help us out. Who should I ask...?");
}

// --- Outside map 220020000 -------------------------------------
if (plr.getMap() === 220020000) {
    npc.sendNext("Okay, then. Inside this room, you'll see a whole lot of plastic barrels lying around. Strike the barrels to knock them down, and see if you can find the lost #bMachine Parts#k inside. You'll need to collect 10 #bMachine Parts#k and then talk to me afterwards. There's a time limit on this, so go!");

    if (instanceProperties().playerCount > 0) {
        npc.sendNextPrev("I'm sorry, but it seems like someone else is inside looking through the barrels. Only one person is allowed in here, so you'll have to wait for your turn.");
    }

    // empty instance – reset + enter
    npc.sendBackNext(
        "The area is empty. I'll send you in now. You have 20 minutes to gather the pieces.",
        true,
        true
    );
    plr.warp(922000000, 0);
}

// --- Inside map 922000000 --------------------------------------
if (plr.getMap() === 922000000) {
    if (plr.itemCount(4031092) < 10) {
        if (npc.sendSelection(
            "Have you taken care of everything? If you wish to leave, I'll let you out. Ready to go? \r\n\r\n#L0##bPlease let me out."
        ) === 0) {
            if (npc.sendYesNo("Hmm... All right. I can let you out, but you'll have to start from the beginning next time. Still wanna leave?")) {
                plr.warp(922000009, 0);
            }
        }
    }

    npc.sendNext("Oh ho, you really brought 10 Machine Parts items, and just in time. All right then! Since you have done so much for the toy factory, l'll give you a great present. Before I do that, however, make sure you have at least one empty slot in your Use tab.");

    if (plr.getFreeSlots(2 /* USE */) < 1) {
        npc.sendOk("Use item inventory is full.");
    }

    plr.giveEXP(140874);
    plr.giveItem(2040708, 1);
    plr.takeItem(4031092, 0, 10, 2);   // slot 0 to remove first instance
    plr.warp(220020000, 0);
}