// import {
//   displayErrorInHtml,
//   displayFatalErrorInHtml,
//   displayPostInHtml,
// } from "./templates.js";
// let ws;

import { getUserInfo } from "./app.js";
import {
  displayErrorInHtml,
  displayFatalErrorInHtml,
  displayPostInHtml,
  errorInHtml,
} from "./templates.js";
let ws;
// function initPostWebSocket() {
//   ws = new WebSocket("ws://localhost:8080/post");
//   ws.onmessage = handlePost;
//   ws.onopen = handleOpen;
//   ws.onclose = handleClose;
// }



function initPostWebSocket() {
  console.log("InitPostWS");
  ws = new WebSocket("ws://localhost:8080/post");
  ws.onmessage = handlePost;
  ws.onopen = handleOpen;
  ws.onclose = handleClose;
}

// function handlePost(event) {
//   try {
//     const data = JSON.parse(event.data);
//     console.log("Message received:", data);
//     if (data.Errors.Status) {
//       displayFatalErrorInHtml(data.Errors);
//     } else {
//       displayPostInHtml(data);
//     }
//   } catch (error) {
//     console.error("Error parsing server response:", error);
//   }
// }

function handlePost(event) {
  try {
    const data = JSON.parse(event.data);
    if (data){
      if (data.Errors.Status) {
        const { username, firstname } = getUserInfo();
        errorInHtml(data.Errors.Code, data.Errors.Message, username, firstname);
      } else {
        displayPostInHtml(data);
      }
    }
  } catch (error) {
    console.error("Error parsing server response post:", error);
    const { username, firstname } = getUserInfo();
    errorInHtml("500", "Error parsing server response", username, firstname);
  }
}


function handleOpen() {
  console.log("Connected to the WebSocket server");
  ws.send(JSON.stringify({ type: "get_all_posts" }));
}

function handleClose(event) {
  if (event.wasClean) {
    console.log(
      `Connection closed cleanly, code=${event.code}, reason=${event.reason}`
    );
  } else {
    console.error("Connection died");
  }
}

const Post = () => {
  const postForm = document.getElementById("postForm");
  if (postForm) {
    // console.log(postForm);
    postForm.addEventListener("submit", (event) => {
      console.log("addEventlistener");
      event.preventDefault();
      const title = document.getElementById("title").value;
      const cat = document.getElementById("category").value.trim().split(" ");
      const content = document.getElementById("content").value;
      // const imageInput = document.getElementById("image");
      // const image = imageInput.files[0];
      verifyAndCreatePost(title, cat, content);
    });
  }
};

function createPost(title, category, content) {
  const message = {
    type: "new_post",
    title: title,
    content: content,
    category: category,
    image: { Name: "noPhoto", Size: 0, Status: false, Type: "noPhoto" },
  };
  isValidImage;
  ws.send(JSON.stringify(message));
}

// verifie la validité des input du post et leur redirection
async function verifyAndCreatePost(title, category, content, image) {
  const errors = validateInputs(title, category, content, image);
  let imageObject = {
    status: false,
    // path: "noPhoto"
  };
  if (Object.keys(errors).length > 0) {
    displayErrors(errors);
    return;
  }

  if (image) {
    const isValid = await isValidImage(image);
    console.log("ISVALID:", isValid, "IMAGE", image);
    if (!isValid) {
      displayErrorInHtml("Invalid image format");
      return;
    } else {
      imageObject = {
        name: image.name,
        size: image.size,
        type: image.type,
        status: true,
      };
    }
  }
  createPost(title, category, content, imageObject);
  resetForm();
}

// verifie si l'image qu'on souhaite upload est correcte
function isValidImage(file) {
  return new Promise((resolve, reject) => {
    if (!file) {
      reject("No file provided");
      return;
    }

    const validTypes = [
      "image/svg+xml",
      "image/jpeg",
      "image/gif",
      "image/png",
    ];
    const reader = new FileReader();

    reader.onloadend = () => {
      const arr = new Uint8Array(reader.result).subarray(0, 4);
      let header = "";
      for (let i = 0; i < arr.length; i++) {
        header += arr[i].toString(16);
      }
      let type = "";
      switch (header) {
        case "89504e47":
          type = "image/png";
          break;
        case "47494638":
          type = "image/gif";
          break;
        case "ffd8ffe0":
        case "ffd8ffe1":
        case "ffd8ffe2":
          type = "image/jpeg";
          break;
        default:
          type = "unknown";
          break;
      }

      if (validTypes.includes(type)) {
        resolve(true);
      } else {
        resolve(false);
      }
    };

    reader.onerror = (error) => {
      reject(error);
    };

    reader.readAsArrayBuffer(file.slice(0, 4));
  });
}

// gere les eventuelles erreur de saisie
function validateInputs(title, category, content, image) {
  const validations = [
    {
      field: "title",
      required: true,
      maxLength: 30,
      errorMessage:
        "A title is required and must be no more than 30 characters long.",
    },
    {
      field: "category",
      required: true,
      maxLength: 100,
      errorMessage: "You must choose at least one category.",
    },
    {
      field: "content",
      required: true,
      maxLength: 500,
      errorMessage:
        "Content is required and must be no more than 500 characters long.",
    },
    {
      field: "image",
      maxSize: 10485760,
      errorMessage: "The image size exceeds 10MB. Please reduce its size.",
    },
  ];

  let errors = {};
  validations.forEach((validation) => {
    const { field, required, maxLength, maxSize, errorMessage } = validation;
    const value = eval(field);

    if (required && !value) {
      errors[field] = errorMessage;
    } else if (value && field !== "category" && value.length > maxLength) {
      errors[field] = errorMessage;
    } else if (field === "image" && image && image.size > maxSize) {
      errors[field] = errorMessage;
    }

    if (field === "category" && required) {
      const trimmedCategory = value
        .map((cat) => cat.trim())
        .filter((cat) => cat !== "");
      if (trimmedCategory.length === 0) {
        errors[field] = "Category cannot be empty or contain only spaces.";
      }
    }
  });

  return errors;
}

// reset les inputs
function resetForm() {
  const modal = document.getElementById("modal");
  if (modal) {
    modal.style.display = "none";
  }
}

function displayErrors(errors) {
  Object.values(errors).forEach((error) => displayErrorInHtml(error));
}

export { Post, initPostWebSocket };
