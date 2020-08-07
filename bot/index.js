const Discord = require('discord.js');
const client = new Discord.Client();
const mysql = require('mysql');
const config = require('./config.json');

/* You should change the lines below to working things, or Idigo won't work */
const websiteUrl = "https://example.com/idigo"; // Change this to your website URL, it MUST NOT end with a slash (/)
const serverId   = serverId;        // Change this to your discord server ID, make it a string, why not

const http = require('http');

let con = mysql.createConnection({
    host:       config.db.hostname,
    user:       config.db.username,
    password:   config.db.password,
    database:   config.db.database
});

client.on('ready', s => {
    console.log(`Connected as ${client.user.tag}`);
    setInterval(() => {
        con.query(`SELECT 1`, (res, err) => {})
    }, 30*1000);
});

function timeDifference(current, previous) {

    let msPerMinute = 60 * 1000;
    let msPerHour = msPerMinute * 60;
    let msPerDay = msPerHour * 24;
    let msPerMonth = msPerDay * 30;
    let msPerYear = msPerDay * 365;

    let elapsed = current - previous;

    if (elapsed < msPerMinute) {
         return Math.round(elapsed/1000) + ' seconds ago';   
    }

    else if (elapsed < msPerHour) {
         return Math.round(elapsed/msPerMinute) + ' minutes ago';   
    }

    else if (elapsed < msPerDay ) {
         return Math.round(elapsed/msPerHour ) + ' hours ago';   
    }

    else if (elapsed < msPerMonth) {
        return 'approximately ' + Math.round(elapsed/msPerDay) + ' days ago';   
    }

    else if (elapsed < msPerYear) {
        return 'approximately ' + Math.round(elapsed/msPerMonth) + ' months ago';   
    }

    else {
        return 'approximately ' + Math.round(elapsed/msPerYear ) + ' years ago';   
    }
}

function rand(min, max) { // min and max included 
    return Math.floor(Math.random() * (max - min + 1) + min);
}
const helpers = {
    embed: {
        footer: {
            text: "• Idigo by sn#6740",
            image: websiteUrl + "/idigo.gif"
        },
        color: {
            danger: "#e74c3c",
            success: "#2ecc71",
            warning: "#f1c40f",
            primary: "#3498db"
        }
    },
    idigo: {
        image: websiteUrl + "/idigo.gif"
    }
}

function makeDefaultEmbed(message) {
    const embed = new Discord.MessageEmbed()
        .setTitle(`Idigo | Pin`)
        .setThumbnail(helpers.idigo.image)
        .setAuthor(`${message.author.tag}`, `https://cdn.discordapp.com/avatars/${message.author.id}/${message.author.avatar}?size=128`)
        .setTimestamp()
        .setFooter(helpers.embed.footer.text, helpers.embed.footer.image)
    return embed;
}

client.on('guildMemberAdd', member => {
    const role = client.guilds.cache.get(serverId).roles.cache.find(x => x.name == "user");
    member.roles.add(role);
});

client.on(`message`, m => {
    if (!m.content.startsWith(config.prefix) || m.author.bot) return;

    const args = m.content.slice(config.prefix.length).trim().split(' ');
    const command = args.shift().toLowerCase();

    switch(command) {
        case "pin":
            // let cooldownT = 30 * 1000, cooldownG = cooldowns.pin.get(m.author.id);
            // if(cooldownG) return m.channel.send(`Please wait ${humd(cooldownG - Date.now(), { round: true })} before running ${command} again`);

            con.query(`SELECT * FROM users WHERE discord = '${m.author.id}'`, (err, res) => {
                if(m.channel.parent.name.toLowerCase() == "pins") res[0] = { license: 1 };
                if(!res[0] || res[0].license != 1) return (function() {
                    const embed = new Discord.MessageEmbed()
                        .setTitle(`Idigo`)
                        .setThumbnail(helpers.idigo.image)
                        .setAuthor(`${m.author.tag}`, `https://cdn.discordapp.com/avatars/${m.author.id}/${m.author.avatar}?size=128`)
                        .setTimestamp()
                        .setFooter(helpers.embed.footer.text, helpers.embed.footer.image)
                        .setColor(helpers.embed.color.danger)
                        .setDescription(`You don't have a license!\nCheck out <#738802873990643732> to get one!`)
                    m.channel.send(embed);
                    m.delete().catch(() => {});
                })();
                con.query(`SELECT * FROM pins WHERE author = '${m.author.id}'`, (err, res) => {
                    if(res != undefined && res.length != 1) return (() => {
                        const embed = getDefaultEmbed(m)
                            .setColor(helpers.embed.color.danger)
                            .setDescription(`You already have an unsed PIN!\nCheck your DMs and use that PIN instead of generating a new one!`)
                        m.channel.send(embed);
                        m.delete().catch(() => {});
                    })();
                    let pin = rand(100000, 999999);
                    m.author.send(`Your PIN is: ${pin}`).then(() => {
                    const embed = makeDefaultEmbed(m)
                        .setColor(helpers.embed.color.success)
                        .setDescription(`Sent you a new PIN on your DMs\nThanks for using Idigo!`);
                    m.channel.send(embed)
                    .catch(err => {
                        const embed = makeDefaultEmbed(m)
                            .setColor(helpers.embed.color.warning)
                            .setDescription(`Couldn't send you a private message!\n\nDo you want me to send the PIN in this channel?\n:warning: Someone could steal your PIN!\n\nPlease, enable your messages or reply with **"yes"**`);
                        m.channel.send(embed);
                        m.delete().catch(() => {});
                        m.channel.awaitMessages((mf => mf.author.id == m.author.id && mf.content == "yes"), {max: 1, time: 10*1000})
                        .then(() => {
                            m.channel.send(`${m.author}\nYour PIN is: ${pin}`);
                            con.query(`INSERT INTO pin (pin, author, channel) VALUES ('${pin}', '${m.author.id}', '${m.channel.id}')`, (err, res) => {});
                        })
                        .catch(() => {});
                    })
                    .then(() => {
                        con.query(`INSERT INTO pin (pin, author, channel) VALUES ('${pin}', '${m.author.id}', '${m.channel.id}')`, (err, res) => {});
                        m.delete().catch(() => {});
                    });
                });
            });
            });
            break;
        }
});

client.login(config.token);


const reqHandler = function (req, res) {
    res.writeHead(200);
    res.write(`Received request`)
    res.end();
    console.log(req.headers);
    const embed = new Discord.MessageEmbed()
        .setTitle(`Idigo | Result`)
        .setThumbnail(helpers.idigo.image)
        .setTimestamp()
        .setFooter(helpers.embed.footer.text, helpers.embed.footer.image)
        .setDescription(`${(req.headers.cheats == "None") ? "User is clean" : "User is not clean\n" + req.headers.cheats.split("|||").join("\n")}`)
        .addFields(
            { name: "OS", value: req.headers.build, inline: true },
            { name: "Recycle", value: `${timeDifference((() => { let date = new Date(); return (date.getTime()/1000).toFixed(0); })(), req.headers.recycle)}`, inline: true }
        )
    if(req.headers.cheats == "None") embed.setColor(helpers.embed.color.success); else embed.setColor(helpers.embed.color.danger);
    client.guilds.cache.get(serverId).channels.cache.get(req.headers.channel).send(embed)
  }
  
  const server = http.createServer(reqHandler);
  server.listen(51335);
