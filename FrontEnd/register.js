url = "http://localhost:8080"
// console.log("jsdlfajsl")
const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');

signUpButton.addEventListener('click', () => {
	container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
	container.classList.remove("right-panel-active");
});

function signIn() {
	if (!validate_signin) {
		return
	}

	console.log("here")
	let enter_form = document.forms["signin-form"];
	data = `username=${enter_form["username"].value}&password=${enter_form["password"].value}`
	var request = {
		credentials: "same-origin",
		method: 'POST',
		headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' },
		body: data
	};
	fetch(url + "/api/signin", request).then(function (response) {
		if (response.status == 200) {
			// location.replace(url)
		} else {
			response.text().then(function (res) {
				window.confirm(JSON.parse(res)["message"])
			})
		}
	}).catch(function (error) {
		console.log("Error: " + error);
	})
}

function signUp() {
	let register_form = document.forms["signup-form"];
	if (!validate_signup(register_form)) {
		return
	}

	// data = `fullname=${register_form["fullname"].value}&username=${register_form["username"].value}&email=${register_form["email"].value}&password=${register_form["password"].value}`
	data = `username=${register_form["username"].value}&email=${register_form["email"].value}&password=${register_form["password"].value}`
	var request = {
		method: 'POST',
		headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' },
		body: data
	};
	fetch(url + "/api/signup", request).then(function (response) {
		stat = response.status
		if (stat == 201) {
			// location.replace(url)
			response.text().then(function (res) {
				console.log(JSON.parse(res)["message"])
			})
		} else {
			response.text().then(function (res) {
				window.confirm(JSON.parse(res)["message"])
			})
		}
	}).catch(function (error) {
		console.log("Error: " + error);
	})
}

function validate_signin() {
	let enter_form = document.forms["signin-form"];
	if (enter_form["username"].value == "") {
		window.confirm("please enter a valid username")
		enter_form["username"].focus();
	}
	else if (enter_form["password"].value == "" || enter_form["password"].value.length <= 5) {
		window.confirm("please enter valid password")
		enter_form["password"].focus();
	}
	else {
		return true
	}
	return false
}

function validate_signup(register_form) {
	if (register_form["fullname"].value == "") {
		window.confirm("please enter your fullname")
		register_form["fullname"].focus();
	}
	else if (register_form["username"].value == "") {
		window.confirm("please enter a username")
		register_form["username"].focus();
	}
	else if (!ValidateEmail(register_form["email"].value)) {
		window.confirm("please enter a valid email")
		register_form["email"].focus();
	}
	else if (register_form["password"].value == '' || register_form["password"].length < 5) {
		window.confirm("please enter a valid password")
		register_form["password"].focus();
	}
	else {
		return true
	}
	return false
}

function ValidateEmail(mail) {
	if (/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/.test(mail)) {
		return (true)
	}
	return (false)
}