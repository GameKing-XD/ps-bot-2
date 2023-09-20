async function fetchGuilds() {
        return fetch("/api/guilds").
                then(res => res.json()).
                then(x => {
                        return Object.entries(x)
                })
}

function h(tag, attributes = {}, children = []) {
        let element;
        if (tag[0] === ".") {
                element = document.createElement('div');
                attributes.class = tag.substr(1);
        } else element = document.createElement(tag);

        if (typeof attributes === "string")
                attributes = {'class': attributes};

        for(let key in attributes) {
                if (!Object.hasOwnProperty.call(attributes, key)) continue;
                if (key[0] === "@") {
                        element.addEventListener(key.substr(1), attributes[key]);
                        continue;
                }
                element.setAttribute(key, attributes[key]);
        }

        children.forEach(child => {
                // Skip over empties
                if (typeof child === "undefined" || child === null) return;
                        if (typeof child === "string") {
                        let childText = document.createTextNode(child);
                        element.appendChild(childText);
                        return;
                }
                element.appendChild(child);
        });
        return element;
}

function e(q) {
	if (q.length === 0) return null;
	let q1 = q.substr(1);
	switch(q[0]) {
		case "#":
			return document.getElementById(q1);
		case '.':
			return document.getElementsByClassName(q1);
		default:
			return document.querySelectorAll(q1);
	}
}



async function updateGuilds()  {
        fetchGuilds().then(guilds => {
                e("#guild-selector").innerHTML = ""
                        e("#guild-selector").appendChild(h("option"))
                guilds.map(guild => h('option', {value: guild[0]}, [guild[1]])).forEach(x => 
                        e("#guild-selector").appendChild(x)
                )


                updateSounds()
        }) 
}

async function updateSounds() {
        const val = e("#guild-selector").value
        e("#sounds").innerHTML = ""
        if (val === "") return;
        fetchSounds(val).then(sounds => {
                sounds.map(s => button(val, s))
                        .forEach(x => {
                                e("#sounds").appendChild(x)
                        })
        })
}

async function fetchSounds(guild) {
        return fetch("/api/sounds/" + guild).then(res => res.json())
}

async function playSound(guild, s) {
        fetch("/api/play", {
                body: JSON.stringify({
                        guild: guild,
                        sound: s,
                }),
                headers: {
                        "Content-Type": "application/json",
                },
                method: "POST",
        })
}



document.addEventListener("readystatechange", () => {
        if (document.readyState !== "interactive") {
                return
        }
        updateGuilds()
        e("#guild-selector").addEventListener("change", () => {
                updateSounds()
        })
})
