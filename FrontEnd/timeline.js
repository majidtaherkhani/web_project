url = "http://localhost:8080"
function post() {
    post_content = document.getElementById("post_textarea").value
    data = `content=${post_content}&parent=""`
    var request = {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' },
        body: data
    };
    fetch(url + "/api/createPost", request).then(function (response) {
        stat = response.status
        if (stat == 201) {
            location.replace(url)
            // get timeline again
            // document.getElementsByClassName("others-tweets")[0].prepend(...)
        } else {
            response.text().then(function (res) {
                window.confirm(JSON.parse(res)["message"])
            })
        }
    }).catch(function (error) {
        console.log("Error: " + error);
    })
}

function profile(user_id) {
    var request = {
        method: 'GET',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' }
    };
    fetch(url + "/api/profile", request).then(function (response) {
        stat = response.status
        if (stat == 200) {
            response.text().then(function (res) {
                profile_json = JSON.parse(res)
                make_profile(profile_json)
            })
            // get timeline again
            // document.getElementsByClassName("others-tweets")[0].prepend(...)
        } else {
            response.text().then(function (res) {
                window.confirm(JSON.parse(res)["message"])
            })
        }
    }).catch(function (error) {
        console.log("Error: " + error);
    })
}

function like(outer_div) {
    console.log(outer_div.parentNode.parentNode.id)
    var request = {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' }
    };
    fetch(url + "/api/like/" + outer_div.parentNode.parentNode.id, request).then(function (response) {
        stat = response.status
        if (stat == 200) {
            outer_div.childNodes[0].classList.toggle("fa-heart")
            outer_div.childNodes[0].classList.toggle("fa-heart-o")
            outer_div.classList.toggle("liked")
            response.text().then(function (res) {
                if (JSON.parse(res)["message"] == "post liked")
                    outer_div.childNodes[2].innerHTML = +outer_div.childNodes[2].innerHTML + 1
                else
                    outer_div.childNodes[2].innerHTML = +outer_div.childNodes[2].innerHTML - 1

            })
        } else {
            response.text().then(function (res) {
                console.log(res)
                window.confirm(JSON.parse(res)["message"])
            })
        }
    }).catch(function (error) {
        console.log("Error: " + error);
    })
}


function bookmark(outer_div) {
    // fetch bookmark ...
    outer_div.classList.toggle("bookmarked")
}

function show_comments(post) {
    commentsModal = document.getElementById("commentsModal")
    commentsModal.style.display = "block"
    // commentsModal.style.opacity = 1
    document.getElementById("container").style.filter = "blur(8px)"
}

function close_comments() {
    document.getElementById("commentsModal").style.display = "none"
    document.getElementById("container").style.filter = ""
}

function make_post(post_id, poster_fullname, poster_id, post_date, post_content, post_comments, post_likes, isLiked, isBookmarked) {
    post = document.createElement("div")
    post.id = post_id
    post.classList.add("other-tweet");
    post.innerHTML = `
        <div class="profile-msg">
            <div class="others-profile">
                <img src="images/no-image.jpg" alt="">
            </div>
            <div class="name-msg">
                <span><p><b>${poster_fullname} @${poster_id}.<small>${post_date}</small></b></p></span>
                <div class="msg">
                    <p>${post_content}</p>
                </div>
            </div>
        </div>
        <div class="your-reaction">
            <div class="comment" onclick="show_comments(post)"><i class="fa fa-comment-o"></i><p>${post_comments}</p></div>
            <div class=\"${"like" + (isLiked ? " liked" : "")}\" onclick="like(this)"><i class=\"${isLiked ? "fa fa-heart" : "fa fa-heart-o"}\"></i><p>${post_likes}</p></div>
            <div class=\"${"bookmark" + (isBookmarked ? " bookmarked" : "")}\" onclick="bookmark(this)"><i class="fa fa-bookmark"></i></div>
        </div>
    `
    console.log("your-reaction" + (isLiked ? " liked" : "") + (isBookmarked ? " bookmarked" : ""))
    return post
}

function load_timeline() {
    // document.getElementById("content-menu").appendChild(make_profile_header("matin_ft", 100, 200, 300, "this is my bio"))
    lorem_ipsum = "hLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

    // document.getElementsByClassName("content-menu")[0].appendChild(make_profile_header("matin_ft", 100, 200, 300, "this is my bio"))
    for (i=1;i <= 4; i++){
        document.getElementById("others-tweets").appendChild(make_post(123456, "matin fotouhi", "matin_ft", "yesterday", lorem_ipsum, 3, 100, true, false))
    }


}

function make_profile(profile) {
    console.log(profile)
}

function make_profile_header(user_id, posts, following, followers, bio) {
    profile = document.createElement("div")
    profile.classList.add("profile")
    profile.innerHTML = `
    <div class="profile-image">
    
        <img src="images/profile-photo.jpg"
            alt="">
    
    </div>
    
    <div class="profile-user-settings">
    
        <h1 class="profile-user-name">${user_id}</h1>
    
        <button class="btn profile-edit-btn">Edit Profile</button>
    
    </div>
    
    <div class="profile-stats">
    
        <ul>
            <li><span class="profile-stat-count">${posts}</span> posts</li>
            <li><span class="profile-stat-count">${followers}</span> followers</li>
            <li><span class="profile-stat-count">${following}</span> following</li>
        </ul>
    
    </div>
    
    <div class="profile-bio">
    
        <p><span class="profile-real-name">${bio}</p>
    
    </div>
    
    <!-- End of profile section -->
    `
    return profile
}


