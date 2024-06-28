import { Comment } from "./comments.js";
import { renderHomePage } from "./app.js";
import { Like } from "./like.js"

export const register = `<div class="inputField">
<div class="loginForm">
    <div class="heading">
        <h1>Welcome buddy</h1>
        <img class="logo" src="../static/icons/Vector.svg" alt="logo">
    </div>
    <form class="inputForm" id="registerForm">
        <input class="input" placeholder="Firstname" type="text" id="firstname" name="firstname">
        <input class="input" placeholder="lastName" type="text" id="lastname" name="lastname">
        <input class="input" placeholder="userName" type="text" id="username" name="username">
        <input class="input" placeholder="Email" type="email" id="email" name="email">
        <input class="input" placeholder="Password" type="password" id="password" name="password">
        <input class="input" placeholder="Age" type="number" id="age" name="age">
        <input class="input" placeholder="Gender     ex: male, female or other" type="text" id="gender" name="gender">
        <button class="submit" type="submit">Sign Up</button>
    </form>
    <div class="acount">Already have an account? <button class="singin">Sign In</button></div>
    </div>
</div>
<div class="Allimage">
<div class="imgField">
    <div class="bento2">
        <img class="bento2img" src="../static/assets/bento2.jpg" alt="image register">
    </div>
    <div class="bento3">
        <img class="bento3img" src="../static/assets/bento3.jpg" alt="image register">
    </div>
</div>
<div class="bento1">
    <img class="bento1img" src="../static/assets/bento1.jpg" alt="image register">
</div>
</div>
`;

export const login = `<div class="inputField">
<div class="loginForm">
    <div class="heading">
        <h1>Welcome back</h1>
        <img class="logo" src="../static/icons/Vector.svg" alt="logo">
    </div>
    <form class="inputForm" id="loginForm" >
        <input class="input" placeholder="Email or Username" type="text" id="emailorusername" name="email">
        <input class="input" placeholder="Password" type="password" id="password" name="password">
        <button class="submit" type="submit">Log in</button>
    </form>
    <div class="acount">Don’t have an account? <button class="singup">Sign Up</button></div>
    </div>
</div>
<div class="Allimage">
<div class="imgField">
    <div class="bento2">
        <img class="bento2img" src="../static/assets/bento2.jpg" alt="image register">
    </div>
    <div class="bento3">
        <img class="bento3img" src="../static/assets/bento3.jpg" alt="image register">
    </div>
</div>
<div class="bento1">
    <img class="bento1img" src="../static/assets/bento1.jpg" alt="image register">
</div>
</div>`;

export const notifHtml = `
    <div class="notif-container">
        <div class="notifAvatar">
            <img style="width: 50px;" src="../static/assets/avatar.svg" alt="">
        </div>
    </div>`;

export const message = `
<div class="message">
    <div class="messageheader">
        <div class="oneusersmessage">
            <div class="userandtyping">
                <div class="userbadgesmessage">E</div>
                <div class="userfisrtandlastnamemessage">Dribble dribble</div>
            </div>
            <div class="istyping"></div>
        </div>
    </div>
    <div class="chat-container">

    </div>
    <div class="messagefooter">
        <form class="messageChat" method="post">
            <input class="messageInput" type="text" title="comment" placeholder="start chat..." required="">
            
                <button type="submit" class="sendChat">
        <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24"
             fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
             class="lucide lucide-send-horizontal">
            <path d="m3 3 3 9-3 9 19-9Z" />
            <path d="M6 12h16" />
        </svg>
    </button>
        </form>
    </div>
</div>
`;
{/* <button></button> */ }

