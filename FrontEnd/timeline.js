url = "http://localhost:8080"
lorem_ipsum = "hLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

function Post() {
    post_content = document.getElementById("post_textarea").value
    data = `content=${post_content}&parent=`
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

function Profile(user_id) {
    if (user_id == "") {
        // console.log(document.getElementsByClassName("fa-home")[0].parentNode)
        document.getElementsByClassName("fa-home")[0].parentNode.classList.remove("option-active")
        document.getElementsByClassName("fa-user-o")[0].parentNode.classList.add("option-active")
    }
    var request = {
        method: 'GET',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' }
    };
    targetURL = user_id == "" ? url + "/api/profile" : url + "/api/otherProfile/" + user_id
    fetch(targetURL, request).then(function (response) {
        stat = response.status
        if (stat == 200) {
            response.text().then(function (res) {
                profile_json = JSON.parse(res)
                console.log(profile_json)
                make_profile(profile_json)
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

function Follow_unfollow() {
    user_id = document.getElementsByClassName("profile-user-name")[0].innerHTML
    var request = {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' }
    };
    fetch(url + "/api/follow/" + user_id, request).then(function (response) {
        stat = response.status
        if (stat == 201) {
            response.text().then(function (res) {
                message = JSON.parse(res)["message"]
                if (message == "followed") {
                    document.getElementsByClassName("profile-edit-btn")[0].innerHTML = "unfollow"
                    document.getElementsByClassName("profile-stat-count")[1].innerHTML = +document.getElementsByClassName("profile-stat-count")[1].innerHTML + 1
                }

                else {
                    document.getElementsByClassName("profile-edit-btn")[0].innerHTML = "follow"
                    document.getElementsByClassName("profile-stat-count")[1].innerHTML -= 1
                }
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

function like(outer_div) {
    // console.log(outer_div.childNodes)
    // outer_div.childNodes[1].innerHTML = +outer_div.childNodes[].innerHTML + 1
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
                    outer_div.childNodes[1].innerHTML = +outer_div.childNodes[1].innerHTML + 1
                else
                    outer_div.childNodes[1].innerHTML = +outer_div.childNodes[1].innerHTML - 1

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

function comment(comment) {
    content = comment.parentNode.childNodes[1].value
    parent = document.getElementById("comments-others-tweets").childNodes[0].id
    data = `content=${content}&parent=${parent}`
    console.log(data)
    var request = {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' },
        body: data
    };
    fetch(url + "/api/createPost", request).then(function (response) {
        stat = response.status
        if (stat == 201) {
            document.getElementsByClassName("commentText")[0].value = ""
            show_comments(current_post_comments)
        } else {
            response.text().then(function (res) {
                window.confirm(JSON.parse(res)["message"])
            })
        }
    }).catch(function (error) {
        console.log("Error: " + error);
    })
}

function Search_user() {
    user_id = document.getElementById("search").value
    Profile(user_id)
}

function bookmark(outer_div) {
    var request = {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' }
    };
    fetch(url + "/api/mark/" + outer_div.parentNode.parentNode.id, request).then(function (response) {
        stat = response.status
        if (stat == 201) {
            // outer_div.childNodes[0].classList.toggle("fa-heart")
            // outer_div.childNodes[0].classList.toggle("fa-heart-o")
            // outer_div.classList.toggle("liked")
            response.text().then(function (res) {
                // if (JSON.parse(res)["message"] == "marked")
                //     outer_div.childNodes[1].innerHTML = +outer_div.childNodes[1].innerHTML + 1
                // else
                //     outer_div.childNodes[1].innerHTML = +outer_div.childNodes[1].innerHTML - 1
                console.log(res)

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
    outer_div.classList.toggle("bookmarked")
}

function show_comments(post) {
    current_post_comments = post
    comments_others_tweets = document.getElementById("comments-others-tweets")
    comments_others_tweets.innerHTML = ""
    comments_others_tweets.appendChild(post.parentNode.parentNode.cloneNode(true))
    var request = {
        method: 'GET',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' }
    };
    fetch(url + "/api/getComments/" + post.parentNode.parentNode.id, request).then(function (response) {
        stat = response.status
        if (stat == 200) {
            response.text().then(function (res) {
                for (post of JSON.parse(res).comment) {
                    comments_others_tweets.appendChild(make_post(post.id, post.fullName, post.creator,
                        post["created-at"], post.content, post.commentNumber, post.likeNumber, post.like, post.mark))
                }
                console.log(JSON.parse(res))
            })
        } else {
            response.text().then(function (res) {
                window.confirm(JSON.parse(res)["message"])
            })
        }
    }).catch(function (error) {
        console.log("Error: " + error);
    })

    // for (i=1;i <= 4; i++){
    //     comments_others_tweets.appendChild(make_post(123456, "matin fotouhi", "matin_ft", "yesterday", lorem_ipsum, 3, 100, true, false))
    // }
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
            <div class="comment" onclick="show_comments(this)"><i class="fa fa-comment-o"></i><p>${post_comments}</p></div>
            <div class=\"${"like" + (isLiked ? " liked" : "")}\" onclick="like(this)"><i class=\"${isLiked ? "fa fa-heart" : "fa fa-heart-o"}\"></i><p>${post_likes}</p></div>
            <div class=\"${"bookmark" + (isBookmarked ? " bookmarked" : "")}\" onclick="bookmark(this)"><i class="fa fa-bookmark"></i></div>
        </div>
    `
    return post
}

function load_timeline() {
    add_timeline_header()
    document.getElementsByClassName("fa-home")[0].parentNode.classList.add("option-active")
    document.getElementsByClassName("fa-user-o")[0].parentNode.classList.remove("option-active")
    other_tweet = document.getElementsByClassName("others-tweets")[0]
    other_tweet.innerHTML = ""
    var request = {
        method: 'GET',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8' }
    };
    fetch("/api/timeline", request).then(function (response) {
        stat = response.status
        if (stat == 200) {
            response.text().then(function (res) {
                timeline = JSON.parse(res)["timeLine"]
                if (timeline != null)
                    for (post of timeline)
                        other_tweet.appendChild(make_post(post.id, post.fullName, post.creator, post["created-at"], post.content, post.commentNumber, post.likeNumber, post.like, post.mark))

            })
        } else {
            response.text().then(function (res) {
                window.confirm(JSON.parse(res)["message"])
            })
        }
    }).catch(function (error) {
        console.log("Error: " + error);
    })
    // document.getElementById("content-menu").appendChild(make_profile_header("matin_ft", 100, 200, 300, "this is my bio"))

    // document.getElementsByClassName("content-menu")[0].appendChild(make_profile_header("matin_ft", 100, 200, 300, "this is my bio"))
    // for (i = 1; i <= 4; i++) {
    //     document.getElementById("others-tweets").appendChild(make_post(123456, "matin fotouhi", "matin_ft", "yesterday", lorem_ipsum, 3, 100, true, false))
    // }
}

function add_timeline_header() {
    document.getElementById("content-menu").innerHTML = `<div class="prefer">
                <span>
                    <a href="">Home</a>
                </span>
                <span>
                    <i class="fa fa-star-o"></i>
                </span>
            </div>

            <div class="you-tweet-other-tweet">
                <div class="your-tweet">
                    <div class="profile-message">
                        <span><img src="images/nani.png" alt=""></span>
                        <span><textarea placeholder="what's happening" rows="5" id="post_textarea"></textarea></span>
                    </div>
                    <div class="add-extra">
                        <!-- <div class="images-more">
                            <span><a href=""><i class="fa fa-picture-o"></i></a></span>
                            <span><a href=""><i class="fa fa-bars"></i></a></span>
                            <span><a href=""><i class="fa fa-smile-o"></i></a></span>
                            <span><a href=""><i class="fa fa-calender-plus-o"></i></a></span>
                        </div> -->
                        <span ><button onclick="Post()">Post</button></span>
                    </div>
                </div>

                <div class="others-tweets" id="others-tweets">



                </div>


            </div>

        </div>`
    
}

function make_profile(profile) {
    you_other_tweet = document.getElementsByClassName("you-tweet-other-tweet")[0]
    you_other_tweet.innerHTML = ""
    you_other_tweet.appendChild(make_profile_header(profile.username, profile.posts == null ? 0 : profile.posts.length, profile.following, profile.followers, profile.bio, profile.follow == null ? true : false, profile.follow))
    for (post of profile.posts) {
        // console.log(post["created-at"])
        // console.log(post.id, post.fullName, profile.username, post.created-at)
        you_other_tweet.appendChild(make_post(post.id, post.fullName, post.creator, post["created-at"], post.content, post.commentNumber, post.likeNumber, post.like, post.mark))
    }

    // console.log(profile)
}

function make_profile_header(user_id, posts, following, followers, bio, me, follow) {
    profile = document.createElement("div")
    profile.classList.add("profile")
    // console.log("Follow_unfollow(\"${user_id}\")")
    profile.innerHTML = `
    <div class="profile-image">
    
        <img src="images/profile-photo.jpg"
            alt="">
    
    </div>
    
    <div class="profile-user-settings">
    
        <h1 class="profile-user-name">${user_id}</h1>
    
        <button class="btn profile-edit-btn" id="profile-edit-btn" onclick="Follow_unfollow()">${me == true ? "edit profile" : (follow == true ? "unfollow" : "follow")}</button>
    
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

function log_out() {
    var request = {
        method: 'POST'
    };
    fetch(url + "/logout", request).then(function (response) {
        location.replace(url)
    }).catch(function (error) {
        console.log("Error: " + error);
    })
}