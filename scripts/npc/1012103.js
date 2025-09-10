var couponCut  = 5150053; // Haircut coupon
var couponDye = 5151036; // Hair dye coupon

// Main menu
npc.sendSelection(
    "I'm the head of this hair salon Natalie. If you have #b#t5150053##k or #b#t5151036##k, allow me to take care of your hairdo. Please choose the one you want.\r\n"
    + "#L0#Haircut (VIP coupon)#l\r\n"
    + "#L1#Dye your hair (VIP coupon)#l"
);
var sel = npc.selection();

if (sel === 0) {
    // Haircut branch
    var hair;
    if (plr.gender() < 1) { // male
        hair = [33040, 30060, 30210, 30140, 30200, 33170, 33100];
    } else { // female
        hair = [31150, 34090, 31300, 31700, 31350, 31740, 34110];
    }

    // Add last digit from current hair to each style
    var lastDigit = plr.hair() % 10;
    for (var i = 0; i < hair.length; i++) {
        hair[i] += lastDigit;
    }

    var choice = npc.askAvatar(
        "I can totally change up your hairstyle and make it look so good. Why don't you change it up a bit? with #b#t5150053##k I'll change it for you. Choose the one to your liking~",
        ...hair
    );

    if (plr.itemCount(couponCut) > 0) {
        plr.giveItem(couponCut, -1);
        plr.setHair(hair[choice]);
        npc.sendBackNext(
            "Check it out!! What do you think? Even I think this one is a work of art! AHAHAHA. Please let me know when you want another haircut, because I'll make you look good each time!",
            false, true
        );
    } else {
        npc.sendBackNext(
            "Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't give you a haircut without it. I'm sorry...",
            false, true
        );
    }

} else if (sel === 1) {
    // Hair dye branch
    var base = Math.floor(plr.hair() / 10) * 10;
    var colors = [base + 0, base + 1, base + 2, base + 4, base + 6];

    var choice = npc.askAvatar(
        "I can totally change your haircolor and make it look so good. Why don't you change it up a bit? With #b#t5151036##k I'll change it for you. Choose the one to your liking.",
        ...colors
    );

    if (plr.itemCount(couponDye) > 0) {
        plr.giveItem(couponDye, -1);
        plr.setHair(colors[choice]);
        npc.sendBackNext(
            "Check it out!! What do you think? Even I think this one is a work of art! AHAHAHA. Please let me know when you want to dye your hair again, because I'll make you look good each time!",
            false, true
        );
    } else {
        npc.sendBackNext(
            "Hmmm...it looks like you don't have our designated coupon...I'm afraid I can't dye your hair without it. I'm sorry...",
            false, true
        );
    }
}