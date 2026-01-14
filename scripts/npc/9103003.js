// Reward NPC for Ludibrium PQ Completion
// This NPC appears at the bonus stage exit

var pass = 4001022;
var key = 4001023;

// Clean up any remaining quest items
if (plr.itemCount(pass) > 0) {
    plr.removeItemsByIDSilent(pass, plr.itemCount(pass));
}
if (plr.itemCount(key) > 0) {
    plr.removeItemsByIDSilent(key, plr.itemCount(key));
}

// Give rewards
var expReward = 8500;
plr.giveEXP(expReward);

// Random reward system - simplified version
var rand = Math.floor(Math.random() * 100);
var rewardID = 0;
var rewardAmount = 1;

if (rand < 5) {
    rewardID = 2000004; // Elixir
    rewardAmount = 10;
} else if (rand < 10) {
    rewardID = 2000002; // White Potion
    rewardAmount = 100;
} else if (rand < 15) {
    rewardID = 2000003; // Blue Potion
    rewardAmount = 100;
} else if (rand < 20) {
    rewardID = 4010000; // Adamantium Ore
    rewardAmount = 15;
} else if (rand < 25) {
    rewardID = 4010001; // Silver Ore
    rewardAmount = 15;
} else if (rand < 30) {
    rewardID = 4010002; // Orihalcon Ore
    rewardAmount = 15;
} else if (rand < 35) {
    rewardID = 4020000; // Garnet
    rewardAmount = 15;
} else if (rand < 40) {
    rewardID = 4020001; // Amethyst
    rewardAmount = 15;
} else if (rand < 45) {
    rewardID = 4020002; // Aquamarine
    rewardAmount = 15;
} else {
    rewardID = 4003000; // Screws
    rewardAmount = 50;
}

if (plr.giveItem(rewardID, rewardAmount)) {
    npc.sendOk("Incredible! You've completed all the stages! Here's your reward: " + rewardAmount + " #t" + rewardID + "# and " + expReward + " EXP!");
} else {
    npc.sendOk("Your inventory is full! Please make space and try again.");
}

// Warp player out
plr.warp(221024500); // Exit map
