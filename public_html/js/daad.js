<?xml version="1.0" encoding="UTF-8"?>
<Module>
<ModulePrefs title="daad"
			 title_url="http://www.unkar.org/"
			 description="２ちゃんねるビューアです。板やスレッドを閲覧することができます。"
			 width="320"
			 height="420"
			 screenshot="http://file.unkar.org/2ch/img/daad.png"
			 thumbnail="http://file.unkar.org/2ch/img/daad-thumb.png"
			 author="tanaton"
			 author_email="heiwaboke+igoogle@gmail.com"
			 author_link="http://www.unkar.org/"
			 author_location="Japan">
<Require feature="opensocial-0.8" />
<Require feature="views" />
<Require feature="settitle" />
<Require feature="setprefs" />
</ModulePrefs>
<UserPref name="reboot_url" default_value="server" datatype="hidden" />
<Content type="html" view="profile,home,canvas"><![CDATA[
<script type="text/javascript">
(function(){

var window = this,
LF = String.fromCharCode(10),
undefined,
emptyString = '',
_doc = document,

daad = window.daad = function(url){
	this.nich = {};
	this.now_obj = undefined;
	this.data = {
		server	: {},
		board	: [],
		thread	: []
	};
	this.access(url);
},

conf = daad.conf = {
	daad		: '1.0.0',
	name		: 'daad',
	convert_url	: 'http://www.unkar.org/convert.php',
	linkurl		: 'http://www.unkar.org/read.html',
	kakiko		: 'http://p2.2ch.net/p2/post_form.php',
	timeout		: 8000,
	change		: 128,
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
	line		: /(\d+)(\-(\d+))?/g,
	sssp		: /(sssp)(\:\/\/img\.2ch\.net\/ico\/[\-_\w\.\/?&]+)/g,
	search		: /([^<]*)(<[^>]*>)?/g
},

img = unkar.img = {
	main_url		: 'http://file.unkar.org/img/unkar/',
	Convert			: function(name){ return this.main_url + this[name]; },
	against			: 'against.gif',
	arrow_down		: 'arrow-down.gif',
	arrow_downleft	: 'arrow-downleft.gif',
	arrow_downright	: 'arrow-downright.gif',
	arrow_left		: 'arrow-left.gif',
	arrow_right		: 'arrow-right.gif',
	arrow_up		: 'arrow-up.gif',
	arrow_upleft	: 'arrow-upleft.gif',
	arrow_upright	: 'arrow-upright.gif',
	back_forth		: 'back-forth.gif',
	bookmark		: 'bookmark.gif',
	bulb			: 'bulb.gif',
	calendar		: 'calendar.gif',
	calendar2		: 'calendar2.gif',
	camera			: 'camera.gif',
	cart			: 'cart.gif',
	caution			: 'caution.gif',
	chart			: 'chart.gif',
	checkmark		: 'checkmark.gif',
	clipboard		: 'clipboard.gif',
	clock			: 'clock.gif',
	closed_folder	: 'closed-folder.gif',
	database		: 'database.gif',
	diskette		: 'diskette.gif',
	document		: 'document.gif',
	double_arrow	: 'double-arrow.gif',
	edit			: 'edit.gif',
	eject			: 'eject.gif',
	exclaim			: 'exclaim.gif',
	fastforward		: 'fastforward.gif',
	favourite		: 'favourite.gif',
	flag			: 'flag.gif',
	folder			: 'folder.png',
	graph			: 'graph.gif',
	grow			: 'grow.gif',
	headphones		: 'headphones.gif',
	home			: 'home.gif',
	hourglass		: 'hourglass.gif',
	info			: 'info.gif',
	key				: 'key.gif',
	loading			: 'loading.gif',
	lock			: 'lock.gif',
	mail			: 'mail.gif',
	move			: 'move.gif',
	music			: 'music.gif',
	news			: 'news.gif',
	note			: 'note.gif',
	open_folder		: 'open-folder.gif',
	paper_clip		: 'paper-clip.gif',
	paper_clip2		: 'paper-clip2.gif',
	pause			: 'pause.gif',
	phone			: 'phone.gif',
	play			: 'play.gif',
	plus			: 'plus.gif',
	print			: 'print.gif',
	question_mark	: 'question-mark.gif',
	quote			: 'quote.gif',
	refresh			: 'refresh.gif',
	rewind			: 'rewind.gif',
	search			: 'search.gif',
	shield			: 'shield.gif',
	skip_back		: 'skip-back.gif',
	skip			: 'skip.gif',
	skull			: 'skull.gif',
	statusbar		: 'statusbar.gif',
	stop			: 'stop.gif',
	template		: 'template.gif',
	text_bigger		: 'text-bigger.gif',
	text_smaller	: 'text-smaller.gif',
	trash			: 'trash.gif',
	two_docs		: 'two-docs.gif',
	twotone			: 'twotone.gif',
	undo			: 'undo.gif',
	user			: 'user.gif',
	vegetable		: 'vegetable.gif',
	x				: 'x.gif',
	zoom_in			: 'zoom-in.gif',
	zoom_out		: 'zoom-out.gif'
},

id = daad.id = {
	canvas		: 'nich',
	server		: 'server',
	board		: 'board',
	thread		: 'thread',
	popup		: 'popup',
	outer		: 'outer',
	ie			: 'ie',
	prefs		: {
		url			: 'reboot_url'
	},
	tab			: function(str){ return str + '-tab'; },
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

addDiv = daad.addDiv = function(canvas, id, klass){
	var elm = _doc.createElement('div');
	elm.id = id;
	if(klass){
		elm.className = klass;
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
	tmp = emptyString,
	i = 0,
	length = array.length;
	array.sort(function(a, b){ return a - b; });
	while(i < length){
		list[list.length] = tmp = array[i];
		while(tmp === array[++i]);
	}
	return list;
},

ajax = daad.ajax = function(path, self, func){
	var get_url,
	params = {},
	timeout = (self.timeout ? self.timeout : 60),
	url = conf.convert_url + '/' + path;
	// 計測開始
	traceLog.start();
	get_url = function(response){
		if(response && response.text){
			traceLog.load();
			self[func](response.text);
		} else {
			dom.uLogErr('接続に失敗しました');
			traceLog.load();
			traceLog.stop();
		}
	};
	// キャッシュ保持時間を指定
	params[gadgets.io.RequestParameters.REFRESH_INTERVAL] = timeout;
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
	reading		: 0,
	change_flag	: 0,

	init: function(){
		this.nich = {};
		this.data = {
			server	: {},
			board	: [],
			thread	: []
		};
		this.thread = undefined;
		this.board = undefined;
		this.server = undefined;
		this.reading = 0;
		this.change_flag = 0;
	},

	access	: function(url){
		var line, key, obj;
		// changeを止める
		this.reading++;
		// hashを最新に更新しておく
		this.hash(url);
		if(line = url.match(regs.url_split)){
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
			this.data.server = new this.saba(id.server);
		}
	},

	renew: function(id){
		if(this[id] !== undefined){
			this[id].update();
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
		var obj, line;
		if(this.reading > 0) return;
		if(line = url.match(regs.url_split)){
			obj = (!viewType.canvas) ? this.now_obj : ((line[3]) ? this.thread : this.board);
			if(obj === undefined) return;
			if(url !== obj.url){
				this.change_flag++;
				this.access(url);
				this.change_flag--;
			}
		}
	},

	p2Kakiko: function(name){
		var line = [],
		url = conf.kakiko,
		obj = this[name];
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
	this.timeout = 3600;
	this.title = '板一覧';
	this.line = [];
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

		txt[txt.length] = '<div id="s0" class="saba">';
		for(i = 0; i < length; i++){
			if((line = list[i].split(l)).length === 2){
				txt[txt.length] = '<img src="' + img.Convert('folder') + '" alt="ロード中" width="16" height="16"><a href="#' + line[0] + '" onclick="$daad.access(\'' + line[0] + '\');">' + line[1] + '</a><br>';
			} else {
				txt[txt.length] = '</div><img src="' + img.Convert('folder') + '" alt="ロード中" width="16" height="16"><a href="javascript:$daad.server.itaView(\'s' + k + '\');" class="tan">' + list[i] + '</a><br><div id="s' + k + '" class="saba">';
				k++;
			}
		}
		this.printData(txt);
	},

	itaView: function(id){
		var e = $(id);
		e.style.display = (e.style.display === 'none') ? 'block' : 'none';
	},

	text: function(){
		return this.line;
	}
},

fnita = fn.ita = function(url){
	this.url = url;
	this.timeout = 10;
	this.data = [];
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
		var start = ['<table><tbody><tr><th><a href="#' + this.url + '" onclick="$daad.board.sort(\'num\', ' + this.sortflag + ');" class="tan" >No</a></th><th>title</th><th><a href="#' + this.url + '" onclick="$daad.board.sort(\'res\', ' + this.sortflag + ');" class="tan" >res</a></th><th><a href="#' + this.url + '" onclick="$daad.board.sort(\'spd\', ' + this.sortflag + ');" class="tan">res/day</a></th>'];
		if(viewType.canvas){
			start[0] += '<th><a href="#' + this.url + '" onclick="$daad.board.sort(\'sin\', ' + this.sortflag + ');" class="tan">since</a></th></tr>';
		} else {
			start[0] += '</tr>';
		}
		txt = start.concat(txt);
		txt[txt.length] = '</tbody></table>';
		return txt;
	},

	style: function(i){
		var txt = '<tr' + ((i % 2) ? ' class="line-color"' : emptyString) + '><td>' + this.data[i].num + '</td><td><a href="#' + this.url + '/' + this.data[i].sin + '" onclick="$daad.access(\'' + this.url + '/' + this.data[i].sin + '\');">' + this.data[i].thread + '</a></td><td class="res">' + this.data[i].res + '</td><td class="spd">' + this.data[i].spd + '</td>';
		if(viewType.canvas){
			txt += '<td class="sin">' + this.data[i].since + '</td></tr>';
		} else {
			txt += '</tr>';
		}
		return txt;
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
	this.timeout = 10;
	this.res = [];
	this.anker = [];
	this.id= [];
	this.searchMatch = 0;
	ajax(url, this, 'ajaxCall');
},
sure = fnsure.prototype = {
	name	: id.thread,

	ankerStyle: function(i){
		var color = (this.anker[i] === undefined) ? emptyString : ' onmouseover="$daad.thread.resPop(' + i + ');" ' + ((this.anker[i].length < 3) ? 'class="ninki"' : 'class="makka"'),
		res = this.res[i];
		return '<dt><a id="l'+i+'" href="#l'+i+'"'+color+'>'+i+'</a> ：<span class="nich"><b>'
				+ res[0] + '</b></span>[' + res[1] + ']：' + res[2] + '</dt><dd>' + res[3] + '</dd>';
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
					ret = '<a href="#' + line + '" onclick="$unkar.access(\'' + line + '\');" onmouseover="unkar.title.name(\'' + line + '\');">' + str + '</a>';
				}
			} else if(regexp = p2.match(reg.ita)){
				// 板だった場合
				if(in_array(regexp[1], conf.filter) !== -1){
					ret = hcheck(p1, str);
				} else {
					line = regexp[1] + regexp[2];
					ret = '<a href="#' + line + '" onclick="$unkar.access(\'' + line + '\');" onmouseover="unkar.title.name(\'' + line + '\');">' + str + '</a>';
				}
			} else {
				ret = hcheck(p1, str);
			}
			return ret;
		},
		idColor = function(str, p2, p1){
			var color = ((id[p1].length >= 5) ? ' class="makka"' : (id[p1].length > 1) ? emptyString : ' class="tan"'),
			text = ' <a href="#l' + i + '"' + color + ' onmouseover="$daad.thread.idPop(\'' + p1 + '\');">ID:</a>' + p1;
			if(viewType.canvas){
				text = p2 + text;
			} else {
				text = '<a href="#l' + i + '" class="tan" onmouseover="$daad.$popUp.print(\'date' + i + '\', \'' + p2 + '\')">date</a>' + text;
			}
			return text;
		},
		line = http.replace(reg.tag, emptyString)
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

	text: function(flag){
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


tab = daad.tab = {
	print: function(name){
		var list = $daad.data[name],
		length = list.length,
		url = emptyString,
		_url = $daad[name].url,
		title = emptyString,
		txt = ['<ul>'],
		obj = {},
		tmp = emptyString,
		i = 0;
		for(; i < length; i++){
			obj = list[i];
			url = obj.url || emptyString;
			title = obj.title || emptyString;
			if(title.length > 10){
				title = title.slice(0, 10) + '...';
			}
			tmp = '<li ondblclick="daad.menu.tab(' + i + ', \'' + name + '\');" oncontextmenu="daad.menu.tab(' + i + ', \'' + name + '\', event);"';
			if(url === _url){
				txt[txt.length] = tmp + ' class="now"><a href="#' + url + '" title="' + obj.title + '">' + title + '</a>';
			} else {
				txt[txt.length] = tmp + '><a href="#' + url + '" title="' + obj.title + '" onclick="$daad.access(\'' + url + '\');">' + title + '</a>';
			}
		}
		$(id.tab(name)).innerHTML = txt.join(emptyString) + '</ul>';
		return length;
	},

	delAll: function(){
		$(id.tab(id.board)).innerHTML = emptyString;
		$(id.tab(id.thread)).innerHTML = emptyString;
		$(id.board).innerHTML = emptyString;
		$(id.thread).innerHTML = emptyString;
		$daad.init();
		$daad.access(id.server);
	},

	delLeft: function(key, name){
		var list = $daad.data[name],
		url = emptyString,
		i = 0;
		// 先に表示を変える
		$daad.hash(list[key].url);
		list[key].print();
		for(; i < key; i++){
			url = list[i].url;
			delete $daad.nich[url];
		}
		list.splice(0, key);
		this.print(name);
	},

	delRight: function(key, name){
		var list = $daad.data[name],
		length = list.length,
		url = emptyString,
		i = key + 1;
		// 先に表示を変える
		$daad.hash(list[key].url);
		list[key].print();
		for(; i < length; i++){
			url = list[i].url;
			delete $daad.nich[url];
		}
		list.splice(key + 1, length - key);
		this.print(name);
	},

	del: function(key, name){
		var obj = $daad.data[name],
		url = obj[key].url,
		_url = $daad[name].url;
		obj.splice(key, 1);
		delete $daad.nich[url];
		this.print(name);
		if(!obj[0]){
			$daad[name] = undefined;
			$(name).innerHTML = emptyString;
			$daad.hash(emptyString);
		} else {
			if((url === _url) && obj[--key]){
				$daad.hash(obj[key].url);
				obj[key].print();
			}
		}
	}
},

popUp = daad.popUp = function(id){
	this.rootID = id;
	if(viewType.canvas){
		addDiv(_doc.body, id);
	} else {
		addDiv(_doc.body, id, 'popup');
	}
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
		self = this,
		id = this.createID(this.level);
		if($(id) === null){
			addDiv(nich, id);
		}
		nich = $(id);
		nich.style.display = 'none';
		nich.style.position = 'absolute';
		if(flag){
			addEvent(nich, 'mouseout', function(e){ return self.saku(e); });
			addEvent(nich, 'click', function(){ return self.remove(); });
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
		self = this,
		flag = (function(target){
			var tid = emptyString,
			val = 0;
			try {
				if(!target || target.nodeType !== 1)
					return -1;
			} catch(error){
				return 0;
			}
			if(target.id != null){
				tid = self.idLevel(target.id);
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

menu = daad.menu = {
	$popUp		: undefined,
	style		: {
		fontSize		: '12px',
		backgroundColor	: '#FFFFFF',
		width			: '200px',
		padding			: '2px',
		border			: 'solid 1px black'
	},

	context: function(){
		var nich = this.$popUp.plus(true);
		style.change(nich, this.style);
		this.option(nich, 'タブを全て削除', function(){ daad.tab.delAll(); });
		nich.appendChild(_doc.createElement('hr'));
		this.option(nich, 'スレッドを立てる(p2)', function(){ $daad.p2Kakiko(id.board); });
		this.option(nich, 'スレッドに書き込む(p2)', function(){ $daad.p2Kakiko(id.thread); });
		this.option(nich, '板一覧表示', function(){ $daad.access(daad.id.server); });
		nich.appendChild(_doc.createElement('hr'));
		this.option(nich, '板一覧更新', function(){ $daad.renew(daad.id.server); });
		this.option(nich, 'スレッド一覧更新', function(){ $daad.renew(daad.id.board); });
		this.option(nich, 'スレッド更新', function(){ $daad.renew(daad.id.thread); });
		this.$popUp.point(nich);
	},

	context2: function(){
		var nich = this.$popUp.plus(true);
		style.change(nich, this.style);
		if($daad.now_obj){
			if($daad.now_obj.name === id.board){
				this.option(nich, 'スレッドを立てる(p2)', function(){ $daad.p2Kakiko(id.board); });
			} else if($daad.now_obj.name === id.thread){
				this.option(nich, 'スレッドに書き込む(p2)', function(){ $daad.p2Kakiko(id.thread); });
			}
			if($daad.now_obj.name !== id.server){
				nich.appendChild(_doc.createElement('hr'));
				this.option(nich, 'お気に入りに追加', function(e){ daad.bookmark.add($daad.now_obj.url, $daad.now_obj.title); });
				nich.appendChild(_doc.createElement('hr'));
			}
		}
		this.option(nich, '板一覧表示', function(){ $daad.access(daad.id.server); });
		this.option(nich, '更新', function(){ home.update(); });
		this.option(nich, '進む', function(){ home.forward(); });
		this.option(nich, '戻る', function(){ home.back(); });
		this.$popUp.point(nich);
	},

	tab: function(key, name, e){
		var nich = this.$popUp.plus(true);
		style.change(nich, this.style);
		this.option(nich, 'お気に入りに追加', function(e){
			var line = $daad.data[name][key];
			daad.bookmark.add(line.url, line.title);
		});
		nich.appendChild(_doc.createElement('hr'));
		this.option(nich, 'タブを全て削除', function(e){ daad.tab.delAll(); });
		this.option(nich, 'このタブの左側を削除', function(e){ daad.tab.delLeft(key, name); });
		this.option(nich, 'このタブの右側を削除', function(e){ daad.tab.delRight(key, name); });
		this.option(nich, 'このタブを削除', function(e){ daad.tab.del(key, name); });
		this.$popUp.point(nich);
	},
	
	option: function(nich, title, click){
		var elm = _doc.createElement('div');
		addEvent(elm, 'mouseover', function(){
			var style = elm.style;
			style.backgroundColor = '#316AC5';
			style.color = '#FFFFFF';
		});
		addEvent(elm, 'mouseout', function(){
			var style = elm.style;
			style.backgroundColor = '#FFFFFF';
			style.color = '#000000';
		});
		addEvent(elm, 'click', click);
		elm.innerHTML = title;
		elm.style.paddingLeft = '15px';
		nich.appendChild(elm);
	},

	createPopUp: function(){
		this.$popUp = new daad.popUp('popmenu');
	}
},

style = daad.style = {
	change: function(nich, css){
		if(nich === null) return false;
		if(css === undefined){
			nich.style.backgroundColor = '#FFFFCC';
			nich.style.border = 'solid 1px black';
		} else {
			extend(nich.style, css);
		}
		return true;
	}
},

title = daad.title = {
	nameList: [],
	timeout: 604800,

	name: function(url){
		var path = url + '?name=title',
		list = this.nameList,
		len = list.length;
		if(len-- > 0) do {
			if(list[len][0] === url){
				break;
			}
		} while(len--);
		if(len < 0){
			ajax(path, this, 'ajaxCall');
		} else {
			$daad.$popUp.print(list[len][0], list[len][1]);
		}
	},
	
	ajaxCall: function(http){
		var list = this.nameList,
		data = [];
		if(http){
			// 1行目がパス、２行目がタイトル
			data = list[list.length] = http.split(LF);
			traceLog.stop();
			$daad.$popUp.print(data[0], data[1]);
		}
	}
},

search = daad.search = {
	main: function(name){
		var obj = $daad[name],
		word = $(id.search(name)).value,
		element, data, reg, sid, i, rep;
		if((obj === undefined) || !word){
			dom.uLogErr('検索失敗。');
			return false;
		}
		element = $(name);
		data = obj.text().join(emptyString);
		reg = new RegExp(word, 'g');
		sid = id.searchID(name);
		i = 0;
		rep = function(str){
			return '<span id="' + sid + (i++) + '" class="search">' + str + '</span>';
		};
		data = data.replace(regs.search, function(str, p1, p2){
			return p1.replace(reg, rep) + p2;
		});
		obj.searchMatch = i;
		element.innerHTML = data;
		if(i === 0){
			obj.searchEnabled = false;
			obj.searchNow = undefined;
			dom.uLogErr('何も見つかりませんでした。');
		} else {
			obj.searchEnabled = true;
			dom.uLog(i + '件見つかりました。');
			this.move(0, name);
		}
		return i;
	},
	
	next: function(name){
		var obj = $daad[name],
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
		this.move(num, name);
		return true;
	},
	
	back: function(name){
		var obj = $daad[name],
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
		this.move(num, name);
		return true;
	},

	move: function(num, name){
		var obj = $daad[name],
		element, top, elm;
		if((obj === undefined) || !(obj.searchEnabled)){
			dom.uLogErr('移動失敗。');
			return false;
		}
		element = $(id.searchID(name) + num),
		top = element.offsetTop,
		elm = $(name);
		if(!browser.msie){
			top -= elm.offsetTop;
		}
		elm.scrollTop = top;
		obj.searchNow = num;
		return true;
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

home = daad.home = {
	back: function(){
		history.back();
	},

	forward: function(){
		history.forward();
	},

	update: function(){
		if($daad.now_obj){
			$daad.now_obj.update();
		} else {
			dom.uLogErr('更新できません');
		}
	},

	home: function(){
		$daad.access(id.server);
	},

	history: function(){
		var elem = $(id.canvas),
		line = {},
		txt = ['<h1>履歴</h1><ul>'],
		i = 0,
		len = $daad.data.board.length;
		for(i = 0; i < len; i++){
			line = $daad.data.board[i];
			txt[txt.length] = '<li><a href="#' + line.url + '" onclick="$daad.access(\'' + line.url + '\');">' + line.title + '</a></li>';
		}
		len = $daad.data.thread.length;
		for(i = 0; i < len; i++){
			line = $daad.data.thread[i];
			txt[txt.length] = '<li><a href="#' + line.url + '" onclick="$daad.access(\'' + line.url + '\');">' + line.title + '</a></li>';
		}
		txt[txt.length] = '</ul>';
		$daad.$popUp.print('history', txt.join(emptyString));
	}
},

// 共通のメソッド
fncanvas = extend([saba, ita, sure], {
	x		: 0,
	y		: 0,
	title	: emptyString,

	scrollkeep: function(){
		var element = viewType.canvas ? $(this.name) : $(id.canvas);
		this.x = element.scrollLeft;
		this.y = element.scrollTop;
	},
	
	scroll: function(){
		var element = viewType.canvas ? $(this.name) : $(id.canvas);
		element.scrollLeft = this.x;
		element.scrollTop = this.y;
	},

	update: function(){
		ajax(this.url, this, 'ajaxCall');
	},

	print: function(){
		traceLog.start();
		this.printData(this.text());
	},

	printData: function(txt){
		var name = this.name,
		title = this.title,
		url = this.url,
		self = this;
		(function(){
			if(dom.ready){
				var element = viewType.canvas ? $(name) : $(id.canvas);
				// スクロール状態を保存する
				if($daad[name]) $daad[name].scrollkeep();
				element.innerHTML = txt.join(emptyString);
				gadgets.window.setTitle(title + ' - ' + conf.name);
				// グローバル変数にオブジェクト格納
				$daad[name] = self;
				if(!viewType.canvas){ // canvas以外
					$daad.now_obj = self;
				}
				if($daad.nich[url]){
					// タブを作成
					if(viewType.canvas) tab.print(name);
					// スクロールさせる
					self.scroll();
					self.searchEnabled = false;
				}
				// changeを動作させる
				if($daad.reading > 0) $daad.reading--;
				prefs.set(daad.id.prefs.url, url);
				traceLog.stop();
			} else {
				setTimeout(arguments.callee, 0);
				return;
			}
		})();
	}
}),

// ページロード完了時の処理
pageLoad = daad.pageLoad = function(){
	$(id.canvas).style.height = '100%';
	dom.stdout = $(id.outer);

	// popup初期化
	$daad.createPopUp();
	menu.createPopUp();
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
		if(viewType.canvas){
			menu.context();
		} else {
			menu.context2();
		}
	};

	// イベントハンドラをセットする
	addEvent(_doc, 'mousemove', mousefunc);

	// 右クリックメニュー
	if(viewType.canvas){
		addEvent($(id.server), 'contextmenu', contextevent);
		addEvent($(id.board), 'contextmenu', contextevent);
		addEvent($(id.thread), 'contextmenu', contextevent);
	} else {
		addEvent($(id.canvas), 'contextmenu', contextevent);
		addEvent($('back'), 'click', function(){ home.back(); });
		addEvent($('forward'), 'click', function(){ home.forward(); });
		addEvent($('update'), 'click', function(){ home.update(); });
		addEvent($('home'), 'click', function(){ home.home(); });
		addEvent($('history'), 'click', function(){ home.history(); });
	}

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

//$daad = window.$daad = new daad(location.hash.slice(1));
$daad = window.$daad = new daad(prefs.getString(daad.id.prefs.url));

// ページの構築が完了したらloadを呼び出す
gadgets.util.registerOnLoadHandler(pageLoad);

})();
</script>
]]></Content>
<Content type="html" view="profile,home"><![CDATA[
<style type="text/css">
<!--
	* {font-family:"ＭＳ Ｐゴシック","ＭＳ ゴシック";}
	body {
		color:black;background-color:#efefef;
		font-size:12px;line-height:115%;
	}
	h1 {color:red;font-size:13px;font-weight:400;}
	tbody {font-size:12px;line-height:115%;}
	dd {margin: 0px 0px 12px 18px;}
	#nich {width:100%;height:95%;overflow:scroll;}
	.title {font-size:13px;color:red;}
	.nich {color:green;}
	.res {color:#B02B2C;}
	.spd {color:#006E2E;}
	.popup {
		font-size: 10px;line-height:105%;
		border-style: solid;border-color: #000000;border-width: 1px;
		padding: 3px;filter: alpha(opacity=90);
		-moz-opacity:0.9;opacity:0.9;white-space:pre;
	}
	.popup dl {margin:0px;}
	.popup dt {margin:0px;}
	.popup dd {margin:0px 0px 7px 10px;}
	.saba {display:none;margin-left:5px;}
	.line-color {background-color:#F5F5F5;}
	a:link,a:visited {color:blue;}
	a:active {color:red;}
	a:hover {color:#660099;}
	a:link.tan,a:visited.tan,a:active.tan,a:hover.tan {color:#333333;text-decoration: underline;}
	a:link.makka,a:visited.makka,a:active.makka,a:hover.makka {color:#FF0000;text-decoration: underline;}
	a:link.ninki,a:visited.ninki,a:active.ninki,a:hover.ninki {color:#AF00CF;text-decoration: underline;}
-->
</style>
<div id="menu">
<img src="http://file.unkar.org/img/unkar/arrow-left.gif" alt="戻る" width="16" height="16" id="back" />
<img src="http://file.unkar.org/img/unkar/arrow-right.gif" alt="進む" width="16" height="16" id="forward" />
<img src="http://file.unkar.org/img/unkar/refresh.gif" alt="更新" width="16" height="16" id="update" />
<img src="http://file.unkar.org/img/unkar/home.gif" alt="HOME" width="16" height="16" id="home" />
<img src="http://file.unkar.org/img/unkar/clipboard.gif" alt="履歴" width="16" height="16" id="history" />
<span id="outer"></span>
</div>
<div id="nich"></div>
]]></Content>
<Content type="html" view="canvas"><![CDATA[
<style type="text/css">
<!--
	* {font-family:"ＭＳ Ｐゴシック","ＭＳ ゴシック";}
	body {
		color:#000000;font-size:14px;
		line-height:120%;margin:0;padding:0;
	}
	#menu {width:100%;height:auto;margin:0;padding:0;}
	#server {
		font-size:14px;float:left;width:15%;height:100%;
		vertical-align:top;overflow:scroll;
	}
	#board {
		float:left;width:85%;height:30%;overflow:scroll;
	}
	#board-tab {
		float:left;line-height:13px;width:85%;
		height:auto;white-space:nowrap;font-size:12px;
	}
	#board-menu {
		float:left;width:85%;height:18px;
		white-space:nowrap;font-size:12px;
	}
	#board-search {border: solid 1px #000000;}
	#board table {
		width:100%;font-size:14px;
		overflow:hidden;white-space:nowrap;
	}
	.line-color {background-color:#F5F5F5;}
	#board th {font-size:12px;background-color:#F0F0E8;}
	#thread {
		float:left;background-color:#EFEFEF;
		font-size:16px;line-height:115%;
		width:85%;height:60%;overflow:scroll;
	}
	#thread-tab {
		float:left;line-height:13px;background-color:#EFEFEF;
		font-size:12px;width:85%;height:auto;white-space:nowrap;
	}
	#thread-menu {
		float:left;width:85%;height:18px;
		white-space:nowrap;font-size:12px;
	}
	#thread-search {border: solid 1px #000000;}
	#popup {white-space:pre;}
	.nich {color:green;}
	.res {color:#B02B2C;font-weight:700;}
	.sin {color:#3F4C6B;}
	.spd {color:#006E2E;font-weight:700;}
	.title {font-size:18px;color:red;}
	.navi {color:#3300FF;text-decoration:underline;}
	.search {background-color:#FFFF00;}
	.saba {display:none;margin-left:20px;}
	h1 {
		color:red;font-size:18px;font-weight:400;
		margin:10px 0px 14px 3px;
	}
	h2 {
		font-size:16px;font-weight:400;
		margin:0;padding:0;
	}
	form, input {
		display:inline;vertical-align:top;
		margin:0;padding:0;
	}
	b {font-weight:700;}
	dd {margin:0 0 20px 32px;}
	dl {margin:0 0 0 5px;}
	ul {list-style:none;margin:0;padding:0;}
	li {
		background-color:#EBE5C3;
		float:left;width:auto;
		margin:0px 2px 2px 2px;padding:2px;
	}
	li.now {background-color:#FFEEA8;}
	#thread-tab li a,
	#board-tab li a {color:#000000;text-decoration:none;}
	li a:hover {background-color:#E0EFF9;}
	a:link,a:visited {color:blue;}
	a:link.makka,a:visited.makka,a:active.makka,a:hover.makka {color:red;text-decoration:underline;}
	a:link.ninki,a:visited.ninki,a:active.ninki,a:hover.ninki {color:#AF00CF;text-decoration:underline;}
	a:link.tan,a:visited.tan,a:active.tan,a:hover.tan {color:#333;text-decoration:underline;}
	.popup dl,.popup dt,dt {margin:0;}
-->
</style>
<div id="menu">
<img src="http://file.unkar.org/img/unkar/info.gif" alt="メニュー" width="16" height="16" onclick="daad.menu.contextmenu();">　
<img src="http://file.unkar.org/img/unkar/trash.gif" alt="タブを全て削除" width="16" height="16" onclick="daad.tab.delAll();">　
<span id="outer"></span>
</div>
<div id="nich">
<div id="server"></div>
<div id="board-menu">
<form onsubmit="daad.search.main(daad.id.board);return false;">
<input type="text" id="board-search" size="40" onkeyup="daad.search.main(daad.id.board);">　
</form>
<img src="http://file.unkar.org/img/unkar/zoom-in.gif" alt="↓" width="16" height="16" onclick="daad.search.next(daad.id.board);">　
<img src="http://file.unkar.org/img/unkar/zoom-out.gif" alt="↑" width="16" height="16" onclick="daad.search.back(daad.id.board);">　
<img src="http://file.unkar.org/img/unkar/refresh.gif" alt="スレッド一覧更新" width="16" height="16" onclick="$daad.renew(daad.id.board);">　
<img src="http://file.unkar.org/img/unkar/edit.gif" alt="スレ立て" width="16" height="16" onclick="$daad.p2Kakiko(daad.id.board);">　
</div>
<div id="board-tab"></div>
<div id="board"></div>
<div id="thread-menu">
<form onsubmit="daad.search.main(daad.id.thread);return false;">
<input type="text" id="thread-search" size="40" onkeyup="daad.search.main(daad.id.thread);">　
</form>
<img src="http://file.unkar.org/img/unkar/zoom-in.gif" alt="↓" width="16" height="16" onclick="daad.search.next(daad.id.thread);">　
<img src="http://file.unkar.org/img/unkar/zoom-out.gif" alt="↑" width="16" height="16" onclick="daad.search.back(daad.id.thread);">　
<img src="http://file.unkar.org/img/unkar/refresh.gif" alt="スレッド一覧更新" width="16" height="16" onclick="$daad.renew(daad.id.thread);">　
<img src="http://file.unkar.org/img/unkar/edit.gif" alt="書き込み" width="16" height="16" onclick="$daad.p2Kakiko(daad.id.thread);">　
</div>
<div id="thread-tab"></div>
<div id="thread"></div>
</div>
]]></Content>
</Module>
