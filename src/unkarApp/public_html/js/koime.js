(function(){

var window = this,
LF = String.fromCharCode(10),
undefined,
emptyString = '',
_doc = document,

koime = window.koime = function(url){
	this.nich = {};
	this.now_obj = undefined;
	this.data = {
		server	: {},
		board	: [],
		thread	: []
	};
	this.access(url);
},

conf = koime.conf = {
	koime		: '0.9.0',
	name		: 'koime',
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

mouse = koime.mouse = {
	x			: 0,
	y			: 0
},

regs = koime.regs = {
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

img = koime.img = {
	load			: 'http://file.unkar.org/img/unkar/loading.gif',
	load2			: 'http://file.unkar.org/img/loading.gif',
	folder			: 'http://file.unkar.org/img/unkar/folder.png',
	cross			: 'http://file.unkar.org/img/unkar/cross.png',
	arrow_down		: 'http://file.unkar.org/img/unkar/arrow_down.png',
	arrow_up		: 'http://file.unkar.org/img/unkar/arrow_up.png',
	arrow_left		: 'http://file.unkar.org/img/unkar/arrow_left.png',
	arrow_right		: 'http://file.unkar.org/img/unkar/arrow_right.png',
	comment_add		: 'http://file.unkar.org/img/unkar/comment_add.png',
	comment_edit	: 'http://file.unkar.org/img/unkar/comment_edit.png',
	file			: 'http://file.unkar.org/img/unkar/file.png',
	info			: 'http://file.unkar.org/img/unkar/info.png',
	magnify			: 'http://file.unkar.org/img/unkar/magnify.png',
	star			: 'http://file.unkar.org/img/unkar/star.png',
	trash			: 'http://file.unkar.org/img/unkar/trash.png',
	file_delete		: 'http://file.unkar.org/img/unkar/file_delete.png',
	download		: 'http://file.unkar.org/img/unkar/download.png',
	italist			: 'http://file.unkar.org/img/unkar/italist.png',
	attach			: 'http://file.unkar.org/img/unkar/attach.png',
	disc			: 'http://file.unkar.org/img/unkar/disc.png'
},

id = koime.id = {
	canvas		: 'nich',
	server		: 'server',
	board		: 'board',
	thread		: 'thread',
	popup		: 'popup',
	outer		: 'outer',
	prefs		: {
		url			: 'reboot_url'
	},
	menu		: function(str){ return str + '-menu'; },
	search		: function(str){ return str + '-search'; },
	searchNext	: function(str){ return str + '-search-next'; },
	searchBack	: function(str){ return str + '-search-back'; },
	searchID	: function(id) { return 'sNo-' + id + '-'; }
},

dom = koime.dom = {
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
browser = koime.browser = {
	safari		: (UA.indexOf('webkit') !== -1),
	opera		: (UA.indexOf('opera') !== -1),
	msie		: (UA.indexOf('msie') !== -1) && (UA.indexOf('opera') === -1),
	mozilla		: (UA.indexOf('mozilla') !== -1) && !/(compatible|webkit)/.test(UA)
},

addEvent = koime.addEvent = (function(){
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

stopEvent = koime.stopEvent = function(e){
	if(e.stopPropagation){
		e.stopPropagation();
	} else {
		e.cancelBubble = true;
	}
},

extend = koime.extend = function(dest, source){
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

addDiv = koime.addDiv = function(canvas, id, klass){
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

uniq = koime.uniq = function(array){
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

ajax = unkar.ajax = function(path, self, func){
	var xml = null,
	e = emptyString,
	timeout, lastmod, timerID, e,
	url = conf.convert_url + '/' + path;
	// 計測開始
	traceLog.start();
	try {
		xml = new XMLHttpRequest();
	} catch(e){
		try {
			xml = new ActiveXObject("Msxml2.XMLHTTP");
		} catch(e){
			try {
				xml = new ActiveXObject("Microsoft.XMLHTTP");
			} catch(e){
				xml = null;
			}
		}
	}
	if(xml){
		timeout = function(){
			xml.abort();
			dom.uLogErr('接続がタイムアウトしました');
		};
		lastmod = (self.lastModified || 'Mon, 26 Jul 1997 05:00:00 GMT');
		timerID = setTimeout(timeout, conf.timeout);
		xml.onreadystatechange = function(){
			if(xml.readyState === 4){
				if(xml.status === 200){
					clearTimeout(timerID);
					traceLog.load();
					// 最終更新時間の取得
					self.lastModified = xml.getResponseHeader('Last-Modified');
					self[func](xml);
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
fn = koime.prototype = {
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
		location.hash = '#' + url;
	},

	change: function(url){
		var obj, line;
		if(this.reading > 0) return;
		if(line = url.match(regs.url_split)){
			obj = this.now_obj;
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
		this.$popUp = new koime.popUp(id.popup);
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
				txt[txt.length] = '<img src="' + img.folder + '" alt="ロード中" width="16" height="16"><a href="#' + line[0] + '" onclick="$koime.access(\'' + line[0] + '\');">' + line[1] + '</a><br>';
			} else {
				txt[txt.length] = '</div><img src="' + img.folder + '" alt="ロード中" width="16" height="16"><a href="javascript:$koime.server.itaView(\'s' + k + '\');" class="tan">' + list[i] + '</a><br><div id="s' + k + '" class="saba">';
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
			data[i].spd = (86400 / ((now - data[i].sin) / data[i].res)) | 0; // 整数変換
			txt[i] = this.style(i);
		}
		this.printData(this.header(txt));
	},

	header: function(txt){
		var start = ['<table><tbody><tr><th><a href="#' + this.url + '" onclick="$koime.board.sort(\'num\', ' + this.sortflag + ');" class="tan" >No</a></th><th>title</th><th><a href="#' + this.url + '" onclick="$koime.board.sort(\'res\', ' + this.sortflag + ');" class="tan" >res</a></th><th><a href="#' + this.url + '" onclick="$koime.board.sort(\'spd\', ' + this.sortflag + ');" class="tan">res/day</a></th><th><a href="#' + this.url + '" onclick="$koime.board.sort(\'sin\', ' + this.sortflag + ');" class="tan">since</a></th></tr>'];
		txt = start.concat(txt);
		txt[txt.length] = '</tbody></table>';
		return txt;
	},

	style: function(i){
		return '<tr' + ((i % 2) ? ' class="line-color"' : emptyString) + '><td>' + this.data[i].num + '</td><td><a href="#' + this.url + '/' + this.data[i].sin + '" onclick="$koime.access(\'' + this.url + '/' + this.data[i].sin + '\');">' + this.data[i].thread + '</a></td><td class="res">' + this.data[i].res + '</td><td class="spd">' + this.data[i].spd + '</td><td class="sin">' + this.data[i].since + '</td></tr>';
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
		var color = (this.anker[i] === undefined) ? emptyString : ' onmouseover="$koime.thread.resPop(' + i + ');" ' + ((this.anker[i].length < 3) ? 'class="ninki"' : 'class="makka"'),
		res = this.res[i];
		return '<dt><a id="l'+i+'" href="#l'+i+'"'+color+'>'+i+'</a> ：<span class="nich"><b>'
				+ res[0] + '</b></span>[' + res[1] + ']：' + res[2] + '</dt><dd>' + res[3] + '</dd>';
	},
	
	anchorPop: function(i){
		if(this.res[i] === undefined) return emptyString;
		$koime.$popUp.print('a' + i, '<dl>' + this.ankerStyle(i) + '</dl>');
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
		$koime.$popUp.print('c' + num[0], this.resChain(num));
	},

	resPop: function(i){
		$koime.$popUp.print('r' + i, this.resTree(i));
	},
	
	idPop: function(id){
		$koime.$popUp.print(id, this.idTree(id));
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
				return '<a href="#l'+p3+'" onmouseover="$koime.thread.split(\''+p3+p4+'\');">'+p1+p3+p4+'</a>';
			}
			return '<a href="#l'+p3+'" onmouseover="$koime.thread.anchorPop('+p3+');">'+p1+p3+'</a>';
		},
		urlLink = function(str, p1, p2){
			var line, regexp;
			if(regexp = p2.match(reg.sure)){
				// スレッドだった場合
				line = regexp[1] + '/' + regexp[2];
				return '<a href="#' + line + '" onclick="$koime.access(\'' + line + '\');" onmouseover="koime.title.name(\'' + line + '\');">' + str + '</a>';
			} else if(regexp = p2.match(reg.ita)){
				// 板だった場合
				if(in_array(regexp[1], conf.filter) !== -1){
					if((p1.indexOf('h') === -1)){
						return '<a href="h' + str + '" target="_blank">' + str + '</a>';
					}
					return '<a href="' + str + '" target="_blank">' + str + '</a>';
				}
				line = regexp[1] + regexp[2];
				return '<a href="#' + line + '" onclick="$koime.access(\'' + line + '\');" onmouseover="koime.title.name(\'' + line + '\');">' + str + '</a>';
			} else if(p1.indexOf('h') === -1){
				// hが付いていなかった場合
				return '<a href="h' + str + '" target="_blank">' + str + '</a>';
			} else {
				return '<a href="' + str + '" target="_blank">' + str + '</a>';
			}
		},
		idColor = function(str, p2, p1){
			var color = ((id[p1].length >= 5) ? ' class="makka"' : (id[p1].length > 1) ? emptyString : ' class="tan"'),
			text = '<a href="#l' + i + '" class="tan" onmouseover="$koime.$popUp.print(\'date' + i + '\', \'' + p2 + '\')">date</a> <a href="#l' + i + '"' + color + ' onmouseover="$koime.thread.idPop(\'' + p1 + '\');">ID:</a>' + p1;
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
			res[i][3] = res[i][3].replace(reg.id2, '<a href="#l' + i + '" class="tan" onmouseover="$koime.thread.idPop(\'$1\');">ID:</a>$1')
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

tab = koime.tab = {
	print: function(name){
		var list = $koime.data[name],
		length = list.length,
		url = emptyString,
		_url = $koime[name].url,
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
			tmp = '<li ondblclick="koime.menu.tab(' + i + ', \'' + name + '\');" oncontextmenu="koime.menu.tab(' + i + ', \'' + name + '\', event);"';
			if(url === _url){
				txt[txt.length] = tmp + ' class="now"><a href="#' + url + '" title="' + obj.title + '">' + title + '</a>';
			} else {
				txt[txt.length] = tmp + '><a href="#' + url + '" title="' + obj.title + '" onclick="$koime.access(\'' + url + '\');">' + title + '</a>';
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
		$koime.init();
		$koime.access(id.server);
	},

	delLeft: function(key, name){
		var list = $koime.data[name],
		url = emptyString,
		i = 0;
		// 先に表示を変える
		$koime.hash(list[key].url);
		list[key].print();
		for(; i < key; i++){
			url = list[i].url;
			delete $koime.nich[url];
		}
		list.splice(0, key);
		this.print(name);
	},

	delRight: function(key, name){
		var list = $koime.data[name],
		length = list.length,
		url = emptyString,
		i = key + 1;
		// 先に表示を変える
		$koime.hash(list[key].url);
		list[key].print();
		for(; i < length; i++){
			url = list[i].url;
			delete $koime.nich[url];
		}
		list.splice(key + 1, length - key);
		this.print(name);
	},

	del: function(key, name){
		var obj = $koime.data[name],
		url = obj[key].url,
		_url = $koime[name].url;
		obj.splice(key, 1);
		delete $koime.nich[url];
		this.print(name);
		if(!obj[0]){
			$koime[name] = undefined;
			$(name).innerHTML = emptyString;
			$koime.hash(emptyString);
		} else {
			if((url === _url) && obj[--key]){
				$koime.hash(obj[key].url);
				obj[key].print();
			}
		}
	}
},

popUp = koime.popUp = function(id){
	this.rootID = id;
	addDiv(_doc.body, id, 'popup');
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

menu = koime.menu = {
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
		if($koime.now_obj){
			if($koime.now_obj.name === id.board){
				this.option(nich, 'スレッドを立てる(p2)', function(){ $koime.p2Kakiko(id.board); });
			} else if($koime.now_obj.name === id.thread){
				this.option(nich, 'スレッドに書き込む(p2)', function(){ $koime.p2Kakiko(id.thread); });
			}
			if($koime.now_obj.name !== id.server){
				nich.appendChild(_doc.createElement('hr'));
				this.option(nich, 'お気に入りに追加', function(e){ koime.bookmark.add($koime.now_obj.url, $koime.now_obj.title); });
				nich.appendChild(_doc.createElement('hr'));
			}
		}
		this.option(nich, '板一覧表示', function(){ $koime.access(koime.id.server); });
		this.option(nich, '更新', function(){ home.update(); });
		this.option(nich, '進む', function(){ home.forward(); });
		this.option(nich, '戻る', function(){ home.back(); });
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
		this.$popUp = new koime.popUp('popmenu');
	}
},

style = koime.style = {
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

title = koime.title = {
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
			$koime.$popUp.print(list[len][0], list[len][1]);
		}
	},
	
	ajaxCall: function(http){
		var list = this.nameList,
		data = [];
		if(http){
			// 1行目がパス、２行目がタイトル
			data = list[list.length] = http.split(LF);
			traceLog.stop();
			$koime.$popUp.print(data[0], data[1]);
		}
	}
},

search = koime.search = {
	main: function(name){
		var obj = $koime[name],
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
		var obj = $koime[name],
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
		var obj = $koime[name],
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
		var obj = $koime[name],
		element, top, elm;
		if((obj === undefined) || !(obj.searchEnabled)){
			dom.uLogErr('移動失敗。');
			return false;
		}
		element = $(id.searchID(name) + num),
		top = element.offsetTop,
		elm = $(name);
		top -= elm.offsetTop;
		elm.scrollTop = top;
		obj.searchNow = num;
		return true;
	}
},

traceLog = koime.traceLog = {
	startTime: 0,
	loadTime: 0,

	start: function(){
		var txt = '<img src="' + img.load2 + '" alt="ロード中" width="16" height="16">ロード中';
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

home = koime.home = {
	back: function(){
		history.back();
	},

	forward: function(){
		history.forward();
	},

	update: function(){
		if($koime.now_obj){
			$koime.now_obj.update();
		} else {
			dom.uLogErr('更新できません');
		}
	},

	home: function(){
		$koime.access(id.server);
	},

	history: function(){
		var elem = $(id.canvas),
		line = {},
		txt = ['<h1>履歴</h1><ul>'],
		i = 0,
		len = $koime.data.board.length;
		for(i = 0; i < len; i++){
			line = $koime.data.board[i];
			txt[txt.length] = '<li><a href="#' + line.url + '" onclick="$koime.access(\'' + line.url + '\');">' + line.title + '</a></li>';
		}
		len = $koime.data.thread.length;
		for(i = 0; i < len; i++){
			line = $koime.data.thread[i];
			txt[txt.length] = '<li><a href="#' + line.url + '" onclick="$koime.access(\'' + line.url + '\');">' + line.title + '</a></li>';
		}
		txt[txt.length] = '</ul>';
		$koime.$popUp.print('history', txt.join(emptyString));
	}
},

// 共通のメソッド
fncanvas = extend([saba, ita, sure], {
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
				var element = $(id.canvas);
				// スクロール状態を保存する
				if($koime[name]) $koime[name].scrollkeep();
				element.innerHTML = txt.join(emptyString);
				_doc.title = title + ' - ' + conf.name;
				// グローバル変数にオブジェクト格納
				$koime[name] = self;
				$koime.now_obj = self;
				if($koime.nich[url]){
					// スクロールさせる
					self.scroll();
					self.searchEnabled = false;
				}
				// changeを動作させる
				if($koime.reading > 0) $koime.reading--;
				traceLog.stop();
			} else {
				setTimeout(arguments.callee, 0);
				return;
			}
		})();
	}
}),

// ページロード完了時の処理
pageLoad = koime.pageLoad = function(){
	$(id.canvas).style.height = '100%';
	dom.stdout = $(id.outer);

	// popup初期化
	$koime.createPopUp();
	menu.createPopUp();
	// 構築完了
	dom.ready = true;

	var mousefunc = function(){
		var mouse = koime.mouse;
		mouse.x = e.pageX;
		mouse.y = e.pageY;
	},
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
	setInterval(function(){
		$koime.change(location.hash.slice(1));
	}, conf.change);
	return true;
},

$koime = window.$koime = new koime(location.hash.slice(1));

// ページの構築が完了したらloadを呼び出す
if(_doc.addEventListener){ // IE以外
	addEvent(_doc, 'DOMContentLoaded', pageLoad);
} else { // その他
	addEvent(window, 'load', pageLoad);
}

})();









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
<img src="http://file.unkar.org/img/modoru.png" alt="戻る" id="back" />
<img src="http://file.unkar.org/img/susumu.png" alt="進む" id="forward" />
<img src="http://file.unkar.org/img/kousin.png" alt="更新" id="update" />
<img src="http://file.unkar.org/img/ita.png" alt="HOME" id="home" />
<img src="http://file.unkar.org/img/rireki.png" alt="履歴" id="history" />
<span id="outer"></span>
</div>
<div id="nich"></div>