export const home = `
<div>
    <div class="header">
        <div class="logohome"><svg width="4vw" height="6vh" viewBox="0 0 126 179" fill="none"
            xmlns="http://www.w3.org/2000/svg">
            <path
                d="M104.498 25.815C109.754 16.715 116.502 7.981 125.966 0C107.346 7.25 89.946 17.36 74.416 29.92C41.726 54.07 16.806 89.19 5.866 131.6C5.266 133.93 4.716 136.26 4.196 138.58V138.6C4.086 139.13 3.966 139.66 3.856 140.18C3.8 140.45 3.752 140.72 3.697 140.991C3.517 141.874 3.337 142.76 3.169 143.64C3.148 143.755 3.13 143.866 3.109 143.98C3.033 144.39 2.966 144.8 2.893 145.21C2.841 145.48 2.794 145.75 2.749 146.011C2.688 146.359 2.632 146.702 2.573 147.048C2.257 148.921 1.971 150.793 1.715 152.664C1.663 153.036 1.611 153.407 1.559 153.781C1.498 154.223 1.448 154.662 1.396 155.102C1.37 155.308 1.34 155.515 1.315 155.722C1.317 155.721 1.318 155.72 1.32 155.719V155.722C0.89 159.232 0.570001 162.732 0.360001 166.222L0.317001 167.126C1.45 167.202 2.349 168.136 2.349 169.288C2.349 170.49 1.374 171.465 0.171997 171.465C0.150997 171.465 0.132 171.459 0.111 171.459L0.090004 171.902L0 173.752V178.772H4.26C4.31 178.622 4.38 178.472 4.44 178.322C4.94 177.162 5.46 175.982 6.03 174.832C8.3 170.152 11.14 165.772 14.42 163.702C7.94 155.362 9.17 144.262 10.06 139.892C10.071 139.825 10.085 139.765 10.099 139.703C10.236 139.094 10.372 138.492 10.518 137.882C10.888 136.282 11.268 134.692 11.678 133.102C16.798 113.282 25.058 95.112 35.798 78.992C24.728 100.032 20.118 118.032 17.648 134.772C18.078 134.582 18.48 134.323 18.928 134.182C54.094 123.102 69.615 105.258 79.559 85.102H51.238L82.907 77.794C84.159 74.862 85.331 71.894 86.466 68.902H69.958L88.02 64.734C92.227 53.299 96.153 41.616 101.998 30.382H84.718L104.498 25.815Z"
                    fill="#116EFF" />
                </svg>
        </div>
        <div class="logout">logout</div>
    </div>
    <div class="appcontainer">
        <div class="profileAndOnline">
            <div class="profileprofile">
                <div class="profileid">D</div>
                <div>
                    <div class="profilecontent">
                        <div class="profileusertag">
                            <div id="myname">Dribble</div>
                                <div id="myusername">@dribble</d>
                                    </div>
                                </div>
                                <div class="profileposts">
                                    <div>Posts</div>
                                    <div>1.2k</div>
                                </div>
                            </div>
                        </div>
                        <div class="profileLikes">
                            <div>Likes</div>
                            <div>1.2k</div>
                        </div>
                    </div>
                    <div class="usersLeft">
                        <div class="usertitleLeft">Users</div>
                        <div class="usersbackgroundLeft">
                            
                        </div>
                        <div class="userLeft">
                        </div>
                    </div>
                </div>
                <div class="posts" id="displayPosts">
                </div>
                <div class="postsAndMessage">
                    <div class="addPost">
                        <p>Add new post....</p>
                        <svg role="img" xmlns="http://www.w3.org/2000/svg" width="30px" height="30px"
                            viewBox="0 0 24 24" aria-labelledby="plusIconTitle" stroke="#F4F7F9" stroke-width="1"
                            stroke-linecap="square" stroke-linejoin="miter" fill="none" color="#F4F7F9">
                            <title id="plusIconTitle">Plus</title>
                            <path d="M20 12L4 12M12 4L12 20" />
                        </svg>
                    </div>
                    </div>
                </div>
            </div>
        </div>
        </div>
        <div id="modal" class="modal">
        </div>
        `;

function renderOnlineOrOfflineTemplate(user) {
  if (user.online) {
    return `<div class="oneusersLeft" id=${user.user_id}>
                <span class="onlineLeft"></span>
                <div class="userbadgesLeft">${user.username[0]}</div>
                <div class="userfisrtandlastnameLeft" id=${user.user_id}>${user.firstname} ${user.lastname}</div>
                <div class="notificationLeft">2</div>
            </div>`;
  } else {
    return `<div class="oneusersLeft">
                <div class="userbadgesLeft">${user.username[0]}</div>
                <div class="userfisrtandlastnameLeft" id=${user.user_id}>${user.firstname} ${user.lastname}</div>
                <div class="notificationLeft">2</div>
            </div>`;
  }
}

