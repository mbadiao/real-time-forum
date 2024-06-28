import { message } from "./templates.js";

export var chunk;



export function HandleMessage() {
    document.addEventListener('click', function (event) {
        if (event.target.classList.contains('userfisrtandlastnameLeft')) {
            const id = event.target.id;
            console.log('Selected user ID:', id);
            const oldDiv = document.body.getElementsByClassName('message')[0];

            let name = event.target.textContent

            if (oldDiv) {
                oldDiv.remove();
            }
            const div = document.createElement('div');
            div.innerHTML = message;
            const postNmessage = document.getElementsByClassName('postsAndMessage')[0]
            postNmessage.appendChild(div)
            const namePlace = document.getElementsByClassName('userfisrtandlastnamemessage')[0]
            namePlace.textContent = name
            messagesWebsocket(id)
        }
    });
}

function messagesWebsocket(id) {
    Typing(id)
    let chatBody = document.getElementsByClassName('chat-container')[0];
    chatBody.scrollTop = chatBody.scrollHeight;

    let wsMessage = new WebSocket(`ws://localhost:8080/messages?receiver=${id}`);

    let form = document.getElementsByClassName('messageChat')[0];
    let message = document.getElementsByClassName('messageInput')[0];


    wsMessage.onopen = function () {

    }

    form.addEventListener('submit', function (e) {
        e.preventDefault();
        // console.log('test');
        if (message.value === '') {
            alert('veuillez remplir le formulaire');
            return;
        }
        wsMessage.send(message.value);
        // handleNotifs(id)
        // console.log(message.value);
        message.value = '';
    });
    wsMessage.onmessage = function (event) {
        throttle(MessageDistribution(chatBody, event.data, id), 200)
    }
}
function broadCast(user, messages, scroll) {
    // console.log(messages);
    // messages.reverse();
    let chatBody = document.getElementsByClassName('chat-container')[0]

    messages.forEach(message => {
        const newMessage = document.createElement('div');
        newMessage.className = message.receiver === user ? 'chat-bubble sender' : 'chat-bubble receiver';
        newMessage.textContent = message.content
        const format = document.createElement('div')
        format.className = message.receiver === user ? 'format sent' : 'format received'
        format.style.color = 'white'
        if (message.receiver === user) {
            format.textContent = message.formatReceiver
        } else {
            format.textContent = message.formatSender
        }
        newMessage.appendChild(format)
        if (scroll < 1) {
            chatBody.appendChild(newMessage);
        } else {
            chatBody.prepend(newMessage)
        }
    });
}

function MessageDistribution(chatBody, data, id) {
    let scroll = 1
    clearMessages();

    let tab = JSON.parse(data);
    chunk = chunkArray(tab, 10);
    // console.log(chunk);
    if (scroll === 1) {
        if (chunk[scroll - 1]) {

            broadCast(id, chunk[scroll - 1], scroll);
        }
        scroll++
    }

    chatBody.scrollTop = chatBody.scrollHeight;
    chatBody.addEventListener('scroll', function () {
        if (chatBody.scrollTop === 0) {
            // console.log((scroll <= chunk.length));
            if (scroll <= chunk.length) {
                setTimeout(() => {
                    if (chunk[scroll - 1]) {
                        broadCast(id, chunk[scroll - 1], scroll);
                    }
                    scroll++
                }, 1000);
            }
        }
    })
}

function chunkArray(array, chunkSize) {
    array.reverse()
    const chunkedArray = [];
    for (let i = 0; i < array.length; i += chunkSize) {
        const chunk = array.slice(i, i + chunkSize);
        chunkedArray.push(chunk);
    }
    return chunkedArray;
}

function throttle(func, delay) {
    let timer = null;
    return (...args) => {
        if (timer === null) {
            func(...args);
            timer = setTimeout(() => {
                timer = null;
            }, delay);
        }
    };
}

function clearMessages() {
    const messageDivs = document.getElementsByClassName('chat-bubble');
    if (messageDivs) {
        while (messageDivs.length > 0) {
            messageDivs[0].remove();
        }
    }
}

let bubble = `
<div class="typing-bubble">
  <span class="dot"></span>
  <span class="dot"></span>
  <span class="dot"></span>
</div>`

function Typing(id) {
    let wsTyping = new WebSocket("ws://localhost:8080/typing-progress")
    let form = document.getElementsByClassName('messageChat')[0];
    form.addEventListener('keydown', function () {
        wsTyping.send(JSON.stringify({ "id": id, "typing": 'true' }));
    });

    let user = document.getElementsByClassName('userfisrtandlastnamemessage')[0]
    let isTyping = document.getElementsByClassName('istyping')[0]
    let typingTimeout;

    form.addEventListener('keyup', function () {
        clearTimeout(typingTimeout);
        typingTimeout = setTimeout(function () {
            wsTyping.send(JSON.stringify({ "id": id, "typing": 'false' }));
        }, 1000);
    });

    wsTyping.onmessage = function (event) {
        if (event.data === user.textContent) {
            let chatBbubble = document.getElementsByClassName('typing-bubble')[0]
            if (!chatBbubble) {
                let chatContainer = document.getElementsByClassName('chat-container')[0]
                const newMessage = document.createElement('div');
                newMessage.className = 'typing-chat';
                newMessage.innerHTML = bubble
                chatContainer.appendChild(newMessage)
                chatContainer.scrollTop = chatContainer.scrollHeight;
            }

            const data = event.data;
            const splitData = data.split(' ');
            const firstname = splitData[0];
            isTyping.textContent = `${firstname} is typing...`;
        } else {
            let typingChat = document.getElementsByClassName('typing-chat')[0]
            if (typingChat) {
                typingChat.remove()
            }
            isTyping.textContent = '';
        }
    }
}