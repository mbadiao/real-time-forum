import { errorInHtml, displayPostInHtml } from "./templates.js";
import { getUserInfo } from "./app.js";

const ws = new WebSocket("ws://localhost:8080/like");

ws.onopen = () => {
    console.log("Connected to the WebSocket server like");
};

ws.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);
        if (data) {
            if (data.Errors && data.Errors.Status) {
                const { username, firstname } = getUserInfo();
                errorInHtml(data.Errors.Code, data.Errors.Message, username, firstname);
                return;
            }
            console.log('here');
            displayPostInHtml(data)
        }
    } catch (error) {
        console.error("Error parsing server response like:", error);
        const { username, firstname } = getUserInfo();
        errorInHtml("500", "Error parsing server response", username, firstname);
    }
};

ws.onclose = (event) => {
    if (event.wasClean) {
        console.log(
            `Connection closed cleanly, code=${event.code}, reason=${event.reason}`
        );
    } else {
        console.error("Connection died");
    }
};


export function Like() {
    document.querySelectorAll(".sendLikeWs").forEach(sendLikeButton => {
        sendLikeButton.addEventListener("click", (event) => {
            event.preventDefault();
            const postId = sendLikeButton.getAttribute('like-postid');
            const userId = sendLikeButton.getAttribute('like-userid');
            const like = {
                postid: Number(postId),
                userid: Number(userId)
            };
            console.log("EVENEMENT LIKE");
            ws.send(JSON.stringify(like));
        });
    });
}