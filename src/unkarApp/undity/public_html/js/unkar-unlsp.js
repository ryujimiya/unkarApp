
var unlsp = {
	back_flag		: false,
	forward_flag	: false
};

jQuery().live('pageshow', function(){
	// 新しいページが読み込まれた場合
	unlsp.back_flag = true;
	un_initPageAnchor();
});

jQuery().live('pagehide', function(){
	// 今までのページの表示を終了する場合
	unlsp.forward_flag = true;
});

jQuery(document.body).bind('swiperight', function(){
	// 戻る
	if(unlsp.back_flag){
		history.back();
	}
});

jQuery(document.body).bind('swipeleft', function(){
	// 進む
	if(unlsp.forward_flag){
		history.forward();
	}
});

