// Pietri - LudiPQ Stage NPC
// Used for stage progression

var mapId = plr.mapID();
var pass = 4001022; // Pass of Dimension
var key = 4001023; // Key of Dimension

// Stage 1: Collect 25 passes
if (mapId == 922010100) {
    var requiredPasses = 25;
    if (plr.itemCount(pass) >= requiredPasses) {
        if (npc.sendYesNo("Good job! You have collected " + requiredPasses + " #t" + pass + "#s. Would you like to move to the next stage?")) {
            plr.removeItemsByID(pass, requiredPasses);
            var field = plr.inst;
            field.properties()["clear"] = true;
            field.showEffect("quest/party/clear");
            field.playSound("Party1/Clear");
            field.portalEffect("gate");
            npc.sendOk("The portal to the next stage is now open!");
        }
    } else {
        npc.sendOk("You need to collect #b" + requiredPasses + " #t" + pass + "#s#k to proceed. You currently have #b" + plr.itemCount(pass) + "#k.");
    }
}

// Stage 2: Collect 15 passes
else if (mapId == 922010200) {
    var requiredPasses = 15;
    if (plr.itemCount(pass) >= requiredPasses) {
        if (npc.sendYesNo("Good job! You have collected " + requiredPasses + " #t" + pass + "#s. Would you like to move to the next stage?")) {
            plr.removeItemsByID(pass, requiredPasses);
            var field = plr.inst;
            field.properties()["clear"] = true;
            field.showEffect("quest/party/clear");
            field.playSound("Party1/Clear");
            field.portalEffect("gate");
            npc.sendOk("The portal to the next stage is now open!");
        }
    } else {
        npc.sendOk("You need to collect #b" + requiredPasses + " #t" + pass + "#s#k to proceed. You currently have #b" + plr.itemCount(pass) + "#k.");
    }
}

// Stage 3: Collect 32 passes (answer to the question)
else if (mapId == 922010300) {
    var requiredPasses = 32;
    if (plr.itemCount(pass) >= requiredPasses) {
        if (npc.sendYesNo("Good job! You have collected " + requiredPasses + " #t" + pass + "#s. Would you like to move to the next stage?")) {
            plr.removeItemsByID(pass, requiredPasses);
            var field = plr.inst;
            field.properties()["clear"] = true;
            field.showEffect("quest/party/clear");
            field.playSound("Party1/Clear");
            field.portalEffect("gate");
            npc.sendOk("The portal to the next stage is now open!");
        }
    } else {
        npc.sendOk("You need to collect #b" + requiredPasses + " #t" + pass + "#s#k to proceed. You currently have #b" + plr.itemCount(pass) + "#k.");
    }
}

// Stage 4: Collect 6 passes
else if (mapId == 922010400) {
    var requiredPasses = 6;
    if (plr.itemCount(pass) >= requiredPasses) {
        if (npc.sendYesNo("Good job! You have collected " + requiredPasses + " #t" + pass + "#s. Would you like to move to the next stage?")) {
            plr.removeItemsByID(pass, requiredPasses);
            var field = plr.inst;
            field.properties()["clear"] = true;
            field.showEffect("quest/party/clear");
            field.playSound("Party1/Clear");
            field.portalEffect("gate");
            npc.sendOk("The portal to the next stage is now open!");
        }
    } else {
        npc.sendOk("You need to collect #b" + requiredPasses + " #t" + pass + "#s#k to proceed. You currently have #b" + plr.itemCount(pass) + "#k.");
    }
}

// Stage 5: Collect 24 passes
else if (mapId == 922010500) {
    var requiredPasses = 24;
    if (plr.itemCount(pass) >= requiredPasses) {
        if (npc.sendYesNo("Good job! You have collected " + requiredPasses + " #t" + pass + "#s. Would you like to move to the next stage?")) {
            plr.removeItemsByID(pass, requiredPasses);
            var field = plr.inst;
            field.properties()["clear"] = true;
            field.showEffect("quest/party/clear");
            field.playSound("Party1/Clear");
            field.portalEffect("gate");
            npc.sendOk("The portal to the next stage is now open!");
        }
    } else {
        npc.sendOk("You need to collect #b" + requiredPasses + " #t" + pass + "#s#k to proceed. You currently have #b" + plr.itemCount(pass) + "#k.");
    }
}

// Stage 6: Jump quest stage
else if (mapId == 922010600) {
    npc.sendOk("Navigate through the platforms to reach the portal to the next stage!");
}

// Stage 7: Collect 3 passes
else if (mapId == 922010700) {
    var requiredPasses = 3;
    if (plr.itemCount(pass) >= requiredPasses) {
        if (npc.sendYesNo("Good job! You have collected " + requiredPasses + " #t" + pass + "#s. Would you like to move to the next stage?")) {
            plr.removeItemsByID(pass, requiredPasses);
            var field = plr.inst;
            field.properties()["clear"] = true;
            field.showEffect("quest/party/clear");
            field.playSound("Party1/Clear");
            field.portalEffect("gate");
            npc.sendOk("The portal to the next stage is now open!");
        }
    } else {
        npc.sendOk("You need to collect #b" + requiredPasses + " #t" + pass + "#s#k to proceed. You currently have #b" + plr.itemCount(pass) + "#k.");
    }
}

// Stage 8: Platform puzzle
else if (mapId == 922010800) {
    npc.sendOk("Stand on the correct platforms and I will check your answer when the party leader speaks to me!");
}

// Stage 9: Boss stage - Collect 1 key
else if (mapId == 922010900) {
    var requiredKeys = 1;
    if (plr.itemCount(key) >= requiredKeys) {
        if (npc.sendYesNo("Incredible! You defeated Alishar and obtained the #t" + key + "#! Would you like to proceed to the bonus stage?")) {
            plr.removeItemsByID(key, requiredKeys);
            var field = plr.inst;
            field.properties()["clear"] = true;
            field.showEffect("quest/party/clear");
            field.playSound("Party1/Clear");
            field.portalEffect("gate");
            npc.sendOk("The portal to the bonus stage is now open!");
        }
    } else {
        npc.sendOk("Defeat Alishar and bring me the #b#t" + key + "##k to proceed!");
    }
}

// Bonus stage
else if (mapId == 922011000) {
    npc.sendOk("Congratulations on completing the Ludibrium Party Quest! Break as many boxes as you can for rewards!");
}
