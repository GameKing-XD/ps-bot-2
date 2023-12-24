export function h(tag, attributes = {}, children = []) {
        let element;
        if (tag[0] === '.') {
                element = document.createElement('div');
                attributes.class = tag.slice(1);
        } else {
                element = document.createElement(tag);
        }

        if (typeof attributes === 'string') {
                attributes = { class: attributes };
        }

        for (let key in attributes) {
                if (!Object.hasOwnProperty.call(attributes, key)) continue;
                if (key[0] === '@') {
                        element.addEventListener(key.slice(1), attributes[key]);
                        continue;
                }
                element.setAttribute(key, attributes[key]);
        }

        children.forEach((child) => {
                // Skip over empties
                if (typeof child === 'undefined' || child === null) {
                        return;
                }

                if (typeof child === 'string') {
                        let childText = document.createTextNode(child);
                        element.appendChild(childText);
                        return;
                }
                element.appendChild(child);
        });
        return element;
}

export function e(q) {
        if (q.length === 0) return null;
        let q1 = q.slice(1);
        switch (q[0]) {
                case '#':
                        return document.getElementById(q1);
                case '.':
                        return document.getElementsByClassName(q1);
                default:
                        return document.querySelectorAll(q1);
        }
}
