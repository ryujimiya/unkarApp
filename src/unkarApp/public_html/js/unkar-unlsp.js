
var unlsp = {
	back_flag		: false,
	forward_flag	: false
};

jQuery().live('pageshow', function(){
	// �V�����y�[�W���ǂݍ��܂ꂽ�ꍇ
	unlsp.back_flag = true;
	un_initPageAnchor();
});

jQuery().live('pagehide', function(){
	// ���܂ł̃y�[�W�̕\�����I������ꍇ
	unlsp.forward_flag = true;
});

jQuery(document.body).bind('swiperight', function(){
	// �߂�
	if(unlsp.back_flag){
		history.back();
	}
});

jQuery(document.body).bind('swipeleft', function(){
	// �i��
	if(unlsp.forward_flag){
		history.forward();
	}
});

