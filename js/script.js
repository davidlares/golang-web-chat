$( document ).ready(function() {

	var conexion_final;
	var user_name;

	$('#form_registro').on('submit', function(e) {
		user_name = $('#user_name').val()
		e.preventDefault();
		$.ajax({
			type: "post",
			url: "/validate",
			data: {
				'user_name' : user_name
			},
			success: function(data){
				validate_response(data)
			}
		});
	});
	function validate_response(data){
		obj = JSON.parse(data);
		if (obj.valid === true){
			create_conection();
		}else{
			location.reload();
		}
	}

	function create_conection(){
		var conexion = new WebSocket("ws://localhost:8000/ws/" + user_name);
		conexion_final = conexion;
		conexion.onopen = function(){
			conexion.onmessage = function(response){
				console.log(response)
				val = $("#chat_area").val();
		   	$("#chat_area").val(val + "\n" + response.data);
			}
		}
		$("#registro").hide();
   		$("#container_chat").show();
	}

    $('#form_message').on('submit', function(e) {
    	e.preventDefault();
    	conexion_final.send($('#msg').val());
    	$('#msg').val("")
    });
});
