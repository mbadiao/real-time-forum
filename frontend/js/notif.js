import { notifHtml } from "./templates.js";

export function wsNotifs() {
    let wsNotifications = new WebSocket('ws://localhost:8080/notifications')
    const logout = document.querySelector(".logout");
    // console.log(logout);
    if (logout) {
        logout.addEventListener("click", function () {
            // console.log('logout a ete bel et bien apuyee');
            wsNotifications.close();
        });
    }

    wsNotifications.onopen = function (event) {

    }

    wsNotifications.onmessage = function (event) {
        let maincontainer = document.getElementsByClassName('mainContainer home')[0]

        let notifDiv = document.createElement('div')
        notifDiv.className = "notifContainContainer"
        notifDiv.innerHTML = notifHtml
        notifDiv.style.animation = "slideDown 2s"; // Apply the animation

        let notifContainer = document.createElement('div')

        notifContainer.innerHTML = `
            <div class="notif-message">
                <div class="notif-sender">
                    ${event.data}
                </div>
                <div class="notif-text">
                    vous a envoyer un message
                </div>
            </div>
        `
        notifDiv.appendChild(notifContainer)

        maincontainer.appendChild(notifDiv)

        setTimeout(function () {
            notifDiv.style.animation = "slideUp 2s forwards";
            notifDiv.addEventListener('animationend', () => {
                maincontainer.removeChild(notifDiv);
            });
        }, 4000);
    }
    wsNotifications.onclose = function (event) {
        // alert('ws for notifs is out of service,', event.data)
    }
}

{/* <div class="notif-message">
<div class="notif-sender">
    
</div>
<div class="notif-text">
    vous a envoyer un message
</div>
</div> */}



