$(document).ready(function(){

		
		$('#apps').on('click', function(){			
		
			$('html, body').animate({
				scrollTop: $(".app-cards").offset().top-100
			}, 1000);
			
		});
		
		
		$('#features').on('click', function(){			
			
			$('html, body').animate({
				scrollTop: $("#features-screen").offset().top
			}, 1000);
			
		});
				

  });
