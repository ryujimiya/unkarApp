(function(){

var LF = String.fromCharCode(10),
CRLF = String.fromCharCode(13) + LF,
emptyString = '',
_doc = document,

daad = window.daad = function(url){
	$daad = this;
	this.nich = {};
	this.now_obj = undefined;
	this.data = {
		server	: {},
		board	: [],
		thread	: [],
		tsearch	: [],
		history	: {}
	};
	this.access(url);
},

conf = daad.conf = {
	daad		: '0.1.0',
	name		: 'daad',
	convert_url	: 'http://www.unkar.org/convert.php',
	linkurl		: 'http://www.unkar.org/read.html',
	kakiko		: 'http://p2.2ch.net/p2/post_form.php',
	timeout		: 15000,
	change		: 64,
	maginX		: 25, // px
	maginY		: 5, // px
	filter		: [
		'www.2ch.net',
		'info.2ch.net',
		'find.2ch.net',
		'm.2ch.net',
		'stats.2ch.net',
		'movie.2ch.net',
		'img.2ch.net'
	]
},

mouse = daad.mouse = {
	x			: 0,
	y			: 0
},

regs = daad.regs = {
	sure		: /^(\w+\.2ch\.net|\w+\.bbspink\.com)\/test\/read\.\w+[\/#](\w+\/\d{9,10})(\/[l,\-\d]+)?/,
	ita			: /^(\w+\.2ch\.net|\w+\.bbspink\.com)(\/\w+)/,
	unkar		: /^(\w+\.2ch\.net|\w+\.bbspink\.com)(\/\w+)(\/\w+)/,
	id			: / ID:([\w\+\/]+)/,
	id2			: /ID:([\w\+\/]+)/g,
	id3			: /^(.+) ID:([\w\+\/]+)/,
	http		: /(s?h?ttps?):\/\/([\-_.!~*'()\w;\/?:\@&=+\$,%#]+)/g,
	be			: / BE:(\d+)\-(.+)/,
	ank			: /(&gt;(&gt;)?)(\d+)([\-,\d]*)/g,
	ank2		: /^(&gt;(&gt;)?)?(\d+)([\-,\d]*)$/,
	tag			: /<\/?[aA].*?>/g,
	url_split	: /^(\w+\.2ch\.net|\w+\.bbspink\.com)(\/\w+)(\/\d{9,10})?/,
	deldat		: /(\d+)\.dat<>(.*\s\((\d+)\))/,
	tsearch		: /^(\w+\.2ch\.net|\w+\.bbspink\.com)\/(\w+)\/(\d+)\<>(.*)/,
	line		: /(\d+)(\-(\d+))?/g,
	sssp		: /(sssp)(\:\/\/img\.2ch\.net\/ico\/[\-_\w\.\/?&]+)/g,
	search		: /([^<]*)(<[^>]*>)?/g
},

img = daad.img = {
	main_url		: 'http://file.unkar.org/img/unkar/',
	Convert			: function(name){ return this.main_url + this[name]; },
//	against			: 'against.gif',
//	arrow_down		: 'arrow-down.gif',
//	arrow_downleft	: 'arrow-downleft.gif',
//	arrow_downright	: 'arrow-downright.gif',
//	arrow_left		: 'arrow-left.gif',
//	arrow_right		: 'arrow-right.gif',
//	arrow_up		: 'arrow-up.gif',
//	arrow_upleft	: 'arrow-upleft.gif',
//	arrow_upright	: 'arrow-upright.gif',
//	back_forth		: 'back-forth.gif',
//	bookmark		: 'bookmark.gif',
//	bulb			: 'bulb.gif',
//	calendar		: 'calendar.gif',
//	calendar2		: 'calendar2.gif',
//	camera			: 'camera.gif',
//	cart			: 'cart.gif',
//	caution			: 'caution.gif',
//	chart			: 'chart.gif',
//	checkmark		: 'checkmark.gif',
//	clipboard		: 'clipboard.gif',
//	clock			: 'clock.gif',
//	closed_folder	: 'closed-folder.gif',
//	database		: 'database.gif',
//	diskette		: 'diskette.gif',
//	document		: 'document.gif',
//	double_arrow	: 'double-arrow.gif',
//	edit			: 'edit.gif',
//	eject			: 'eject.gif',
//	exclaim			: 'exclaim.gif',
//	fastforward		: 'fastforward.gif',
//	favourite		: 'favourite.gif',
//	flag			: 'flag.gif',
//	graph			: 'graph.gif',
//	grow			: 'grow.gif',
//	headphones		: 'headphones.gif',
//	home			: 'home.gif',
//	hourglass		: 'hourglass.gif',
//	info			: 'info.gif',
//	key				: 'key.gif',
//	lock			: 'lock.gif',
//	mail			: 'mail.gif',
//	move			: 'move.gif',
//	music			: 'music.gif',
//	news			: 'news.gif',
//	note			: 'note.gif',
//	open_folder		: 'open-folder.gif',
//	paper_clip		: 'paper-clip.gif',
//	paper_clip2		: 'paper-clip2.gif',
//	pause			: 'pause.gif',
//	phone			: 'phone.gif',
//	play			: 'play.gif',
//	plus			: 'plus.gif',
//	print			: 'print.gif',
//	question_mark	: 'question-mark.gif',
//	quote			: 'quote.gif',
//	refresh			: 'refresh.gif',
//	rewind			: 'rewind.gif',
//	search			: 'search.gif',
//	shield			: 'shield.gif',
//	skip_back		: 'skip-back.gif',
//	skip			: 'skip.gif',
//	skull			: 'skull.gif',
//	statusbar		: 'statusbar.gif',
//	stop			: 'stop.gif',
//	template		: 'template.gif',
//	text_bigger		: 'text-bigger.gif',
//	text_smaller	: 'text-smaller.gif',
//	trash			: 'trash.gif',
//	two_docs		: 'two-docs.gif',
//	twotone			: 'twotone.gif',
//	undo			: 'undo.gif',
//	user			: 'user.gif',
//	vegetable		: 'vegetable.gif',
//	x				: 'x.gif',
//	zoom_in			: 'zoom-in.gif',
//	zoom_out		: 'zoom-out.gif',
	folder			: 'folder.png',
	loading			: 'loading.gif'
},

id = daad.id = {
	canvas		: 'nich',
	server		: 'server',
	board		: 'board',
	thread		: 'thread',
	tsearch		: 'tsearch',
	history		: 'history',
	popup		: 'popup',
	outer		: 'outer',
	ie			: 'ie',
	name		: 'name',
	mail		: 'mail',
	body		: 'body',
	prefs		: {
		url			: 'reboot_url'
	},
	menu		: function(str){ return str + '-menu'; },
	search		: function(str){ return str + '-search'; },
	searchNext	: function(str){ return str + '-search-next'; },
	searchBack	: function(str){ return str + '-search-back'; },
	searchID	: function(id) { return 'sNo-' + id + '-'; }
},

dom = daad.dom = {
	ready: false,
	stdout: {},
	
	uLog: function(str){
		this.stdout.innerHTML = str;
	},

	uLogErr: function(str){
		this.stdout.innerHTML = '<span style="color:#FF0000;">' + str + '</span>';
	}
},

$ = function(id){
	return _doc.getElementById(id);
},

UA = navigator.userAgent.toLowerCase(),
browser = daad.browser = {
	safari		: (UA.indexOf('webkit') !== -1),
	opera		: (UA.indexOf('opera') !== -1),
	msie		: (UA.indexOf('msie') !== -1) && (UA.indexOf('opera') === -1),
	mozilla		: (UA.indexOf('mozilla') !== -1) && !/(compatible|webkit)/.test(UA)
},

VT = gadgets.views.getCurrentView().getName().toLowerCase(),
viewType = daad.viewType = {
	canvas		: (gadgets.views.ViewType.CANVAS.toLowerCase() === VT),
	home		: (gadgets.views.ViewType.HOME.toLowerCase() === VT),
	profile		: (gadgets.views.ViewType.PROFILE.toLowerCase() === VT),
	preview		: (gadgets.views.ViewType.PREVIEW.toLowerCase() === VT)
},

/*
serviceType = daad.serviceType = {
	mixi		: (opensocial.Environment.getDomain().indexOf('mixi') !== -1),
	google		: (opensocial.Environment.getDomain().indexOf('google') !== -1)
},
*/

addEvent = daad.addEvent = (function(){
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
})(),

stopEvent = daad.stopEvent = function(e){
	if(e.stopPropagation){
		e.stopPropagation();
	} else {
		e.cancelBubble = true;
	}
},

extend = daad.extend = function(dest, source){
	if(!(dest instanceof Array)) dest = [dest];
	var i = 0,
	len = dest.length,
	property = emptyString;
	for(; i < len; i++){
		for(property in source){
			dest[i][property] = source[property];
		}
	}
	return dest;
},

addDiv = daad.addDiv = function(canvas, id, className){
	var elm = _doc.createElement('div');
	elm.id = id;
	if(className){
		elm.className = className;
	}
	canvas.appendChild(elm);
},

in_array = function(search, array){
	for(var i = 0, len = array.length; i < len; i++){
		if(search === array[i]){
			return i;
		}
	}
	return -1;
},

uniq = daad.uniq = function(array){
	var list = [],
	tmp = [],
	i = 0,
	length = array.length;
	for(i = 0; i < length; i++){
		if(tmp[array[i]] === undefined){
			list[list.length] = array[i];
			tmp[array[i]] = 1;
		}
	}
	return list;
},

ajax = daad.ajax = function(path, that, func){
	var params = {},
	timeout = (that.timeout ? that.timeout : 60),
	url = conf.convert_url + '/' + path,
	get_url = function(response){
		if(response && response.text){
			traceLog.load();
			that[func](response.text);
		} else {
			dom.uLogErr('接続に失敗しました');
			traceLog.load();
			traceLog.stop();
		}
	};
	// キャッシュ保持時間を指定
	params[gadgets.io.RequestParameters.REFRESH_INTERVAL] = timeout;
	// 計測開始
	traceLog.start();
	gadgets.io.makeRequest(url, get_url, params);
},

get2ch = daad.get2ch = function(path, that, func){
	var params = {},
	headers = {},
	line = [],
	timeout = (that.timeout ? that.timeout : 60),
	url = emptyString,
	get_url = function(response){
		if(!response) return false;
		if(response.rc !== 200){
			traceLog.load();
			traceLog.stop();
			dom.uLogErr('unkarで探します');
			ajax(path, that, func);
		} else {
			traceLog.load();
			if(response.text){
				that[func](response.text);
			} else {
				traceLog.stop();
			}
		}
	};
	// キャッシュ保持時間を指定
	params[gadgets.io.RequestParameters.REFRESH_INTERVAL] = timeout;
	if(line = path.match(regs.unkar)){
		url = 'http://' + line[1] + line[2] + '/dat' + line[3] + '.dat';
		headers['Accept'] = 'text/html';
		headers['Accept-Encoding'] = 'gzip';
		headers['Accept-Language'] = 'ja,en';
		headers['User-Agent'] = 'Monazilla/1.00 (' + conf.name + ')';
		params[gadgets.io.RequestParameters.HEADERS] = headers;
	} else {
		dom.uLogErr('解析できないパス');
		return;
	}
	// 計測開始
	traceLog.start();
	gadgets.io.makeRequest(url, get_url, params);
},

nowTimestamp = function(){
	return +new Date();
},

// 管理オブジェクト
fn = daad.prototype = {
	$popUp		: undefined,
	thread		: undefined,
	board		: undefined,
	server		: undefined,
	tsearch		: undefined,
	history		: undefined,
	reading		: 0,
	change_flag	: 0,

	init: function(){
		this.nich = {};
		this.data = {
			server	: {},
			board	: [],
			thread	: [],
			tsearch	: [],
			history	: {}
		};
		this.thread = undefined;
		this.board = undefined;
		this.server = undefined;
		this.tsearch = undefined;
		this.history = undefined;
		this.reading = 0;
		this.change_flag = 0;
	},

	access	: function(url){
		var data, line, key, obj;
		// changeを止める
		this.reading++;
		// hashを最新に更新しておく
		this.hash(url);
		data = url.split('/');
		if(data[0] === id.tsearch){
			// 検索結果
			obj = this.data.tsearch;
			if(data.length == 4){
				if(this.nich[url]){
					this.nich[url].print();
				} else {
					obj[obj.length] = this.nich[url] = new this.tken(data);
				}
			}
		} else if(data[0] === id.history){
			// 履歴
			url = id.history;
			this.data.history = this.nich[url] = new this.rireki(url);
		} else if(line = url.match(regs.url_split)){
			key = line[1] + line[2] + (line[3] ? line[3] : emptyString);
			if(this.nich[key]){
				this.nich[key].print();
			} else {
				if(line[3]){
					obj = this.data.thread;
					obj[obj.length] = this.nich[key] = new this.sure(key);
				} else {
					obj = this.data.board;
					obj[obj.length] = this.nich[key] = new this.ita(key);
				}
			}
		} else {
			// 板一覧
			url = id.server;
			this.data.server = this.nich[url] = new this.saba(url);
		}
	},

	renew: function(){
		obj = this.now_obj;
		if(obj !== undefined){
			obj.update();
		} else {
			dom.uLogErr('更新失敗。');
		}
	},

	hash: function(url){
		if(this.change_flag) return;
		if(browser.msie || browser.opera){
			var ie = $(id.ie),
			element = undefined;
			if(ie === null) return;
			element = ie.contentDocument || ie.contentWindow.document;
			element.open();
			element.close();
			element.location.hash = '#' + url;
		} else {
			location.hash = '#' + url;
		}
	},

	change: function(url){
		var obj, data;
		if(this.reading > 0) return;
		if(url === emptyString) return;
		obj = this.now_obj;
		if(obj === undefined) return;
		if(url === obj.url) return;
		this.change_flag++;
		this.access(url);
		this.change_flag--;
	},

	delAll: function(){
		$(id.canvas).innerHTML = emptyString;
		this.init();
		this.access(id.server);
	},

	delLeft: function(key, name){
		var list = this.data[name],
		url = emptyString,
		i = 0;
		// 先に表示を変える
		this.hash(list[key].url);
		list[key].print();
		for(; i < key; i++){
			url = list[i].url;
			delete this.nich[url];
		}
		list.splice(0, key);
	},

	delRight: function(key, name){
		var list = this.data[name],
		length = list.length,
		url = emptyString,
		i = key + 1;
		// 先に表示を変える
		this.hash(list[key].url);
		list[key].print();
		for(; i < length; i++){
			url = list[i].url;
			delete this.nich[url];
		}
		list.splice(key + 1, length - key);
	},

	del: function(key, name){
		var obj = this.data[name],
		url = obj[key].url,
		_url = this[name].url;
		obj.splice(key, 1);
		delete this.nich[url];
		if(!obj[0]){
			this[name] = undefined;
			$(id.canvas).innerHTML = emptyString;
			this.hash(emptyString);
		} else {
			if((url === _url) && obj[--key]){
				this.hash(obj[key].url);
				obj[key].print();
			}
		}
	},

	p2Kakiko: function(){
		var line = [],
		url = conf.kakiko,
		obj = this.now_obj;
		if(obj){
			if(line = obj.url.match(regs.url_split)){
				url += '?host=' + line[1] + '&bbs=' + line[2].slice(1);
				if(line[3]){
					url += '&key=' + line[3].slice(1);
				} else {
					url += '&newthread=true';
				}
				window.open(url);
				return true;
			}
		}
		dom.uLogErr('失敗。');
		return false;
	},
	
	createPopUp: function(){
		this.$popUp = new daad.popUp(id.popup);
	}
},

fnsaba = fn.saba = function(url){
	this.url = url;
	this.title = '板一覧';
	this.line = [];
	this.timeout = 3600;
	ajax(url, this, 'ajaxCall');
},
saba = fnsaba.prototype = {
	name	: id.server,

	ajaxCall: function(http){
		var i = 0,
		list = http.split(LF),
		txt = this.line,
		length = list.length - 1,
		line = [],
		l = '<>',
		k = 1;

		txt[txt.length] = '<div id="s0" style="display:none;margin-left:5px;">';
		for(i = 0; i < length; i++){
			if((line = list[i].split(l)).length === 2){
				txt[txt.length] = '<img src="' + img.Convert('folder') + '" alt="ロード中" width="16" height="16"><a href="#' + line[0] + '">' + line[1] + '</a><br>';
			} else {
				txt[txt.length] = '</div><img src="' + img.Convert('folder') + '" alt="ロード中" width="16" height="16"><a href="javascript:$daad.server.itaView(\'s' + k + '\');" class="tan">' + list[i] + '</a><br><div id="s' + k + '" style="display:none;margin-left:5px;">';
				k++;
			}
		}
		this.printData(txt);
	},

	itaView: function(id){
		var e = $(id);
		e.style.display = (e.style.display === 'none') ? 'block' : 'none';
	},

	update: function(){
		ajax(this.url, this, 'ajaxCall');
	},

	text: function(){
		return this.line;
	}
},

fnita = fn.ita = function(url){
	this.url = url;
	this.data = [];
	this.timeout = 10;
	ajax(url, this, 'ajaxCall');
},
ita = fnita.prototype = {
	name	: id.board,
	sortflag: 1,

	ajaxCall: function(http){
		var list = http.split(LF),
		size = list.length - 2,
		txt = [],
		line = [],
		data = this.data,
		t = [],
		i = 0,
		tmp = 0,
		d = new Date(),
		now = (+d) / 1000,
		reg = regs.deldat;

		// 板名取得
		this.title = list[size];

		for(; i < size; i++){
			if(!(line = list[i].match(reg))){
				data[i] = {
					num: i + 1,
					res: 0,
					thread: "スレッドが壊れているみたい",
					sin: 0,
					since: "故障",
					spd: 0
				};
				continue;
			}
			data[i] = {};
			data[i].num = i + 1;
			data[i].res = line[3];
			data[i].thread = line[2];
			data[i].sin = line[1];
			d.setTime(line[1] * 1000);
			t[0] = d.getFullYear();
			if((t[1] = d.getMonth() + 1) < 10)	t[1] = '0' + t[1];
			if((t[2] = d.getDate()) < 10)		t[2] = '0' + t[2];
			if((t[3] = d.getHours()) < 10)		t[3] = '0' + t[3];
			if((t[4] = d.getMinutes()) < 10)	t[4] = '0' + t[4];
			if((t[5] = d.getSeconds()) < 10)	t[5] = '0' + t[5];
			data[i].since = t[0] + '/' + t[1] + '/' + t[2] + ' ' + t[3] + ':' + t[4] + ':' + t[5];
			tmp = (now - data[i].sin) / data[i].res;
			data[i].spd = (86400 / tmp) | 0;
			txt[i] = this.style(i);
		}
		this.printData(this.header(txt));
	},

	header: function(txt){
		var start = ['<ul><li>並び替え：<a href="#' + this.url + '" onclick="$daad.board.sort(\'num\', ' + this.sortflag + ');" class="tan" >番号</a>：<a href="#' + this.url + '" onclick="$daad.board.sort(\'res\', ' + this.sortflag + ');" class="tan" >レス数</a>：<a href="#' + this.url + '" onclick="$daad.board.sort(\'spd\', ' + this.sortflag + ');" class="tan">勢い</a>：<a href="#' + this.url + '" onclick="$daad.board.sort(\'sin\', ' + this.sortflag + ');" class="tan">日時</a></li>'];
		txt = start.concat(txt);
		txt[txt.length] = '</ul>';
		return txt;
	},

	style: function(i){
		return '<li' + ((i % 2) ? ' class="line-color"' : emptyString) + '>' + this.data[i].num + '<a href="#' + this.url + '/' + this.data[i].sin + '">' + this.data[i].thread + '</a><br><span class="resdate">レス:' + this.data[i].res + ' 勢い:' + this.data[i].spd + ' 日時:' + this.data[i].sin + '</span></li>';
	},

	update: function(){
		ajax(this.url, this, 'ajaxCall');
	},

	text: function(){
		var length = this.data.length,
		txt = [],
		i = 0;
		for(; i < length; i++){
			txt[i] = this.style(i);
		}
		return this.header(txt);
	},

	sort: function(target, type){
		traceLog.start();
		var seed = [],
		data = Array.apply(null, this.data),
		length = data.length,
		i = 0,
		txt = [];

		if(data[0][target] === undefined) return false;
		for(i = 0; i < length; i++){
			seed[i] = [data[i][target], data[i]];
		}
		if(type){
			seed.sort(function(a, b){ return b[0] - a[0]; });
			this.sortflag = 0;
		} else {
			seed.sort(function(a, b){ return a[0] - b[0]; });
			this.sortflag = 1;
		}
		for(i = 0; i < length; i++){
			this.data[i] = seed[i][1];
			txt[i] = this.style(i);
		}
		this.printData(this.header(txt));
	}
},

fnsure = fn.sure = function(url){
	this.url = url;
	this.res = [];
	this.anker = [];
	this.id= [];
	this.timeout = 30;
	get2ch(url, this, 'ajaxCall');
},
sure = fnsure.prototype = {
	name	: id.thread,

	ankerStyle: function(i){
		var color = (this.anker[i] === undefined) ? emptyString : ' onmouseover="$daad.thread.resPop(' + i + ');" ' + ((this.anker[i].length < 3) ? 'class="ninki"' : 'class="makka"'),
		res = this.res[i];
		return '<dt class="resdate"><a id="l'+i+'" href="javascript:void(0)"'+color+'>'+i+'</a>:<span class="nich"><b>'
				+ res[0] + '</b></span>[' + res[1] + ']<br>' + res[2] + '</dt>'
				+ ((res[3].length > 256) ? '<dd class="mini">' : '<dd>') + res[3] + '</dd>';
	},
	
	anchorPop: function(i){
		if(this.res[i] === undefined) return emptyString;
		$daad.$popUp.print('a' + i, '<dl>' + this.ankerStyle(i) + '</dl>');
	},
	
	split: function(str){
		var num = [],
		max = 0,
		min = 0;
		str = str.replace(regs.line, function(match, p1, p2, p3){
			if(p2){
				// 巨大な数値の場合に備え4桁に切りそろえる
				min = parseInt(p1.substr(0, 4), 10);
				max = parseInt(p3.substr(0, 4), 10);
				while(min <= max){
					num[num.length] = min++;
				}
			} else {
				num[num.length] = parseInt(p1.substr(0, 4), 10);
			}
			return emptyString;
		});
		if(!num) return false;
		num = uniq(num);
		$daad.$popUp.print('c' + num[0], this.resChain(num));
	},

	resPop: function(i){
		$daad.$popUp.print('r' + i, this.resTree(i));
	},
	
	idPop: function(id){
		$daad.$popUp.print(id, this.idTree(id));
	},

	resChain: function(array){
		if(!(array instanceof Array)) return emptyString;
		var i = 0,
		l = array.length,
		list = [],
		line = 0;
		while(i < l){
			line = array[i++];
			if(this.res[line] === undefined) continue;
			list[list.length] = this.ankerStyle(line);
		}
		if(!list) return emptyString;
		return '抽出数(' + list.length + ')<dl>' + list.join(emptyString) + '</dl>';
	},

	resTree: function(number){
		return this.resChain(this.anker[number]);
	},

	idTree: function(id){
		return this.resChain(this.id[id]);
	},

	ajaxCall: function(http){
		this.res = [];
		this.anker = [];
		this.id= [];

		var txt = [],
		regexp = [],
		res = this.res,
		anker = this.anker,
		id = this.id,
		i = 1,
		split = '<>',
		abone = '壊れています',
		check = [],
		reg = regs,
		ankerNumberColor = function(str, p1, p2, p3, p4){
			if(anker[p3] === undefined){
				anker[p3] = [i];
			} else {
				if(check[p3] !== true){
					// >>1>>1>>1等をカウントしてしまうのを防ぐ
					anker[p3][anker[p3].length] = i;
				}
			}
			check[p3] = true;
			if(p4){
				return '<a href="#l'+p3+'" onmouseover="$daad.thread.split(\''+p3+p4+'\');">'+p1+p3+p4+'</a>';
			}
			return '<a href="#l'+p3+'" onmouseover="$daad.thread.anchorPop('+p3+');">'+p1+p3+'</a>';
		},
		hcheck = function(protocol, url){
			if((protocol.indexOf('h') === -1)){
				// hが付いていなかった場合
				return '<a href="h' + url + '" target="_blank">' + url + '</a>';
			}
			return '<a href="' + url + '" target="_blank">' + url + '</a>';
		},
		urlLink = function(str, p1, p2){
			var line, regexp, ret;
			if(regexp = p2.match(reg.sure)){
				// スレッドだった場合
				if(in_array(regexp[1], conf.filter) !== -1){
					ret = hcheck(p1, str);
				} else {
					line = regexp[1] + '/' + regexp[2];
					ret = '<a href="#' + line + '" onclick="$daad.access(\'' + line + '\');" onmouseover="daad.title.name(\'' + line + '\');">' + str + '</a>';
				}
			} else if(regexp = p2.match(reg.ita)){
				// 板だった場合
				if(in_array(regexp[1], conf.filter) !== -1){
					ret = hcheck(p1, str);
				} else {
					line = regexp[1] + regexp[2];
					ret = '<a href="#' + line + '" onclick="$daad.access(\'' + line + '\');" onmouseover="daad.title.name(\'' + line + '\');">' + str + '</a>';
				}
			} else {
				ret = hcheck(p1, str);
			}
			return ret;
		},
		idColor = function(str, p2, p1){
			var color = ((id[p1].length >= 5) ? ' class="makka"' : (id[p1].length > 1) ? emptyString : ' class="tan"'),
			text = ' <a href="#l' + i + '"' + color + ' onmouseover="$daad.thread.idPop(\'' + p1 + '\');">ID:</a>' + p1;
			text = '日時:' + p2 + text;
			return text;
		},
		line =	http.replace(reg.tag, emptyString)
					.replace(reg.http, urlLink)
					.replace(reg.sssp, '<img src="http$2" alt="2chicon">')
					.split(LF),
		size = line.length;

		for(; i < size; i++){
			if((res[i] = line[i-1].split(split)).length !== 5){
				res[i] = [abone, abone, abone, abone, abone];
			}
			if(regexp = res[i][2].match(reg.id)){
				if(id[regexp[1]] === undefined){
					id[regexp[1]] = [i];
				} else {
					id[regexp[1]][id[regexp[1]].length] = i;
				}
			}
			res[i][0] = res[i][0].replace(reg.ank2, ankerNumberColor);
			res[i][1] = res[i][1].replace(reg.ank2, ankerNumberColor);
			res[i][3] = res[i][3].replace(reg.id2, '<a href="#l' + i + '" class="tan" onmouseover="$daad.thread.idPop(\'$1\');">ID:</a>$1')
								 .replace(reg.ank, ankerNumberColor);
			check = [];
		}
		this.title = res[1][4];
		txt[0] = '<h1>' + this.title + '</h1><dl>';
		for(i = 1; i < size; i++){
			res[i][2] = res[i][2].replace(reg.id3, idColor)
								 .replace(reg.be, ' <a href="http://be.2ch.net/test/p.php?i=$1" target="_blank">?$2</a>');
			txt[i] = this.ankerStyle(i);
		}
		txt[i] = '</dl>';
		this.printData(txt);
	},

	update: function(){
		get2ch(this.url, this, 'ajaxCall');
	},

	text: function(){
		var size = this.res.length,
		txt = [],
		i = 1;
		txt[0] = '<h1>' + this.title + '</h1><dl>';
		for(; i < size; i++){
			txt[i] = this.ankerStyle(i);
		}
		txt[i] = '</dl>';
		return txt;
	}
},

fntken = fn.tken = function(data){
	this.num = parseInt(data[2], 10);
	this.page = parseInt(data[3], 10);
	this.title = data[1] + ' (' + this.page + 'ページ目)';
	this.data = [];
	this.url = data.join('/');
	this.timeout = 3600 * 24;
	ajax(this.url, this, 'ajaxCall');
},
tken = fntken.prototype = {
	name	: id.tsearch,

	ajaxCall: function(http){
		var list = http.split(LF),
		size = list.length - 2,
		txt = [],
		line = [],
		data = this.data,
		t = [],
		i = 0,
		d = new Date(),
		reg = regs.tsearch;

		// ヒット数取得
		line = list[size].split('<>');
		this.hit = parseInt(line[0], 10);
		this.title = line[1];

		for(; i < size; i++){
			line = list[i].match(reg);
			data[i] = {};
			data[i].num = i + 1;
			data[i].thread = line[4];
			data[i].sin = line[3];
			data[i].url = line[1] + '/' + line[2] + '/' + line[3];
			d.setTime(line[3] * 1000);
			t[0] = d.getFullYear();
			if((t[1] = d.getMonth() + 1) < 10)	t[1] = '0' + t[1];
			if((t[2] = d.getDate()) < 10)		t[2] = '0' + t[2];
			if((t[3] = d.getHours()) < 10)		t[3] = '0' + t[3];
			if((t[4] = d.getMinutes()) < 10)	t[4] = '0' + t[4];
			if((t[5] = d.getSeconds()) < 10)	t[5] = '0' + t[5];
			data[i].since = t[0] + '/' + t[1] + '/' + t[2] + ' ' + t[3] + ':' + t[4] + ':' + t[5];
			txt[i] = this.style(i);
		}
		this.printData(this.header(txt));
	},

	header: function(txt){
		var page = this.num * (this.page - 1),
		max = page + this.num,
		start = ['<h1>検索結果</h1>' + this.title + ' に一致するスレッド<br>' + this.hit + '件中 ' + page + '～' + ((max < this.hit) ? max : this.hit) + '件目<br><ul>'];
		txt = start.concat(txt);
		txt[txt.length] = '</ul>';
		return txt;
	},

	style: function(i){
		return '<li' + ((i % 2) ? ' class="line-color"' : emptyString) + '>' + this.data[i].num + '<a href="#' + this.data[i].url + '">' + this.data[i].thread + '</a><br><span class="resdate"> 日時:' + this.data[i].since + '</span></li>';
	},

	update: function(){
		ajax(this.url, this, 'ajaxCall');
	},

	text: function(){
		var length = this.data.length,
		txt = [],
		i = 0;
		for(; i < length; i++){
			txt[i] = this.style(i);
		}
		return this.header(txt);
	}
},

fnrireki = fn.rireki = function(url){
	this.url = url;
	this.title = '履歴';
	this.line = [];
	this.main();
},
rireki = fnrireki.prototype = {
	name	: id.history,

	main: function(){
		traceLog.start();
		var line = {},
		txt = ['<hr>板履歴<ul>'],
		i = 0,
		len = $daad.data.board.length;
		for(i = 0; i < len; i++){
			line = $daad.data.board[i];
			txt[txt.length] = '<li class="rireki' + ((i % 2) ? ' line-color' : emptyString) + '"><a href="#' + line.url + '">' + line.title + '</a></li>';
		}
		txt[txt.length] = '</ul><hr>スレッド履歴<ul>';
		len = $daad.data.thread.length;
		for(i = 0; i < len; i++){
			line = $daad.data.thread[i];
			txt[txt.length] = '<li class="rireki' + ((i % 2) ? ' line-color' : emptyString) + '"><a href="#' + line.url + '">' + line.title + '</a></li>';
		}
		txt[txt.length] = '</ul><hr>検索履歴<ul>';
		len = $daad.data.tsearch.length;
		for(i = 0; i < len; i++){
			line = $daad.data.tsearch[i];
			txt[txt.length] = '<li class="rireki' + ((i % 2) ? ' line-color' : emptyString) + '"><a href="#' + line.url + '">' + line.title + '</a></li>';
		}
		txt[txt.length] = '</ul>';
		this.line = txt;
		this.printData(txt);
	},

	update: function(){
		this.main();
	},

	text: function(){
		return this.line;
	}
},


popUp = daad.popUp = function(id){
	this.rootID = id;
	addDiv(_doc.body, id, daad.id.popup);
	this.root= $(id);
	this.level = 1;
	this.block = emptyString;
},
unpopUp = popUp.prototype = {
	createID: function(id){
		return this.rootID + id;
	},

	idLevel: function(id){
		return id.slice(this.rootID.length);
	},

	plus: function(flag){
		var nich = this.root,
		that = this,
		id = this.createID(this.level);
		if($(id) === null){
			addDiv(nich, id);
		}
		nich = $(id);
		nich.style.display = 'none';
		nich.style.position = 'absolute';
		if(flag){
			addEvent(nich, 'mouseout', function(e){ return that.saku(e); });
			addEvent(nich, 'click', function(){ return that.remove(); });
		}
		this.level++;
		return nich;
	},

	saku: function(element){
		// 移動先
		var t = element.relatedTarget || element.toElement,
		// イベント発生元
		e = element.currentTarget || element.srcElement,
		currentId = this.idLevel(e.id),
		that = this,
		flag = (function(target){
			var tid = emptyString,
			val = 0;
			try {
				if(!target || target.nodeType !== 1){
					return -1;
				}
			} catch(error){
				return 0;
			}
			if(target.id != null){
				tid = that.idLevel(target.id);
				if(tid == currentId){
					return 0;
				} else if((tid != emptyString) && (!isNaN(tid))){
					val = parseInt(tid);
					if(currentId >= val){
						return val + 1;
					}
					return 0;
				}
			}
			target = target.parentNode;
			if(target){
				return arguments.callee(target);
			}
			return -1;
		})(t);
		if(flag === -1){
			this.remove();
		} else if(flag > 0){
			this.cut(flag);
		}
		return flag;
	},

	remove: function(){
		var root = this.root,
		length = root.childNodes.length;
		while(length > 0){
			root.removeChild(root.childNodes[--length]);
		}
		this.level = 1;
		this.block = emptyString;
		return true;
	},

	cut: function(i){
		var length = this.level,
		root = this.root;
		this.level = i;
		this.block = emptyString;
		for(; i < length; i++){
			root.removeChild($(this.createID(i)));
		}
		return true;
	},

	print: function(key, data){
		if(data === emptyString) return false;
		if(this.block === key){
			this.move(this.level - 1); // 多重防止
			return false;
		}
		this.block = key;
		var nich = this.plus(true);
		nich.innerHTML = data;
		style.change(nich);
		this.point(nich);
		return true;
	},

	point: function(nich){
		var x = 0, y = 0, style = nich.style;
		style.zIndex = this.level;
		if((x = mouse.x - conf.maginX) < 0) x = 0;
		if((y = mouse.y - conf.maginY) < 0) y = 0;
		style.left = x + 'px';
		style.top = y + 'px';
		style.display = 'block';
		if((y = mouse.y - nich.offsetHeight - conf.maginY) < 0) y = 0;
		style.display = 'none';
		style.top = y + 'px';
		style.display = 'block';
	},

	move: function(i){
		this.point($(this.createID(i)));
	},

	clean: function(){
		_doc.body.removeChild(this.root);
	}
},

style = daad.style = {
	change: function(nich, css){
		if(nich === null) return false;
		if(css === undefined){
			nich.style.paddingLeft = '2px';
			nich.style.backgroundColor = '#FFFFCC';
			nich.style.border = 'solid 1px black';
		} else {
			extend(nich.style, css);
		}
		return true;
	}
},

title = daad.title = {
	nameList: {},

	name: function(url){
		var path = url + '?name=title',
		list = this.nameList;
		if(list[url] === undefined){
			ajax(path, this, 'ajaxCall');
		} else {
			$daad.$popUp.print(url, list[url]);
		}
	},
	
	ajaxCall: function(res){
		var list = this.nameList,
		data = [];
		if(res){
			// 1行目がパス、２行目がタイトル
			data = res.split(LF);
			list[data[0]] = data[1];
			traceLog.stop();
			$daad.$popUp.print(data[0], data[1]);
		}
	}
},

search = daad.search = {
	main: function(){
		var obj = $daad.now_obj,
		word = $(id.search(id.canvas)).value,
		element, data, reg, sid, i, rep;
		if((obj === undefined) || !word){
			dom.uLogErr('検索失敗。');
			return false;
		}
		element = $(id.canvas);
		data = obj.text().join(emptyString);
		reg = new RegExp(word, 'g');
		sid = id.searchID(id.canvas);
		i = 0;
		rep = function(str){
			return '<span id="' + sid + (i++) + '" class="search">' + str + '</span>';
		};
		data = data.replace(regs.search, function(str, p1, p2){
			return p1.replace(reg, rep) + p2;
		});
		obj.searchMatch = i;
		obj.searchNow = 0;
		element.innerHTML = data;
		if(i === 0){
			obj.searchEnabled = false;
			dom.uLogErr('何も見つかりませんでした。');
		} else {
			obj.searchEnabled = true;
			dom.uLog(i + '件見つかりました。');
			this.move(0, id.canvas);
		}
		return i;
	},
	
	next: function(){
		var obj = $daad.now_obj,
		now, match, num;
		if((obj === undefined) || !(obj.searchEnabled)){
			dom.uLogErr('検索失敗。');
			return false;
		}
		now = obj.searchNow;
		match = obj.searchMatch;
		num = 0;
		if(now !== undefined){
			now++;
			if(now < match){
				num = now;
			} else {
				dom.uLogErr('下まで検索したので上に戻りました。');
			}
		}
		this.move(num, id.canvas);
		return true;
	},
	
	back: function(){
		var obj = $daad.now_obj,
		now, match, num;
		if((obj === undefined) || !(obj.searchEnabled)){
			dom.uLogErr('検索失敗。');
			return false;
		}
		now = obj.searchNow;
		match = obj.searchMatch;
		num = match - 1;
		if(now !== undefined){
			now--;
			if(now >= 0){
				num = now;
			} else {
				dom.uLogErr('上まで検索したので下に戻りました。');
			}
		}
		this.move(num, id.canvas);
		return true;
	},

	move: function(num){
		var obj = $daad.now_obj,
		element, top, elm;
		if((obj === undefined) || !(obj.searchEnabled)){
			dom.uLogErr('移動失敗。');
			return false;
		}
		element = $(id.searchID(id.canvas) + num);
		elm = $(id.canvas);
		elm.scrollTop = element.offsetTop;
		obj.searchNow = num;
		return true;
	}
},

write = daad.write = {
	main: function(path){
		var name = $(daad.id.name).value,
		mail = $(daad.id.name).value,
		body = $(daad.id.name).value;
		if(path.match(regs.unkar)){
			this.thread(path, name, mail, body);
		} else {
			daad.dom.uLogErr('スレ立てには対応してません。');
		}
	},

	thread: function(path, name, mail, body){
		var params = {},
		line = [],
		h = [],
		url = emptyString,
		msg = emptyString;
		if(line = path.match(regs.unkar)){
			url = 'http://' + line[1] + line[2] + '/test/bbs.cgi';
		} else {
			dom.uLogErr('解析できないパス');
			return;
		}
		name = escape(name);
		mail = escape(mail);
		body = escape(body);
		msg = 'submit=%8F%91%82%AB%8D%9E%82%DE&FROM=' + name + '&mail=' + mail + '&MESSAGE=' + body + '&bbs=' + line[2] + '&key=' + line[3] + '&time=1';
		h[0] = 'POST /test/bbs.cgi HTTP/1.1';
		h[1] = 'Host: ' + line[1];
		h[2] = 'User-Agent: Monazilla/1.00 (' + conf.name + ')';
		h[3] = 'Referer: http://' + line[1] + '/test/read.cgi/' + line[2] + '/' + line[3] + '/l50';
		h[4] = 'Content-Type: application/x-www-form-urlencoded';
		h[5] = 'Content-Length: ' + msg.length;
		h[6] = emptyString;
		h[7] = emptyString;
		h[8] = msg;
		params[gadgets.io.RequestParameters.METHOD] = gadgets.io.MethodType.POST;
		params[gadgets.io.RequestParameters.POST_DATA] = h.join(CRLF);
		gadgets.io.makeRequest(url, function(response){
			if(!response) return;
			if(line = response.text.replace(/Set\-Cookie: (.+)\n/)){
				h[6] = 'Cookie: ' + line[1] + '; NAME=""; MAIL=""; suka=pontan';
				params[gadgets.io.RequestParameters.POST_DATA] = h.join(CRLF);
				gadgets.io.makeRequest(url, function(res){
					daad.dom.uLog('書き込み成功');
				}, params);
			} else {
				daad.dom.uLogErr('書き込み失敗');
			}
		}, params);
	}
},

traceLog = daad.traceLog = {
	startTime: 0,
	loadTime: 0,

	start: function(){
		var txt = '<img src="' + img.Convert('loading') + '" alt="ロード中" width="16" height="16">ロード中';
		this.startTime = nowTimestamp();
		if(dom.ready){
			dom.uLog(txt);
		}
	},

	load: function(){
		this.loadTime = nowTimestamp() - this.startTime;
	},

	stop: function(){
		dom.uLog((nowTimestamp() - this.startTime) + 'ms (通信 ' + this.loadTime + 'ms)');
		this.loadTime = 0;
	}
},

menu = daad.menu = {
	style		: {
		fontSize		: '10px',
		backgroundColor	: '#FFFFFF',
		padding			: '2px',
		border			: 'solid 1px black'
	},

	context: function(){
		var nich = $daad.$popUp.plus(true),
		obj = $daad.now_obj;
		style.change(nich, this.style);
		if(obj !== undefined){
			if(obj.name === id.board){
				this.option(nich, 'スレッドを立てる(p2)', function(){ $daad.p2Kakiko(id.board); });
			} else if(obj.name === id.thread){
				this.option(nich, 'スレッドに書き込む(p2)', function(){ $daad.p2Kakiko(id.thread); });
				this.option(nich, 'スレッドに書き込む(google) テスト中', function(e){
					daad.stopEvent(e);
					daad.menu.twrite();
				});
			} else if(obj.name === id.history){
				this.option(nich, '履歴を全て削除', function(){ $daad.delAll(); });
			}
			nich.appendChild(_doc.createElement('hr'));
			this.option(nich, '単語検索', function(e){
				daad.stopEvent(e);
				daad.menu.search();
			});
			if(obj.searchEnabled){
				this.option(nich, '次を検索', function(){ daad.search.next(); });
				this.option(nich, '前を検索', function(){ daad.search.back(); });
			}
			nich.appendChild(_doc.createElement('hr'));
		}
		this.option(nich, '板一覧表示', function(){ $daad.access(id.server); });
		this.option(nich, '履歴', function(){ $daad.access(id.history); });
		this.option(nich, '更新', function(){ $daad.renew(); });
		$daad.$popUp.point(nich);
	},

	search: function(){
		var nich = $daad.$popUp.plus(true);
		style.change(nich, this.style);
		nich.innerHTML = '単語検索フォーム<form onsubmit="daad.search.main();return false;"><input type="text" id="' + daad.id.search(daad.id.canvas) + '" size="20" onclick="daad.stopEvent(arguments[0]);"></form>';
		$daad.$popUp.point(nich);
	},

	tsearch: function(){
		var nich = $daad.$popUp.plus(true);
		style.change(nich, this.style);
		nich.innerHTML = 'スレッド検索フォーム<form onsubmit="$daad.access(daad.id.tsearch + \'/\' + encodeURI(document.getElementById(daad.id.search(daad.id.tsearch)).value) + \'/20/1\');return false;"><input type="text" id="' + daad.id.search(daad.id.tsearch) + '" size="20" onclick="daad.stopEvent(arguments[0]);"></form>';
		$daad.$popUp.point(nich);
	},

	twrite: function(){
		var nich = $daad.$popUp.plus(true),
		text = [];
		style.change(nich, this.style);
		text[text.length] = '書き込みフォーム';
		text[text.length] = '<form onsubmit="daad.write.main(\'' + $daad.now_obj.url + '\');return false;">';
		text[text.length] = '名前：<input type="text" id="' + daad.id.name + '" size="20" onclick="daad.stopEvent(arguments[0]);">';
		text[text.length] = 'メール：<input type="text" id="' + daad.id.mail + '" size="20" onclick="daad.stopEvent(arguments[0]);">';
		text[text.length] = '<textarea rows="4" cols="40" id="' + daad.id.body + '" onclick="daad.stopEvent(arguments[0]);"></textarea>';
		text[text.length] = '<input type="submit" onclick="daad.stopEvent(arguments[0]);">';
		text[text.length] = '</form>';
		nich.innerHTML = text.join('<br>');
		$daad.$popUp.point(nich);
	},

	option: function(nich, title, click){
		var elm = _doc.createElement('span');
		addEvent(elm, 'click', click);
		elm.innerHTML = title;
		elm.style.paddingLeft = '5px';
		nich.appendChild(elm);
		nich.appendChild(_doc.createElement('br'));
	}
},

// 共通のメソッド
fncanvas = extend([saba, ita, sure, tken, rireki], {
	x		: 0,
	y		: 0,
	title	: emptyString,

	scrollkeep: function(){
		var element = $(id.canvas);
		this.x = element.scrollLeft;
		this.y = element.scrollTop;
	},
	
	scroll: function(){
		var element = $(id.canvas);
		element.scrollLeft = this.x;
		element.scrollTop = this.y;
	},

	print: function(){
		traceLog.start();
		this.printData(this.text());
	},

	printData: function(txt){
		var title = this.title,
		url = this.url,
		that = this;
		(function(){
			if(dom.ready){
				var element = $(id.canvas);
				// スクロール状態を保存する
				if($daad.now_obj){
					$daad.now_obj.scrollkeep();
				}
				element.innerHTML = txt.join(emptyString);
				gadgets.window.setTitle(title + ' - ' + conf.name);
				// グローバル変数にオブジェクト格納
				$daad.now_obj = that;
				$daad[that.name] = that;
				// スクロールさせる
				that.scroll();
				that.searchEnabled = false;
				// changeを動作させる
				if($daad.reading > 0){
					$daad.reading--;
				}
				prefs.set(daad.id.prefs.url, url);
				traceLog.stop();
			} else {
				setTimeout(arguments.callee, 0);
				return;
			}
		})();
	}
}),

home = daad.home = {
	back: function(){
		window.history.back();
	},

	forward: function(){
		window.history.forward();
	},

	update: function(){
		if($daad.now_obj){
			$daad.now_obj.update();
		} else {
			dom.uLogErr("更新失敗");
		}
	},

	home: function(){
		$daad.access(id.server);
	},

	history: function(){
		$daad.access(id.history);
	}
},

// ページロード完了時の処理
pageLoad = daad.pageLoad = function(){
	dom.stdout = $(id.outer);

	// popup初期化
	$daad.createPopUp();
	// 構築完了
	dom.ready = true;

	var mousefunc = (function(){
		var mouse = daad.mouse;
		if(browser.msie){
			return function(){
				mouse.x = event.x + (_doc.body.scrollLeft || _doc.documentElement.scrollLeft);
				mouse.y = event.y + (_doc.body.scrollTop || _doc.documentElement.scrollTop);
			}
		} else {
			return function(e){
				mouse.x = e.pageX;
				mouse.y = e.pageY;
			}
		}
	})(),
	contextevent = function(){
		menu.context();
	};

	// イベントハンドラをセットする
	addEvent(_doc, 'mousemove', mousefunc);

	// 右クリックメニュー
	addEvent($(id.canvas), 'contextmenu', contextevent);

	addEvent($('back'), 'click', function(){ home.back(); });
	addEvent($('forward'), 'click', function(){ home.forward(); });
	addEvent($('update'), 'click', function(){ home.update(); });
	addEvent($('home'), 'click', function(){ home.home(); });
	addEvent($('history'), 'click', function(){ home.history(); });

	// ブラウザの戻る進む処理
	setInterval((function(){
		if(browser.msie || browser.opera){
			var elm = _doc.createElement('iframe');
			elm.id = id.ie;
			elm.style.display = 'none';
			_doc.body.appendChild(elm);
			return function(){
				var ie = elm.contentDocument || elm.contentWindow.document;
				$daad.change(ie.location.hash.slice(1));
			};
		} else {
			return function(){
				$daad.change(location.hash.slice(1));
			};
		}
	})(), conf.change);
	return true;
},

// オブジェクト生成
prefs = daad.prefs = new gadgets.Prefs(),

$daad = window.$daad = (function(){
	var obj = {};
	if(viewType.home || viewType.profile){
		obj = new daad(prefs.getString(daad.id.prefs.url));
		// ページの構築が完了したらloadを呼び出す
		gadgets.util.registerOnLoadHandler(daad.pageLoad);
	}
	return obj;
})();

})();
