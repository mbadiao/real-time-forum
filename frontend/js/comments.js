import { getCookie } from "./register.js";
import { errorInHtml } from "./templates.js";
import { getUserInfo } from "./app.js";

const ws = new WebSocket("ws://localhost:8080/comment");

ws.onopen = () => {
  console.log("Connected to the WebSocket server");
};

ws.onmessage = (event) => {
  try {
    const data = JSON.parse(event.data);
    if (data){
      if (data.Errors && data.Errors.Status) {
        const { username, firstname } = getUserInfo();
        errorInHtml(data.Errors.Code, data.Errors.Message, username, firstname);
        return;
      }
      console.log("Received data from the server:", data)
      if (data.comment) {
        const { postid} = data.comment;
        updatePostComments(postid, data);
      }
    }
  } catch (error) {
    console.error("Error parsing server response:", error);
    const { username, firstname } = getUserInfo();
    errorInHtml("500", "Error parsing server response comment", username, firstname);
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

function updatePostComments(postId, newComment) {
    const {comment, user} = newComment;
  const postElement = document.getElementById(`post-${postId}`);
  const commentsSection = postElement.querySelector(".badgeUser.comment.off");
  if (commentsSection) {
    const commentElement = document.createElement("div");
    commentElement.classList.add("comment");
    commentElement.innerHTML = `
      <div class="badge">${user.charAt(0)}</div>
      <div>
        <div class="name">${user}</div>
        <div class="text">${comment.content}</div>
      </div>
    `;
    if (commentsSection.firstChild) {
        commentsSection.insertBefore(commentElement, commentsSection.firstChild);
      } else {
        commentsSection.appendChild(commentElement);
      }
  }
  postElement.querySelector(".commentInput").value = "";
}

export function Comment(postId) {
  const commentForm = document.querySelector(`.sendCommentWs${postId}`);
  const postIds = commentForm.getAttribute("data-postid");
  const userIds = commentForm.getAttribute("data-userid");
  const contents = commentForm.previousElementSibling.value;
  verifyAndCreateComment(postIds, userIds, contents);
}

function verifyAndCreateComment(postId, userId, content) {
  const validation = isValidComment(content);
  if (validation.isValid) {
    createComment(postId, userId, content);
  } else {
    console.error(validation.errorMessage);
  }
}

function isValidComment(content) {
  const maxLength = 500;
  if (content && content.length <= maxLength) {
    return { isValid: true, errorMessage: "" };
  } else {
    return {
      isValid: false,
      errorMessage:
        "Content is required and must be no more than 500 characters long.",
    };
  }
}

function createComment(postId, userId, content) {
  const cookie = getCookie("ForumCookie");
  const comment = {
    postid: Number(postId),
    userid: Number(userId),
    content: content,
    cookie: cookie,
  };
  console.log(JSON.stringify(comment));
  ws.send(JSON.stringify(comment));
}
