
import {
    getCookie,
} from "./register.js";

import { renderOnlineOrOfflineTemplate } from "./templates.js";


export function userStatus() {
    const socket = new WebSocket('ws://localhost:8080/userstatus');
    const usersList = document.querySelector('.userLeft');
    
    socket.onopen = function (event) {
        const cookie = getCookie("ForumCookie")
        socket.send(JSON.stringify(cookie));
        console.log("cookie envoyer");
    };

    socket.onmessage = function (event) {
        const users = JSON.parse(event.data);
        console.log("reception de donnees", users);
        usersList.innerHTML = '';
        for (const username in users) {
            const user = users[username];
            // const li = document.createElement('li');
            // li.textContent = `${user.username} (${user.online ? 'online' : 'offline'})`;
            // usersList.appendChild(li);
            usersList.innerHTML+= renderOnlineOrOfflineTemplate(user)
        }
    };

    socket.onclose = function (event) {
        console.log('WebSocket is closed now.');
    };

    socket.onerror = function (error) {
        console.error('WebSocket Error: ', error);
    };
    const logout = document.querySelector(".logout");
    logout.onclick = function() {
        // socket.send(JSON.stringify({type: "logout"}));
        socket.close();
    };
}