function displayPostInHtml(data) {
  const displaypost = document.getElementById("displayPosts");
  if (displaypost) {
    displaypost.innerHTML = "";
    if (!data) {
      return;
    }
    if (data.Posts) {
      data.Posts.forEach((post) => {
        const div = document.createElement("div");
        div.classList.add("postContainer");
        div.innerHTML = `
          <div class="container" id="post-${post.PostID}">
            <div class="badgeUser">
              <div class="userinfo">
                <div class="badge">${post.Author.Firstname.charAt(0)}</div>
                <div class="fullname">
                  <div class="name">${post.Author.Firstname} ${post.Author.Lastname}</div>
                  <div class="username">@${post.Author.Username}</div>
                </div>
                <div class="date">${post.Formated_date}</div>
              </div>
            </div>
            <div class="postcontaint">
              <div class="title">${post.Title}</div>
              <div class="text">${post.Content}</div>
              <div class="allTags">
                ${post.Categories.map(category => `<div class="tags">#${category}</div>`).join("")}
              </div>
              <div class="allNumbers">
                <div class="numberComment">
                  <svg role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="3.5vw" height="2.5vh" aria-labelledby="chatIconTitle" stroke="#F4F7F9" fill="none" color="#F4F7F9">
                    <title id="chatIconTitle">Chat</title>
                    <path d="M8.82388455,18.5880577 L4,21 L4.65322944,16.4273939 C3.00629211,15.0013 2,13.0946628 2,11 C2,6.581722 6.4771525,3 12,3 C17.5228475,3 22,6.581722 22,11 C22,15.418278 17.5228475,19 12,19 C10.8897425,19 9.82174472,18.8552518 8.82388455,18.5880577 Z"/>
                  </svg>
                  <div class="nbrCommnent">${post.Comments_nbr}</div>
                </div>
                <div class="numberLike">
                  <svg role="img" xmlns="http://www.w3.org/2000/svg" width="3.5vw" height="2.5vh" viewBox="0 0 24 24" aria-labelledby="favouriteIconTitle" class='sendLikeWs' like-postid="${post.PostID}" like-userid="${post.UserID}" stroke="${post.Like_status ? "#FF8587" : "#F4F7F9"}" fill="none">
                    <title id="favouriteIconTitle">Favourite</title>
                    <path d="M12,21 L10.55,19.7051771 C5.4,15.1242507 2,12.1029973 2,8.39509537 C2,5.37384196 4.42,3 7.5,3 C9.24,3 10.91,3.79455041 12,5.05013624 C13.09,3.79455041 14.76,3 16.5,3 C19.58,3 22,5.37384196 22,8.39509537 C22,12.1029973 18.6,15.1242507 13.45,19.7149864 L12,21 Z"/>
                  </svg>
                  <div class="nbrLikes">${post.Like_nbr}</div>
                </div>
              </div>
            </div>
            <div class="comments">
              <input class="commentInput" type="text" title="comment" placeholder="add comment...">
              <svg id="sendComment" 
                data-postid="${post.PostID}" 
                data-userid="${post.UserID}" 
                class="sendCommentWs  sendCommentWs${post.PostID}" 
                xmlns="http://www.w3.org/2000/svg" 
                width="24" height="24" 
                viewBox="0 0 24 24"
                fill="none" 
                stroke="currentColor" 
                stroke-width="2" 
                stroke-linecap="round" 
                stroke-linejoin="round"
                class="lucide lucide-send-horizontal">
                <path d="m3 3 3 9-3 9 19-9Z" />
                <path d="M6 12h16" />
              </svg>
            </div>
            <div class="badgeUser comment off" style="display: none;"></div>
          </div>
        `;
        displaypost.appendChild(div);
        const commentsContainer = div.querySelector(".badgeUser.comment.off");
        if (post.Comments) {
          post.Comments.forEach((comment) => {
            const commentElement = document.createElement("div");
            commentElement.classList.add("comment");
            commentElement.innerHTML = `
              <div class="badge">${comment.Author.Firstname.charAt(0)}</div>
              <div>
                <div class="name">${comment.Author.Firstname} ${comment.Author.Lastname}</div>
                <div class="text">${comment.Content}</div>
              </div>
            `;
            commentsContainer.appendChild(commentElement);
          });
        }
        const postElement = document.getElementById(`post-${post.PostID}`);
        postElement.addEventListener("click", () => togglePost(post.PostID));
        const sendComment = postElement.querySelector(`.sendCommentWs${post.PostID}`);
        sendComment.addEventListener("click", () => Comment(post.PostID));
      });
    }
    const modalmystere = document.getElementById('modal');
    if (modalmystere) {
      modalmystere.style.display = 'none';
    }
    // Attach the Like event listeners after the posts are rendered
    Like();
  }
}


function togglePost(postId) {
  const postElement = document.getElementById(`post-${postId}`);
  const commentsSection = postElement.querySelector(".badgeUser.comment.off");
  if (postElement) {
    document
      .querySelectorAll(".container.postmodal")
      .forEach((modal) => modal.classList.remove("postmodal"));
    const modal = document.getElementById("modal");
    postElement.classList.add("postmodal");
    modal.style.display = "flex";
    commentsSection.style.display = "block";
    window.onclick = function (event) {
      if (event.target === modal) {
        modal.style.display = "none";
        commentsSection.style.display = "none";
        postElement.classList.remove("postmodal");
      }
    };
  }
}

// affiche les potentiels erreur du modal createpost
function displayErrorInHtml(message) {
  const modal = document.getElementById("modal-Element");
  const div = document.createElement("div");
  div.classList.add("error-modal");
  div.innerHTML = message;
  modal.appendChild(div);
  setTimeout(() => {
    modal.removeChild(div);
  }, 3000);
}

// affiche les Erreurs 400 500 404
function displayFatalErrorInHtml(data) {
  const displaypost = document.getElementById("displayPosts");
  if (displaypost) {
    const div = document.createElement("div");
    div.classList.add("postContainer");
    div.innerHTML = `
            <p>${data.Code}</p>
            <p>${data.Message}</p>
        `;
    displaypost.appendChild(div);
  }
}

function errorInHtml(code, status, username, firstname) {
  let errorHtmlContent = `
      <div class="errorpage">
          <div class="message">
              <p>Upppssss....</p>
              <div class="errorstatus">
                  <div class="statuscode">${code}</div>
                  <div>${status}</div>
              </div>
          </div>
          <div class="imageerror">
              <div class="error">
                  <img src="../static/assets/error.jpg" alt="">
              </div>
              <div class="errorcontainer">
                  <div class="error1">
                      <img src="../static/assets/error1.jpg" alt="">
                  </div>
                  <div class="error3">
                      <img src="../static/assets/error3.jpg" alt="">
                  </div>
              </div>
          </div>
      </div>
      <button class="homeerror" data-username="${username}" data-firstname="${firstname}">Home</button>`;

  const toastElements = document.querySelectorAll(".toast");
  const mainContainerElements = document.querySelectorAll(".mainContainer");

  if (toastElements) {
    toastElements.forEach((element) => element.remove());
  }

  if (mainContainerElements) {
    mainContainerElements.forEach((element) => element.remove());
  }

  document.body.insertAdjacentHTML("beforeend", errorHtmlContent);

  document.querySelector(".homeerror").addEventListener("click", function () {
    document
      .querySelectorAll(".errorpage")
      .forEach((element) => element.remove());
    document
      .querySelectorAll(".homeerror")
      .forEach((element) => element.remove());
    document.body.insertAdjacentHTML("beforeend", squeleton);

    // Afficher la page d'accueil
    renderHomePage(username, firstname);
  });
}

const squeleton = `
<div class="toast">
</div>
<div class="mainContainer home">
</div>`;

export {
  displayErrorInHtml,
  displayFatalErrorInHtml,
  displayPostInHtml,
  renderOnlineOrOfflineTemplate,
  errorInHtml,
};

// <div class="chat-bubble sender">
// <div>Hello! How can I help you today?</div>
// </div>
// <div class="chat-bubble receiver">
// <div>Hi! I need some assistance with my order.</div>
// </div>
