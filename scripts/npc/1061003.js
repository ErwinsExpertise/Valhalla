// Start conversation – Mr. Wetbottom in Sleepywood VIP Sauna
const QUEST_ID = 1000600;
const BOOK_ID  = 4031016;

// Helper strings
function quali() {
    if (plr.getLevel() < 30) return "notReady";
    const d = plr.questData(QUEST_ID);
    if (d === "") return "available";
    if (d === "e") return "done";
    if (d === "4" && plr.itemCount(BOOK_ID) >= 1) return "reward";
    return "inProgress";
}

const qStat = quali();

// Welcome line
if (qStat === "notReady") {
    npc.sendOk("Welcome to the VIP sauna of the #m105040300# Hotel. Actually I need some help here ...");
} else if (qStat === "available") {
    npc.sendNext("Good job getting here. Actually I have a favor to ask you. If you accept it, I'll give you a piece of clothing that you'll need in return. I think you are more than capable of doing it. Even if you don't care, just please listen to my story first.");
    npc.sendBackNext("I have a son that can do me no wrong. But one day he took a book of mine that is very dear to me and left. That book ... hmmm ... I can't give you the full detail on it, but it is a very very important book to me...", true, true);
    if (npc.sendYesNo("If you get me that book back safely, I'll give you a comfortable article of clothing, perfect for saunas like this. What do you think? Will you find my son and take back that book?")) {
        plr.startQuest(QUEST_ID);
        plr.setQuestData(QUEST_ID, "1");
        npc.sendOk("Ohhh ... thank you so much. It won't be easy locating my son in this humongous island, the Victoria Island. I'm guessing that he may be in a passage made of trees near the forest at #m101000000#, because that's his favorite spot ... best of luck to you!");
    } else {
        npc.sendOk("I see ... I guess you're busy with things here and there ... but I'll definitely reward you handsomely for your work so if you ever change your mind, please let me know.");
    }
} else if (qStat === "done") {
    npc.sendOk("I'm so glad I got this book back safely. It's my number one treasure, you know. Am I not worried about #p1061004#? The fairies are taking care of him alright, so I'm not worried one bit.");
} else if (qStat === "reward") {
    plr.setQuestData(QUEST_ID, "e");
    plr.removeItemsByID(BOOK_ID, 1);
    let robe = (plr.job() % 1000 < 500) ? 1050018 : 1051017;
    plr.giveItem(robe, 1);
    plr.giveEXP(500);
    plr.giveMesos(10000);
    npc.sendOk("This is the book! The book I was looking for! WHEW! Thank you so much! Here, the piece of clothing, like I promised.\r\n\r\n#v"+robe+"#");
} else {   // inProgress
    npc.sendOk("You haven't found my book yet. Please come back when you have. My son is located in a passage made of trees near the forest at Ellinia, because that's his favorite spot ... best of luck to you!");
}