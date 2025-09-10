// Dr. Lenu stateless NPC
const ITEM_REG = 5152010;
const ITEM_VIP = 5152013;

var faceId = plr.getFace();
var gender = plr.getGender();

// greeting + menu
var sel = npc.sendMenu(
    "Why hello there! I'm Dr. Lenu, in charge of the cosmetic lenses here at the Henesys Plastic Surgery Shop! With #b#t5152010##k or #b#t5152013##k, you can have the kind of look you've always wanted! All you have to do is find the cosmetic lens that most fits you, then let us take care of the rest. Now, what would you like to use?",
    "#bCosmetic Lenses at Henesys (Reg coupon)#l",
    "#bCosmetic Lenses at Henesys (VIP coupon)#l"
);

if (sel === 0) {        // Regular coupon
    if (npc.sendYesNo("If you use the regular coupon, you'll be awarded a random pair of cosmetic lenses. Are you going to use #b#t5152010##k and really make the change to your eyes?")) {
        if (!plr.itemCount(ITEM_REG)) {
            npc.sendNext("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..");
        }
        plr.takeItemById(ITEM_REG, 1);

        // build random face
        var baseColor = [100, 200, 300, 400, 500, 600, 700];
        var teye = faceId % 100;
        teye += gender < 1 ? 20000 : 21000;
        var newFace = baseColor[Math.floor(Math.random() * baseColor.length)] + teye;

        plr.setFace(newFace);
        npc.sendNext("Tada~! Check it out!! What do you think? I really think your eyes look sooo fantastic now~~! Please come again ~");
    }
} else if (sel === 1) { // VIP coupon
    // prepare preview list
    var baseColor = [100, 200, 300, 400, 500, 600, 700];
    var teye = faceId % 100;
    teye += gender < 1 ? 20000 : 21000;
    var faces = [];
    for (var i = 0; i < baseColor.length; i++) faces.push(teye + baseColor[i]);

    var chosen = npc.sendAvatar("With our specialized machine, you can see the results of your potential treatment in advance. What kind of lens would you like to wear? Choose the style of your liking...", faces);
    if (chosen >= 0) {
        if (!plr.itemCount(ITEM_VIP)) {
            npc.sendNext("I'm sorry, but I don't think you have our cosmetic lens coupon with you right now. Without the coupon, I'm afraid I can't do it for you..");
        }
        plr.takeItemById(ITEM_VIP, 1);
        plr.setFace(faces[chosen]);
        npc.sendNext("Tada~! Check it out!! What do you think? I really think your eyes look sooo fantastic now~~! Please come again ~");
    }
}