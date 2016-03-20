$(document).ready(function() {
	initEntries();
});

function initEntries() {
	$('div.docular-dir').click(function() {
		var url = $(this).attr('data-url');
		$.post(url, function(data) {
			$("body").html(data);
			initEntries();
		});
	});
	$('div.docular-file').click(function() {
		var url = $(this).attr('data-url');
		window.open(url, '_blank');
	});
	$('div.docular-maff').click(function() {
		var url = $(this).attr('data-url');
		window.open(url, '_blank');
	});
}
