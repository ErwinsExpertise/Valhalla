// Pietri - LudiPQ Stage NPC (party2_play)
// Shows help for non-leaders, handles stage progression for party leaders

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

// Show help text for non-party leaders
if (!plr.isPartyLeader()) {
    if (mapId == 922010100) {
        npc.sendOk("Here is information about the 1st stage. You'll see monsters at different points on the map. These monsters have an item called #b#t4001022##k, which opens the door to another dimension. With it, you can take a step closer to the top of the Tower of Eos, where the door to the other dimension will open and, finally, you will find the culprit for everything.\r\nDefeat the monsters, collect #b25 #t4001022#s#k and give it to your group leader, who will in turn give it to me. This will take you to the next stage. Good luck!");
    } else if (mapId == 922010200) {
        npc.sendOk("Here is information about the 2nd stage. You'll see crates all over the map. Break a box and you will be sent to another map or rewarded with a #t4001022#. Search each box, collect #b15 #t4001022#s#k and bring them all to me. Gather the collected #t4001022#s, give them all to your group leader, who, in turn, will give them to me.\r\nBy the way, even if you are sent to another place, you can find another box there. So, don't just leave the strange place you went to. If you just leave there, you can't go back and you'll need to start the quest from the beginning. Good luck!");
    } else if (mapId == 922010300) {
        npc.sendOk("Here is information about the 3rd stage. Here you will see a bunch of monsters and boxes. If you defeat the monsters, they will drop #b#t4001022##k, just like monsters from the other dimension. If you break the box, a monster will appear and it will also give #b#t4001022##k.\r\nThe number of #b#t4001022#s#k you need to collect will be determined by the answer to the question I will ask the leader of your group. The answer to the question will determine the number of #b#t4001022#s#k you will need to collect. Once I ask the group leader the question, they can discuss the answer with the members. Good luck!");
    } else if (mapId == 922010400) {
        npc.sendOk("Here is information about the 4th stage. Here you will find a black space created by the dimensional rift. Inside, you'll find a monster called #b#o9300008##k hiding in the darkness. For this reason, you will barely be able to see it if you don't have your eyes wide open. Defeat the monsters and collect #b6 #t4001022#s#k.\r\nLike I said, #b#o9300008##k cannot be seen unless your eyes are wide open. It's a different kind of monster that stealthily merges into the darkness. Good luck!");
    } else if (mapId == 922010500) {
        npc.sendOk("Here is information about the 5th stage. Here you will find many spaces and, inside them, you will find some monsters. Your duty is to collect with the group #b24 #t4001022#s#k. This is the description: There will be cases where you need to be of a certain profession, or you cannot collect #b#t4001022##k. Therefore, be careful. Here's a clue. There is a monster called #b#o9300013##k that is unbeatable. Only a rogue can get to the other side of the monster. There is also a route that only witches can take. Finding out is up to you. Good luck!");
    } else if (mapId == 922010600) {
        npc.sendOk("Here is the information about the 6th stage. Here, you'll see boxes with numbers written on them, and if you stand on top of the correct box by pressing the UP ARROW, you'll be transported to the next box. I'll give the party leader a clue on how to get past this stage #bonly twice#k and it's the leader's duty to remember the clue and take the right step, one at a time.\r\nOnce you reach the top, you'll find the portal to the next stage. When everyone in your party has passed through the portal, the stage is complete. Everything will depend on remembering the correct boxes. Good luck!");
    } else if (mapId == 922010700) {
        npc.sendOk("Here is information about the 7th stage. Here you will find a ridiculously powerful monster called #b#o9300010##k. Defeat the monster and find the #b#t4001022##k needed to proceed to the next stage. Please collect #b3 #t4001022#s#k.\r\nTo finish off the monster, defeat it from afar. The only way to attack would be from a long distance, but... oh yes, be careful, #o9300010# is very dangerous. You will definitely get hurt if you are not careful. Good luck!");
    } else if (mapId == 922010800) {
        npc.sendOk("Here is information about the 8th stage. Here you will find many platforms to climb. #b5#k of them will be connected to the #bportal that leads to the next stage#k. To pass, place #b5 of your party members on the correct platform#k.\r\nA word of warning: You will need to stand firmly in the center of the platform for your answer to count as correct. Also remember that only 5 members can stay on the platform. When this happens, the group leader must #bclick me twice to know if the answer is correct or not#k. Good luck!");
    } else if (mapId == 922010900) {
        npc.sendOk("Here is the information about the 9th stage. Now is your chance to finally get your hands on the real culprit. Go right and you'll see a monster. Defeat it to find a monstrous #b#o9300012##k appearing out of nowhere. He will be very agitated by your group's presence, be careful.\r\nYour task is to defeat him, collect the #b#t4001023##k he has and bring it to me. If you manage to take the key away from the monster, there is no way the dimensional door can be opened again. I have faith in you. Good luck!");
    } else if (mapId == 922011000) {
        npc.sendOk("Welcome to the bonus stage. I can't believe you actually defeated #b#o9300012##k! Incredible! But we don't have much time, so I'll get right to the point. There are many boxes here. Your task is to break the boxes within the time limit and get the items inside. If you're lucky, you might even snag a great item here and there. If this doesn't excite you, I don't know what will. Good luck!");
    }
}
// Party leader - handle stage progression
else {
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
        npc.sendOk("Stand on the correct platforms and I will check your answer when you speak to me again!");
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
    // Bonus stage - same help text for leader
    else if (mapId == 922011000) {
        npc.sendOk("Welcome to the bonus stage. I can't believe you actually defeated #b#o9300012##k! Incredible! But we don't have much time, so I'll get right to the point. There are many boxes here. Your task is to break the boxes within the time limit and get the items inside. If you're lucky, you might even snag a great item here and there. If this doesn't excite you, I don't know what will. Good luck!");
    }
}
