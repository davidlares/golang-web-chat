$( document ).ready(function() {

	var conn;
	var username;

	$('#form').on('submit', function(e) {
		username = $('#username').val()
		e.preventDefault();
		$.ajax({
			type: "post",
			url: "/validate",
			data: {
				'user_name' : username
			},
			success: function(data){
				validate_response(data)
			}
		});
	});
	function validate_response(data){
		obj = JSON.parse(data);
		if (obj.isvalid === true){
			create_conection();
		}else{
			location.reload();
		}
	}

	function create_conection(){
		var connection = new WebSocket("ws://localhost:8000/chat/" + username);
		conn = connection;
		connection.onopen = function(){
			connection.onmessage = function(response){
				console.log(response)
				val = $("#area").val();
		   	$("#area").val(val + "\n" + response.data);
			}
		}
		$("#registro").hide();
   		$("#container_chat").show();
	}

    $('#form_message').on('submit', function(e) {
    	e.preventDefault();
    	conn.send($('#msg').val());
    	$('#msg').val("")
    });
});
