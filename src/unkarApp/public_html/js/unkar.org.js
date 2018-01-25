// うんかーTOP
// Closure Compilerでの圧縮を前提にすること！

var _doc = document;
var UA = navigator.userAgent.toLowerCase();
var emptyString = '';
var $utop = {};
var jQueryObj = window['$'];

/**
 * @constructor
 */
var utop = function(){
	this.createTable();
	this.tmpData = [];
	this.viewData = [];
	this.allUrls = [];
	this.allReqs = 0;
	this.allReqsSec = 0;
	this.startTime = 0;
	this.allTimes = 0;
	this.rangeReqs = 0;
	this.rangeReqsSec = 0;
	var ws;
	if('WebSocket' in window){
		ws = WebSocket;
	} else if('MozWebSocket' in window){
		ws = MozWebSocket;
	}
	this.ws = new ws('ws://' + utop.conf.wshost + ':' + utop.conf.wsport + '/' + utop.conf.wspath);
	this.ws.onclose = function(){
		// 接続が閉じられた際は特に何もしない
	};
	this.ws.onmessage = (function(that, func){
		return function(msg){
			func.call(that, msg);
		};
	})(this, this.message);
};
utop.conf = {
	timeRange	: 30,
	timeRangeMs	: 30 * 1000,
	wshost		: 'unkar.org',
	wsport		: '12345',
	wspath		: 'unkartop'
};

utop.id = {
	table		: 'utop-accessview',
	info		: 'utop-accessinfo',
	allrange	: 'utop-allrange',
	range		: 'utop-range',
	startbutton	: 'ws-start-button',
	root		: 'surelist'
};

utop.utils = {};
utop.utils.nowTimestamp = Date.now || function(){
	return +new Date();
};

utop.prototype.message = function(msg){
	var obj = jQueryObj['parseJSON'](msg.data);
	var len = obj.length;
	var list = {
		data: {},
		timestamp: utop.utils.nowTimestamp(),
		count: len
	};
	var url = emptyString;
	var time_range_ms = utop.conf.timeRangeMs;
	var tmpData = this.tmpData;
	var tmpArr = this.viewData = [];
	var tmpArrLen = 0;
	var tmpObj = {};
	var tmpline = {};
	var key = emptyString;
	var time;
	var tmp;
	var i;

	for(i = 0; i < len; i++){
		url = obj[i]['board'] + '/' + obj[i]['thread'];
		if(list.data[url] === undefined){
			list.data[url] = obj[i];
			list.data[url].count = 1;
		} else {
			list.data[url].count++;
		}
	}
	tmpData[tmpData.length] = list;

	this.allReqs += list.count;
	if(this.startTime == 0){
		this.startTime = list.timestamp - 1000;
	}
	this.allTimes = list.timestamp - this.startTime;
	i = this.allTimes / 1000;
	i = Math.round((this.allReqs / i) * 10);
	this.allReqsSec = i / 10;
	i = 0;
	len = tmpData.length;
	while(len > i){
		if((list.timestamp - tmpData[i].timestamp) > time_range_ms){
			i++;
		} else {
			break;
		}
	}
	if(i > 0){
		tmpData = tmpData.slice(i);
	}
	len = tmpData.length;
	for(i = 0; i < len; i++){
		tmp = tmpData[i].data;
		time = tmpData[i].timestamp;
		for(key in tmp){
			tmpline = tmp[key];
			if(tmpObj[key] === undefined){
				// オブジェクトをdeep cloneする
				tmpObj[key] = {
					count: tmpline.count,
					res: tmpline['res'],
					title: tmpline['title'],
					board: tmpline['board'],
					thread: tmpline['thread'],
					bname: tmpline['boardname'],
					lastdate: tmpline['lastdate'],
					addtime: tmpline['addtime'],
					timestamp: time
				};
			} else {
				tmpObj[key].count += tmpline.count;
			}
		}
	}
	this.rangeReqs = 0;
	for(key in tmpObj){
		tmpArr[tmpArrLen++] = tmpObj[key];
		this.rangeReqs += tmpObj[key].count;
	}
	i = Math.round((this.rangeReqs / utop.conf.timeRange) * 10);
	this.rangeReqsSec = i / 10;

	// ソートする
	tmpArr.sort(function(a, b){
		var len = b.count - a.count;
		return ((len !== 0) ? len : (b.timestamp - a.timestamp));
	});
	this.updateTable();
};

utop.prototype.updateTable = function(){
	var arr = this.viewData;
	var sec = utop.conf.timeRange;
	var url = emptyString;
	var t = [];
	var d = new Date();
	jQueryObj('#surelist dt')['each'](function(index){
		if(arr[index] === undefined) return false;
		var it = arr[index];
		url = '/r/' + it.board + '/' + it.thread;
		this['innerHTML'] = '<a href="' + url + '">' + it.title + '</a>';
		return true;
	});
	jQueryObj('#surelist dd')['each'](function(index){
		if(arr[index] === undefined) return false;
		var it = arr[index];
		url = '/r/' + it.board;
		d.setTime((it.thread | 0) * 1000);
		t[0] = d.getFullYear();
		if((t[1] = d.getMonth() + 1) < 10)	t[1] = '0' + t[1];
		if((t[2] = d.getDate()) < 10)		t[2] = '0' + t[2];
		if((t[3] = d.getHours()) < 10)		t[3] = '0' + t[3];
		if((t[4] = d.getMinutes()) < 10)	t[4] = '0' + t[4];
		if((t[5] = d.getSeconds()) < 10)	t[5] = '0' + t[5];
		var since = t[0] + '/' + t[1] + '/' + t[2] + ' ' + t[3] + ':' + t[4] + ':' + t[5];
		this['innerHTML'] = since + '　<a href="' + url + '">' + it.bname + '</a>' + ((it.lastdate > (utop.utils.nowTimestamp() / 1000)) ? '　(dat落ち)' : '');
		return true;
	});
	jQueryObj('#' + utop.id.allrange)['text'](this.allReqs + ' reqs (  ' + this.allReqsSec + '/sec)');
	jQueryObj('#' + utop.id.range)['text'](this.rangeReqs + ' reqs (  ' + this.rangeReqsSec + '/sec)');
};

utop.prototype.createTable = function(){
	var text = [];
	text[text.length] = '<table id="' + utop.id.info + '">';
	text[text.length] = '<tbody>';
	text[text.length] = '<tr><td>このページを開いてからのリクエスト数</td><td id="' + utop.id.allrange + '"></td></tr>';
	text[text.length] = '<tr><td>直近' + utop.conf.timeRange + '秒のリクエスト数</td><td id="' + utop.id.range + '"></td></tr>';
	text[text.length] = '</tbody>';
	text[text.length] = '</table>';

	jQueryObj('#' + utop.id.root)['append'](text.join(emptyString));
};

jQueryObj(function(){
	if(('MozWebSocket' in window) || ('WebSocket' in window)){
		$utop = new utop();
	}
	jQueryObj('.showbutton')['click'](function(){
		var jobj = jQueryObj(this);
		jobj['parent']('.cate-area')['children']('.italist-hide')['slideToggle']('fast');
		jobj['hide']('fast');
	});
	var a = jQueryObj('<a/>')['attr']('href', '/r')['text']('うんかーJSモード')['click'](function(ev){
		if(ev.preventDefault){
			ev.preventDefault();
		} else {
			ev.returnValue = false;
		}
		document.cookie = 'unkarjs=1;path=/r;expires=' + (new Date(utop.utils.nowTimestamp() + 60 * 60 * 24 * 365 * 1000)).toUTCString();
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

