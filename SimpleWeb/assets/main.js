$(document).ready(function() {
    $.ajax({
        url: "/show"
    }).then(function(data) {
        if(data.nickname !== "") {
            $('#message').append(data.message);
        }
    });
});
