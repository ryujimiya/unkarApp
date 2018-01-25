(function(){

var LF = String.fromCharCode(10),
emptyString = '',
_doc = document,

podunko = window.podunko = function(url){
	$podunko = this;
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

conf = podunko.conf = {
	podunko		: '0.1.0',
	name		: 'PodUnko',
	convert_url	: 'http://unkar.org/convert.php',
	linkurl		: 'http://unkar.org/read.html',
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

mouse = podunko.mouse = {
	x			: 0,
	y			: 0
},

regs = podunko.regs = {
	sure		: /^(\w+\.2ch\.net|\w+\.bbspink\.com)\/test\/read\.\w+[\/#](\w+\/\d{9,10})(\/[l,\-\d]+)?/,
	ita			: /^(\w+\.2ch\.net|\w+\.bbspink\.com)(\/\w+)/,
	id			: / ID:([\w\+\/]+)/,
	id2			: /ID:([\w\+\/]+)/g,
	id3			: /^(.+) ID:([\w\+\/]+)/,
	http		: /(s?h?ttps?):\/\/([\-_.!~*'()\w;\/?:\@&=+\$,%#]+)/g,
	be			: / BE:(\d+)\-(.+)/,
	ank			: /(&gt;(&gt;)?)(\d+)([\-,\d]*)/g,
	ank2		: /^(&gt;(&gt;)?)?(\d+)([\-,\d]*)$/,
	tag			: /<\/?[aA].*?>/g,
	url_split	: /^(\w+)(\/\d{9,10})?/,
	deldat		: /(\d+)\.dat<>(.*\s\((\d+)\))/,
	tsearch		: /^(\w+\.2ch\.net|\w+\.bbspink\.com)\/(\w+)\/(\d+)\<>(.*)/,
	line		: /(\d+)(\-(\d+))?/g,
	sssp		: /(sssp)(\:\/\/img\.2ch\.net\/ico\/[\-_\w\.\/?&]+)/g,
	search		: /([^<]*)(<[^>]*>)?/g
},

img = podunko.img = {
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

id = podunko.id = {
	canvas		: 'nich',
	server		: 'server',
	board		: 'board',
	thread		: 'thread',
	tsearch		: 'tsearch',
	history		: 'history',
	popup		: 'popup',
	outer		: 'outer',
	menu		: function(str){ return str + '-menu'; },
	search		: function(str){ return str + '-search'; },
	searchNext	: function(str){ return str + '-search-next'; },
	searchBack	: function(str){ return str + '-search-back'; },
	searchID	: function(id) { return 'sNo-' + id + '-'; }
},

dom = podunko.dom = {
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
browser = podunko.browser = {
	iphone		: (UA.indexOf('iphone') !== -1),
	ipad		: (UA.indexOf('ipad') !== -1),
	ipod		: (UA.indexOf('ipod') !== -1)
},

addEvent = podunko.addEvent = function(elm, type, func, flag){
	flag = (flag) ? true : false;
	elm.addEventListener(type, func, flag);
},

stopEvent = podunko.stopEvent = function(e){
	e.preventDefault();
	e.stopPropagation();
},

extend = podunko.extend = function(dest, source){
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

addDiv = podunko.addDiv = function(canvas, id, className){
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

uniq = podunko.uniq = function(array){
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

ajax = podunko.ajax = function(path, that, func){
	var xml = new XMLHttpRequest(),
	e = emptyString,
	timeout, lastmod, timerID,
	url = conf.convert_url + '/' + path;
	// 計測開始
	traceLog.start();
	
	if(xml){
		timeout = function(){
			xml.abort();
			dom.uLogErr('接続がタイムアウトしました');
		};
		lastmod = (that.lastModified || 'Mon, 26 Jul 1997 05:00:00 GMT');
		timerID = setTimeout(timeout, conf.timeout);
		xml.onreadystatechange = function(){
			if(xml.readyState === 4){
				if(xml.status === 200){
					clearTimeout(timerID);
					traceLog.load();
					// 最終更新時間の取得
					that.lastModified = xml.getResponseHeader('Last-Modified');
					that[func](xml);
				} else if(xml.status === 304){
					clearTimeout(timerID);
					traceLog.load();
					traceLog.stop();
				}
			}
		};
		xml.open('GET', url, true);
		// 最終更新時間があるならセットする
		xml.setRequestHeader('If-Modified-Since', lastmod);
		xml.send(emptyString);
	}
},

nowTimestamp = function(){
	return +new Date();
},

// 管理オブジェクト
fn = podunko.prototype = {
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
			key = line[1] + (line[2] ? line[2] : emptyString);
			if(this.nich[key]){
				this.nich[key].print();
			} else {
				if(line[2]){
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
		location.hash = '#' + url;
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

	createPopUp: function(){
		this.$popUp = new podunko.popUp(id.popup);
	}
},

fnsaba = fn.saba = function(url){
	this.url = url;
	this.title = '板一覧';
	this.line = [];
	ajax(url, this, 'ajaxCall');
},
saba = fnsaba.prototype = {
	name	: id.server,

	ajaxCall: function(http){
		var i = 0,
		list = http.responseText.split(LF),
		txt = this.line,
		length = list.length - 1,
		line = [],
		board = [],
		l = '<>',
		sla = '/',
		k = 1;

		txt[txt.length] = '<div id="s0" style="display:none;margin-left:5px;">';
		for(i = 0; i < length; i++){
			if((line = list[i].split(l)).length === 2){
				if((board = line[0].split(sla)).length === 2){
					txt[txt.length] = '<img src="' + img.Convert('folder') + '" alt="ロード中" width="16" height="16"><a href="#' + board[1] + '">' + line[1] + '</a><br>';
				}
			} else {
				txt[txt.length] = '</div><img src="' + img.Convert('folder') + '" alt="ロード中" width="16" height="16"><a href="javascript:$podunko.server.itaView(\'s' + k + '\');" class="tan">' + list[i] + '</a><br><div id="s' + k + '" style="display:none;margin-left:5px;">';
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
	ajax(url, this, 'ajaxCall');
},
ita = fnita.prototype = {
	name	: id.board,
	sortflag: 1,

	ajaxCall: function(http){
		var list = http.responseText.split(LF),
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
		var start = ['<ul><li>並び替え：<a href="#' + this.url + '" onclick="$podunko.board.sort(\'num\', ' + this.sortflag + ');" class="tan" >番号</a>：<a href="#' + this.url + '" onclick="$podunko.board.sort(\'res\', ' + this.sortflag + ');" class="tan" >レス数</a>：<a href="#' + this.url + '" onclick="$podunko.board.sort(\'spd\', ' + this.sortflag + ');" class="tan">勢い</a>：<a href="#' + this.url + '" onclick="$podunko.board.sort(\'sin\', ' + this.sortflag + ');" class="tan">日時</a></li>'];
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
	ajax(url, this, 'ajaxCall');
},
sure = fnsure.prototype = {
	name	: id.thread,

	ankerStyle: function(i){
		var color = (this.anker[i] === undefined) ? emptyString : ' onclick="podunko.stopEvent(arguments[0]);$podunko.thread.resPop(' + i + ');" ' + ((this.anker[i].length < 3) ? 'class="ninki"' : 'class="makka"'),
		res = this.res[i];
		return '<dt class="resdate"><a id="l'+i+'" href="javascript:void(0)"'+color+'>'+i+'</a>:<span class="nich"><b>'
				+ res[0] + '</b></span>[' + res[1] + ']<br>' + res[2] + '</dt>'
				+ ((res[3].length > 256) ? '<dd class="mini">' : '<dd>') + res[3] + '</dd>';
	},
	
	anchorPop: function(i){
		if(this.res[i] === undefined) return emptyString;
		$podunko.$popUp.print('a' + i, '<dl>' + this.ankerStyle(i) + '</dl>');
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
		$podunko.$popUp.print('c' + num[0], this.resChain(num));
	},

	resPop: function(i){
		$podunko.$popUp.print('r' + i, this.resTree(i));
	},
	
	idPop: function(id){
		$podunko.$popUp.print(id, this.idTree(id));
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
				return '<a href="javascript:void(0)" onclick="podunko.stopEvent(arguments[0]);$podunko.thread.split(\''+p3+p4+'\');">'+p1+p3+p4+'</a>';
			}
			return '<a href="javascript:void(0)" onclick="podunko.stopEvent(arguments[0]);$podunko.thread.anchorPop('+p3+');">'+p1+p3+'</a>';
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
					line = '/' + regexp[2];
					ret = '<a href="#' + line + '" ontouchstart="podunko.title.name(\'' + line + '\');">' + str + '</a>';
				}
			} else if(regexp = p2.match(reg.ita)){
				// 板だった場合
				if(in_array(regexp[1], conf.filter) !== -1){
					ret = hcheck(p1, str);
				} else {
					line = regexp[2];
					ret = '<a href="#' + line + '" ontouchstart="podunko.title.name(\'' + line + '\');">' + str + '</a>';
				}
			} else {
				ret = hcheck(p1, str);
			}
			return ret;
		},
		idColor = function(str, p2, p1){
			var color = ((id[p1].length >= 5) ? ' class="makka"' : (id[p1].length > 1) ? emptyString : ' class="tan"'),
			text = ' <a href="javascript:void(0)"' + color + ' onclick="podunko.stopEvent(arguments[0]);$podunko.thread.idPop(\'' + p1 + '\');">ID:</a>' + p1;
			text = '日時:' + p2 + text;
			return text;
		},
		line = http.responseText.replace(reg.tag, emptyString)
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
			res[i][3] = res[i][3].replace(reg.id2, '<a href="javascript:void(0)" class="tan" onclick="podunko.stopEvent(arguments[0]);$podunko.thread.idPop(\'$1\');">ID:</a>$1')
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
		ajax(this.url, this, 'ajaxCall');
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
	ajax(this.url, this, 'ajaxCall');
},
tken = fntken.prototype = {
	name	: id.tsearch,

	ajaxCall: function(http){
		var list = http.responseText.split(LF),
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
		len = $podunko.data.board.length;
		for(i = 0; i < len; i++){
			line = $podunko.data.board[i];
			txt[txt.length] = '<li class="rireki' + ((i % 2) ? ' line-color' : emptyString) + '"><a href="#' + line.url + '">' + line.title + '</a></li>';
		}
		txt[txt.length] = '</ul><hr>スレッド履歴<ul>';
		len = $podunko.data.thread.length;
		for(i = 0; i < len; i++){
			line = $podunko.data.thread[i];
			txt[txt.length] = '<li class="rireki' + ((i % 2) ? ' line-color' : emptyString) + '"><a href="#' + line.url + '">' + line.title + '</a></li>';
		}
		txt[txt.length] = '</ul><hr>検索履歴<ul>';
		len = $podunko.data.tsearch.length;
		for(i = 0; i < len; i++){
			line = $podunko.data.tsearch[i];
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


popUp = podunko.popUp = function(id){
	this.rootID = id;
	addDiv(_doc.body, id, podunko.id.popup);
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
		style.left = '1px';
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

style = podunko.style = {
	change: function(nich, css){
		if(nich === null) return false;
		if(css === undefined){
			nich.style.paddingLeft = '2px';
			nich.style.width = 'auto';
			nich.style.backgroundColor = '#FFFFCC';
			nich.style.border = 'solid 1px black';
		} else {
			extend(nich.style, css);
		}
		return true;
	}
},

title = podunko.title = {
	nameList: {},

	name: function(url){
		var path = url + '?name=title',
		list = this.nameList;
		if(list[url] === undefined){
			ajax(path, this, 'ajaxCall');
		} else {
			$podunko.$popUp.print(url, list[url]);
		}
	},
	
	ajaxCall: function(res){
		var list = this.nameList,
		data = [];
		if(res){
			// 1行目がパス、２行目がタイトル
			data = res.responseText.split(LF);
			list[data[0]] = data[1];
			traceLog.stop();
			$podunko.$popUp.print(data[0], data[1]);
		}
	}
},

search = podunko.search = {
	main: function(){
		var obj = $podunko.now_obj,
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
		var obj = $podunko.now_obj,
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
		var obj = $podunko.now_obj,
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
		var obj = $podunko.now_obj,
		element, top, elm;
		if((obj === undefined) || !(obj.searchEnabled)){
			dom.uLogErr('移動失敗。');
			return false;
		}
		element = $(id.searchID(id.canvas) + num),
		top = element.offsetTop,
		elm = _doc.body;
		if(top - 100 > 0){
			top -= 100;
		} else {
			top = 0;
		}
		elm.scrollTop = top;
		obj.searchNow = num;
		return true;
	}
},

traceLog = podunko.traceLog = {
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

menu = podunko.menu = {
	style		: {
		fontSize		: '10px',
		backgroundColor	: '#FFFFFF',
		width			: '150px',
		padding			: '2px',
		border			: 'solid 1px black'
	},

	context: function(){
		var nich = $podunko.$popUp.plus(true),
		obj = $podunko.now_obj;
		style.change(nich, this.style);
		if(obj !== undefined){
			if(obj.name === id.history){
				this.option(nich, '履歴を全て削除', function(){ $podunko.delAll(); });
			}
			nich.appendChild(_doc.createElement('hr'));
			this.option(nich, '単語検索', function(e){
				podunko.stopEvent(e);
				podunko.menu.search();
			});
			if(obj.searchEnabled){
				this.option(nich, '次を検索', function(){ podunko.search.next(); });
				this.option(nich, '前を検索', function(){ podunko.search.back(); });
			}
			nich.appendChild(_doc.createElement('hr'));
		}
		this.option(nich, '板一覧表示', function(){ $podunko.access(id.server); });
		this.option(nich, '履歴', function(){ $podunko.access(id.history); });
		this.option(nich, '更新', function(){ $podunko.renew(); });
		$podunko.$popUp.point(nich);
	},

	search: function(){
		var nich = $podunko.$popUp.plus(true);
		style.change(nich, this.style);
		nich.innerHTML = '単語検索フォーム<form onsubmit="podunko.search.main();return false;"><input type="text" id="' + podunko.id.search(podunko.id.canvas) + '" size="20" onclick="podunko.stopEvent(arguments[0]);"></form>';
		$podunko.$popUp.point(nich);
	},

	tsearch: function(){
		var nich = $podunko.$popUp.plus(true);
		style.change(nich, this.style);
		nich.innerHTML = 'スレッド検索フォーム<form onsubmit="$podunko.access(podunko.id.tsearch + \'/\' + encodeURI(document.getElementById(podunko.id.search(podunko.id.tsearch)).value) + \'/20/1\');return false;"><input type="text" id="' + podunko.id.search(podunko.id.tsearch) + '" size="20" onclick="podunko.stopEvent(arguments[0]);"></form>';
		$podunko.$popUp.point(nich);
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
	top		: 0,
	title	: emptyString,

	scrollkeep: function(){
		var touchy = 0;
		if(podunko.mouse.y > 240){
			touchy = podunko.mouse.y - 240;
		}
		this.top = touchy;
	},

	scroll: function(){
		var element = _doc.body;
		element.scrollTop = this.top;
		podunko.mouse.y = this.top + 240;
	},

	print: function(){
		traceLog.start();
		this.printData(this.text());
	},

	printData: function(txt){
		var title = this.title,
		that = this;
		(function(){
			if(dom.ready){
				var element = $(id.canvas);
				// スクロール状態を保存する
				if($podunko.now_obj){
					$podunko.now_obj.scrollkeep();
				}
				element.innerHTML = txt.join(emptyString);
				window.title = title + ' - ' + conf.name;
				// グローバル変数にオブジェクト格納
				$podunko.now_obj = that;
				$podunko[that.name] = that;
				// スクロールさせる
				that.scroll();
				that.searchEnabled = false;
				// changeを動作させる
				if($podunko.reading > 0){
					$podunko.reading--;
				}
				traceLog.stop();
			} else {
				setTimeout(arguments.callee, 0);
				return;
			}
		})();
	}
}),

// ページロード完了時の処理
pageLoad = podunko.pageLoad = function(){
	dom.stdout = $(id.outer);

	// popup初期化
	$podunko.createPopUp();
	// 構築完了
	dom.ready = true;

	var touchfunc= function(e){
		var mouse = podunko.mouse;
		if(e.touches.length > 0){
			mouse.x = e.touches[0].pageX;
			mouse.y = e.touches[0].pageY;
		}
	},
	menufunc = function(){
		menu.context();
	};

	// イベントハンドラをセットする
	addEvent(_doc, 'touchstart', touchfunc, true);
	addEvent(_doc, 'touchend', touchfunc, true);
	addEvent(_doc, 'gesturestart', menufunc);
	addEvent(_doc, 'contextmenu', menufunc);

	addEvent($('arrowdown'), 'click', function(e){
		var element = _doc.body;
		element.scrollTop = $(id.canvas).offsetHeight;
	});
	addEvent($('home'), 'click', function(e){
		$podunko.access(id.server);
	});
	addEvent($('tsearch'), 'click', function(e){
		menu.tsearch();
	});
	addEvent($('trash'), 'click', function(e){
		$podunko.delAll();
	});

	// ブラウザの戻る進む処理
	setInterval(function(){
		$podunko.change(location.hash.slice(1));
	}, conf.change);
	return true;
},

$podunko = undefined;

$podunko = window.$podunko = new podunko(location.hash.slice(1));
// ページの構築が完了したらloadを呼び出す
addEvent(_doc, 'DOMContentLoaded', pageLoad);

})();
