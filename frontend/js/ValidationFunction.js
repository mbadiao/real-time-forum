// Fonction pour vérifier si un email est valide
function validateEmailOrUsername(emailorusername){
  const regex1 = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;;
  const regex2 = /^[a-zA-Z0-9]{1,30}$/;
  if (!regex1.test(emailorusername) && !regex2.test(emailorusername)) {
    Toast("Invalid Mail or Username");
    return false;
  }
  return true;
}

// Fonction pour vérifier si les champs sont remplis 
const validateTextinput = (textInput, errorMessage) => {
  const regex = /^[a-zA-Z0-9 ]{1,20}$/;
  if (!regex.test(textInput)) {
    Toast(errorMessage);
    return false;
  }
  return true;
};


function validateEmail(email) {
  const regex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  if (!regex.test(email)) {
    Toast("Invalid Mail Format")
    return false
  }
  return true
}



const validateGender = (textInput, errorMessage) => {
  if (textInput!="male" && textInput!="female" && textInput!="other") {
    Toast(errorMessage);
    return false;
  }
  return true;
};

// Fonction pour vérifier si le mot de passe est valide
const validatePassword = (password) => {
  const regex = /\w| |\d{8,}/;
  if (!regex.test(password)) {
    Toast("Invalid password format 8 characters minimum");
    return false;
  }
  return true;
};

// Fonction pour vérifier si l'age est valide
const validInputAge = (age) => {
  const regex = /^[0-9]{1,3}$/;
  if (!regex.test(age)) {
    Toast("Invalid Age");
    return false;
  }
  return true;
};
// Fonction pour vérifier si tous les champs de la page d'inscription sont valides  
const allRegisterFieldAreValidated = (payload) => {
  const { firstname, lastname, username, email, password, age, gender } = payload;
  return validateTextinput(firstname, "Invalid First Name") && validateTextinput(lastname, "Invalid Last Name") && validateTextinput(username, "invalid User Name") && validateEmail(email) && validatePassword(password) && validInputAge(age) && validateGender(gender, "invalid gender");
}

// Fonction pour vérifier si les champs sont remplis et si l'email est valide
const allLonginFieldAreValidated = (payload) => {
  const { emailorusername, password } = payload;
  const validemailorusername = validateEmailOrUsername(emailorusername)
  const validpassword = validatePassword(password)
  return validemailorusername && validpassword
}

// La  fonction Toast pour afficher un message dans le front-end
const Toast = (message) => {
  const toast = document.querySelector(".toast")
  toast.textContent = message;
  toast.classList.add("show");
    setTimeout(() => {
      toast.textContent = '';
      toast.classList.remove("show");
    }, 3000);
};

export { allRegisterFieldAreValidated, allLonginFieldAreValidated, Toast };
