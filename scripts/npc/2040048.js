// Pietri - LudiPQ Stage NPC
// Used for stage progression

var mapId = plr.mapID();
var pass = 4001022; // Pass of Dimension
var key = 4001023; // Key of Dimension

// Helper function to complete a stage
function completeStage(itemID, requiredAmount) {
    if (plr.itemCount(itemID) >= requiredAmount) {
        if (npc.sendYesNo("Good job! You have collected " + requiredAmount + " #t" + itemID + "#s. Would you like to move to the next stage?")) {
            plr.removeItemsByID(itemID, requiredAmount);
            var field = plr.inst;
            field.properties()["clear"] = true;
            field.showEffect("quest/party/clear");
            field.playSound("Party1/Clear");
            field.portalEffect("gate");
            npc.sendOk("The portal to the next stage is now open!");
        }
    } else {
        npc.sendOk("You need to collect #b" + requiredAmount + " #t" + itemID + "#s#k to proceed. You currently have #b" + plr.itemCount(itemID) + "#k.");
    }
}

// Stage 1: Collect 25 passes
if (mapId == 922010100) {
    completeStage(pass, 25);
}

// Stage 2: Collect 15 passes
else if (mapId == 922010200) {
    completeStage(pass, 15);
}

// Stage 3: Collect 32 passes (answer to the question)
else if (mapId == 922010300) {
    completeStage(pass, 32);
}

// Stage 4: Collect 6 passes
else if (mapId == 922010400) {
    completeStage(pass, 6);
}

// Stage 5: Collect 24 passes
else if (mapId == 922010500) {
    completeStage(pass, 24);
}

// Stage 6: Jump quest stage
else if (mapId == 922010600) {
    npc.sendOk("Navigate through the platforms to reach the portal to the next stage!");
}

// Stage 7: Collect 3 passes
else if (mapId == 922010700) {
    completeStage(pass, 3);
}

// Stage 8: Platform puzzle
else if (mapId == 922010800) {
    npc.sendOk("Stand on the correct platforms and I will check your answer when the party leader speaks to me!");
}

// Stage 9: Boss stage - Collect 1 key
else if (mapId == 922010900) {
    if (plr.itemCount(key) >= 1) {
        if (npc.sendYesNo("Incredible! You defeated Alishar and obtained the #t" + key + "#! Would you like to proceed to the bonus stage?")) {
            plr.removeItemsByID(key, 1);
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
