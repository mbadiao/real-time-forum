import { login, register, home } from "./templates.js";
import { initPostWebSocket, Post } from "./posts.js";
import { userStatus } from "./showOnlineUSer.js";
import { HandleMessage } from "./messages.js";
import {
  getLoginDataForm,
  getRegisterDataForm,
  getCookie,
  setCookie,
} from "./register.js";
import { wsNotifs } from "./notif.js";
// import './comments.js'

let userInfo = {};
const mainContainer = document.querySelector(".mainContainer");
const createRegister = () => {
  mainContainer.classList.remove("home");
  mainContainer.innerHTML = "<div class='register'></div>";
};

const renderLogin = () => {
  createRegister();
  const auth = document.querySelector(".register");
  auth.innerHTML = login;
  attachEventListeners();
  getLoginDataForm();
};

const renderRegister = () => {
  createRegister();
  const auth = document.querySelector(".register");
  auth.innerHTML = register;
  attachEventListeners();
  getRegisterDataForm();
};

const attachEventListeners = () => {
  const singin = document.querySelector(".singin");
  const singup = document.querySelector(".singup");

  if (singup) {
    singup.addEventListener("click", (event) => {
      event.preventDefault();
      renderRegister();
    });
  }

  if (singin) {
    singin.addEventListener("click", (event) => {
      event.preventDefault();
      renderLogin();
    });
  }
};

const createModalTemplate = () => {
  return `
            <div class="modal-content" id="modal-Element">
                <span class="close">&times;</span>
                <h2>Add Post</h2>
                <form id="postForm">
                    <div class="form-group">
                        <input placeholder="Title" class="titlepost" type="text" id="title" name="title" required>
                    </div>
                    <div class="form-group">
                        <input placeholder="Content" class="contentpost" type="text" id="content" name="content" required>
                    </div>
                    <div class="form-group">
                        <input placeholder="Category" class="categorypost" type="text" id="category" name="category" required>
                    </div>
                    <button id="submit" class="submitpost" type="submit">Submit</button>
                </form>
            </div>
        `;
};

const renderHomePage = (username, firstname) => {
  const mainContainer = document.querySelector(".mainContainer");
  mainContainer.classList.add("home");
  mainContainer.innerHTML = home;
  initPostWebSocket();
  wsNotifs();
  const logout = document.querySelector(".logout");
  if (logout) {
    logout.addEventListener("click", function (event) {
      event.preventDefault();
      setCookie("ForumCookie", "", -1);
      renderLogin();
    });
  }
  userStatus();

  const modal = document.getElementById("modal");
  const btn = document.querySelector(".addPost");
  btn.onclick = function () {
    modal.innerHTML = createModalTemplate();
    modal.style.display = "flex";
    Post(modal);
    const span = document.querySelector(".close");
    span.onclick = function () {
      modal.innerHTML = "";
      modal.style.display = "none";
    };
    window.onclick = function (event) {
      if (event.target === modal) {
        modal.innerHTML = "";
        modal.style.display = "none";
      }
    };
  };

  document.querySelector(".profileid").textContent = username[0]
  document.getElementById("myname").textContent = firstname
  document.getElementById("myusername").textContent = "@" + username
};

const renderLoginOrHomePages = () => {
  const cookie = getCookie("ForumCookie");

  fetch("http://localhost:8080/checksession", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ cookievalue: cookie }),
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.processedValue !== "ok") {
        renderLogin();
      } else {
        userInfo = {
          username: data.username,
          firstname: data.firstname
        };
        renderHomePage(data.username, data.firstname);
      }
    })
    .catch((error) => {
      console.error("Error front end cookie", error);
    });
};

function getUserInfo() {
  return userInfo;
}

HandleMessage();
renderLoginOrHomePages();
export { renderHomePage, renderLoginOrHomePages, getUserInfo };
