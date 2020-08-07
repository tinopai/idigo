<?php

$error = "No error set";
if(isset($_GET['e']) && $_GET['e'] == "NOJAVAW") {
    $error = "No Minecraft instance found";
}

?>
<!DOCTYPE html>
<html lang="en">
<!-- CSS only -->
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css" integrity="sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk" crossorigin="anonymous">

<!-- JS, Popper.js, and jQuery -->
<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/js/bootstrap.min.js" integrity="sha384-OgVRvuATP1z7JjHLkuOU7Xw704+h835Lr+6QL9UvYjZE3Ipu6Tp75j7Bh/kR0JKI" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.0.0/animate.min.css"/>
<head>
    <meta charset="UTF-8">
</head>
<style>
* {
    overflow: hidden; /* Hide scrollbars */
}
.no-select {
    -webkit-touch-callout: none; /* iOS Safari */
    -webkit-user-select: none; /* Safari */
     -khtml-user-select: none; /* Konqueror HTML */
       -moz-user-select: none; /* Old versions of Firefox */
        -ms-user-select: none; /* Internet Explorer/Edge */
            user-select: none; /* Non-prefixed version, currently
                                  supported by Chrome, Edge, Opera and Firefox */
}


@keyframes pleasewait {
    0% {
        color: #7c7881;
    }
    50% {
        color: #c0bcc5;
    }
    100% {
        color: #7c7881;
    }
}

@keyframes lsdBtn {
    0% {
        border-color: #eb4d4b;
        background-color: #eb4d4b;
    }
    25% {
        border-color: #9b59b6;
        background-color: #9b59b6;
    }
    50% {
        border-color: #eb4d4b;
        background-color: #eb4d4b;
    }
    75% {
        border-color: #9b59b6;
        background-color: #9b59b6;
    }
    100% {
        border-color: #eb4d4b;
        background-color: #eb4d4b;
    }
}

.more-center {
    display: block;
    margin-left: auto;
    margin-right: auto;
}

.actualFooter {
    display: fixed;
    text-align: center;
    bottom: -5%;
}
</style>
<body style="background-color: #121212;" class="no-select">
    <div class="text-center more-center animate__animated animate__fadeIn" style="margin-top: 6%;">
        <img draggable="false" src="https://discord.monster/api/idigo/idigo.gif" style="display:block;margin-left:auto;margin-right:auto;margin-top:auto;margin-bottom: auto;" width="128" height="128">
        <h1 style="text-align: center;" class="text-light">Error</h1>
        <h3 style="animation: pleasewait 2s infinite;"><?php echo $error; ?></h3>
        <button id="exitButton" class="btn btn-outline-primary text-center more-center" style="width: 35%;margin-top:3%;animation: lsdBtn 6s infinite;color:white;background-color: transparent;box-shadow:none;" onclick="exit()">Exit</button>
    </div>
</body>
</html>

<script>
function openWebsite(url) {
    astilectron.sendMessage(`Client::OpenWebsite|${url}`, function(message) {
        console.log("received " + message)
    });
}

function exit() {
    astilectron.sendMessage(`Client::Exit`)
}
</script>