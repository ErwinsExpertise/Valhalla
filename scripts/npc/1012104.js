// Brittany the assistant hair NPC – REG coupons (random result)
const couponHaircut = 5150052
const couponDye = 5151035
const maleHaircuts = [30310,30330,30060,30150,30410,30210,30140,30120,30200,30560,30510,30610,30470]
const femaleHaircuts = [31150,31310,31300,31160,31100,31410,31030,31080,31070,31610,31350,31510,31740]

// Intro
npc.sendBackNext(
    "I'm Brittany the assistant. If you have #b#t5150052##k or #b#t5151035##k by any chance, then how about letting me change your hairdo?",
    false, true
)

// Menu
npc.sendSelection(
    "What would you like to do today?\r\n" +
    "#L0##bHaircut (REG coupon)#k#l\r\n" +
    "#L1##bDye your hair (REG coupon)#k#l"
)
var choice = npc.selection()

if (choice === 0) {
    // Haircut (random style within gender, keep color)
    var pool = (plr.gender() < 1) ? maleHaircuts : femaleHaircuts
    var newStyle = pool[Math.floor(Math.random() * pool.length)]
    newStyle += plr.hair() % 10

    if (!npc.sendYesNo("If you use the REG coupon, your hair will change RANDOMLY. Use #b#t5150052##k and change your hairstyle?")) {
        npc.sendOk("See you another time!")
    } else if (plr.itemCount(couponHaircut) >= 1) {
        plr.removeItemsByID(couponHaircut, 1)
        plr.setHair(newStyle)
        npc.sendOk("Hey, here's the mirror. What do you think of your new haircut? Come back later when you want to change it up again!")
    } else {
        npc.sendOk("Hmmm... are you sure you have our designated coupon? Sorry, no haircut without it.")
    }

} else if (choice === 1) {
    // Dye (random color, keep base style)
    var base = Math.floor(plr.hair() / 10) * 10
    var colors = [base + 0, base + 1, base + 2, base + 3, base + 4, base + 5]
    var newStyle = colors[Math.floor(Math.random() * colors.length)]

    if (!npc.sendYesNo("If you use the REG coupon, your hair color will change RANDOMLY. Use #b#t5151035##k and change it up?")) {
        npc.sendOk("See you another time!")
    } else if (plr.itemCount(couponDye) >= 1) {
        plr.removeItemsByID(couponDye, 1)
        plr.setHair(newStyle)
        npc.sendOk("Hey, here's the mirror. What do you think of your new hair color? Come back later when you want to change it up again!")
    } else {
        npc.sendOk("Hmmm... are you sure you have our designated coupon? Sorry, we can’t dye your hair without it.")
    }

} else {
    npc.sendOk("See you another time!")
}