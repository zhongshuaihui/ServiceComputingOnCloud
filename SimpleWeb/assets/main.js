$(document).ready(function() {
    $.ajax({
        url: "/api/show"
    }).then(function(data) {
        if(data.nickname !== "") {
            $('#message').append(data.message);
        }
    });
});