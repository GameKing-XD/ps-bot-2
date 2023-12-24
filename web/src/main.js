import {h,e} from './js/fw.js'
import {Fzf} from 'fzf';

import 'bootswatch/dist/darkly/bootstrap.min.css';

const GUILDS_URL="/api/guilds";
const PLAYSOUND_URL="/api/play";
const SOUNDS_URL="/api/sounds/";
const GUILD_DROPDOWN_SELECTOR="#guild-selector";
const SOUNDS_PARENT_SELECTOR="#sounds";

async function fetchGuilds() {
        return fetch("/api/guilds").
                then(res => res.json()).
                then(x => Object.entries(x));
}



function updateGuilds()  {
        fetchGuilds(GUILDS_URL).then(guilds => {
                e(GUILD_DROPDOWN_SELECTOR).innerHTML = ""
                        e(GUILD_DROPDOWN_SELECTOR).appendChild(h("option"))
                guilds.map(guild => h('option', {value: guild[0]}, [guild[1]])).forEach(x => 
                        e(GUILD_DROPDOWN_SELECTOR).appendChild(x)
                )

                updateSounds();
        }) 
}

function updateSounds() {
        const val = e(GUILD_DROPDOWN_SELECTOR).value
        e(SOUNDS_PARENT_SELECTOR).innerHTML = ""
        if (val === "") return;
        fetchSounds(val).then(sounds => {
                const searchEngine = new Fzf(sounds);
                const soundParent = e(SOUNDS_PARENT_SELECTOR);
                const searchbar = h("input", {
                        "class": 'form-control mb-4', 
                        'placeholder': 'Search...', 
                        'id': 'sound-searchbar', 

                        'autocorrect': "off",
                        'autocapitalize': "none",
                        "@input": ev => {
                                const results = searchEngine.find(ev.target.value);

                                const soundEntries = h('div', {'id': 'sound-entries', 'class': 'row'}, 
                                        results.map(({item}) => button(val, item)),
                                );
                                e("#sound-entries")?.remove();
                                soundParent.appendChild(soundEntries);
                        }
                });
                e(SOUNDS_PARENT_SELECTOR).appendChild(searchbar);

                const soundEntries = h('div', {'id': 'sound-entries', 'class': 'row'}, 
                        sounds.map((item) => button(val, item)),
                );
                e("#sound-entries")?.remove();
                soundParent.appendChild(soundEntries);



        })
}


function button(val, s) {
        return h(
                'div', 
                {'class': 'col-6 mb-1 col-md-3'}, 
                [
                        h(
                                'button', 
                                {'class': 'btn btn-primary d-block w-100', "@click": () => playSound(val, s)}, 
                                [s]
                        )
                ]
        )
}

async function fetchSounds(guild) {
        return fetch(SOUNDS_URL + guild).then(res => res.json())
}

function playSound(guild, s) {
        fetch(PLAYSOUND_URL, {
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



updateGuilds()
e(GUILD_DROPDOWN_SELECTOR).addEventListener("change", () => {
        updateSounds()
});

