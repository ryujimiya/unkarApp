// うんかー検索
var jQueryObj = window['$'];
var us = {};
us.utils = {};
us.utils.nowTimestamp = Date.now || function(){
	return +new Date();
};

jQueryObj(function(){
	var a = jQueryObj('<a/>')['attr']('href', '/r')['text']('うんかーJSモード')['click'](function(ev){
		if(ev.preventDefault){
			ev.preventDefault();
		} else {
			ev.returnValue = false;
		}
		document.cookie = 'unkarjs=1;path=/r;expires=' + (new Date(us.utils.nowTimestamp() + 60 * 60 * 24 * 365 * 1000)).toUTCString();
		location.href = '/r';
	});
	jQueryObj('#gNav')['append'](jQueryObj('<li/>')['append'](a));
	jQueryObj('#searchform')['submit'](function(ev){
		var val = jQueryObj('#searchtext')['val']();
		if(ev.preventDefault){
			ev.preventDefault();
		} else {
			ev.returnValue = false;
		}
		if(match = val.match(/(\/\w+\/\d{9,10}(?:\/[-,\d]*)?)$/)){
			location.href = '/r' + match[1];
		} else {
			var sf = jQueryObj('#searchform');
			var act = sf['attr']('action');
			function append(id, val){
				var ret = '';
				var v = jQueryObj('#search-select-' + id)['val']();
				if(v !== val){
					ret = '&' + id + '=' + v;
				}
				return ret;
			};
			act += '?q=' + val;
			act += append('board', '');
			act += append('type', 'number');
			act += append('order', 'desc');
			location.href = act;
		}
	});
});

