// 縮小ツール：Google Closure Compiler
// 圧縮レベル：ADVANCED_OPTIMIZATIONS

/**
 * @type {string}
 * @const
 */
var LF = String.fromCharCode(10);

/**
 * @type {string}
 * @const
 */
var emptyString = '';

var UA = navigator.userAgent.toLowerCase();
var dummyFunction = function(){};
var _doc = document;
var _href = location.href;
var _hash = location.hash;

var unkarl = {
	unkarl		: '1.8.3',
	appurl		: 'http://' + location.host + '/r/',
	itaurl		: 'http://' + location.host + '/convert.php',
	querystr	: '?name=title',
	timeout		: 8000
};

/**
 * @const
 */
unkarl.id = {
	contents		: 'content',
	popup			: 'popup',
	headerlist_sub	: 'uNav'
};

/**
 * @type {Object.<string, RegExp>}
 * @const
 */
unkarl.regs = {
	http		: /((?=[hst])(?:http(?:s)?|ttp(?:s)?|shttp)):\/\/([\-_.!~*'()\w;\/?:\@&=+\$,%#\|]+)/,
	c2ch 		: /^c\.2ch\.net\/test\/\-\/(\w+)(\/\d{9,10})?/,
	sure 		: /^(?:\w+\.(?=[2b])(?:2ch\.net|bbspink\.com))\/test\/read\.\w+[\/#](\w+\/\d{9,10})(\/[l,\-\d]+)?/,
	ita			: /^(?:\w+\.(?=[2b])(?:2ch\.net|bbspink\.com))\/(\w+)/,
	unkar		: /\/(\w+\/\d{9,10})/,
	tag			: /<[^>]*>/g,
	youtube		: /youtube\.(?:com|jp|co\.jp)\/watch(?:_videos)?\?.*v(?:ideo_ids)?=([\-\w]+)/,
	youtube_min	: /youtu\.be\/([\-\w]+)/,
	nico2_min	: /nico(?:video\.jp\/\w+\/|\.ms)(?:(?:sm|nm|lv|co)?\d+)/,
	nico2		: /^(?:(?:www\.)?nicovideo\.jp\/watch\/|nico\.ms)((?:sm|nm)?\d+)/,
	nico2_live	: /^(?:live\.nicovideo\.jp\/(?:watch|gate)\/|nico\.ms)(lv\d+)/,
	nico2_com	: /^(?:com\.nicovideo\.jp\/community\/|nico\.ms)(co\d+)/,
	res			: /^(?:&gt;(?:&gt;)?|>>?)(\d+)([-,\d]*)/,
	line		: /(\d{1,4})\-(\d{1,4})/,
	dtlist		: /^(\d{1,4})(?:.* ID:([+\/\w!]+))?/,
	id			: /ID:([+\/\w!]+)/,
	num			: /^\d{1,4}$/,
	linkcount	: /&gt;(?:&gt;)?(\d{1,4})/g,
	image_ext	: /\.(?:(?:tif?|gi)f|jp(?:eg?|g)|p(?:ng|sd)|a(?:rt|i)|bmp|ico)$/
};

unkarl.browser = {
	safari		: (UA.indexOf('webkit') !== -1),
	chrome		: (UA.indexOf('chrome') !== -1),
	opera		: (UA.indexOf('opera') !== -1),
	mozilla		: (UA.indexOf('mozilla') !== -1) && !/(compatible|webkit)/.test(UA),
	msie		: (UA.indexOf('msie') !== -1) && (UA.indexOf('opera') === -1),
	version		: (window.opera ?
					(opera.version().replace(/\d$/, emptyString) - 0)
					: parseFloat((/(?:ie |fox\/|ome\/|ion\/)(\d+\.\d)/.exec(UA) || [,0])[1]))
};

unkarl.mouse = {
	x			: 0,
	y			: 0
};

unkarl.idReferenceList = {};
unkarl.resReferenceList = {};
unkarl.nameList = [];
unkarl.nameListDummy = ['dummy', '取得中…', emptyString, emptyString, emptyString];
unkarl.resPopObj = {};

unkarl.now = Date.now || function(){
	return +new Date();
};

unkarl.addEvent = (function(){
	if(_doc.addEventListener){
		return function(elm, type, func){
			elm.addEventListener(type, func, false);
		};
	} else if(_doc.attachEvent){
		return function(elm, type, func){
			elm.attachEvent('on' + type, func);
		};
	} else {
		return function(elm, type, func){
			elm['on' + type] = func;
		};
	}
})();

unkarl.delEvent = (function(){
	if(_doc.removeEventListener){
		return function(elm, type, func){
			elm.removeEventListener(type, func, false);
		};
	} else if(_doc.detachEvent){
		return function(elm, type, func){
			elm.detachEvent('on' + type, func);
		};
	} else {
		return function(elm, type, func){
			elm['on' + type] = null;
		};
	}
})();

unkarl.mousePoint = function(e){
	if(e.touches && (e.touches.length > 0)){
		unkarl.mouse.x = e.touches[0].pageX;
		unkarl.mouse.y = e.touches[0].pageY;
	} else if(unkarl.browser.msie){
		unkarl.mouse.x = event.x + (_doc.body.scrollLeft || _doc.documentElement.scrollLeft);
		unkarl.mouse.y = event.y + (_doc.body.scrollTop || _doc.documentElement.scrollTop);
	} else {
		unkarl.mouse.x = e.pageX;
		unkarl.mouse.y = e.pageY;
	}
};

unkarl.ajax = function(path, that, func){
	var xml;
	var e;
	var timeout;
	var timerID;

	try {
		xml = new XMLHttpRequest();
	} catch(e){
		try {
			xml = new ActiveXObject('Msxml2.XMLHTTP');
		} catch(e){
			try {
				xml = new ActiveXObject('Microsoft.XMLHTTP');
			} catch(e){
				xml = null;
			}
		}
	}
	if(xml){
		timeout = function(){
			unkarl.nameList[path] = undefined;
			xml.abort();
		},
		timerID = setTimeout(timeout, unkarl.timeout);
		xml.onreadystatechange = function(){
			if(xml.readyState === 4){
				clearTimeout(timerID);
				func.call(that, xml, path);
			}
		};
		xml.open('GET', path, true);
		xml.send(emptyString);
	}
};

unkarl.canonicalRes = function(list){
	var clist = [];
	var clistlen = 0;
	var cstr = emptyString;
	var first = 0;	// 連番の開始値
	var second = 0;	// 期待値
	var value = 0;
	var key = 0;
	var tmpstr = emptyString;
	var len;

	list = unkarl.uniq(list);
	list.sort(function(a, b){ return a - b; });
	len = list.length;

	if(len === 0){
		// 空だった場合
		return cstr;
	}
	for(; key < len; key++){
		value = list[key];
		if(first === 0){
			first = value;
			// 期待値を更新
			second = (first + 1);
		} else {
			if(value === second){
				// 期待値通りだったので期待値を更新
				second++;
			} else {
				// 期待値ではない場合
				if((first + 1) === second){
					// 連番が始まる前に崩壊
					tmpstr = first;
				} else {
					// 連番ができた
					tmpstr = first + '-' + (second - 1);
				}
				clist[clistlen++] = tmpstr;
				// 開始値と期待値を更新
				first = value;
				second = (value + 1);
			}
		}
	}
	if(first === value){
		// 開始値を更新してすぐにループを抜けた場合
		tmpstr = value;
	} else {
		tmpstr = first + '-' + value;
	}
	clist[clistlen++] = tmpstr;
	if(clistlen > 0){
		cstr = clist.join(',');
	}
	return cstr;
};

unkarl.uniq = function(array){
	var l = [];
	var llen = 0;
	var tmp = [];
	var it = 0;
	var i = 0;
	var len = array.length;

	for(i = 0; i < len; i++){
		it = array[i];
		if(tmp[it] === undefined){
			l[llen++] = it;
			tmp[it] = 1;
		}
	}
	return l;
};

unkarl.createChangeLink = function(){
	var a = _doc.createElement('a');
	a.href = _href;
	a.className = 'unkarjs-button';
	a.innerHTML = 'うんかーJSモードに移行';
	unkarl.addEvent(a, 'click', function(event){
		if(event.preventDefault){
			event.preventDefault();
		} else {
			event.returnValue = false;
		}
		_doc.cookie = 'unkarjs=1;path=/r;expires=' + (new Date(unkarl.now() + 60 * 60 * 24 * 365 * 1000)).toUTCString();
		location.hash = emptyString;
		location.reload(true);
	});
	return a;
};

unkarl.pageLoad = function(){
	var line = [];
	var href = _href;
	var li = {};
	var tmpelm = {};
	var hash = emptyString;
	var reghash = /^#!?\//;
	var oplist;
	var i;
	var len;

	if(_hash){
		if(reghash.test(_hash)){
			hash = _hash.replace(reghash, emptyString);
			location.replace(unkarl.appurl + hash);
			return;
		}
	}

	unkarl.initPageAnchor();
	// レスポップアップ処理
	unkarl.resPopObj = new unkarl.resPop();

	tmpelm = _doc.getElementById(unkarl.id.headerlist_sub);
	if(tmpelm !== null){
		li = _doc.createElement('li');
		li.appendChild(unkarl.createChangeLink());
		tmpelm.appendChild(li);

		if(_doc.getElementsByClassName){
			// ナビゲーションにも追加
			oplist = Array.prototype.slice.call(_doc.getElementsByClassName('optionlist'));
			len = oplist.length;
			for(i = 0; i < len; i++){
				oplist[i].appendChild(_doc.createTextNode(' '));
				oplist[i].appendChild(unkarl.createChangeLink());
			}
		}
	}
};

unkarl.initPageAnchor = function(){
	var dl = _doc.getElementById(unkarl.id.contents);
	var dt_list_live = {};
	var dd_list_live = {};
	var dt_list = [];
	var dd_list = [];
	var id_list = unkarl.idReferenceList = {};
	var res_list = unkarl.resReferenceList = {};
	var m = [];
	/** @type {number} */
	var length;
	/** @type {number} */
	var i;
	/** @type {string} */
	var text;
	/** @type {number} */
	var num;
	/** @type {string} */
	var id;
	var res_count_func_flag = [];
	var res_count_func = function(str, p1){
		var rl = res_list;
		if(rl[p1] === undefined){
			rl[p1] = [num];
		} else if(res_count_func_flag[p1] !== true){
			// >>1>>1>>1等をカウントしてしまうのを防ぐ
			rl[p1][rl[p1].length] = num;
		}
		res_count_func_flag[p1] = true;
	};

	// nodelist変換
	if(	(unkarl.browser.msie && (unkarl.browser.version < 9)) ||
		(unkarl.browser.safari && (unkarl.browser.version < 4)) ||
		(unkarl.browser.opera && (unkarl.browser.version < 9.5)) ||
		(unkarl.browser.mozilla && (unkarl.browser.version < 3.5))){
		dt_list_live = dl.getElementsByTagName('dt');
		dd_list_live = dl.getElementsByTagName('dd');
		length = dt_list_live.length;
		for(i = 0; i < length; i++){
			dt_list[i] = dt_list_live[i];
			dd_list[i] = dd_list_live[i];
		}
	} else {
		dt_list = Array.prototype.slice.call(dl.getElementsByTagName('dt'));
		dd_list = Array.prototype.slice.call(dl.getElementsByTagName('dd'));
	}

	length = dt_list.length;
	for(i = 0; i < length; i++){
		text = dt_list[i].innerHTML;
		if(m = text.replace(unkarl.regs.tag, emptyString).match(unkarl.regs.dtlist)){
			num = +m[1];
			id = m[2];
			if(id){
				if(id_list[id] === undefined){
					id_list[id] = [num];
				} else {
					id_list[id][id_list[id].length] = num;
				}
			}
			// dtもカウントする
			text.replace(unkarl.regs.linkcount, res_count_func);
			// dtがあったらddも検索
			text = dd_list[i].innerHTML;
			text.replace(unkarl.regs.linkcount, res_count_func);
		}
		res_count_func_flag = [];
	}
};


/**
 * @constructor
 * @param {string} id
 * @param {string=} croot
 * @param {string=} cpop
 * @param {number=} mx
 * @param {number=} my
 */
unkarl.popUp = function(id, croot, cpop, mx, my){
	var tmp;
	var that = this;
	this.level_ = 1;
	this.block_ = emptyString;
	this.rootId_ = id;
	this.regPop = new RegExp('^' + this.rootId_ + '([0-9]+)$');	// privateではない
	this.cPop_ = cpop || emptyString;
	this.mx_ = mx || 0;
	this.my_ = my || 0;
	this.deleteLevel_ = function(e){
		that.deleteNode(e);
	};
	this.deleteAll_ = function(){
		that.removeTree();
	};

	this.rootElem_ = _doc.createElement('div');
	this.rootElem_.id = this.rootId_;
	if(croot){
		this.rootElem_.className = croot;
	}
	// DOMツリーに組み込む
	_doc.body.appendChild(this.rootElem_);
};

unkarl.popUp.prototype.createSpace_ = function(){
	var nich = _doc.getElementById(this.createId_(this.level_));

	if(nich === null){
		nich = _doc.createElement('div');
		nich.id = this.createId_(this.level_);
		nich.className = this.cPop_;
		nich.style.position = 'absolute';
		this.rootElem_.appendChild(nich);
	}
	unkarl.addEvent(nich, 'mouseout', this.deleteLevel_);
	unkarl.addEvent(nich, 'click', this.deleteAll_);
	this.level_++;
	return nich;
};

unkarl.popUp.prototype.createId_ = function(level){
	return this.rootId_ + level;
};

unkarl.popUp.prototype.append = function(key, data){
	var nich;
	if(this.block_ === key){
		this.movePoint_(this.createId_(this.level_ - 1)); // 多重防止
		return false;
	}
	this.block_ = key;
	nich = this.createSpace_();
	nich.innerHTML = data;
	this.setMousePoint_(nich);
	return true;
};

unkarl.popUp.prototype.deleteNode = function(element){
	this.deleteNodeSub(
		element.relatedTarget || element.toElement,
		(element.currentTarget || element.srcElement).id
	);
};

/**
 * @param {Element} target
 * @param {string} currentId
 */
unkarl.popUp.prototype.deleteNodeSub = function(target, currentId){
	var tid;
	var tid_list;
	var tid_num;
	var cid_list;
	var cid_num;

	if(target == null || target.nodeType !== 1){	// null or undefined or nodetype !== 1
		this.removeTree();
		return;
	} else if(target.id != null){
		tid = target.id;
		if(tid == currentId){
			return;
		} else if(tid_list = this.regPop.exec(tid)){
			tid_num = +tid_list[1];
			if(cid_list = this.regPop.exec(currentId)){
				cid_num = +cid_list[1];
				if(cid_num >= tid_num){
					this.cutBranch(tid_num + 1);
				}
			}
			return;
		}
	}
	this.deleteNodeSub(target.parentNode, currentId);
};

unkarl.popUp.prototype.removeTree = function(){
	this.cutBranch(1);
};

unkarl.popUp.prototype.cutBranch = function(i){
	var length = this.level_;
	var root = this.rootElem_;
	var elem;

	this.level_ = i;
	this.block_ = emptyString;
	for(; i < length; i++){
		elem = _doc.getElementById(this.createId_(i));
		unkarl.delEvent(elem, 'mouseout', this.deleteLevel_);
		unkarl.delEvent(elem, 'click', this.deleteAll_);
		root.removeChild(elem);
	}
};

unkarl.popUp.prototype.setMousePoint_ = function(nich){
	var x = 0;
	var y = 0;
	var style = nich.style;

	x = unkarl.mouse.x + this.mx_;
	if(x < 0 || window.innerWidth < 800){
		x = 0;
	}
	y = unkarl.mouse.y - nich.offsetHeight + this.my_;
	if(y < 0){
		y = 0;
	}
	style.zIndex = this.level_;
	style.left = x + 'px';
	style.top = y + 'px';
};

unkarl.popUp.prototype.movePoint_ = function(id){
	this.setMousePoint_(_doc.getElementById(id));
};


/**
 * @constructor
 */
unkarl.resPop = function(){
	this.dataCache_ = [];
	// ポップアップオブジェクト生成
	this.popobj = new unkarl.popUp(
		unkarl.id.popup,
		'popup-root',	//{fontSize: '14px'},
		'popup-branch',	//{backgroundColor: '#FFFFCC', border: 'solid 1px black'},
		-25,
		-5
	);
	this.createEvent_();
};

unkarl.resPop.prototype.createEvent_ = function(){
	// イベント生成
	var that = this;
	var func = this.mouseEventAnalyze_;

	unkarl.addEvent(_doc, 'mouseover', function(e){
		var source = e.target || e.srcElement;

		if(source.tagName.toUpperCase() === 'A'){
			func.call(that, e, source);
		}
	});
};

unkarl.resPop.prototype.mouseEventAnalyze_ = function(e, elm){
	var text = elm.innerHTML;
	var line = [];
	var url = emptyString;
	var path = emptyString;
	var href;
	var id;
	var link;

	unkarl.mousePoint(e);
	if(unkarl.regs.num.test(text)){
		if(unkarl.resReferenceList[text] !== undefined){
			this.resChain_(unkarl.resReferenceList[text].join(','));
		}
	} else if(id = elm.getAttribute('data-id')){
		if(unkarl.idReferenceList[id] !== undefined){
			this.resChain_(unkarl.idReferenceList[id].join(','), id);
		}
	} else if(line = text.match(unkarl.regs.res)){
		if(line[2]){
			this.resChain_(line[1] + line[2]);
		} else {
			this.resAnchorOnce_(line[1]);
		}
	} else if(line = text.match(unkarl.regs.c2ch)){
		if(line[2] !== undefined){
			this.getPageTitle_(line[1] + line[2]);
		} else {
			this.getPageTitle_(line[1]);
		}
	} else if(line = text.match(unkarl.regs.http)){
		// URL系
		path = line[2];
		if(line = path.match(unkarl.regs.sure)){
			this.getPageTitle_(line[1]);
		} else if(line = path.match(unkarl.regs.ita)){
			this.getPageTitle_(line[1]);
		} else if((line = path.match(unkarl.regs.youtube)) || (line = path.match(unkarl.regs.youtube_min))){
			this.popobj.append('youtube' + line[1], '<iframe width="560" height="340" src="https://www.youtube.com/embed/' + line[1] + '?hd=1&rel=0" style="border:none;"></iframe>');
 		} else if(unkarl.regs.nico2_min.test(path)){
			if(line = path.match(unkarl.regs.nico2)){
				url = 'http://ext.nicovideo.jp/thumb/' + line[1];
			} else if(line = path.match(unkarl.regs.nico2_live)){
				url = 'http://live.nicovideo.jp/embed/' + line[1];
			} else if(line = path.match(unkarl.regs.nico2_com)){
				url = 'http://ext.nicovideo.jp/thumb_community/' + line[1];
			} else {
				url = emptyString;
			}
			if(url !== emptyString){
				this.popobj.append('nico2' + line[1], '<iframe width="312" height="176" src="' + url + '" style="border:none;"></iframe>');
			}
		} /*else if(!elm.getAttribute('data-imglink')){
			elm.setAttribute('data-imglink', '1');
			href = elm.href;
			if(unkarl.regs.image_ext.test(href)){
				if(unkarl.browser.msie){
					elm.removeAttribute('target');
					elm.href = 'javascript:(function(){window.open("' + href + '");})();';
				} else {
					if(unkarl.browser.opera){
						link =
							'<html><head></head><body><p style="display:none"><iframe></iframe></p><script type="text/javascript">' + "\n" +
							'var ele = document.getElementsByTagName("iframe")[0];' +
							'ele.contentWindow.document.write(\'<script type="text/javascript">' + "\n" +
								'window.parent.location.href="' + href + '";' + "\n" +
							'</script>\');</script></body></html>';
					} else {
						link = '<html><head><meta http-equiv="Refresh" content="0; url=' + href + '"></head><body></body></html>';
					}
					elm.href = 'data:text/html; charset=utf-8,' + encodeURIComponent(link);
				}
			}
		}
		*/
	}
};

unkarl.resPop.prototype.getPageTitle_ = function(url){
	var path = unkarl.itaurl + '/' + url + unkarl.querystr;
	var list = unkarl.nameList;

	if(list[path] !== undefined){
		this.printPageTitle_(list[path]);
	} else {
		list[path] = unkarl.nameListDummy;
		unkarl.ajax(path, this, this.nameCallBack_);
	}
};

unkarl.resPop.prototype.nameCallBack_ = function(http, path){
	var list = unkarl.nameList;
	var data = [];

	if(http){
		// 1行目がパス、2行目がタイトル、3行目がサーバ、4行目がレス番号、5行目がレス
		data = list[path] = http.responseText.split(LF);
		this.printPageTitle_(data);
	}
};

unkarl.resPop.prototype.printPageTitle_ = function(array){
	this.popobj.append(array[0], array[1]);
};

unkarl.resPop.prototype.lineCache_ = function(num){
	var dt = _doc.getElementById('l' + num);
	var dd = _doc.getElementById('b' + num);

	if(dt === null || dd === null){
		return false;
	}
	this.dataCache_[num] = '<dt class="' + dt.className + '">' + dt.innerHTML + '</dt><dd class="' + dd.className + '">' + dd.innerHTML + '</dd>';
	return true;
};

unkarl.resPop.prototype.resAnchorOnce_ = function(line){
	var num = +line;

	if(this.dataCache_[num] === undefined){
		if(!this.lineCache_(num)) return false;
	}
	this.addRes_('r' + num, '<dl>' + this.dataCache_[num] + '</dl>', [num]);
};

/**
 * @param {string} str
 * @param {string=} id
 */
unkarl.resPop.prototype.resChain_ = function(str, id){
	var list = str.split(',');
	var num = emptyString;
	var resarray = [];
	var resarraylen = 0;
	var len = list.length;
	var line = [];
	var i = 0;
	var max = 0;
	var min = 0;
	var reg = unkarl.regs.num;

	for(; i < len; i++){
		if(reg.test(list[i])){
			max = +list[i];
			if(this.dataCache_[max] === undefined){
				if(!this.lineCache_(max)) continue;
			}
			num += this.dataCache_[max];
			resarray[resarraylen++] = max;
		} else if(line = list[i].match(unkarl.regs.line)){
			min = +line[1];
			max = +line[2];
			for(; min <= max; min++){
				if(this.dataCache_[min] === undefined){
					if(!this.lineCache_(min)) continue;
				}
				num += this.dataCache_[min];
				resarray[resarraylen++] = min;
			}
		}
	}
	if(num !== emptyString){
		this.addRes_('c' + resarray.join('-'), '<dl>' + num + '</dl>', resarray, id);
	}
};

/**
 * @param {string} key
 * @param {string} data
 * @param {Array.<number>} resarray
 * @param {string=} id
 */
unkarl.resPop.prototype.addRes_ = function(key, data, resarray, id){
	var text = emptyString;
	var line = [];
	var href = _href;
	var url = emptyString;
	var tree_url = emptyString;
	var cstr = emptyString;

	if(resarray.length > 0){
		if(line = href.match(unkarl.regs.unkar)){
			url = unkarl.appurl + line[1];
			if(id !== undefined){
				tree_url = url + '/Tree:ID:' + id;
				url += '/ID:' + id;
			} else {
				tree_url = url + '/Tree:' + resarray[0];
				cstr = unkarl.canonicalRes(resarray);
				if(cstr !== emptyString){
					url += '/' + cstr;
				}
			}
			text = '<p>レス抽出(' + resarray.length + '件) [<a href="' + tree_url + '">ツリー抽出</a>] [<a href="' + url + '">Permalink</a>]</p>' + data;
		} else {
			text = data;
		}
	} else {
		text = data;
	}
	this.popobj.append(key, text);
};


unkarl.pageLoad();
