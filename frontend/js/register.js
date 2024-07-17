import {
  allLonginFieldAreValidated,
  allRegisterFieldAreValidated,
  Toast,
} from "./ValidationFunction.js";
import { renderHomePage } from "./app.js";

function setCookie(name, value, days, path, domain, secure) {
  let expires = "";
  if (days) {
    let date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    expires = "; expires=" + date.toUTCString();
  }
  let cookie = name + "=" + (value || "") + expires + "; path=" + (path || "/");
  if (domain) {
    cookie += "; domain=" + domain;
  }
  if (secure) {
    cookie += "; secure";
  }
  document.cookie = cookie;
}

function getCookie(name) {
  let nameEQ = name + "=";
  let ca = document.cookie.split(";");
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == " ") c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}

const handleSubmit = (data, action) => {
  WebSocketConnection(`ws://localhost:8080/${action}`, data);
};

// Fonction pour recuperer les données du formulaire de connexion
const getLoginDataForm = () => {
  const loginForm = document.getElementById("loginForm");
  if (!loginForm) {
    return;
  }
  loginForm.addEventListener("submit", function (event) {
    event.preventDefault();
    const emailorusername = document.getElementById("emailorusername").value;
    const password = document.getElementById("password").value;
    const data = {
      emailorusername: emailorusername,
      password: password,
    };
    if (!allLonginFieldAreValidated(data)) {
      return;
    }
    handleSubmit(JSON.stringify(data), "login");
  });
};

// Ce fichier contient les fonctions qui permettent de récupérer les données du formulaire d'inscription
const getRegisterDataForm = () => {
  const registerForm = document.getElementById("registerForm");
  if (!registerForm) {
    return;
  }
  registerForm.addEventListener("submit", function (event) {
    event.preventDefault();
    const firstname = document.getElementById("firstname").value;
    const lastname = document.getElementById("lastname").value;
    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    const age = document.getElementById("age").value;
    const gender = document.getElementById("gender").value;

    const data = {
      firstname: firstname,
      lastname: lastname,
      username: username,
      email: email,
      password: password,
      age: age,
      gender: gender,
    };

    if (!allRegisterFieldAreValidated(data)) {
      return;
    }
    handleSubmit(JSON.stringify(data), "register");
  });
};

// Fonction pour envoyer les données du formulaire d'inscription au serveur
const WebSocketConnection = (url, data) => {
  const socket = new WebSocket(url);
  socket.onopen = function () {
    console.log("WebSocket connection opened");
    socket.send(data);
  };

  socket.onmessage = function (event) {
    const { cookie, errorRegister, errorLogin, username, firstname } = JSON.parse(event.data);
    if (errorRegister) {
      Toast(errorRegister);
      return;
    }
    if (errorLogin) {
      Toast(errorLogin);
      return;
    }
    if (cookie) {
      setCookie("ForumCookie", cookie, 1, "/", "localhost", true);
      renderHomePage(username, firstname);
    }
  };

  socket.onclose = function () {
    console.log("WebSocket connection closed");
  };

  socket.onerror = function (error) {
    console.error("WebSocket error:", error);
  };
};


export {
  getLoginDataForm,
  getRegisterDataForm,
  WebSocketConnection,
  getCookie,
  setCookie,
};
